package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	devconfig "github.com/waishnav/mcp-webcoder/internal/config"
	"github.com/waishnav/mcp-webcoder/internal/locales"
)

// --- All available languages ---
var langCodes = []string{
	"auto", "en", "pl", "de", "fr", "es", "it", "pt", "ru", "uk",
	"zh", "ja", "ko", "ar", "hi", "tr", "nl", "sv", "no", "da",
	"fi", "cs", "sk", "hu", "ro", "bg", "el", "he", "th", "vi",
	"id", "ms", "af", "bn", "ca", "et", "fa", "ga", "hr", "lt",
	"lv", "mt", "sl", "sr", "sw", "ta", "ur", "zu",
}

var shellOptions = []string{"auto", "powershell", "cmd", "bash", "sh"}

// guiWidgets holds references to all translatable widgets.
type guiWidgets struct {
	window       fyne.Window
	titleLabel   *widget.Label
	hostEntry    *widget.Entry
	portEntry    *widget.Entry
	rootsEntry   *widget.Entry
	urlEntry     *widget.Entry
	shellSelect  *widget.Select
	langSelect   *widget.Select
	passEntry    *widget.Entry
	confirmEntry *widget.Entry
	browseBtn    *widget.Button
	saveBtn      *widget.Button
	testBtn      *widget.Button
	statusLabel  *widget.Label
	progressBar  *widget.ProgressBarInfinite
	cfg          *devconfig.Config
	formLabels   []*widget.Label // labels inside the form we track
}

// applyLang updates all widget texts to the currently active locale.
func (gw *guiWidgets) applyLang() {
	lang := gw.langSelect.Selected
	if lang == "" || lang == "auto" {
		lang = "en"
	}
	locales.SetLocale(lang)

	gw.window.SetTitle(locales.T("gui.title"))
	gw.titleLabel.SetText(locales.T("gui.title"))

	// Update form labels (we re-create the form, so store references)
	gw.browseBtn.SetText(locales.T("gui.browse"))
	gw.saveBtn.SetText(locales.T("gui.save"))
	gw.testBtn.SetText(locales.T("gui.test"))

	// Update placeholders
	gw.rootsEntry.SetPlaceHolder(locales.T("gui.roots_placeholder"))
	gw.urlEntry.SetPlaceHolder(locales.T("gui.url_placeholder"))
	gw.passEntry.SetPlaceHolder(locales.T("gui.password_placeholder"))
	gw.confirmEntry.SetPlaceHolder(locales.T("gui.confirm_placeholder"))

	// Update password validator message
	gw.passEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		if len(s) < 16 {
			return errors.New(locales.T("gui.password_empty_ok"))
		}
		return nil
	}
	gw.portEntry.Validator = func(s string) error {
		var p int
		if _, err := fmt.Sscanf(s, "%d", &p); err != nil || p < 1 || p > 65535 {
			return errors.New(locales.T("gui.port_invalid"))
		}
		return nil
	}

	// Rebuild form to update labels
	gw.rebuildForm()
}

