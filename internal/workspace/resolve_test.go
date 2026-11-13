package workspace

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveInRootAcceptsPathsInsideRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src", "pkg"), 0o700); err != nil {
		t.Fatal(err)
	}
	for input, want := range map[string]string{
		"":                         root,
		".":                        root,
		"src":                      filepath.Join(root, "src"),
		"src/pkg/../pkg/file.go":   filepath.Join(root, "src", "pkg", "file.go"),
		"new/dir/file.txt":         filepath.Join(root, "new", "dir", "file.txt"),
		filepath.Join(root, "src"): filepath.Join(root, "src"),
	} {
		got, err := ResolveInRoot(root, input)
		if err != nil {
			t.Errorf("ResolveInRoot(%q) error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("ResolveInRoot(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestResolveInRootRejectsTraversal(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "ws")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	inputs := []string{
		"..",
		"../",
		"../secret.txt",
		"src/../../secret.txt",
		filepath.Join(parent, "secret.txt"),
		parent,
		root + "-sibling",
		"a\x00b",
	}
	if runtime.GOOS == "windows" {
		inputs = append(inputs, `\Windows`, "/Windows", "C:Windows", `\\server\share\x`)
	}
	for _, input := range inputs {
		if got, err := ResolveInRoot(root, input); err == nil {
			t.Errorf("ResolveInRoot(%q) = %q, want error", input, got)
		}
	}
}

func TestResolveInRootRejectsLinkEscape(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "ws")
	outside := filepath.Join(parent, "outside")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	makeDirLink(t, outside, filepath.Join(root, "escape"))
	for _, input := range []string{"escape", "escape/secret.txt", "escape/new/file.txt"} {
		if got, err := ResolveInRoot(root, input); err == nil || !strings.Contains(err.Error(), "link") {
			t.Errorf("ResolveInRoot(%q) = %q, %v; want link escape error", input, got, err)
		}
	}

	// A link that stays inside the root remains usable.
	if err := os.Mkdir(filepath.Join(root, "real"), 0o700); err != nil {
		t.Fatal(err)
	}
	makeDirLink(t, filepath.Join(root, "real"), filepath.Join(root, "alias"))
	if _, err := ResolveInRoot(root, "alias/file.txt"); err != nil {
		t.Errorf("internal link rejected: %v", err)
	}
}

func TestResolveInRootAllowsRootReachedThroughLink(t *testing.T) {
	parent := t.TempDir()
	real := filepath.Join(parent, "real")
	if err := os.Mkdir(real, 0o700); err != nil {
		t.Fatal(err)
	}
	linkedRoot := filepath.Join(parent, "linked")
	makeDirLink(t, real, linkedRoot)
	if _, err := ResolveInRoot(linkedRoot, "file.txt"); err != nil {
		t.Fatalf("workspace opened through a link rejected its own files: %v", err)
	}
}

func TestIsPathInsideRootCaseSensitivityMatchesPlatform(t *testing.T) {
	root := filepath.FromSlash("/srv/app")
	upper := filepath.FromSlash("/srv/APP/file")
	if got := IsPathInsideRoot(upper, root); got != pathsAreCaseInsensitive {
		t.Fatalf("IsPathInsideRoot(%q, %q) = %v on %s", upper, root, got, runtime.GOOS)
	}
	if IsPathInsideRoot(filepath.FromSlash("/srv/application"), root) {
		t.Fatal("sibling with shared prefix treated as inside root")
	}
}

// makeDirLink creates a directory symlink, falling back to an NTFS junction on
// Windows where symlinks need Developer Mode or elevation. Junctions are the
// link type an unprivileged process can actually create there.
func makeDirLink(t *testing.T, target, link string) {
	t.Helper()
	symlinkErr := os.Symlink(target, link)
	if symlinkErr == nil {
		return
	}
	if runtime.GOOS == "windows" {
		if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err == nil {
			return
		} else {
			t.Skipf("cannot create symlink (%v) or junction (%v: %s)", symlinkErr, err, output)
		}
	}
	t.Skipf("cannot create symlink: %v", symlinkErr)
}
