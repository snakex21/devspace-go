package tools

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snakex21/devspace-go/internal/procutil"
	"github.com/snakex21/devspace-go/internal/workspace"
)

// resolveToolPath maps a tool path argument to an absolute path inside the
// workspace. Every tool must use it instead of joining paths itself, so that
// "..", absolute paths and links cannot escape the workspace root.
func resolveToolPath(wsRoot, inputPath string) (string, *mcp.CallToolResult) {
	absPath, err := workspace.ResolveInRoot(wsRoot, inputPath)
	if err != nil {
		return "", errorResult(err)
	}
	return absPath, nil
}

func errorResult(err error) *mcp.CallToolResult {
	result := &mcp.CallToolResult{}
	result.SetError(err)
	return result
}

const (
	defaultReadLimitLines = 2000
	maxReadLimitLines     = 5000
	maxToolOutputBytes    = 256 * 1024
	maxSearchFileBytes    = 2 * 1024 * 1024
	maxSearchMatches      = 500
	maxReadOutputBytes    = 4 * 1024 * 1024
	maxEditFileBytes      = 64 * 1024 * 1024
)

var skippedDirs = map[string]bool{
	".git": true, ".hg": true, ".svn": true,
	"node_modules": true, "dist": true, "build": true,
	".next": true, ".turbo": true, ".cache": true,
	".devspace": true, ".devspace-state": true, ".webcoder": true, ".webcoder-state": true,
	"tools": true,
}

var unixSkippedDirs = map[string]bool{
	".android":  true,
	".bun":      true,
	".config":   true,
	".gradle":   true,
	".hermes":   true,
	".local":    true,
	".npm":      true,
	".opencode": true,
	".var":      true,
}

var configuredShell = "auto"

// SetShell configures the shell used by the bash tool.
// Accepted values: auto, powershell, cmd, bash, sh.
func SetShell(shell string) {
	configuredShell = strings.ToLower(strings.TrimSpace(shell))
	if configuredShell == "" {
		configuredShell = "auto"
	}
}

// ReadInput represents the input for the read tool.
type ReadInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Path        string `json:"path" jsonschema:"File path to read, relative to the workspace root."`
	Offset      int    `json:"offset,omitempty" jsonschema:"1-indexed line number to start reading from."`
	Limit       int    `json:"limit,omitempty" jsonschema:"Maximum number of lines to read."`
}

// ReadOutput represents the output for the read tool.
type ReadOutput struct {
	Result string `json:"result" jsonschema:"File contents."`
}

// ReadFile reads a file and returns its content.
func ReadFile(ctx context.Context, req *mcp.CallToolRequest, input ReadInput, wsRoot string) (*mcp.CallToolResult, ReadOutput, error) {
	absPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, ReadOutput{}, nil
	}

	info, err := os.Stat(absPath)
	if err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("file not found: %s", input.Path))
		return result, ReadOutput{}, nil
	}

	if info.IsDir() {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("path is a directory, not a file: %s", input.Path))
		return result, ReadOutput{}, nil
	}
	if info.Size() > maxSearchFileBytes && input.Limit == 0 {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("file is too large (%d bytes). Read with offset/limit or inspect with shell", info.Size()))
		return result, ReadOutput{}, nil
	}

	file, err := os.Open(absPath)
	if err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to open file: %v", err))
		return result, ReadOutput{}, nil
	}
	defer file.Close()

	var content string

	if input.Offset <= 0 {
		input.Offset = 1
	}
	if input.Limit <= 0 {
		input.Limit = defaultReadLimitLines
	}
	if input.Limit > maxReadLimitLines {
		input.Limit = maxReadLimitLines
	}

	if input.Offset > 0 || input.Limit > 0 {
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		lineNum := 0
		linesRead := 0
		var builder strings.Builder
		truncatedLines := false
		truncatedBytes := false

		for scanner.Scan() {
			lineNum++
			if input.Offset > 0 && lineNum < input.Offset {
				continue
			}
			if input.Limit > 0 && linesRead >= input.Limit {
				// Only report truncation if there really is another line.
				truncatedLines = true
				break
			}
			line := scanner.Bytes()
			// 5000 lines of up to 1 MiB each could otherwise build a
			// multi-gigabyte response.
			if builder.Len()+len(line)+1 > maxReadOutputBytes {
				truncatedBytes = true
				break
			}
			if linesRead > 0 {
				builder.WriteByte('\n')
			}
			builder.Write(line)
			linesRead++
		}

		content = builder.String()
		switch {
		case truncatedBytes:
			content += fmt.Sprintf("\n\n[truncated after %d lines (%d bytes); use offset=%d to read more]", linesRead, builder.Len(), input.Offset+linesRead)
		case truncatedLines:
			content += fmt.Sprintf("\n\n[truncated after %d lines; use offset/limit to read more]", input.Limit)
		}
		if scanner.Err() != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("failed to read file: %v", scanner.Err()))
			return result, ReadOutput{}, nil
		}
	} else {
		data, err := io.ReadAll(file)
		if err != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("failed to read file: %v", err))
			return result, ReadOutput{}, nil
		}
		content = string(data)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: content},
		},
	}, ReadOutput{Result: content}, nil
}

