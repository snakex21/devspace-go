package shells

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fakeSystem is a machine description for detection tests. Nothing here
// touches the real filesystem, so the tests assert the same thing on every
// runner regardless of what is installed on it.
type fakeSystem struct {
	goos  string
	path  map[string]string
	files []string
	env   map[string]string
}

func (f fakeSystem) prober() prober {
	return prober{
		lookPath: func(exe string) (string, error) {
			if path, ok := f.path[exe]; ok {
				return path, nil
			}
			return "", errors.New("not found in PATH")
		},
		exists: func(path string) bool {
			for _, file := range f.files {
				if samePath(file, path) {
					return true
				}
			}
			return false
		},
		getenv: func(key string) string { return f.env[key] },
		goos:   f.goos,
	}
}

func windowsEnv() map[string]string {
	return map[string]string{
		"SystemRoot":   `C:\Windows`,
		"ProgramFiles": `C:\Program Files`,
		"LOCALAPPDATA": `C:\Users\dev\AppData\Local`,
	}
}

// samePath compares paths the way the fixtures mean them. Windows fixtures are
// written with backslashes, but filepath.Join on a POSIX runner joins with a
// forward slash, so a literal comparison would miss a file that is really there.
func samePath(a, b string) bool {
	return strings.EqualFold(normalizeSeparators(a), normalizeSeparators(b))
}

