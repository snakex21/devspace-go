package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadCapsResponseSize(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "wide.txt")
	line := strings.Repeat("x", 512*1024) + "\n" // 512 KiB lines
	if err := os.WriteFile(path, []byte(strings.Repeat(line, 40)), 0o600); err != nil {
		t.Fatal(err)
	}
	result, output, _ := ReadFile(context.Background(), nil, ReadInput{Path: "wide.txt", Limit: 5000}, root)
	if result.IsError {
		t.Fatal(resultText(result))
	}
	if len(output.Result) > maxReadOutputBytes+256 {
		t.Fatalf("read returned %d bytes, cap is %d", len(output.Result), maxReadOutputBytes)
	}
	if !strings.Contains(output.Result, "use offset=") {
		t.Fatalf("missing continuation hint: %q", output.Result[len(output.Result)-120:])
	}
}

func TestReadDoesNotClaimTruncationAtExactLimit(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "three.txt"), []byte("a\nb\nc\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, output, _ := ReadFile(context.Background(), nil, ReadInput{Path: "three.txt", Limit: 3}, root)
	if output.Result != "a\nb\nc" {
		t.Fatalf("result = %q", output.Result)
	}
	_, output, _ = ReadFile(context.Background(), nil, ReadInput{Path: "three.txt", Limit: 2}, root)
	if !strings.HasPrefix(output.Result, "a\nb\n\n[truncated after 2 lines") {
		t.Fatalf("result = %q", output.Result)
	}
	_, output, _ = ReadFile(context.Background(), nil, ReadInput{Path: "three.txt", Offset: 2, Limit: 1}, root)
	if !strings.HasPrefix(output.Result, "b\n\n[truncated") {
		t.Fatalf("offset result = %q", output.Result)
	}
}

func TestEditRejectsHugeFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "huge.bin")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// Sparse where supported: no 64 MiB write needed.
	if err := file.Truncate(maxEditFileBytes + 1); err != nil {
		t.Fatal(err)
	}
	file.Close()
	result, _, _ := EditFile(context.Background(), nil, EditInput{Path: "huge.bin", Edits: []EditBlock{{OldText: "a", NewText: "b"}}}, root)
	if !result.IsError || !strings.Contains(resultText(result), "too large") {
		t.Fatalf("result = %s", resultText(result))
	}
}
