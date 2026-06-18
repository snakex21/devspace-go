.PHONY: all windows linux macos-intel macos-mchip clean

ROOT = ../..
BLD = $(ROOT)/build
CMD = $(ROOT)/cmd

all: clean windows linux macos-intel macos-mchip
	@echo "=== Done! Binaries in $(BLD)/ ==="

# --- Platform targets ---
windows: $(BLD)/windows/devspace.exe $(BLD)/windows/devspace-gui.exe
linux: $(BLD)/linux/devspace $(BLD)/linux/devspace-gui
macos-intel: $(BLD)/macos-intel/devspace $(BLD)/macos-intel/devspace-gui
macos-mchip: $(BLD)/macos-mchip/devspace $(BLD)/macos-mchip/devspace-gui

# --- Windows ---
$(BLD)/windows/devspace.exe:
	@mkdir -p $(BLD)/windows
	GOOS=windows GOARCH=amd64 go build -o $@ $(CMD)/devspace/
	@echo "  -> $@"

$(BLD)/windows/devspace-gui.exe:
	@mkdir -p $(BLD)/windows
	GOOS=windows GOARCH=amd64 go build -o $@ $(CMD)/devspace-gui/
	@echo "  -> $@"

# --- Linux ---
$(BLD)/linux/devspace:
	@mkdir -p $(BLD)/linux
	GOOS=linux GOARCH=amd64 go build -o $@ $(CMD)/devspace/
	@echo "  -> $@"

$(BLD)/linux/devspace-gui:
	@mkdir -p $(BLD)/linux
	GOOS=linux GOARCH=amd64 go build -o $@ $(CMD)/devspace-gui/
	@echo "  -> $@"

# --- macOS Intel ---
$(BLD)/macos-intel/devspace:
	@mkdir -p $(BLD)/macos-intel
	GOOS=darwin GOARCH=amd64 go build -o $@ $(CMD)/devspace/
	@echo "  -> $@"

$(BLD)/macos-intel/devspace-gui:
	@mkdir -p $(BLD)/macos-intel
	GOOS=darwin GOARCH=amd64 go build -o $@ $(CMD)/devspace-gui/
	@echo "  -> $@"

# --- macOS M-chip ---
$(BLD)/macos-mchip/devspace:
	@mkdir -p $(BLD)/macos-mchip
	GOOS=darwin GOARCH=arm64 go build -o $@ $(CMD)/devspace/
	@echo "  -> $@"

$(BLD)/macos-mchip/devspace-gui:
	@mkdir -p $(BLD)/macos-mchip
	GOOS=darwin GOARCH=arm64 go build -o $@ $(CMD)/devspace-gui/
	@echo "  -> $@"

clean:
	rm -rf $(BLD)