// WriteInput represents the input for the write tool.
type WriteInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Path        string `json:"path" jsonschema:"File path to write, relative to the workspace root."`
	Content     string `json:"content" jsonschema:"Complete new file content."`
}

// WriteOutput represents the output for the write tool.
type WriteOutput struct {
	Result string `json:"result" jsonschema:"Write result message."`
}

// MkdirInput represents a directory creation operation.
type MkdirInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Path        string `json:"path" jsonschema:"Directory path to create, relative to the workspace root."`
}

// MkdirOutput represents the output for mkdir.
type MkdirOutput struct {
	Result string `json:"result" jsonschema:"Directory creation result message."`
}

// MoveInput represents a file or directory move/rename operation.
type MoveInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	SourcePath  string `json:"sourcePath" jsonschema:"Existing file or directory path, relative to the workspace root."`
	TargetPath  string `json:"targetPath" jsonschema:"Destination file or directory path, relative to the workspace root."`
	Overwrite   bool   `json:"overwrite,omitempty" jsonschema:"Overwrite destination if it already exists. Default false."`
}

// MoveOutput represents the output for move.
type MoveOutput struct {
	Result string `json:"result" jsonschema:"Move result message."`
}

// WriteFile creates or overwrites a file.
func WriteFile(ctx context.Context, req *mcp.CallToolRequest, input WriteInput, wsRoot string) (*mcp.CallToolResult, WriteOutput, error) {
	absPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, WriteOutput{}, nil
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to create directory: %v", err))
		return result, WriteOutput{}, nil
	}

	if err := os.WriteFile(absPath, []byte(input.Content), 0644); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to write file: %v", err))
		return result, WriteOutput{}, nil
	}

	result := fmt.Sprintf("Wrote %s (%d bytes).", input.Path, len(input.Content))
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, WriteOutput{Result: result}, nil
}

// MakeDirectory creates a directory and all missing parents inside a workspace.
func MakeDirectory(ctx context.Context, req *mcp.CallToolRequest, input MkdirInput, wsRoot string) (*mcp.CallToolResult, MkdirOutput, error) {
	absPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, MkdirOutput{}, nil
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to create directory: %v", err))
		return result, MkdirOutput{}, nil
	}

	result := fmt.Sprintf("Created directory %s.", input.Path)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: result}},
	}, MkdirOutput{Result: result}, nil
}

