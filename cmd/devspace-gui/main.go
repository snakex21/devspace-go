package main

import (
	"fmt"
	"os"

	"github.com/snakex21/devspace-go/internal/config"
)

func main() {
	// Quick check if already configured — but still allow reconfiguration
	cfg := config.LoadConfig()
	if len(cfg.AllowedRoots) > 0 {
		fmt.Println("✅ MCP WebCoder jest już skonfigurowany.")
		fmt.Printf("   Roots: %v\n", cfg.AllowedRoots)
		fmt.Println("   Uruchamiam GUI do ewentualnej zmiany...")
		fmt.Println()
	} else {
		fmt.Println("⚙️  MCP WebCoder — pierwsza konfiguracja")
		fmt.Println()
	}

	runGUI()
	os.Exit(0)
}
