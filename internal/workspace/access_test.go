package workspace

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
)

// accessFixture lays out
//
//	parent/ws                 open workspace
//	parent/shared             read-only extra root (lib.go, .git/config, .ssh/id_rsa, sub/)
//	parent/rw                 read-write extra root
//	parent/outside            not a root (secret.txt)
//	parent/home               fake home with .ssh and AppData credential folders
//	parent/sys                fake Windows folder (system deny rule)
//	parent/sys/explicit       read-only extra root placed inside sys
type accessFixture struct {
	parent, ws, shared, rw, outside, home, sys, explicit string
	cfg                                                  *config.Config
	loc                                                  locations
}

func newAccessFixture(t *testing.T) *accessFixture {
	t.Helper()
	parent := t.TempDir()
	f := &accessFixture{
		parent:  parent,
		ws:      filepath.Join(parent, "ws"),
		shared:  filepath.Join(parent, "shared"),
		rw:      filepath.Join(parent, "rw"),
		outside: filepath.Join(parent, "outside"),
		home:    filepath.Join(parent, "home"),
		sys:     filepath.Join(parent, "sys"),
	}
	f.explicit = filepath.Join(f.sys, "explicit")
	files := map[string]string{
		filepath.Join(f.ws, "inside.txt"):                                               "inside",
		filepath.Join(f.shared, "lib.go"):                                               "package lib",
		filepath.Join(f.shared, "sub", "util.go"):                                       "package sub",
		filepath.Join(f.shared, ".git", "config"):                                       "[remote] url=https://token@host",
		filepath.Join(f.shared, ".ssh", "id_rsa"):                                       "PRIVATE",
		filepath.Join(f.shared, "keys", "id_ed25519"):                                   "PRIVATE",
		filepath.Join(f.shared, ".git-credentials"):                                     "https://u:p@host",
		filepath.Join(f.rw, "notes.txt"):                                                "notes",
		filepath.Join(f.outside, "secret.txt"):                                          "TOP-SECRET-MARKER",
		filepath.Join(f.home, ".ssh", "id_rsa"):                                         "PRIVATE",
		filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials", "blob"): "cred",
		filepath.Join(f.home, "projects", "app", "main.go"):                             "package main",
		filepath.Join(f.sys, "system.ini"):                                              "sys",
		filepath.Join(f.explicit, "hosts"):                                              "127.0.0.1",
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(f.ws, 0o700); err != nil {
		t.Fatal(err)
	}
	cfg := config.DefaultConfig()
	cfg.ConfigDir = filepath.Join(parent, "app", ".devspace")
	cfg.StateDir = filepath.Join(parent, "app", ".devspace-state")
	for _, dir := range []string{cfg.ConfigDir, cfg.StateDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	cfg.OutsideAccess.Roots = []config.ExtraRootConfig{
		{Path: f.shared, Access: "read"},
		{Path: f.rw, Access: "readwrite"},
		{Path: f.explicit},
	}
	f.cfg = cfg
	f.loc = locations{
		goos:         "windows",
		home:         f.home,
		appData:      filepath.Join(f.home, "AppData", "Roaming"),
		localAppData: filepath.Join(f.home, "AppData", "Local"),
		systemRoot:   f.sys,
	}
	return f
}

func (f *accessFixture) policy(t *testing.T) *AccessPolicy {
	t.Helper()
	policy, problems := newAccessPolicy(f.cfg, f.loc)
	if len(problems) > 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
	return policy
}

func (f *accessFixture) guard(t *testing.T, origin Origin) *Guard {
	return f.policy(t).Guard(origin, "test")
}

func requireDenied(t *testing.T, g *Guard, wsRoot, input string, op Op) {
	t.Helper()
	got, outside, err := g.Resolve(wsRoot, input, op)
	if err == nil {
		t.Fatalf("Resolve(%q, %s) = %q (outside=%v), want denial", input, op, got, outside)
	}
}

func requireAllowed(t *testing.T, g *Guard, wsRoot, input string, op Op) string {
	t.Helper()
	got, _, err := g.Resolve(wsRoot, input, op)
	if err != nil {
		t.Fatalf("Resolve(%q, %s) error: %v", input, op, err)
	}
	return got
}

func TestGuardKeepsWorkspaceBehaviour(t *testing.T) {
	f := newAccessFixture(t)
	for _, g := range []*Guard{nil, f.guard(t, OriginLocal), f.guard(t, OriginRemote)} {
		for _, op := range []Op{OpRead, OpWrite, OpMove, OpExec} {
			got, outside, err := g.Resolve(f.ws, "inside.txt", op)
			if err != nil || outside || got != filepath.Join(f.ws, "inside.txt") {
				t.Fatalf("workspace path: %q %v %v", got, outside, err)
			}
		}
		requireDenied(t, g, f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
		requireDenied(t, g, f.ws, "../outside/secret.txt", OpRead)
	}
}

func TestGuardReadOnlyRootEnforcement(t *testing.T) {
	f := newAccessFixture(t)
	g := f.guard(t, OriginLocal)

	for _, input := range []string{filepath.Join(f.shared, "lib.go"), "../shared/lib.go", "../shared/sub/../lib.go", f.shared} {
		got, outside, err := g.Resolve(f.ws, input, OpRead)
		if err != nil || !outside {
			t.Fatalf("read %q: %q %v %v", input, got, outside, err)
		}
	}
	for _, op := range []Op{OpWrite, OpMove, OpExec} {
		requireDenied(t, g, f.ws, filepath.Join(f.shared, "lib.go"), op)
		requireDenied(t, g, f.ws, filepath.Join(f.shared, "new", "file.go"), op)
	}
	_, _, err := g.Resolve(f.ws, filepath.Join(f.shared, "lib.go"), OpWrite)
	if !errors.Is(err, ErrAccessDenied) || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("write to read-only root: %v", err)
	}

	// Read-write root: everything except moving the root itself.
	for _, op := range []Op{OpRead, OpWrite, OpMove, OpExec} {
		requireAllowed(t, g, f.ws, filepath.Join(f.rw, "notes.txt"), op)
	}
	requireAllowed(t, g, f.ws, filepath.Join(f.rw, "new", "dir", "file.txt"), OpWrite)
	requireDenied(t, g, f.ws, f.rw, OpMove)
}

func TestGuardMostSpecificRootWins(t *testing.T) {
	f := newAccessFixture(t)
	nested := filepath.Join(f.rw, "vendor")
	if err := os.Mkdir(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	f.cfg.OutsideAccess.Roots = append(f.cfg.OutsideAccess.Roots, config.ExtraRootConfig{Path: nested, Access: "read"})
	g := f.guard(t, OriginLocal)
	requireAllowed(t, g, f.ws, filepath.Join(f.rw, "a.txt"), OpWrite)
	requireDenied(t, g, f.ws, filepath.Join(nested, "a.txt"), OpWrite)
	// The nested root cannot be moved away through the outer read-write root.
	requireDenied(t, g, f.ws, nested, OpMove)
}

func TestGuardRejectsTraversal(t *testing.T) {
	f := newAccessFixture(t)
	if err := os.Mkdir(f.shared+"-evil", 0o700); err != nil {
		t.Fatal(err)
	}
	g := f.guard(t, OriginLocal)
	sep := string(filepath.Separator)
	for _, input := range []string{
		"../outside/secret.txt",
		"../shared/../outside/secret.txt",
		"../shared/sub/../../outside/secret.txt",
		f.shared + sep + ".." + sep + "outside" + sep + "secret.txt",
		f.shared + "/../outside/secret.txt",
		f.shared + "-evil",
		filepath.Join(f.shared+"-evil", "x.txt"),
		f.parent,
		"..",
		"../..",
		"a\x00b",
	} {
		requireDenied(t, g, f.ws, input, OpRead)
	}
}

func TestGuardRejectsLinkEscapes(t *testing.T) {
	f := newAccessFixture(t)
	makeDirLink(t, f.outside, filepath.Join(f.shared, "escape"))
	makeDirLink(t, filepath.Join(f.home, ".ssh"), filepath.Join(f.rw, "keys"))
	makeDirLink(t, f.shared, filepath.Join(f.rw, "to-shared"))
	g := f.guard(t, OriginLocal)

	requireDenied(t, g, f.ws, filepath.Join(f.shared, "escape", "secret.txt"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, "escape", "new.txt"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.rw, "keys", "id_rsa"), OpRead)
	// A link from a read-write root into a read-only one is judged by its
	// target: reading works, writing does not.
	requireAllowed(t, g, f.ws, filepath.Join(f.rw, "to-shared", "lib.go"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.rw, "to-shared", "lib.go"), OpWrite)
	// The workspace itself must not escape into a root through a link either:
	// ResolveInRoot reports it as outside and the policy judges the target.
	makeDirLink(t, f.outside, filepath.Join(f.ws, "escape"))
	requireDenied(t, g, f.ws, "escape/secret.txt", OpRead)
}

// TestGuardRejectsJunctionEscape covers NTFS junctions explicitly (see
// TestResolveInRootRejectsJunctionEscape).
func TestGuardRejectsJunctionEscape(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("junctions are Windows-only")
	}
	f := newAccessFixture(t)
	makeJunction(t, f.outside, filepath.Join(f.ws, "junction"))
	makeJunction(t, filepath.Join(f.home, ".ssh"), filepath.Join(f.shared, "junction-ssh"))
	makeJunction(t, f.outside, filepath.Join(f.shared, "junction-out"))
	g := f.guard(t, OriginLocal)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, "junction-ssh", "id_rsa"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, "junction-out", "secret.txt"), OpRead)
	requireDenied(t, g, f.ws, "junction/secret.txt", OpRead)
}