// MovePath moves or renames a file/directory inside a workspace.
func MovePath(ctx context.Context, req *mcp.CallToolRequest, input MoveInput, wsRoot string) (*mcp.CallToolResult, MoveOutput, error) {
	sourceAbs, failure := resolveToolPath(wsRoot, input.SourcePath)
	if failure != nil {
		return failure, MoveOutput{}, nil
	}
	targetAbs, failure := resolveToolPath(wsRoot, input.TargetPath)
	if failure != nil {
		return failure, MoveOutput{}, nil
	}
	// Moving the root away, or overwriting it (RemoveAll), would destroy the
	// whole workspace.
	cleanRoot := filepath.Clean(wsRoot)
	if workspace.IsPathInsideRoot(cleanRoot, sourceAbs) || workspace.IsPathInsideRoot(cleanRoot, targetAbs) {
		return errorResult(fmt.Errorf("cannot move the workspace root or one of its parents")), MoveOutput{}, nil
	}
	if workspace.IsPathInsideRoot(targetAbs, sourceAbs) && !samePath(targetAbs, sourceAbs) {
		return errorResult(fmt.Errorf("cannot move %s into itself", input.SourcePath)), MoveOutput{}, nil
	}

	if _, err := os.Stat(sourceAbs); err != nil {
		result := &mcp.CallToolResult{}
		if os.IsNotExist(err) {
			result.SetError(fmt.Errorf("source not found: %s", input.SourcePath))
		} else {
			result.SetError(fmt.Errorf("failed to stat source: %v", err))
		}
		return result, MoveOutput{}, nil
	}

	if _, err := os.Stat(targetAbs); err == nil {
		if !input.Overwrite {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("target already exists: %s", input.TargetPath))
			return result, MoveOutput{}, nil
		}
		if err := os.RemoveAll(targetAbs); err != nil {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("failed to remove existing target: %v", err))
			return result, MoveOutput{}, nil
		}
	} else if !os.IsNotExist(err) {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to stat target: %v", err))
		return result, MoveOutput{}, nil
	}

	if err := os.MkdirAll(filepath.Dir(targetAbs), 0755); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to create target parent directory: %v", err))
		return result, MoveOutput{}, nil
	}

	if err := os.Rename(sourceAbs, targetAbs); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to move path: %v", err))
		return result, MoveOutput{}, nil
	}

	result := fmt.Sprintf("Moved %s to %s.", input.SourcePath, input.TargetPath)
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: result}},
	}, MoveOutput{Result: result}, nil
}

