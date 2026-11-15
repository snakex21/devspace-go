package tools

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

// extraRootFixture adds a read-only root "shared" and a read-write root "rw"
// next to the workspace from pathFixture, and returns a context carrying a
// guard for origin.
func extraRootFixture(t *testing.T, origin workspace.Origin, remote bool) (ctx context.Context, root, outside, shared, rw string) {
	t.Helper()
	root, outside = pathFixture(t)
	parent := filepath.Dir(root)
	shared = filepath.Join(parent, "shared")
	rw = filepath.Join(parent, "rw")
	files := map[string]string{
		filepath.Join(shared, "lib.go"):           "package lib // SHARED-MARKER",
		filepath.Join(shared, "sub", "util.go"):   "package sub // SHARED-MARKER",
		filepath.Join(shared, ".ssh", "id_rsa"):   "SSH-KEY-MARKER",
		filepath.Join(shared, ".git", "config"):   "url=https://TOKEN-MARKER@host",
		filepath.Join(shared, ".git", "HEAD"):     "ref: refs/heads/main",
		filepath.Join(shared, ".git-credentials"): "https://u:TOKEN-MARKER@host",
		filepath.Join(rw, "notes.txt"):            "notes",
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cfg := config.DefaultConfig()
	cfg.ConfigDir = filepath.Join(parent, "app", ".devspace")
	cfg.StateDir = filepath.Join(parent, "app", ".devspace-state")
	if err := os.MkdirAll(cfg.StateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	cfg.OutsideAccess = config.OutsideAccessConfig{
		Roots: []config.ExtraRootConfig{
			{Path: shared, Access: "read"},
			{Path: rw, Access: "readwrite"},
		},
		Remote: &remote,
	}
	policy, problems := workspace.NewAccessPolicy(cfg)
	if len(problems) > 0 {
		t.Fatalf("policy problems: %v", problems)
	}
	return WithGuard(context.Background(), policy.Guard(origin, "test")), root, outside, shared, rw
}

func TestToolsUseReadOnlyExtraRoot(t *testing.T) {
	ctx, root, _, shared, _ := extraRootFixture(t, workspace.OriginLocal, false)
	lib := filepath.Join(shared, "lib.go")

	r, _, _ := ReadFile(ctx, nil, ReadInput{Path: lib}, root)
	if r.IsError || !strings.Contains(resultText(r), "SHARED-MARKER") {
		t.Fatalf("read from read-only root: %s", resultText(r))
	}
	r, _, _ = ReadFile(ctx, nil, ReadInput{Path: "../shared/lib.go"}, root)
	if r.IsError {
		t.Fatalf("relative read from read-only root: %s", resultText(r))
	}

	ls, _, _ := ListDirectory(ctx, nil, LsInput{Path: shared}, root)
	text := resultText(ls)
	if ls.IsError || !strings.Contains(text, "lib.go") || !strings.Contains(text, ".git/") {
		t.Fatalf("ls: %s", text)
	}
	if strings.Contains(text, ".ssh") || strings.Contains(text, ".git-credentials") {
		t.Fatalf("ls shows protected entries: %s", text)
	}

	grep, _, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "MARKER", Path: shared}, root)
	text = resultText(grep)
	if grep.IsError || !strings.Contains(text, "SHARED-MARKER") {
		t.Fatalf("grep: %s", text)
	}
	if strings.Contains(text, "SSH-KEY-MARKER") || strings.Contains(text, "TOKEN-MARKER") {
		t.Fatalf("grep leaked protected files: %s", text)
	}
	if !strings.Contains(strings.ToLower(text), strings.ToLower(filepath.ToSlash(lib))) {
		t.Fatalf("grep should show absolute paths outside the workspace: %s", text)
	}

	glob, _, _ := FindFiles(ctx, nil, GlobInput{Pattern: "*", Path: shared}, root)
	text = resultText(glob)
	if glob.IsError || !strings.Contains(text, "util.go") {
		t.Fatalf("glob: %s", text)
	}
	if strings.Contains(text, "id_rsa") || strings.Contains(text, "credentials") || strings.Contains(text, "/config") {
		t.Fatalf("glob leaked protected files: %s", text)
	}

	// Every writing tool is refused, and nothing changes on disk.
	w, _, _ := WriteFile(ctx, nil, WriteInput{Path: lib, Content: "pwned"}, root)
	requireToolError(t, "write", w)
	e, _, _ := EditFile(ctx, nil, EditInput{Path: lib, Edits: []EditBlock{{OldText: "lib", NewText: "pwned"}}}, root)
	requireToolError(t, "edit", e)
	m, _, _ := MakeDirectory(ctx, nil, MkdirInput{Path: filepath.Join(shared, "made")}, root)
	requireToolError(t, "mkdir", m)
	mv, _, _ := MovePath(ctx, nil, MoveInput{SourcePath: lib, TargetPath: "stolen.go"}, root)
	requireToolError(t, "move out of read-only root", mv)
	mv, _, _ = MovePath(ctx, nil, MoveInput{SourcePath: "inside.txt", TargetPath: filepath.Join(shared, "x.txt")}, root)
	requireToolError(t, "move into read-only root", mv)
	b, _, _ := RunBash(ctx, nil, BashInput{Command: "echo hi", WorkingDirectory: shared}, root)
	requireToolError(t, "bash cwd in read-only root", b)
	if data, _ := os.ReadFile(lib); !strings.Contains(string(data), "SHARED-MARKER") {
		t.Fatalf("read-only file changed: %q", data)
	}
	if _, err := os.Stat(filepath.Join(shared, "made")); !os.IsNotExist(err) {
		t.Fatalf("mkdir created a directory in a read-only root: %v", err)
	}

	for _, protected := range []string{filepath.Join(shared, ".ssh", "id_rsa"), filepath.Join(shared, ".git", "config"), filepath.Join(shared, ".git-credentials")} {
		r, _, _ := ReadFile(ctx, nil, ReadInput{Path: protected}, root)
		if !r.IsError || strings.Contains(resultText(r), "MARKER") {
			t.Fatalf("read of protected %s: %s", protected, resultText(r))
		}
	}
}

