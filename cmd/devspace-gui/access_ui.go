package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	devconfig "github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/locales"
)

// accessModes in the order the GUI shows them.
var accessModes = []string{devconfig.AccessModeFull, devconfig.AccessModeRoots, devconfig.AccessModeOff}

func accessModeLabel(mode string) string {
	return locales.T("gui.mode_" + mode)
}

func accessModeFromLabel(label string) string {
	for _, mode := range accessModes {
		if accessModeLabel(mode) == label {
			return mode
		}
	}
	return devconfig.AccessModeFull
}

// denyList is the editable list of custom deny entries
// (outsideAccess.deny). Entries must be absolute; duplicates are ignored.
type denyList struct {
	items []string
}

func (d *denyList) add(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if !filepath.IsAbs(path) {
		return errors.New(locales.T("gui.deny_invalid"))
	}
	path = filepath.Clean(path)
	for _, item := range d.items {
		if item == path || (runtime.GOOS == "windows" && strings.EqualFold(item, path)) {
			return nil
		}
	}
	d.items = append(d.items, path)
	return nil
}

func (d *denyList) remove(index int) {
	if index >= 0 && index < len(d.items) {
		d.items = append(d.items[:index], d.items[index+1:]...)
	}
}

// accessSection holds the widgets for outsideAccess.
type accessSection struct {
	window       fyne.Window
	mode         *widget.RadioGroup
	modeValue    string
	remote       *widget.Check
	denyDefaults *widget.Check
	defaultsInfo *widget.Label
	deny         *denyList
	denyRows     *fyne.Container
	denyEntry    *widget.Entry
	denyAdd      *widget.Button
	denyBrowse   *widget.Button
}

func newAccessSection(cfg *devconfig.Config, w fyne.Window, remote *widget.Check) *accessSection {
	a := &accessSection{
		window:    w,
		modeValue: cfg.OutsideAccess.AccessMode(),
		remote:    remote,
		deny:      &denyList{},
		denyRows:  container.NewVBox(),
		denyEntry: widget.NewEntry(),
	}
	for _, entry := range cfg.OutsideAccess.Deny {
		_ = a.deny.add(entry)
	}
	a.mode = widget.NewRadioGroup(nil, func(selected string) {
		if selected != "" {
			a.modeValue = accessModeFromLabel(selected)
		}
	})
	a.denyDefaults = widget.NewCheck("", nil)
	a.denyDefaults.SetChecked(cfg.OutsideAccess.DenyDefaultsEnabled())
	a.defaultsInfo = widget.NewLabel("")
	a.defaultsInfo.Wrapping = fyne.TextWrapWord
	a.denyAdd = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if err := a.deny.add(a.denyEntry.Text); err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.denyEntry.SetText("")
		a.refreshRows()
	})
	a.denyBrowse = widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			_ = a.deny.add(dialogPath(uri))
			a.refreshRows()
		}, a.window)
		folderDialog.Resize(folderDialogSize(a.window.Canvas().Size()))
		folderDialog.Show()
	})
	a.relabel()
	return a
}

// dialogPath turns a Fyne folder URI into a native path.
func dialogPath(uri fyne.ListableURI) string {
	path := strings.TrimPrefix(uri.Path(), "file://")
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	return filepath.FromSlash(path)
}

// relabel applies the current locale, keeping the selected mode.
func (a *accessSection) relabel() {
	options := make([]string, 0, len(accessModes))
	for _, mode := range accessModes {
		options = append(options, accessModeLabel(mode))
	}
	selected := a.modeValue
	a.mode.Options = options
	a.mode.SetSelected(accessModeLabel(selected))
	a.mode.Refresh()
	a.denyDefaults.SetText(locales.T("gui.deny_defaults"))
	a.defaultsInfo.SetText(locales.T("gui.deny_defaults_list"))
	a.denyEntry.SetPlaceHolder(locales.T("gui.deny_placeholder"))
	a.denyAdd.SetText(locales.T("gui.deny_add"))
	a.denyBrowse.SetText(locales.T("gui.browse"))
	a.refreshRows()
}