func (gw *guiWidgets) rebuildForm() {
	form := &widget.Form{}
	form.Append(locales.T("gui.host"), gw.hostEntry)
	form.Append(locales.T("gui.port"), gw.portEntry)
	form.Append(locales.T("gui.roots"), container.NewBorder(nil, nil, nil, gw.browseBtn, gw.rootsEntry))
	form.Append(locales.T("gui.url"), gw.urlEntry)
	form.Append(locales.T("gui.shell"), gw.shellSelect)
	form.Append(locales.T("gui.lang"), gw.langSelect)
	form.Append(locales.T("gui.password"), gw.passEntry)
	form.Append(locales.T("gui.confirm"), gw.confirmEntry)

	buttons := container.NewHBox(
		layout.NewSpacer(),
		gw.testBtn,
		gw.saveBtn,
		layout.NewSpacer(),
	)

	content := container.NewVBox(
		widget.NewLabelWithStyle(locales.T("gui.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		form,
		widget.NewSeparator(),
		buttons,
		gw.progressBar,
		gw.statusLabel,
	)

	scroll := container.NewVScroll(content)
	gw.window.SetContent(scroll)
}

// runGUI launches the desktop GUI configurator.
func runGUI() {
	a := app.NewWithID("com.mcp-webcoder.configurator")
	a.Settings().SetTheme(theme.DarkTheme())

	cfg := devconfig.LoadConfig() // load existing config with defaults

	// Init locales with current language
	lang := cfg.Lang
	if lang == "" || lang == "auto" {
		lang = "en"
	}
	locales.Init(lang)

	w := a.NewWindow(locales.T("gui.title"))
	w.Resize(fyne.NewSize(640, 680))
	w.CenterOnScreen()

	// --- Widgets ---

	titleLabel := widget.NewLabelWithStyle(locales.T("gui.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	hostEntry := widget.NewEntry()
	hostEntry.SetText(cfg.Host)

	portEntry := widget.NewEntry()
	portEntry.SetText(fmt.Sprintf("%d", cfg.Port))
	portEntry.Validator = func(s string) error {
		var p int
		if _, err := fmt.Sscanf(s, "%d", &p); err != nil || p < 1 || p > 65535 {
			return errors.New(locales.T("gui.port_invalid"))
		}
		return nil
	}

	rootsEntry := widget.NewEntry()
	rootsEntry.SetText(strings.Join(cfg.AllowedRoots, "; "))
	rootsEntry.SetPlaceHolder(locales.T("gui.roots_placeholder"))

	browseBtn := widget.NewButtonWithIcon(locales.T("gui.browse"), theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			current := strings.TrimSpace(rootsEntry.Text)
			path := uri.Path()
			path = strings.TrimPrefix(path, "file://")
			if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}
			if current == "" {
				rootsEntry.SetText(path)
			} else if !strings.Contains(current, path) {
				rootsEntry.SetText(current + "; " + path)
			}
		}, w)
	})

	urlEntry := widget.NewEntry()
	urlEntry.SetText(cfg.PublicBaseURL)
	urlEntry.SetPlaceHolder(locales.T("gui.url_placeholder"))

	shellSelect := widget.NewSelect(shellOptions, nil)
	if cfg.Shell == "" {
		cfg.Shell = "auto"
	}
	shellSelect.SetSelected(cfg.Shell)

	langSelect := widget.NewSelect(langCodes, nil)
	if cfg.Lang == "" || cfg.Lang == "auto" {
		langSelect.SetSelected("auto")
	} else {
		langSelect.SetSelected(cfg.Lang)
	}

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder(locales.T("gui.password_placeholder"))
	passEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		if len(s) < 16 {
			return errors.New(locales.T("gui.password_empty_ok"))
		}
		return nil
	}

	confirmEntry := widget.NewPasswordEntry()
	confirmEntry.SetPlaceHolder(locales.T("gui.confirm_placeholder"))

	statusLabel := widget.NewLabel("")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	// --- Assemble guiWidgets ---
	gw := &guiWidgets{
		window:       w,
		titleLabel:   titleLabel,
		hostEntry:    hostEntry,
		portEntry:    portEntry,
		rootsEntry:   rootsEntry,
		urlEntry:     urlEntry,
		shellSelect:  shellSelect,
		langSelect:   langSelect,
		passEntry:    passEntry,
		confirmEntry: confirmEntry,
		browseBtn:    browseBtn,
		statusLabel:  statusLabel,
		progressBar:  progressBar,
		cfg:          cfg,
	}

	// --- Save button ---
	saveBtn := widget.NewButtonWithIcon(locales.T("gui.save"), theme.ConfirmIcon(), func() {
		if strings.TrimSpace(rootsEntry.Text) == "" {
			dialog.ShowError(errors.New(locales.T("gui.roots_empty")), w)
			return
		}
		if passEntry.Text != "" && passEntry.Text != confirmEntry.Text {
			dialog.ShowError(errors.New(locales.T("gui.password_mismatch")), w)
			return
		}
		if passEntry.Text != "" {
			if err := passEntry.Validate(); err != nil {
				dialog.ShowError(err, w)
				return
			}
		}
		if err := portEntry.Validate(); err != nil {
			dialog.ShowError(err, w)
			return
		}

		cfg.Host = hostEntry.Text
		fmt.Sscanf(portEntry.Text, "%d", &cfg.Port)
		cfg.PublicBaseURL = urlEntry.Text
		cfg.AllowedRoots = parseRoots(rootsEntry.Text)
		cfg.Shell = shellSelect.Selected
		cfg.Lang = langSelect.Selected
		cfg.OAuth.OwnerToken = passEntry.Text

		if len(cfg.AllowedRoots) == 0 {
			dialog.ShowError(errors.New(locales.T("gui.roots_parse_error")), w)
			return
		}

		progressBar.Show()
		statusLabel.SetText(locales.T("gui.status_saving"))

		if err := saveGUIConfig(cfg); err != nil {
			progressBar.Hide()
			statusLabel.SetText(locales.T("gui.save_error") + ": " + err.Error())
			dialog.ShowError(fmt.Errorf("%s: %v", locales.T("gui.save_error"), err), w)
			return
		}

		progressBar.Hide()
		statusLabel.SetText(locales.T("gui.status_saved"))
		dialog.ShowInformation(
			locales.T("gui.saved"),
			fmt.Sprintf(locales.T("gui.saved_msg"), cfg.ConfigDir),
			w,
		)

		// Update locales if language changed
		if langSelect.Selected != lang {
			locales.SetLocale(langSelect.Selected)
		}
	})
	gw.saveBtn = saveBtn

	// --- Test button ---
	testBtn := widget.NewButtonWithIcon(locales.T("gui.test"), theme.MediaPlayIcon(), func() {
		statusLabel.SetText(locales.T("gui.status_checking"))
		info := fmt.Sprintf(
			"Host: %s\nPort: %s\nRoots: %s\nURL: %s\n%s: %s\n\n%s: %s\n%s: %s",
			hostEntry.Text,
			portEntry.Text,
			rootsEntry.Text,
			urlEntry.Text,
			locales.T("gui.password"),
			strings.Repeat("*", len(passEntry.Text)),
			locales.T("gui.shell"),
			shellSelect.Selected,
			locales.T("gui.lang"),
			langSelect.Selected,
		)
		dialog.ShowInformation(locales.T("gui.preview_title"), info, w)
		statusLabel.SetText(locales.T("gui.status_ready"))
	})
	gw.testBtn = testBtn

	// --- Language change → immediate update ---
	langSelect.OnChanged = func(selected string) {
		gw.applyLang()
	}

	// Build initial form
	gw.rebuildForm()

	w.ShowAndRun()
}

