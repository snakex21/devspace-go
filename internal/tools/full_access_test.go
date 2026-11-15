package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/workspace"
)

func fullAccessContext(t *testing.T, origin workspace.Origin, remote bool) (ctx context.Context, root, outside string) {
	t.Helper()
	root, outside = pathFixture(t)
	parent := filepath.Dir(root)
	for path, content := range map[string]string{
		filepath.Join(outside, "keys", "id_rsa"):     "SSH-KEY-MARKER",
		filepath.Join(outside, "keys", "id_rsa.pub"): "ssh-rsa PUBLIC",
		filepath.Join(root, ".ssh", "id_ed25519"):    "SSH-KEY-MARKER",
	} {
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
	cfg.OutsideAccess.Remote = &remote
	policy, problems := workspace.NewAccessPolicy(cfg)
	if len(problems) > 0 {
		t.Fatal(problems)
	}
	return WithGuard(context.Background(), policy.Guard(origin, "test")), root, outside
}

func TestToolsFullAccessReadWriteOutside(t *testing.T) {
	for _, origin := range []workspace.Origin{workspace.OriginLocal, workspace.OriginRemote} {
		ctx, root, outside := fullAccessContext(t, origin, true)
		secret := filepath.Join(outside, "secret.txt")

		r, _, _ := ReadFile(ctx, nil, ReadInput{Path: secret}, root)
		if r.IsError || !strings.Contains(resultText(r), "TOP-SECRET-MARKER") {
			t.Fatalf("%s read outside: %s", origin, resultText(r))
		}
		w, _, _ := WriteFile(ctx, nil, WriteInput{Path: filepath.Join(outside, "new", "file.txt"), Content: "written"}, root)
		if w.IsError {
			t.Fatalf("%s write outside: %s", origin, resultText(w))
		}
		e, _, _ := EditFile(ctx, nil, EditInput{Path: secret, Edits: []EditBlock{{OldText: "TOP", NewText: "EDITED"}}}, root)
		if e.IsError {
			t.Fatalf("%s edit outside: %s", origin, resultText(e))
		}
		mv, _, _ := MovePath(ctx, nil, MoveInput{SourcePath: filepath.Join(outside, "new", "file.txt"), TargetPath: "moved-in.txt"}, root)
		if mv.IsError {
			t.Fatalf("%s move into workspace: %s", origin, resultText(mv))
		}
		b, _, _ := RunBash(ctx, nil, BashInput{Command: "echo ok", WorkingDirectory: outside}, root)
		if b.IsError {
			t.Fatalf("%s bash in outside folder: %s", origin, resultText(b))
		}

		g, _, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "MARKER", Path: outside}, root)
		text := resultText(g)
		if strings.Contains(text, "SSH-KEY-MARKER") || !strings.Contains(text, "EDITED-SECRET-MARKER") {
			t.Fatalf("%s grep outside: %s", origin, text)
		}
		ls, _, _ := ListDirectory(ctx, nil, LsInput{Path: filepath.Join(outside, "keys")}, root)
		if text := resultText(ls); strings.Contains(text, "id_rsa\n") || strings.HasSuffix(text, "id_rsa") || !strings.Contains(text, "id_rsa.pub") {
			t.Fatalf("%s ls keys: %s", origin, text)
		}
		// The deny list also applies inside the workspace in full mode.
		r, _, _ = ReadFile(ctx, nil, ReadInput{Path: ".ssh/id_ed25519"}, root)
		requireToolError(t, "read private key in workspace", r)
		g, _, _ = GrepFiles(ctx, nil, GrepInput{Pattern: "SSH-KEY"}, root)
		if strings.Contains(resultText(g), "SSH-KEY-MARKER") {
			t.Fatalf("grep in workspace leaked a private key: %s", resultText(g))
		}
		mv, _, _ = MovePath(ctx, nil, MoveInput{SourcePath: "..", TargetPath: filepath.Join(outside, "stolen")}, root)
		requireToolError(t, "move parent of workspace", mv)
	}
}

func TestToolsFullAccessRemoteDisabled(t *testing.T) {
	ctx, root, outside := fullAccessContext(t, workspace.OriginRemote, false)
	r, _, _ := ReadFile(ctx, nil, ReadInput{Path: filepath.Join(outside, "secret.txt")}, root)
	requireToolError(t, "remote read with remote:false", r)
	ctx, root, outside = fullAccessContext(t, workspace.OriginLocal, false)
	r, _, _ = ReadFile(ctx, nil, ReadInput{Path: filepath.Join(outside, "secret.txt")}, root)
	if r.IsError {
		t.Fatalf("local read with remote:false: %s", resultText(r))
	}
}
