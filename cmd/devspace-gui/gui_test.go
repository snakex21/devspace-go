package main

import (
	"testing"

	"fyne.io/fyne/v2"
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
