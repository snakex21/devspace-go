package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// benchTree creates dirs*files source-like files of ~4 KiB each.
func benchTree(b *testing.B, dirs, files int) string {
	b.Helper()
	root := b.TempDir()
	var content strings.Builder
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&content, "func handler%d(ctx context.Context) error { return process(ctx, %d) }\n", i, i)
	}
	content.WriteString("// needle: TODO_BENCH_MARKER\n")
	data := []byte(content.String())
	for d := 0; d < dirs; d++ {
		dir := filepath.Join(root, fmt.Sprintf("pkg%03d", d))
		if err := os.MkdirAll(dir, 0o700); err != nil {
			b.Fatal(err)
		}
		for f := 0; f < files; f++ {
			if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("file%03d.go", f)), data, 0o600); err != nil {
				b.Fatal(err)
			}
		}
	}
	return root
}

func BenchmarkGrepFiles(b *testing.B) {
	root := benchTree(b, 20, 20) // 400 files, ~3 MiB
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "TODO_BENCH_MARKER"}, root)
		if result.IsError {
			b.Fatal(resultText(result))
		}
	}
}

func BenchmarkFindFiles(b *testing.B) {
	root := benchTree(b, 20, 20)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindFiles(ctx, nil, GlobInput{Pattern: "*.go"}, root)
	}
}

func BenchmarkReadFile(b *testing.B) {
	root := benchTree(b, 1, 1)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadFile(ctx, nil, ReadInput{Path: "pkg000/file000.go"}, root)
	}
}