func TestGuardWindowsCaseVariants(t *testing.T) {
	if !pathsAreCaseInsensitive {
		t.Skip("case-sensitive filesystem")
	}
	f := newAccessFixture(t)
	g := f.guard(t, OriginLocal)
	requireAllowed(t, g, f.ws, strings.ToUpper(filepath.Join(f.shared, "lib.go")), OpRead)
	requireAllowed(t, g, f.ws, strings.ToLower(filepath.Join(f.shared, "lib.go")), OpRead)
	requireDenied(t, g, f.ws, strings.ToUpper(filepath.Join(f.shared, "lib.go")), OpWrite)
	requireDenied(t, g, f.ws, strings.ToUpper(filepath.Join(f.shared, ".SSH", "ID_RSA")), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, ".GIT", "CONFIG"), OpRead)
	requireDenied(t, g, f.ws, strings.ToUpper(filepath.Join(f.outside, "secret.txt")), OpRead)
}

func TestGuardWindowsPathForms(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows path syntax")
	}
	f := newAccessFixture(t)
	g := f.guard(t, OriginLocal)
	lib := filepath.Join(f.shared, "lib.go")
	drive := filepath.VolumeName(f.shared)
	for _, input := range []string{
		`\\devspace-test-no-such-host\share\lib.go`, // UNC: must fail before any network access
		`\\?\` + lib,
		`\\.\` + lib,
		`\\?\UNC\devspace-test-no-such-host\share`,
		drive + "lib.go",               // drive-relative "C:foo"
		strings.TrimPrefix(lib, drive), // rooted "\Users\...": relative to the current drive
		lib + ":stream",
		lib + "::$DATA",
		lib + ".",
		lib + " ",
		filepath.Join(f.shared, "NUL"),
		filepath.Join(f.shared, "con.txt"),
		filepath.Join(f.shared, "COM1"),
		filepath.Join(f.shared, "CONIN$"),
		filepath.Join(f.shared, ".ssh.", "id_rsa"),
	} {
		requireDenied(t, g, f.ws, input, OpRead)
	}
	_, _, err := g.Resolve(f.ws, `\\devspace-test-no-such-host\share\lib.go`, OpRead)
	if err == nil || !strings.Contains(err.Error(), "UNC") {
		t.Fatalf("UNC path not rejected lexically: %v", err)
	}
}

func TestGuardWindowsShortNames(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("8.3 names are Windows-only")
	}
	f := newAccessFixture(t)
	longDir := filepath.Join(f.shared, "a-very-long-directory-name")
	if err := os.MkdirAll(longDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(longDir, "x.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	longOutside := filepath.Join(f.parent, "outside-with-a-long-name")
	if err := os.MkdirAll(longOutside, 0o700); err != nil {
		t.Fatal(err)
	}
	shortInside := shortPathName(t, longDir)
	shortOutside := shortPathName(t, longOutside)
	if strings.EqualFold(shortInside, longDir) || strings.EqualFold(shortOutside, longOutside) {
		t.Skip("8.3 short names are disabled on this volume")
	}
	g := f.guard(t, OriginLocal)
	got := requireAllowed(t, g, f.ws, filepath.Join(shortInside, "x.txt"), OpRead)
	if !strings.EqualFold(got, filepath.Join(longDir, "x.txt")) {
		t.Fatalf("short name not expanded: %q", got)
	}
	requireDenied(t, g, f.ws, filepath.Join(shortInside, "x.txt"), OpWrite)
	requireDenied(t, g, f.ws, shortOutside, OpRead)
}

func TestGuardDenyList(t *testing.T) {
	f := newAccessFixture(t)
	// The fake home becomes a read-write root: its project files are usable,
	// its credential folders are not.
	f.cfg.OutsideAccess.Roots = append(f.cfg.OutsideAccess.Roots, config.ExtraRootConfig{Path: f.home, Access: "readwrite"})
	f.cfg.OutsideAccess.Deny = []string{filepath.Join(f.home, "private")}
	// The Dev Space Go config sits inside a read-write root.
	f.cfg.ConfigDir = filepath.Join(f.rw, "app", ".devspace")
	f.cfg.StateDir = filepath.Join(f.rw, "app", "state")
	if err := os.MkdirAll(f.cfg.StateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	g := f.guard(t, OriginLocal)

	requireAllowed(t, g, f.ws, filepath.Join(f.home, "projects", "app", "main.go"), OpWrite)
	for _, input := range []string{
		filepath.Join(f.home, ".ssh", "id_rsa"),
		filepath.Join(f.home, ".ssh"),
		filepath.Join(f.home, ".ssh", "new_key"),
		filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials", "blob"),
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data"),
		filepath.Join(f.home, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles"),
		filepath.Join(f.home, "private", "diary.txt"),
		filepath.Join(f.shared, ".git", "config"),
		filepath.Join(f.shared, ".git-credentials"),
		filepath.Join(f.shared, "keys", "id_ed25519"),
		filepath.Join(f.shared, ".ssh", "id_rsa"),
		filepath.Join(f.rw, "app", ".devspace", "config.json"),
		filepath.Join(f.rw, "app", "state", "outside-access.log"),
	} {
		requireDenied(t, g, f.ws, input, OpRead)
		requireDenied(t, g, f.ws, input, OpWrite)
	}
	// Other files in .git stay readable.
	if err := os.WriteFile(filepath.Join(f.shared, ".git", "HEAD"), []byte("ref"), 0o600); err != nil {
		t.Fatal(err)
	}
	requireAllowed(t, g, f.ws, filepath.Join(f.shared, ".git", "HEAD"), OpRead)
}

func TestGuardSystemFoldersNeedAnExplicitRoot(t *testing.T) {
	f := newAccessFixture(t)
	// A root around the system folder does not open it ...
	f.cfg.OutsideAccess.Roots = append(f.cfg.OutsideAccess.Roots, config.ExtraRootConfig{Path: f.parent, Access: "read"})
	g := f.guard(t, OriginLocal)
	requireDenied(t, g, f.ws, filepath.Join(f.sys, "system.ini"), OpRead)
	// ... but a root configured inside it does.
	requireAllowed(t, g, f.ws, filepath.Join(f.explicit, "hosts"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.explicit, "hosts"), OpWrite)
	// The broad root still reaches ordinary folders and still hides secrets.
	requireAllowed(t, g, f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.home, ".ssh", "id_rsa"), OpRead)
}

func TestNewAccessPolicyReportsUnusableRoots(t *testing.T) {
	f := newAccessFixture(t)
	file := filepath.Join(f.shared, "lib.go")
	volumeRoot := string(filepath.Separator)
	if runtime.GOOS == "windows" {
		volumeRoot = filepath.VolumeName(f.parent) + `\`
	}
	bad := []config.ExtraRootConfig{
		{Path: "relative/dir"},
		{Path: filepath.Join(f.parent, "missing")},
		{Path: file},
		{Path: volumeRoot},
		{Path: filepath.Join(f.home, ".ssh")},
		{Path: filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials")},
		{Path: f.cfg.StateDir},
		{Path: f.outside, Access: "admin"},
	}
	if runtime.GOOS == "windows" {
		bad = append(bad, config.ExtraRootConfig{Path: `\\devspace-test-no-such-host\share`})
	}
	f.cfg.OutsideAccess.Roots = bad
	policy, problems := newAccessPolicy(f.cfg, f.loc)
	if got := policy.Roots(OriginLocal); len(got) != 0 {
		t.Fatalf("unusable roots accepted: %+v", got)
	}
	// Every bad root except the invalid access level (reported by
	// config.Problems) is explained.
	if len(problems) != len(bad)-1 {
		t.Fatalf("got %d problems, want %d: %v", len(problems), len(bad)-1, problems)
	}
	if !strings.Contains(strings.Join(f.cfg.Problems(), "\n"), `"admin"`) {
		t.Fatalf("config.Problems does not report the invalid access level: %v", f.cfg.Problems())
	}
}

func TestGuardRemoteOrigin(t *testing.T) {
	f := newAccessFixture(t)
	lib := filepath.Join(f.shared, "lib.go")

	policy := f.policy(t)
	if len(policy.Roots(OriginRemote)) != 0 || len(policy.Roots(OriginLocal)) != 3 {
		t.Fatalf("roots: local %v remote %v", policy.Roots(OriginLocal), policy.Roots(OriginRemote))
	}
	_, _, err := policy.Guard(OriginRemote, "read").Resolve(f.ws, lib, OpRead)
	if err == nil || !strings.Contains(err.Error(), "tunnel") {
		t.Fatalf("remote read of extra root: %v", err)
	}
	// Unknown origins are treated as remote.
	requireDenied(t, policy.Guard(Origin("other"), "read"), f.ws, lib, OpRead)
	requireAllowed(t, policy.Guard(OriginLocal, "read"), f.ws, lib, OpRead)

	f.cfg.OutsideAccess.Remote = true
	remote := f.guard(t, OriginRemote)
	requireAllowed(t, remote, f.ws, lib, OpRead)
	// Enabling remote access lifts no other rule.
	requireDenied(t, remote, f.ws, lib, OpWrite)
	requireDenied(t, remote, f.ws, filepath.Join(f.shared, ".ssh", "id_rsa"), OpRead)
	requireDenied(t, remote, f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
}

func TestGuardMoveCannotTouchRootsOrWorkspace(t *testing.T) {
	f := newAccessFixture(t)
	// Workspace inside a read-write root.
	ws := filepath.Join(f.rw, "project")
	if err := os.Mkdir(ws, 0o700); err != nil {
		t.Fatal(err)
	}
	g := f.guard(t, OriginLocal)
	requireAllowed(t, g, ws, filepath.Join(f.rw, "notes.txt"), OpMove)
	requireDenied(t, g, ws, f.rw, OpMove)
	requireDenied(t, g, ws, "..", OpMove)
	// The workspace root itself resolves inside the workspace; the move tool
	// refuses it separately.
}

func TestGuardHiddenFiltersWalks(t *testing.T) {
	f := newAccessFixture(t)
	g := f.guard(t, OriginLocal)
	for path, want := range map[string]bool{
		filepath.Join(f.shared, "lib.go"):             false,
		filepath.Join(f.shared, "sub"):                false,
		filepath.Join(f.shared, ".ssh"):               true,
		filepath.Join(f.shared, ".git", "config"):     true,
		filepath.Join(f.shared, ".git-credentials"):   true,
		filepath.Join(f.shared, "keys", "id_ed25519"): true,
		filepath.Join(f.outside, "secret.txt"):        true,
	} {
		if got := g.Hidden(path); got != want {
			t.Errorf("Hidden(%q) = %v, want %v", path, got, want)
		}
	}
	var nilGuard *Guard
	if !nilGuard.Hidden(filepath.Join(f.shared, "lib.go")) {
		t.Fatal("nil guard must hide everything outside the workspace")
	}
}

func TestAuditLogRecordsOutsideAccess(t *testing.T) {
	f := newAccessFixture(t)
	g := f.guard(t, OriginLocal)
	requireAllowed(t, g, f.ws, filepath.Join(f.shared, "lib.go"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, "lib.go"), OpWrite)
	requireDenied(t, g, f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
	requireAllowed(t, g, f.ws, "inside.txt", OpRead) // workspace: not audited

	file, err := os.Open(filepath.Join(f.cfg.StateDir, auditFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var entries []auditEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "TOP-SECRET-MARKER") || strings.Contains(scanner.Text(), "package lib") {
			t.Fatalf("audit log contains file content: %s", scanner.Text())
		}
		var entry auditEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, entry)
	}
	if len(entries) != 3 {
		t.Fatalf("got %d audit entries, want 3: %+v", len(entries), entries)
	}
	if !entries[0].Allow || entries[1].Allow || entries[2].Allow {
		t.Fatalf("decisions: %+v", entries)
	}
	if entries[0].Origin != "local" || entries[0].Op != "read" || entries[0].Tool != "test" || entries[0].Time == "" {
		t.Fatalf("entry fields: %+v", entries[0])
	}
	if entries[1].Reason == "" || !strings.Contains(entries[1].Reason, "read-only") {
		t.Fatalf("denial reason: %+v", entries[1])
	}

	// Audit can be turned off.
	off := false
	f.cfg.OutsideAccess.Audit = &off
	f.cfg.StateDir = filepath.Join(f.parent, "app", "state-off")
	if err := os.MkdirAll(f.cfg.StateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	requireAllowed(t, f.guard(t, OriginLocal), f.ws, filepath.Join(f.shared, "lib.go"), OpRead)
	if _, err := os.Stat(filepath.Join(f.cfg.StateDir, auditFileName)); !os.IsNotExist(err) {
		t.Fatalf("audit log written although disabled: %v", err)
	}
}

func TestCanonicalizeKeepsMissingTail(t *testing.T) {
	f := newAccessFixture(t)
	got, err := Canonicalize(filepath.Join(f.shared, "missing", "deeper", "file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	realShared, err := Canonicalize(f.shared)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(realShared, "missing", "deeper", "file.txt"); got != want {
		t.Fatalf("Canonicalize = %q, want %q", got, want)
	}
}

func TestBuiltinDenyRulesCoverPlatforms(t *testing.T) {
	for _, goos := range []string{"windows", "darwin", "linux"} {
		rules := builtinDenyRules(locations{goos: goos, home: "/h", appData: "/h/r", localAppData: "/h/l", systemRoot: "/w"})
		var credentials, system int
		for _, rule := range rules {
			if rule.system {
				system++
			} else {
				credentials++
			}
		}
		if credentials < 3 || system < 1 {
			t.Errorf("%s: %d credential and %d system rules", goos, credentials, system)
		}
	}
}