// EditInput represents an edit operation.
type EditInput struct {
	WorkspaceID string      `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Path        string      `json:"path" jsonschema:"File path to edit, relative to the workspace root."`
	Edits       []EditBlock `json:"edits" jsonschema:"Array of edit operations."`
}

// EditBlock represents a single find-and-replace operation.
type EditBlock struct {
	OldText string `json:"oldText" jsonschema:"Exact text to replace. Must match uniquely."`
	NewText string `json:"newText" jsonschema:"Replacement text."`
}

// EditOutput represents the output for the edit tool.
type EditOutput struct {
	Status string `json:"status" jsonschema:"Edit status."`
	Result string `json:"result" jsonschema:"Edit result message."`
}

// EditFile performs find-and-replace edits on a file.
func EditFile(ctx context.Context, req *mcp.CallToolRequest, input EditInput, wsRoot string) (*mcp.CallToolResult, EditOutput, error) {
	absPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, EditOutput{}, nil
	}

	if info, err := os.Stat(absPath); err == nil && info.Size() > maxEditFileBytes {
		return errorResult(fmt.Errorf("file is too large to edit (%d bytes, limit %d)", info.Size(), maxEditFileBytes)), EditOutput{}, nil
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		result := &mcp.CallToolResult{}
		if os.IsNotExist(err) {
			result.SetError(fmt.Errorf("file not found: %s", input.Path))
		} else {
			result.SetError(fmt.Errorf("failed to read file: %v", err))
		}
		return result, EditOutput{}, nil
	}

	content := string(data)

	for i, edit := range input.Edits {
		count := strings.Count(content, edit.OldText)
		if count == 0 {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("edit %d: oldText not found in file", i+1))
			return result, EditOutput{}, nil
		}
		if count > 1 {
			result := &mcp.CallToolResult{}
			result.SetError(fmt.Errorf("edit %d: oldText matches %d times, must be unique", i+1, count))
			return result, EditOutput{}, nil
		}

		content = strings.Replace(content, edit.OldText, edit.NewText, 1)
	}

	if content == string(data) {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("no changes made to file"))
		return result, EditOutput{}, nil
	}

	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to write file: %v", err))
		return result, EditOutput{}, nil
	}

	result := fmt.Sprintf("Edited %s: %d edit(s) applied.", input.Path, len(input.Edits))
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, EditOutput{Status: "applied", Result: result}, nil
}

// GrepInput represents the input for the grep tool.
type GrepInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Pattern     string `json:"pattern" jsonschema:"Search pattern."`
	Path        string `json:"path,omitempty" jsonschema:"Optional path scope relative to the workspace root."`
	Include     string `json:"include,omitempty" jsonschema:"Optional include glob."`
}

// GrepOutput represents the output for the grep tool.
type GrepOutput struct {
	Result string `json:"result" jsonschema:"Grep results."`
}

// GrepFiles searches file contents for a pattern.
func GrepFiles(ctx context.Context, req *mcp.CallToolRequest, input GrepInput, wsRoot string) (*mcp.CallToolResult, GrepOutput, error) {
	searchPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, GrepOutput{}, nil
	}

	re, err := regexp.Compile(input.Pattern)
	if err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("invalid regex pattern: %v", err))
		return result, GrepOutput{}, nil
	}

	include := input.Include
	var lines []string
	matches := 0
	skipped := 0
	start := time.Now()

	_ = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if time.Since(start) > 20*time.Second || matches >= maxSearchMatches {
			return filepath.SkipAll
		}
		if info.IsDir() {
			if shouldSkipToolDirectory(path, searchPath, info.Name(), runtime.GOOS) {
				return filepath.SkipDir
			}
			return nil
		}
		if include != "" {
			ok, _ := filepath.Match(include, info.Name())
			if !ok {
				return nil
			}
		}
		if info.Size() > maxSearchFileBytes || looksBinaryPath(path) {
			skipped++
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			text := scanner.Text()
			if re.MatchString(text) {
				rel, _ := filepath.Rel(wsRoot, path)
				lines = append(lines, fmt.Sprintf("%s:%d:%s", filepath.ToSlash(rel), lineNo, text))
				matches++
				if matches >= maxSearchMatches {
					break
				}
			}
		}
		return nil
	})

	output := strings.Join(lines, "\n")
	if output == "" {
		output = fmt.Sprintf("No matches found for pattern: %s", input.Pattern)
	} else if matches >= maxSearchMatches {
		output += fmt.Sprintf("\n\n[truncated after %d matches]", maxSearchMatches)
	}
	if skipped > 0 {
		output += fmt.Sprintf("\n[skipped %d large/binary files]", skipped)
	}
	output = truncateOutput(output)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
	}, GrepOutput{Result: output}, nil
}

// GlobInput represents the input for the glob tool.
type GlobInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Pattern     string `json:"pattern" jsonschema:"File glob pattern."`
	Path        string `json:"path,omitempty" jsonschema:"Optional path scope relative to the workspace root."`
}

// GlobOutput represents the output for the glob tool.
type GlobOutput struct {
	Result string `json:"result" jsonschema:"Glob results."`
}

// FindFiles finds files matching a glob pattern.
func FindFiles(ctx context.Context, req *mcp.CallToolRequest, input GlobInput, wsRoot string) (*mcp.CallToolResult, GlobOutput, error) {
	searchPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, GlobOutput{}, nil
	}

	var matches []string

	filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if ctx.Err() != nil || len(matches) >= maxSearchMatches {
			return filepath.SkipAll
		}
		if info.IsDir() {
			if shouldSkipToolDirectory(path, searchPath, info.Name(), runtime.GOOS) {
				return filepath.SkipDir
			}
			return nil
		}
		if looksBinaryPath(path) || info.Size() > maxSearchFileBytes {
			return nil
		}
		matched, _ := filepath.Match(input.Pattern, info.Name())
		if matched {
			rel, _ := filepath.Rel(wsRoot, path)
			matches = append(matches, filepath.ToSlash(rel))
		}
		return nil
	})

	result := strings.Join(matches, "\n")
	if result == "" {
		result = fmt.Sprintf("No files found matching: %s", input.Pattern)
	} else if len(matches) >= maxSearchMatches {
		result += fmt.Sprintf("\n\n[truncated after %d files]", maxSearchMatches)
	}
	result = truncateOutput(result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, GlobOutput{Result: result}, nil
}

// LsInput represents the input for the ls tool.
type LsInput struct {
	WorkspaceID string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Path        string `json:"path" jsonschema:"Directory path to list, relative to the workspace root."`
}

// LsOutput represents the output for the ls tool.
type LsOutput struct {
	Result string `json:"result" jsonschema:"Directory listing."`
}

// ListDirectory lists the contents of a directory.
func ListDirectory(ctx context.Context, req *mcp.CallToolRequest, input LsInput, wsRoot string) (*mcp.CallToolResult, LsOutput, error) {
	absPath, failure := resolveToolPath(wsRoot, input.Path)
	if failure != nil {
		return failure, LsOutput{}, nil
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("failed to read directory: %v", err))
		return result, LsOutput{}, nil
	}

	var lines []string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			lines = append(lines, entry.Name())
			continue
		}

		prefix := "-"
		if entry.IsDir() {
			prefix = "d"
		}

		size := info.Size()
		sizeStr := formatSize(size)
		modTime := info.ModTime().Format("Jan 02 15:04")
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}

		lines = append(lines, fmt.Sprintf("%s %8s %s %s", prefix, sizeStr, modTime, name))
	}

	result := truncateOutput(strings.Join(lines, "\n"))
	if result == "" {
		result = "Empty directory"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, LsOutput{Result: result}, nil
}

// BashInput represents the input for the bash tool.
type BashInput struct {
	WorkspaceID      string `json:"workspaceId" jsonschema:"Workspace identifier returned by open_workspace."`
	Command          string `json:"command" jsonschema:"Shell command to run."`
	WorkingDirectory string `json:"workingDirectory,omitempty" jsonschema:"Optional working directory relative to the workspace root."`
	Timeout          int    `json:"timeout,omitempty" jsonschema:"Timeout in seconds. Defaults to 30, max 300."`
}

// BashOutput represents the output for the bash tool.
type BashOutput struct {
	Result string `json:"result" jsonschema:"Shell command output."`
}

// RunBash executes a shell command. Uses PowerShell on Windows, bash on Unix.
func RunBash(ctx context.Context, req *mcp.CallToolRequest, input BashInput, wsRoot string) (*mcp.CallToolResult, BashOutput, error) {
	cwd, failure := resolveToolPath(wsRoot, input.WorkingDirectory)
	if failure != nil {
		return failure, BashOutput{}, nil
	}

	timeout := 30
	if input.Timeout > 0 {
		timeout = input.Timeout
		if timeout > 300 {
			timeout = 300
		}
	}

	cmdName, cmdArgs := shellCommandForOS(configuredShell, runtime.GOOS, input.Command)

	output, err := runCommand(ctx, cwd, cmdName, cmdArgs, timeout)
	if err != nil && output == "" {
		result := &mcp.CallToolResult{}
		result.SetError(fmt.Errorf("command failed: %v", err))
		return result, BashOutput{}, nil
	}

	result := output
	if err != nil && output != "" {
		if result != "" {
			result += "\n"
		}
		result += "[stderr] " + err.Error()
	}
	if result == "" {
		result = "(no output)"
	}
	result = truncateOutput(result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, BashOutput{Result: result}, nil
}

func shellCommandForOS(preferredShell, goos, command string) (string, []string) {
	if goos == "windows" {
		switch preferredShell {
		case "cmd", "cmd.exe":
			return "cmd.exe", []string{"/C", command}
		case "powershell", "powershell.exe", "pwsh":
			if preferredShell == "pwsh" {
				return "pwsh", []string{"-NoProfile", "-NonInteractive", "-Command", command}
			}
			return "powershell.exe", []string{"-NoProfile", "-NonInteractive", "-Command", command}
		default:
			return "powershell.exe", []string{"-NoProfile", "-NonInteractive", "-Command", command}
		}
	}

	switch preferredShell {
	case "", "auto":
		preferredShell = "bash"
	case "sh":
		preferredShell = "sh"
	}
	return preferredShell, []string{"-c", command}
}

func shouldSkipToolDirectory(path, searchRoot, name, goos string) bool {
	if skippedDirs[name] {
		return true
	}
	if goos == "windows" || !unixSkippedDirs[name] {
		return false
	}
	return filepath.Clean(path) != filepath.Clean(searchRoot)
}

// --- helpers ---

func runCommand(ctx context.Context, cwd, name string, args []string, timeoutSec int) (string, error) {
	timeout := time.Duration(timeoutSec) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(timeoutCtx, name, args...)
	cmd.Dir = cwd
	// Keep at most maxToolOutputBytes in memory; a chatty command must not be
	// able to exhaust the server's RAM. One writer for both streams makes
	// os/exec share a single pipe, preserving interleaving like CombinedOutput.
	output := &cappedBuffer{limit: maxToolOutputBytes}
	cmd.Stdout = output
	cmd.Stderr = output
	// A background process started by the command (e.g. "server &" or
	// "start /b") inherits the output pipe. Without WaitDelay, Wait would block
	// until that process exits, hanging the tool call until the timeout and
	// then forever. Stop collecting output shortly after the shell exits.
	cmd.WaitDelay = shellWaitDelay
	// On timeout or client disconnect, kill the shell together with every
	// process it started, not just the shell itself.
	tree := procutil.Prepare(cmd, false)
	if err := tree.Start(); err != nil {
		return "", err
	}
	err := cmd.Wait()
	tree.Release()

	switch {
	case errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) && ctx.Err() == nil:
		err = fmt.Errorf("command timed out after %s and was terminated", timeout)
	case ctx.Err() != nil:
		err = fmt.Errorf("command canceled: %w", ctx.Err())
	case errors.Is(err, exec.ErrWaitDelay):
		// The shell itself exited successfully; only its background children
		// still hold the output pipe. They keep running.
		err = nil
		output.note("\n[background processes are still running; their further output is not captured]")
	}
	return output.String(), err
}

const shellWaitDelay = 2 * time.Second

// cappedBuffer is an io.Writer that keeps the first limit bytes and counts
// the rest. os/exec calls Write from one goroutine per pipe; both streams
// share one pipe here, but a mutex keeps it safe if that ever changes.
type cappedBuffer struct {
	mu      sync.Mutex
	buf     []byte
	limit   int
	dropped int64
	suffix  string
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if room := b.limit - len(b.buf); room > 0 {
		if len(p) <= room {
			b.buf = append(b.buf, p...)
			return len(p), nil
		}
		b.buf = append(b.buf, p[:room]...)
		b.dropped += int64(len(p) - room)
		return len(p), nil
	}
	b.dropped += int64(len(p))
	return len(p), nil
}

func (b *cappedBuffer) note(text string) {
	b.mu.Lock()
	b.suffix += text
	b.mu.Unlock()
}

func (b *cappedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	text := string(b.buf)
	if b.dropped > 0 {
		text += fmt.Sprintf("\n\n[output truncated after %d bytes; %d more bytes discarded]", b.limit, b.dropped)
	}
	return text + b.suffix
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGT"[exp])
}

func samePath(a, b string) bool {
	return workspace.IsPathInsideRoot(a, b) && workspace.IsPathInsideRoot(b, a)
}

func escapePS(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func truncateOutput(s string) string {
	if len(s) <= maxToolOutputBytes {
		return s
	}
	return s[:maxToolOutputBytes] + fmt.Sprintf("\n\n[output truncated after %d bytes]", maxToolOutputBytes)
}

func looksBinaryPath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".exe", ".dll", ".so", ".dylib", ".bin", ".obj", ".o", ".a",
		".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf",
		".zip", ".tar", ".gz", ".7z", ".rar",
		".db", ".sqlite", ".sqlite3", ".wasm":
		return true
	default:
		return false
	}
}
