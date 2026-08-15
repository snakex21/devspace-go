package tools

import (
	"reflect"
	"testing"
)

func TestShellCommandForOS(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		goos     string
		wantName string
		wantArgs []string
	}{
		{
			name:     "linux auto uses bash",
			shell:    "auto",
			goos:     "linux",
			wantName: "bash",
			wantArgs: []string{"-c", "pwd"},
		},
		{
			name:     "linux explicit sh",
			shell:    "sh",
			goos:     "linux",
			wantName: "sh",
			wantArgs: []string{"-c", "pwd"},
		},
		{
			name:     "windows auto stays powershell",
			shell:    "auto",
			goos:     "windows",
			wantName: "powershell.exe",
			wantArgs: []string{"-NoProfile", "-NonInteractive", "-Command", "pwd"},
		},
		{
			name:     "windows cmd stays cmd",
			shell:    "cmd",
			goos:     "windows",
			wantName: "cmd.exe",
			wantArgs: []string{"/C", "pwd"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotName, gotArgs := shellCommandForOS(test.shell, test.goos, "pwd")
			if gotName != test.wantName {
				t.Fatalf("shell name = %q, want %q", gotName, test.wantName)
			}
			if !reflect.DeepEqual(gotArgs, test.wantArgs) {
				t.Fatalf("shell args = %#v, want %#v", gotArgs, test.wantArgs)
			}
		})
	}
}

func TestUnixToolSkipsDoNotChangeWindows(t *testing.T) {
	if shouldSkipToolDirectory(`/workspace/.gradle`, `/workspace`, ".gradle", "windows") {
		t.Fatal("Windows unexpectedly skips Unix-only directory")
	}
	if !shouldSkipToolDirectory(`/workspace/.gradle`, `/workspace`, ".gradle", "linux") {
		t.Fatal("Linux should skip Unix-only directory below the search root")
	}
	if shouldSkipToolDirectory(`/workspace/.gradle`, `/workspace/.gradle`, ".gradle", "linux") {
		t.Fatal("Linux should allow an explicitly selected Unix-only directory")
	}
}
