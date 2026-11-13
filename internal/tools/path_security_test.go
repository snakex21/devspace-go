package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// pathFixture creates parent/{ws,outside} with a secret outside the workspace.
func pathFixture(t *testing.T) (root, outside string) {
	t.Helper()
	parent := t.TempDir()
	root = filepath.Join(parent, "ws")
	outside = filepath.Join(parent, "outside")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("TOP-SECRET-MARKER"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "inside.txt"), []byte("inside"), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, outside
}

func resultText(result *mcp.CallToolResult) string {
	var parts []string
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			parts = append(parts, text.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func requireToolError(t *testing.T, name string, result *mcp.CallToolResult) {
	t.Helper()
	if result == nil || !result.IsError {
		t.Fatalf("%s: expected a tool error, got %#v", name, result)
	}
	if strings.Contains(resultText(result), "TOP-SECRET-MARKER") {
		t.Fatalf("%s leaked outside content: %s", name, resultText(result))
	}
}

func TestToolsRejectPathsOutsideWorkspace(t *testing.T) {
	root, outside := pathFixture(t)
	ctx := context.Background()
	escapes := []string{"../outside", filepath.Join(outside), "../outside/secret.txt"}

	for _, escape := range escapes {
		r, _, _ := ReadFile(ctx, nil, ReadInput{Path: filepath.Join(escape, "secret.txt")}, root)
		requireToolError(t, "read "+escape, r)

		r2, _, _ := ListDirectory(ctx, nil, LsInput{Path: escape}, root)
		requireToolError(t, "ls "+escape, r2)

		r3, _, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "TOP-SECRET", Path: escape}, root)
		requireToolError(t, "grep "+escape, r3)

		r4, _, _ := FindFiles(ctx, nil, GlobInput{Pattern: "*.txt", Path: escape}, root)
		requireToolError(t, "glob "+escape, r4)

		r5, _, _ := MakeDirectory(ctx, nil, MkdirInput{Path: filepath.Join(escape, "made")}, root)
		requireToolError(t, "mkdir "+escape, r5)

		r6, _, _ := WriteFile(ctx, nil, WriteInput{Path: filepath.Join(escape, "written.txt"), Content: "x"}, root)
		requireToolError(t, "write "+escape, r6)

		r7, _, _ := EditFile(ctx, nil, EditInput{Path: filepath.Join(escape, "secret.txt"), Edits: []EditBlock{{OldText: "TOP", NewText: "pwned"}}}, root)
		requireToolError(t, "edit "+escape, r7)

		r8, _, _ := MovePath(ctx, nil, MoveInput{SourcePath: "inside.txt", TargetPath: filepath.Join(escape, "moved.txt")}, root)
		requireToolError(t, "move target "+escape, r8)

		r9, _, _ := MovePath(ctx, nil, MoveInput{SourcePath: filepath.Join(escape, "secret.txt"), TargetPath: "stolen.txt"}, root)
		requireToolError(t, "move source "+escape, r9)

		r10, _, _ := RunBash(ctx, nil, BashInput{Command: "echo hi", WorkingDirectory: escape}, root)
		requireToolError(t, "bash cwd "+escape, r10)
	}

	for _, unexpected := range []string{"made", "written.txt", "moved.txt"} {
		if _, err := os.Stat(filepath.Join(outside, unexpected)); err == nil {
			t.Fatalf("tool created %s outside the workspace", unexpected)
		}
	}
	if data, _ := os.ReadFile(filepath.Join(outside, "secret.txt")); string(data) != "TOP-SECRET-MARKER" {
		t.Fatalf("outside file modified: %q", data)
	}
	if _, err := os.Stat(filepath.Join(root, "inside.txt")); err != nil {
		t.Fatalf("inside file was moved: %v", err)
	}
}

func TestToolsAcceptAbsolutePathInsideWorkspace(t *testing.T) {
	root, _ := pathFixture(t)
	// Before the fix, an absolute path passed validation but was then joined
	// onto the root a second time (root/root/...), so reads failed.
	result, output, _ := ReadFile(context.Background(), nil, ReadInput{Path: filepath.Join(root, "inside.txt")}, root)
	if result.IsError || output.Result != "inside" {
		t.Fatalf("absolute read failed: %s", resultText(result))
	}
}

func TestMoveRefusesToMoveOrOverwriteWorkspaceRoot(t *testing.T) {
	root, _ := pathFixture(t)
	ctx := context.Background()
	for _, input := range []MoveInput{
		{SourcePath: "inside.txt", TargetPath: ".", Overwrite: true},
		{SourcePath: "inside.txt", TargetPath: "", Overwrite: true},
		{SourcePath: ".", TargetPath: "renamed"},
	} {
		result, _, _ := MovePath(ctx, nil, input, root)
		if !result.IsError {
			t.Fatalf("move %+v succeeded", input)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "inside.txt")); err != nil {
		t.Fatalf("workspace content destroyed: %v", err)
	}
}

func TestMoveRefusesToMoveDirectoryIntoItself(t *testing.T) {
	root, _ := pathFixture(t)
	if err := os.Mkdir(filepath.Join(root, "dir"), 0o700); err != nil {
		t.Fatal(err)
	}
	result, _, _ := MovePath(context.Background(), nil, MoveInput{SourcePath: "dir", TargetPath: "dir/sub/dir"}, root)
	if !result.IsError {
		t.Fatal("moving a directory into itself succeeded")
	}
	if _, err := os.Stat(filepath.Join(root, "dir")); err != nil {
		t.Fatalf("directory lost: %v", err)
	}
}

func TestToolsRejectLinkEscape(t *testing.T) {
	root, outside := pathFixture(t)
	link := filepath.Join(root, "link")
	if err := os.Symlink(outside, link); err != nil {
		if runtime.GOOS != "windows" {
			t.Skipf("cannot create symlink: %v", err)
		}
		// Unprivileged Windows processes can still create junctions.
		if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, outside).CombinedOutput(); err != nil {
			t.Skipf("cannot create junction: %v %s", err, output)
		}
	}
	result, _, _ := ReadFile(context.Background(), nil, ReadInput{Path: "link/secret.txt"}, root)
	requireToolError(t, "read through link", result)
	result2, _, _ := WriteFile(context.Background(), nil, WriteInput{Path: "link/new.txt", Content: "x"}, root)
	requireToolError(t, "write through link", result2)
	if _, err := os.Stat(filepath.Join(outside, "new.txt")); err == nil {
		t.Fatal("write followed a link outside the workspace")
	}
}
