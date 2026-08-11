package main

import (
	"fmt"

	"github.com/snakex21/devspace-go/internal/config"
	"github.com/snakex21/devspace-go/internal/locales"
)

func main() {
	// Quick check if already configured — but still allow reconfiguration
	cfg := config.LoadConfig()
	locales.Init(cfg.Lang)
	if len(cfg.AllowedRoots) > 0 {
		fmt.Printf("✅ %s\n", locales.T("gui.configured"))
		fmt.Printf("   %s %v\n", locales.T("gui.roots"), cfg.AllowedRoots)
		fmt.Println()
	} else {
		fmt.Printf("⚙️  %s\n", locales.T("gui.first_config"))
		fmt.Println()
	}

	runGUI(cfg)
}
