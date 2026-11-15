package workspace

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
)

// fullFixture reuses accessFixture but switches to the default mode
// (mode unset = full) and adds a few files the small deny list covers.
func fullFixture(t *testing.T) *accessFixture {
	t.Helper()
	f := newAccessFixture(t)
	f.cfg.OutsideAccess = config.OutsideAccessConfig{} // defaults: full, remote on, deny defaults on
	files := map[string]string{
		filepath.Join(f.home, ".ssh", "id_ed25519.pub"):                                                             "ssh-ed25519 AAAA",
		filepath.Join(f.home, ".ssh", "config"):                                                                     "Host x",
		filepath.Join(f.home, ".config", "gh", "hosts.yml"):                                                         "oauth_token: x",
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data"):         "db",
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Network", "Cookies"): "db",
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Bookmarks"):          "{}",
		filepath.Join(f.outside, ".env"):                                                                            "SECRET=1",
		filepath.Join(f.outside, ".npmrc"):                                                                          "registry=x",
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func TestFullModeIsTheDefault(t *testing.T) {
	f := fullFixture(t)
	policy := f.policy(t)
	if policy.Mode() != config.AccessModeFull || !policy.FullAccess(OriginLocal) || !policy.FullAccess(OriginRemote) {
		t.Fatalf("default mode %q, local %v, remote %v", policy.Mode(), policy.FullAccess(OriginLocal), policy.FullAccess(OriginRemote))
	}
}

func TestFullModeAllowsOutsidePathsLocallyAndRemotely(t *testing.T) {
	f := fullFixture(t)
	for _, origin := range []Origin{OriginLocal, OriginRemote} {
		g := f.guard(t, origin)
		for _, op := range []Op{OpRead, OpWrite, OpExec, OpMove} {
			requireAllowed(t, g, f.ws, filepath.Join(f.outside, "secret.txt"), op)
			requireAllowed(t, g, f.ws, "../outside/secret.txt", op)
			requireAllowed(t, g, f.ws, filepath.Join(f.parent, "new-dir", "new.txt"), op)
		}
		// Not blocked in full mode: .env, .npmrc, .ssh/config, public keys,
		// browser bookmarks, "system" folders and .git/config.
		for _, path := range []string{
			filepath.Join(f.outside, ".env"),
			filepath.Join(f.outside, ".npmrc"),
			filepath.Join(f.home, ".ssh", "config"),
			filepath.Join(f.home, ".ssh", "id_ed25519.pub"),
			filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Bookmarks"),
			filepath.Join(f.sys, "system.ini"),
			filepath.Join(f.shared, ".git", "config"),
		} {
			requireAllowed(t, g, f.ws, path, OpWrite)
		}
		// The workspace root (or a parent of it) still cannot be moved.
		requireDenied(t, g, f.ws, f.parent, OpMove)
	}
}

func TestFullModeDefaultDenyList(t *testing.T) {
	f := fullFixture(t)
	g := f.guard(t, OriginLocal)
	denied := []string{
		filepath.Join(f.home, ".ssh", "id_rsa"),
		filepath.Join(f.shared, "keys", "id_ed25519"),
		filepath.Join(f.shared, ".git-credentials"),
		filepath.Join(f.home, ".config", "gh", "hosts.yml"),
		filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials", "blob"),
		filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Protect", "S-1-5"),
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data"),
		filepath.Join(f.home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Network", "Cookies"),
		filepath.Join(f.cfg.ConfigDir, "config.json"),
		filepath.Join(f.cfg.StateDir, "outside-access.log"),
	}
	for _, path := range denied {
		requireDenied(t, g, f.ws, path, OpRead)
		requireDenied(t, g, f.ws, path, OpWrite)
	}

	// The list also applies inside the workspace, because in full mode the
	// workspace itself can be a home folder.
	homeWs := f.home
	requireDenied(t, g, homeWs, ".ssh/id_rsa", OpRead)
	requireAllowed(t, g, homeWs, ".ssh/config", OpRead)
	if !g.Hidden(filepath.Join(f.home, ".ssh", "id_rsa")) || g.Hidden(filepath.Join(f.home, ".ssh", "config")) {
		t.Fatal("walk filter does not follow the full-mode deny list")
	}
	if g.WalkFilter(false) == nil {
		t.Fatal("full mode needs a walk filter inside the workspace too")
	}

	// Custom deny entries add to the list.
	f.cfg.OutsideAccess.Deny = []string{filepath.Join(f.outside, "private")}
	requireDenied(t, f.guard(t, OriginLocal), f.ws, filepath.Join(f.outside, "private", "x.txt"), OpRead)
}

func TestFullModeDenyListCanBeCleared(t *testing.T) {
	f := fullFixture(t)
	off := false
	f.cfg.OutsideAccess.DenyDefaults = &off
	f.cfg.OutsideAccess.Deny = []string{}
	g := f.guard(t, OriginLocal)
	for _, path := range []string{
		filepath.Join(f.home, ".ssh", "id_rsa"),
		filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials", "blob"),
		filepath.Join(f.cfg.ConfigDir, "config.json"),
	} {
		requireAllowed(t, g, f.ws, path, OpRead)
	}
	if g.WalkFilter(false) != nil {
		t.Fatal("no deny rules: walks need no filter")
	}
	if len(f.policy(t).DeniedPaths()) != 0 {
		t.Fatalf("deny paths left: %v", f.policy(t).DeniedPaths())
	}
}

func TestFullModeRemoteCanBeDisabled(t *testing.T) {
	f := fullFixture(t)
	off := false
	f.cfg.OutsideAccess.Remote = &off
	policy := f.policy(t)
	if policy.FullAccess(OriginRemote) || !policy.FullAccess(OriginLocal) {
		t.Fatal("remote:false must keep full access local-only")
	}
	_, _, err := policy.Guard(OriginRemote, "read").Resolve(f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
	if err == nil || !strings.Contains(err.Error(), "tunnel") {
		t.Fatalf("remote outside read: %v", err)
	}
	requireAllowed(t, policy.Guard(OriginRemote, "read"), f.ws, "inside.txt", OpRead)
	requireAllowed(t, policy.Guard(OriginLocal, "read"), f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
}

func TestFullModeKeepsPathHygiene(t *testing.T) {
	f := fullFixture(t)
	g := f.guard(t, OriginLocal)
	sep := string(filepath.Separator)
	inputs := []string{"a\x00b"}
	if runtime.GOOS == "windows" {
		drive := filepath.VolumeName(f.outside)
		secret := filepath.Join(f.outside, "secret.txt")
		inputs = append(inputs,
			sep+sep+"devspace-test-no-such-host"+sep+"share"+sep+"x",
			sep+sep+"?"+sep+secret,
			sep+sep+"."+sep+secret,
			drive+"foo",
			strings.TrimPrefix(secret, drive),
			secret+":stream",
			secret+"::$DATA",
			secret+".",
			filepath.Join(f.outside, "NUL"),
			filepath.Join(f.outside, "com1.txt"),
		)
	}
	for _, input := range inputs {
		requireDenied(t, g, f.ws, input, OpRead)
	}
}

func TestFullModeLinksAreJudgedByTarget(t *testing.T) {
	f := fullFixture(t)
	makeDirLink(t, filepath.Join(f.home, ".ssh"), filepath.Join(f.outside, "keys"))
	g := f.guard(t, OriginLocal)
	requireDenied(t, g, f.ws, filepath.Join(f.outside, "keys", "id_rsa"), OpRead)
	requireAllowed(t, g, f.ws, filepath.Join(f.outside, "keys", "config"), OpRead)
	// Inside the workspace as well.
	makeDirLink(t, filepath.Join(f.home, ".ssh"), filepath.Join(f.ws, "keys"))
	requireDenied(t, g, f.ws, "keys/id_rsa", OpRead)
}

func TestOffModeConfinesToWorkspace(t *testing.T) {
	f := fullFixture(t)
	f.cfg.OutsideAccess.Mode = config.AccessModeOff
	f.cfg.OutsideAccess.Roots = []config.ExtraRootConfig{{Path: f.shared}}
	policy := f.policy(t)
	if policy.FullAccess(OriginLocal) || len(policy.Roots(OriginLocal)) != 0 {
		t.Fatal("off mode grants access")
	}
	g := policy.Guard(OriginLocal, "read")
	requireDenied(t, g, f.ws, filepath.Join(f.outside, "secret.txt"), OpRead)
	requireDenied(t, g, f.ws, filepath.Join(f.shared, "lib.go"), OpRead)
	requireAllowed(t, g, f.ws, "inside.txt", OpWrite)
	if g.WalkFilter(false) != nil {
		t.Fatal("off mode must not filter workspace walks")
	}
}

func TestInvalidModeFallsBackToOff(t *testing.T) {
	f := fullFixture(t)
	f.cfg.OutsideAccess.Mode = "everything"
	if mode := f.policy(t).Mode(); mode != config.AccessModeOff {
		t.Fatalf("invalid mode resolved to %q", mode)
	}
	if !strings.Contains(strings.Join(f.cfg.Problems(), "\n"), "everything") {
		t.Fatal("invalid mode not reported")
	}
}

func TestRegistryOpensWorkspacesOutsideAllowedRootsOnlyInFullMode(t *testing.T) {
	f := fullFixture(t)
	f.cfg.AllowedRoots = []string{f.ws}
	open := func(t *testing.T, origin Origin, path string) error {
		t.Helper()
		registry := NewRegistry(f.cfg, nil)
		registry.SetAccessPolicy(f.policy(t))
		ctx, err := registry.OpenWorkspaceAs(origin, path, ModeCheckout, "")
		if err != nil {
			return err
		}
		if _, err := registry.GetWorkspaceAs(ctx.Workspace.ID, origin); err != nil {
			t.Fatalf("reopen by id: %v", err)
		}
		return nil
	}

	for _, origin := range []Origin{OriginLocal, OriginRemote} {
		if err := open(t, origin, f.outside); err != nil {
			t.Fatalf("full mode, %s: %v", origin, err)
		}
	}
	// A workspace on a credential store would bypass the deny list for its
	// own files; refuse it.
	if err := open(t, OriginLocal, filepath.Join(f.home, "AppData", "Roaming", "Microsoft", "Credentials")); err == nil {
		t.Fatal("workspace inside a credential store accepted")
	}
	if err := open(t, OriginLocal, "relative/dir"); err == nil {
		t.Fatal("relative workspace path accepted")
	}

	off := false
	f.cfg.OutsideAccess.Remote = &off
	if err := open(t, OriginRemote, f.outside); err == nil {
		t.Fatal("remote:false still opens workspaces outside allowedRoots remotely")
	}
	if err := open(t, OriginLocal, f.outside); err != nil {
		t.Fatalf("local open with remote:false: %v", err)
	}
	f.cfg.OutsideAccess.Mode = config.AccessModeRoots
	if err := open(t, OriginLocal, f.outside); err == nil {
		t.Fatal("roots mode opens workspaces outside allowedRoots")
	}

	// A workspace opened locally in full mode is not handed to a tunnel
	// request when remote access is off.
	f.cfg.OutsideAccess.Mode = ""
	registry := NewRegistry(f.cfg, nil)
	registry.SetAccessPolicy(f.policy(t))
	ctx, err := registry.OpenWorkspaceAs(OriginLocal, f.outside, ModeCheckout, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.GetWorkspaceAs(ctx.Workspace.ID, OriginRemote); err == nil {
		t.Fatal("remote request got a workspace outside allowedRoots")
	}
}