func (a *accessSection) refreshRows() {
	a.denyRows.RemoveAll()
	if len(a.deny.items) == 0 {
		a.denyRows.Add(widget.NewLabel(locales.T("gui.deny_empty")))
	}
	for i, item := range a.deny.items {
		index := i
		remove := widget.NewButtonWithIcon(locales.T("gui.deny_remove"), theme.DeleteIcon(), func() {
			a.deny.remove(index)
			a.refreshRows()
		})
		a.denyRows.Add(container.NewBorder(nil, nil, nil, remove, widget.NewLabel(item)))
	}
	a.denyRows.Refresh()
}

// formItems returns the rows the configuration form shows for this section.
func (a *accessSection) formItems(extraRoots fyne.CanvasObject) []*widget.FormItem {
	return []*widget.FormItem{
		widget.NewFormItem(locales.T("gui.access_mode"), a.mode),
		widget.NewFormItem("", a.remote),
		widget.NewFormItem(locales.T("gui.extra_roots"), extraRoots),
		widget.NewFormItem(locales.T("gui.deny_title"), container.NewVBox(a.denyDefaults, a.defaultsInfo)),
		widget.NewFormItem(locales.T("gui.deny_custom"), container.NewVBox(
			a.denyRows,
			container.NewBorder(nil, nil, nil, container.NewHBox(a.denyBrowse, a.denyAdd), a.denyEntry),
		)),
	}
}

// apply copies the section into cfg. A path still typed in the entry is
// added as well, so it is not silently lost on save.
func (a *accessSection) apply(cfg *devconfig.OutsideAccessConfig) error {
	if err := a.deny.add(a.denyEntry.Text); err != nil {
		return err
	}
	a.denyEntry.SetText("")
	cfg.Mode = a.modeValue
	remote := a.remote.Checked
	cfg.Remote = &remote
	denyDefaults := a.denyDefaults.Checked
	cfg.DenyDefaults = &denyDefaults
	cfg.Deny = append([]string{}, a.deny.items...)
	return nil
}

// mergeGUIConfig writes the values the GUI edits into an existing
// config.json, keeping every other key (set by hand or by a newer
// version) untouched. Inside outsideAccess the GUI owns mode, remote,
// roots, deny and denyDefaults; other keys such as audit are kept.
func mergeGUIConfig(existing []byte, cfg *devconfig.Config) ([]byte, error) {
	document := map[string]json.RawMessage{}
	if len(existing) > 0 {
		existing = stripUTF8BOM(existing)
		if err := json.Unmarshal(existing, &document); err != nil {
			return nil, fmt.Errorf("existing config.json is not valid JSON (%v); fix or remove it before saving", err)
		}
	}
	set := func(target map[string]json.RawMessage, key string, value any) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		target[key] = data
		return nil
	}
	values := map[string]any{
		"host":          cfg.Host,
		"port":          cfg.Port,
		"allowedRoots":  cfg.AllowedRoots,
		"publicBaseUrl": cfg.PublicBaseURL,
		"stateDir":      devconfig.PortablePath(cfg.StateDir),
		"worktreeRoot":  devconfig.PortablePath(cfg.WorktreeRoot),
		"agentDir":      devconfig.PortablePath(cfg.AgentDir),
		"toolMode":      cfg.ToolMode,
		"toolNaming":    cfg.ToolNaming,
		"shell":         cfg.Shell,
		"lang":          cfg.Lang,
		"widgets":       cfg.Widgets,
		"skillsEnabled": cfg.SkillsEnabled,
	}
	for key, value := range values {
		if err := set(document, key, value); err != nil {
			return nil, err
		}
	}

	outside := map[string]json.RawMessage{}
	if raw, ok := document["outsideAccess"]; ok {
		_ = json.Unmarshal(raw, &outside)
	}
	access := cfg.OutsideAccess
	deny := access.Deny
	if deny == nil {
		deny = []string{}
	}
	for key, value := range map[string]any{
		"mode":         access.AccessMode(),
		"remote":       access.RemoteEnabled(),
		"denyDefaults": access.DenyDefaultsEnabled(),
		"deny":         deny,
	} {
		if err := set(outside, key, value); err != nil {
			return nil, err
		}
	}
	if len(access.Roots) > 0 {
		if err := set(outside, "roots", access.Roots); err != nil {
			return nil, err
		}
	} else {
		delete(outside, "roots")
	}
	if err := set(document, "outsideAccess", outside); err != nil {
		return nil, err
	}
	return json.MarshalIndent(document, "", "  ")
}

func stripUTF8BOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}
