package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func searchFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"src/main.go":               "package main\n// FINDME one\nfunc main() {}\n",
		"src/util.txt":              "nothing\nFINDME two\n",
		"node_modules/dep/index.js": "FINDME vendored\n",
		"build/out.go":              "FINDME generated\n",
		"image.png":                 "FINDME binary\n",
	}
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestGrepOutputFormatAndFilters(t *testing.T) {
	root := searchFixture(t)
	ctx := context.Background()
	_, output, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "FINDME"}, root)
	got := output.Result
	for _, want := range []string{"src/main.go:2:// FINDME one", "src/util.txt:2:FINDME two", "[skipped 1 large/binary files]"} {
		if !strings.Contains(got, want) {
			t.Errorf("grep output missing %q:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"vendored", "generated"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("grep searched a skipped directory (%s):\n%s", unwanted, got)
		}
	}
	_, output, _ = GrepFiles(ctx, nil, GrepInput{Pattern: "FINDME", Include: "*.go"}, root)
	if strings.Contains(output.Result, "util.txt") || !strings.Contains(output.Result, "main.go") {
		t.Errorf("include filter ignored:\n%s", output.Result)
	}
}

func TestSearchInsideExplicitlyRequestedSkippedDirectory(t *testing.T) {
	root := searchFixture(t)
	ctx := context.Background()
	_, grepOut, _ := GrepFiles(ctx, nil, GrepInput{Pattern: "FINDME", Path: "build"}, root)
	if !strings.Contains(grepOut.Result, "build/out.go:1:FINDME generated") {
		t.Fatalf("grep in path=build found nothing:\n%s", grepOut.Result)
	}
	_, globOut, _ := FindFiles(ctx, nil, GlobInput{Pattern: "*.js", Path: "node_modules"}, root)
	if !strings.Contains(globOut.Result, "node_modules/dep/index.js") {
		t.Fatalf("glob in path=node_modules found nothing:\n%s", globOut.Result)
	}
}

func TestGlobMatchesNamesAndSkipsBinaries(t *testing.T) {
	root := searchFixture(t)
	_, output, _ := FindFiles(context.Background(), nil, GlobInput{Pattern: "*"}, root)
	lines := strings.Split(output.Result, "\n")
	if len(lines) != 2 || !strings.Contains(output.Result, "src/main.go") || !strings.Contains(output.Result, "src/util.txt") {
		t.Fatalf("glob output:\n%s", output.Result)
	}
}

func TestGrepFindsLongLines(t *testing.T) {
	root := t.TempDir()
	long := strings.Repeat("a", 200*1024) + "NEEDLE"
	if err := os.WriteFile(filepath.Join(root, "long.txt"), []byte("x\n"+long+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A second file after the long one must still be scanned with a sane buffer.
	if err := os.WriteFile(filepath.Join(root, "z.txt"), []byte("NEEDLE\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, output, _ := GrepFiles(context.Background(), nil, GrepInput{Pattern: "NEEDLE"}, root)
	if !strings.Contains(output.Result, "long.txt:2:") || !strings.Contains(output.Result, "z.txt:1:NEEDLE") {
		t.Fatalf("output (%d bytes) missing matches", len(output.Result))
	}
}