func parseRoots(input string) []string {
	input = strings.ReplaceAll(input, "\r\n", ";")
	input = strings.ReplaceAll(input, "\n", ";")
	parts := strings.Split(input, ";")
	var roots []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			roots = append(roots, trimmed)
		}
	}
	return roots
}

func saveGUIConfig(cfg *devconfig.Config) error {
	if err := os.MkdirAll(cfg.ConfigDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	configPath := filepath.Join(cfg.ConfigDir, "config.json")
	configData := map[string]interface{}{
		"host":          cfg.Host,
		"port":          cfg.Port,
		"allowedRoots":  cfg.AllowedRoots,
		"publicBaseUrl": cfg.PublicBaseURL,
		"stateDir":      cfg.StateDir,
		"worktreeRoot":  cfg.WorktreeRoot,
		"agentDir":      cfg.AgentDir,
		"toolMode":      cfg.ToolMode,
		"toolNaming":    cfg.ToolNaming,
		"shell":         cfg.Shell,
		"lang":          cfg.Lang,
		"widgets":       cfg.Widgets,
		"skillsEnabled": cfg.SkillsEnabled,
		"logLevel":      cfg.Logging.Level,
		"logFormat":     cfg.Logging.Format,
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("write config.json: %w", err)
	}

	// Save or remove auth.json
	authPath := filepath.Join(cfg.ConfigDir, "auth.json")
	if cfg.OAuth.OwnerToken != "" {
		authData := map[string]string{"ownerToken": cfg.OAuth.OwnerToken}
		authBytes, err := json.MarshalIndent(authData, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal auth: %w", err)
		}
		if err := os.WriteFile(authPath, authBytes, 0600); err != nil {
			return fmt.Errorf("write auth.json: %w", err)
		}
	} else {
		os.Remove(authPath)
	}

	return nil
}
