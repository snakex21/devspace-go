package main

import (
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"

	devconfig "github.com/snakex21/devspace-go/internal/config"
)

func TestFolderDialogSizeTracksParentWindow(t *testing.T) {
	tests := []struct {
		name   string
		parent fyne.Size
		want   fyne.Size
	}{
		{name: "default window", parent: fyne.NewSize(640, 680), want: fyne.NewSize(640, 510)},
		{name: "medium window", parent: fyne.NewSize(1000, 800), want: fyne.NewSize(750, 600)},
		{name: "maximized window", parent: fyne.NewSize(2000, 1200), want: fyne.NewSize(1200, 800)},
		{name: "small window", parent: fyne.NewSize(500, 400), want: fyne.NewSize(500, 400)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := folderDialogSize(test.parent); got != test.want {
				t.Fatalf("folderDialogSize(%v) = %v, want %v", test.parent, got, test.want)
			}
		})
	}
}

func TestParseAndFormatExtraRoots(t *testing.T) {
	abs := t.TempDir()
	withSpace := filepath.Join(abs, "my folder")
	text := "ro " + abs + "\r\n\r\n# comment\nrw " + withSpace + "\n" + abs + "\nreadwrite " + abs
	roots, err := parseExtraRoots(text)
	if err != nil {
		t.Fatal(err)
	}
	want := []devconfig.ExtraRootConfig{
		{Path: abs, Access: devconfig.AccessRead},
		{Path: withSpace, Access: devconfig.AccessReadWrite},
		{Path: abs, Access: devconfig.AccessRead},
		{Path: abs, Access: devconfig.AccessReadWrite},
	}
	if len(roots) != len(want) {
		t.Fatalf("got %+v", roots)
	}
	for i := range want {
		if roots[i] != want[i] {
			t.Fatalf("root %d = %+v, want %+v", i, roots[i], want[i])
		}
	}
	again, err := parseExtraRoots(formatExtraRoots(roots))
	if err != nil || len(again) != len(roots) {
		t.Fatalf("round trip: %+v %v", again, err)
	}
	for _, bad := range []string{"ro relative/path", "rw", "C:"} {
		if _, err := parseExtraRoots(bad); err == nil {
			t.Errorf("parseExtraRoots(%q) accepted a non-absolute path", bad)
		}
	}
}