func normalizeSeparators(path string) string {
	return filepath.ToSlash(filepath.Clean(strings.ReplaceAll(path, `\`, "/")))
}

// requireWindowsHost guards the cases that need the host filepath package to
// understand backslashes and drive letters. Those branches only ever run on
// Windows, and the windows job in CI exercises them.
func requireWindowsHost(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "windows" {
		t.Skip("windows path semantics; covered by the windows CI job")
	}
}

func assertIDs(t *testing.T, got []Shell, want ...string) {
	t.Helper()
	ids := IDs(got)
	if strings.Join(ids, ",") != strings.Join(want, ",") {
		t.Fatalf("detected %v, want %v", ids, want)
	}
}

func find(t *testing.T, shells []Shell, id string) Shell {
	t.Helper()
	for _, sh := range shells {
		if sh.ID == id {
			return sh
		}
	}
	t.Fatalf("%s was not detected in %v", id, IDs(shells))
	return Shell{}
}

// TestDetectFindsTheShellsOnAWindowsDevBox mirrors a real machine: Git for
// Windows on PATH, ash from a w64devkit install, no pwsh and no zsh.
func TestDetectFindsTheShellsOnAWindowsDevBox(t *testing.T) {
	system := fakeSystem{
		goos: "windows",
		path: map[string]string{
			"powershell.exe": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			"cmd.exe":        `C:\Windows\system32\cmd.exe`,
			"bash.exe":       `C:\Program Files\Git\usr\bin\bash.exe`,
			"sh.exe":         `C:\Program Files\Git\usr\bin\sh.exe`,
			"ash.exe":        `D:\Tools\w64devkit\bin\ash.exe`,
		},
		env: windowsEnv(),
	}

	found := detect(system.prober())

	assertIDs(t, found, "powershell", "cmd", "bash", "sh", "ash")
	if bash := find(t, found, "bash"); bash.Kind != KindPosix {
		t.Fatalf("bash kind is %q, want %q", bash.Kind, KindPosix)
	}
	if ash := find(t, found, "ash"); ash.Path != `D:\Tools\w64devkit\bin\ash.exe` {
		t.Fatalf("ash path is %q", ash.Path)
	}
}

// TestDetectFindsGitShellsThatAreNotOnPath is the common Windows case: Git is
// installed but its bin directories were never added to PATH.
func TestDetectFindsGitShellsThatAreNotOnPath(t *testing.T) {
	system := fakeSystem{
		goos: "windows",
		path: map[string]string{"powershell.exe": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`},
		files: []string{
			`C:\Program Files\Git\bin\bash.exe`,
			`C:\Program Files\Git\usr\bin\sh.exe`,
		},
		env: windowsEnv(),
	}

	found := detect(system.prober())

	assertIDs(t, found, "powershell", "bash", "sh")
	if bash := find(t, found, "bash"); !samePath(bash.Path, `C:\Program Files\Git\bin\bash.exe`) {
		t.Fatalf("bash path is %q", bash.Path)
	}
}

// TestDetectKeepsTheHistoricalAutoOrder guards the upgrade path: a machine
// with everything installed must still default to PowerShell on Windows.
func TestDetectKeepsTheHistoricalAutoOrder(t *testing.T) {
	system := fakeSystem{
		goos: "windows",
		path: map[string]string{
			"powershell.exe": `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
			"pwsh.exe":       `C:\Program Files\PowerShell\7\pwsh.exe`,
			"cmd.exe":        `C:\Windows\system32\cmd.exe`,
			"bash.exe":       `C:\Program Files\Git\bin\bash.exe`,
		},
		env: windowsEnv(),
	}

	chosen, err := resolve(detect(system.prober()), "auto")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if chosen.ID != "powershell" {
		t.Fatalf("auto chose %q, want powershell", chosen.ID)
	}
}

func TestDetectOnUnixPrefersBash(t *testing.T) {
	system := fakeSystem{
		goos: "linux",
		path: map[string]string{
			"bash": "/bin/bash",
			"sh":   "/bin/sh",
			"zsh":  "/usr/bin/zsh",
		},
	}

	found := detect(system.prober())

	assertIDs(t, found, "bash", "sh", "zsh")
	for _, sh := range found {
		if sh.Kind != KindPosix {
			t.Fatalf("%s kind is %q, want %q", sh.ID, sh.Kind, KindPosix)
		}
	}
}

func TestDetectReportsNothingOnABareMachine(t *testing.T) {
	found := detect(fakeSystem{goos: "linux"}.prober())
	if len(found) != 0 {
		t.Fatalf("detected %v on a machine with no shells", IDs(found))
	}
	if _, err := resolve(found, "auto"); err == nil {
		t.Fatal("resolve accepted a machine with no shells")
	}
}

// TestDetectMarksWindowsWslBash covers the System32 bash.exe launcher, which
// runs inside WSL and does not see the Windows working directory.
func TestDetectMarksWindowsWslBash(t *testing.T) {
	requireWindowsHost(t)

	system := fakeSystem{
		goos:  "windows",
		path:  map[string]string{"bash.exe": `C:\Windows\System32\bash.exe`},
		env:   windowsEnv(),
		files: nil,
	}

	bash := find(t, detect(system.prober()), "bash")
	if bash.Note != "WSL" {
		t.Fatalf("note is %q, want WSL", bash.Note)
	}
	if !strings.Contains(bash.Label(), "WSL") {
		t.Fatalf("label %q does not mention WSL", bash.Label())
	}
}

func TestResolveAcceptsUntidyConfigurationValues(t *testing.T) {
	available := []Shell{
		{ID: "powershell", Kind: KindPowerShell, Path: `C:\powershell.exe`},
		{ID: "bash", Kind: KindPosix, Path: `C:\bash.exe`},
	}

	for _, preference := range []string{"PowerShell.EXE", "  powershell  ", "POWERSHELL"} {
		chosen, err := resolve(available, preference)
		if err != nil {
			t.Fatalf("resolve(%q): %v", preference, err)
		}
		if chosen.ID != "powershell" {
			t.Fatalf("resolve(%q) chose %q", preference, chosen.ID)
		}
	}
}

func TestResolveExplainsAShellThatIsNotInstalled(t *testing.T) {
	available := []Shell{{ID: "powershell", Kind: KindPowerShell, Path: `C:\powershell.exe`}}

	_, err := resolve(available, "bash")
	if err == nil {
		t.Fatal("resolve accepted a shell that is not installed")
	}
	if !strings.Contains(err.Error(), "not installed") || !strings.Contains(err.Error(), "powershell") {
		t.Fatalf("error does not say what is available: %v", err)
	}
}

func TestResolveSeparatesUnknownFromUnavailable(t *testing.T) {
	available := []Shell{{ID: "bash", Kind: KindPosix, Path: "/bin/bash"}}

	_, err := resolve(available, "fish")
	if err == nil {
		t.Fatal("resolve accepted an unknown shell")
	}
	if !strings.Contains(err.Error(), "unknown shell") {
		t.Fatalf("error should call fish unknown: %v", err)
	}
}

func TestResolveHonoursAnExplicitPath(t *testing.T) {
	requireWindowsHost(t)

	system := fakeSystem{
		goos:  "windows",
		files: []string{`C:\Program Files\Git\bin\bash.exe`},
		env:   windowsEnv(),
	}

	sh, ok := resolveExplicitPath(system.prober(), `C:\Program Files\Git\bin\bash.exe`)
	if !ok {
		t.Fatal("an existing explicit path was rejected")
	}
	if sh.ID != "bash" || sh.Kind != KindPosix {
		t.Fatalf("explicit path resolved to %q/%q", sh.ID, sh.Kind)
	}

	if _, ok := resolveExplicitPath(system.prober(), `C:\nope\bash.exe`); ok {
		t.Fatal("a missing explicit path was accepted")
	}
	if _, ok := resolveExplicitPath(system.prober(), "bash"); ok {
		t.Fatal("a plain id was treated as a path")
	}
}

func TestOptionsOffersAutoAndKeepsAMissingConfiguredValue(t *testing.T) {
	available := []Shell{
		{ID: "powershell", Kind: KindPowerShell, Path: `C:\powershell.exe`},
		{ID: "bash", Kind: KindPosix, Path: `C:\bash.exe`},
	}

	if got := options(available, "auto"); strings.Join(got, ",") != "auto,powershell,bash" {
		t.Fatalf("options are %v", got)
	}
	if got := options(available, "zsh"); strings.Join(got, ",") != "auto,powershell,bash,zsh" {
		t.Fatalf("a configured value that is not installed was dropped: %v", got)
	}
}

func TestKindForNameCoversEveryShellFamily(t *testing.T) {
	cases := map[string]Kind{
		"powershell.exe": KindPowerShell,
		"pwsh":           KindPowerShell,
		"cmd.exe":        KindCmd,
		"bash.exe":       KindPosix,
		"sh":             KindPosix,
		"ash.exe":        KindPosix,
		"dash":           KindPosix,
		"zsh":            KindPosix,
	}

	for name, want := range cases {
		if got := kindForName(name); got != want {
			t.Fatalf("kindForName(%q) = %q, want %q", name, got, want)
		}
	}
}
