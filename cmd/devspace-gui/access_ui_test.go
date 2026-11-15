package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	devconfig "github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/locales"
)

func TestDenyListAddRemove(t *testing.T) {
	locales.Init("en")
	abs := t.TempDir()
	d := &denyList{}
	if err := d.add("relative/path"); err == nil {
		t.Fatal("relative deny entry accepted")
	}
	for _, path := range []string{abs, abs, filepath.Join(abs, "a"), "  ", filepath.Join(abs, "b") + string(filepath.Separator)} {
		if err := d.add(path); err != nil {
			t.Fatal(err)
		}
	}
	if len(d.items) != 3 {
		t.Fatalf("items: %v", d.items)
	}
	d.remove(1)
	d.remove(10)
	if len(d.items) != 2 || d.items[1] != filepath.Join(abs, "b") {
		t.Fatalf("after remove: %v", d.items)
	}
}

// TestAccessSectionRoundTrip drives the GUI widgets: the section shows the
// loaded config, edits are applied, saved through mergeGUIConfig and read
// back by config.LoadConfig exactly as a hand-edited file would be.
func TestAccessSectionRoundTrip(t *testing.T) {
	test.NewTempApp(t)
	locales.Init("en")
	locales.SetLocale("en")
	dir := t.TempDir()
	keep := filepath.Join(dir, "keep")
	off := false
	cfg := devconfig.DefaultConfig()
	cfg.ConfigDir = dir
	cfg.AllowedRoots = []string{dir}
	cfg.OutsideAccess = devconfig.OutsideAccessConfig{Deny: []string{keep}, Remote: &off}

	w := test.NewTempWindow(t, widget.NewLabel(""))
	remote := widget.NewCheck("", nil)
	remote.SetChecked(cfg.OutsideAccess.RemoteEnabled())
	section := newAccessSection(cfg, w, remote)

	if section.modeValue != devconfig.AccessModeFull || section.mode.Selected != accessModeLabel(devconfig.AccessModeFull) {
		t.Fatalf("mode shown: %q / %q", section.modeValue, section.mode.Selected)
	}
	if !section.denyDefaults.Checked || len(section.deny.items) != 1 || remote.Checked {
		t.Fatalf("loaded state: defaults %v deny %v remote %v", section.denyDefaults.Checked, section.deny.items, remote.Checked)
	}
	if !strings.Contains(section.defaultsInfo.Text, "SSH private keys") {
		t.Fatalf("built-in list not shown: %q", section.defaultsInfo.Text)
	}
	if len(section.formItems(widget.NewLabel(""))) != 5 {
		t.Fatal("unexpected form rows")
	}

	// Edit: roots mode, defaults off, one entry removed, one typed, one added.
	section.mode.SetSelected(accessModeLabel(devconfig.AccessModeRoots))
	section.denyDefaults.SetChecked(false)
	section.deny.remove(0)
	added := filepath.Join(dir, "added")
	section.denyEntry.SetText(added)
	section.denyAdd.OnTapped()
	typed := filepath.Join(dir, "typed")
	section.denyEntry.SetText(typed)
	remote.SetChecked(true)
	if err := section.apply(&cfg.OutsideAccess); err != nil {
		t.Fatal(err)
	}

	// The relabel after a language switch keeps the selection.
	locales.SetLocale("pl")
	section.relabel()
	if section.modeValue != devconfig.AccessModeRoots || section.mode.Selected != accessModeLabel(devconfig.AccessModeRoots) {
		t.Fatalf("relabel lost the mode: %q", section.mode.Selected)
	}
	locales.SetLocale("en")

	// An existing file with hand-written values the GUI does not edit.
	existing := `{"allowedHosts":["a.example"],"logging":{"level":"debug","toolCalls":false},"futureKey":{"x":1},` +
		`"outsideAccess":{"audit":false,"roots":[{"path":"/old"}],"mode":"full"}}`
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("\xef\xbb\xbf"+existing), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := saveGUIConfig(cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"allowedHosts", "logging", "futureKey"} {
		if _, ok := doc[key]; !ok {
			t.Fatalf("hand-written key %s lost: %s", key, data)
		}
	}
	text := string(data)
	for _, want := range []string{`"mode": "roots"`, `"remote": true`, `"denyDefaults": false`, `"audit": false`} {
		if !strings.Contains(text, want) {
			t.Fatalf("saved config lacks %s:\n%s", want, text)
		}
	}
	if strings.Contains(text, `"/old"`) {
		t.Fatalf("GUI roots should replace the file's roots:\n%s", text)
	}

	t.Setenv("DEVSPACE_CONFIG_DIR", dir)
	loaded := devconfig.LoadConfig()
	access := loaded.OutsideAccess
	if access.AccessMode() != devconfig.AccessModeRoots || access.DenyDefaultsEnabled() || !access.RemoteEnabled() || access.AuditEnabled() {
		t.Fatalf("loaded back: %+v", access)
	}
	if len(access.Deny) != 2 || access.Deny[0] != added || access.Deny[1] != typed {
		t.Fatalf("deny entries: %v", access.Deny)
	}
	if loaded.Logging.Level != "debug" || loaded.Logging.ToolCalls {
		t.Fatalf("hand-written logging changed: %+v", loaded.Logging)
	}

	// Saving again with an empty list writes "deny": [] (cleared), not a
	// missing key.
	cfg.OutsideAccess.Deny = nil
	if err := saveGUIConfig(cfg); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if !strings.Contains(string(data), `"deny": []`) {
		t.Fatalf("cleared deny list not explicit:\n%s", data)
	}
}

func TestMergeGUIConfigRejectsInvalidExistingFile(t *testing.T) {
	if _, err := mergeGUIConfig([]byte("{not json"), devconfig.DefaultConfig()); err == nil {
		t.Fatal("invalid existing config silently overwritten")
	}
	data, err := mergeGUIConfig(nil, devconfig.DefaultConfig())
	if err != nil || !strings.Contains(string(data), `"mode": "full"`) {
		t.Fatalf("fresh config: %s %v", data, err)
	}
}