func TestToolsUseReadWriteExtraRoot(t *testing.T) {
	ctx, root, _, _, rw := extraRootFixture(t, workspace.OriginLocal, false)
	target := filepath.Join(rw, "dir", "new.txt")

	w, _, _ := WriteFile(ctx, nil, WriteInput{Path: target, Content: "hello"}, root)
	if w.IsError {
		t.Fatalf("write: %s", resultText(w))
	}
	e, _, _ := EditFile(ctx, nil, EditInput{Path: target, Edits: []EditBlock{{OldText: "hello", NewText: "bye"}}}, root)
	if e.IsError {
		t.Fatalf("edit: %s", resultText(e))
	}
	mv, _, _ := MovePath(ctx, nil, MoveInput{SourcePath: target, TargetPath: filepath.Join(rw, "moved.txt")}, root)
	if mv.IsError {
		t.Fatalf("move: %s", resultText(mv))
	}
	if data, err := os.ReadFile(filepath.Join(rw, "moved.txt")); err != nil || string(data) != "bye" {
		t.Fatalf("moved file: %q %v", data, err)
	}
	mv, _, _ = MovePath(ctx, nil, MoveInput{SourcePath: rw, TargetPath: "rw-stolen"}, root)
	requireToolError(t, "move of the extra root itself", mv)

	command := "pwd"
	if runtime.GOOS == "windows" {
		command = "(Get-Location).Path"
	}
	b, _, _ := RunBash(ctx, nil, BashInput{Command: command, WorkingDirectory: rw}, root)
	if b.IsError || !strings.Contains(strings.ToLower(resultText(b)), strings.ToLower(filepath.Base(rw))) {
		t.Fatalf("bash in read-write root: %s", resultText(b))
	}
}

func TestToolsRefuseExtraRootsForRemoteRequests(t *testing.T) {
	ctx, root, _, shared, rw := extraRootFixture(t, workspace.OriginRemote, false)
	r, _, _ := ReadFile(ctx, nil, ReadInput{Path: filepath.Join(shared, "lib.go")}, root)
	requireToolError(t, "remote read", r)
	if !strings.Contains(resultText(r), "tunnel") {
		t.Fatalf("remote denial should explain the tunnel rule: %s", resultText(r))
	}
	w, _, _ := WriteFile(ctx, nil, WriteInput{Path: filepath.Join(rw, "x.txt"), Content: "x"}, root)
	requireToolError(t, "remote write", w)
	// The workspace itself works as before.
	r, _, _ = ReadFile(ctx, nil, ReadInput{Path: "inside.txt"}, root)
	if r.IsError {
		t.Fatalf("remote workspace read: %s", resultText(r))
	}

	ctx, root, _, shared, _ = extraRootFixture(t, workspace.OriginRemote, true)
	r, _, _ = ReadFile(ctx, nil, ReadInput{Path: filepath.Join(shared, "lib.go")}, root)
	if r.IsError {
		t.Fatalf("remote read with outsideAccess.remote: %s", resultText(r))
	}
}

func TestToolsStillRejectPathsOutsideEveryRoot(t *testing.T) {
	ctx, root, outside, _, _ := extraRootFixture(t, workspace.OriginLocal, false)
	for _, escape := range []string{"../outside/secret.txt", filepath.Join(outside, "secret.txt"), "../shared/../outside/secret.txt"} {
		r, _, _ := ReadFile(ctx, nil, ReadInput{Path: escape}, root)
		requireToolError(t, "read "+escape, r)
		g, _, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "TOP-SECRET", Path: filepath.Dir(escape)}, root)
		requireToolError(t, "grep "+escape, g)
	}
}
