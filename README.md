# MCP WebCoder

**Give ChatGPT & Claude secure access to your local machine. Turn any MCP host into your coding partner.**

MCP WebCoder is a self-hosted MCP server that lets AI assistants read, edit, search, and run code in your real local projects — your files, your tools, your terminal — without uploading anything to a third party. You run it on your machine, expose it through a tunnel you control, and optionally secure it with a password.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [What AI Can Do](#what-ai-can-do)
- [Configuration](#configuration)
- [Tunnel (Remote Access)](#tunnel-remote-access)
- [Shell Support](#shell-support)
- [Security](#security)
- [Building from Source](#building-from-source)
- [Platform Support](#platform-support)
- [Project Structure](#project-structure)

### 🌍 Translations

| Language | | Language | | Language | |
|---|---|---|---|---|---|
| [Afrikaans](readme/af.md) | [العربية](readme/ar.md) | [Български](readme/bg.md) | [বাংলা](readme/bn.md) | [Català](readme/ca.md) |
| [Čeština](readme/cs.md) | [Dansk](readme/da.md) | [Deutsch](readme/de.md) | [Ελληνικά](readme/el.md) | [English](readme/en.md) |
| [Español](readme/es.md) | [Eesti](readme/et.md) | [فارسی](readme/fa.md) | [Suomi](readme/fi.md) | [Français](readme/fr.md) |
| [Gaeilge](readme/ga.md) | [עברית](readme/he.md) | [हिन्दी](readme/hi.md) | [Hrvatski](readme/hr.md) | [Magyar](readme/hu.md) |
| [Indonesia](readme/id.md) | [Italiano](readme/it.md) | [日本語](readme/ja.md) | [한국어](readme/ko.md) | [Lietuvių](readme/lt.md) |
| [Latviešu](readme/lv.md) | [Melayu](readme/ms.md) | [Malti](readme/mt.md) | [Nederlands](readme/nl.md) | [Norsk](readme/no.md) |
| [Polski](readme/pl.md) | [Português](readme/pt.md) | [Română](readme/ro.md) | [Русский](readme/ru.md) | [Slovenčina](readme/sk.md) |
| [Slovenščina](readme/sl.md) | [Српски](readme/sr.md) | [Svenska](readme/sv.md) | [Kiswahili](readme/sw.md) | [தமிழ்](readme/ta.md) |
| [ไทย](readme/th.md) | [Türkçe](readme/tr.md) | [Українська](readme/uk.md) | [اردو](readme/ur.md) | [Tiếng Việt](readme/vi.md) |
| [简体中文](readme/zh.md) | [isiZulu](readme/zu.md) |

---

## Quick Start

### 1. Download
Pick your platform from [Releases](../../releases) or build from source:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configure (GUI or text)
```bash
mcp-webcoder-gui              # Desktop configurator (GUI)
mcp-webcoder init             # Text-based configurator
```

### 3. Run
```bash
mcp-webcoder                  # Starts server. Auto-detects config.
```

This also auto-starts a Cloudflare Tunnel if `cloudflared` is found in `tools/`.

### 4. Connect your MCP client
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Or locally: `http://127.0.0.1:7676/mcp`

---

## Installation

No Node.js, no npm, no Python. Single binary.

| Platform | Download |
|---|---|
| **Windows** | `mcp-webcoder.exe` + `mcp-webcoder-gui.exe` |
| **Linux** | `mcp-webcoder` (GUI: compile natively) |
| **macOS Intel** | `mcp-webcoder` (GUI: compile natively) |
| **macOS M-chip** | `mcp-webcoder` (GUI: compile natively) |

Requires **Go 1.23+** only if building from source.

---

## What AI Can Do

Once connected, the AI can open one of your approved project folders as a workspace:

- **Read, write, and edit** files inside the workspace
- **Create directories and move/rename files** safely inside the workspace
- **Search code** with regex and inspect directories
- **Run shell commands** (PowerShell on Windows, bash on Unix)
- **Discover project instructions** from `AGENTS.md` / `CLAUDE.md`
- **Auto-configure** with portable `.webcoder/config.json`

10 MCP tools: `open_workspace`, `read`, `write`, `mkdir`, `move`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuration

All config lives **in the same folder as the executable** (portable):

```
.webcoder/
├── config.json       ← allowed roots, port, shell, language, auth
└── auth.json         ← owner password (optional)
```

### config.json
```json
{
  "host": "127.0.0.1",
  "port": 7676,
  "allowedRoots": ["C:/projects"],
  "publicBaseUrl": "http://127.0.0.1:7676",
  "shell": "auto",
  "lang": "auto",
  "toolMode": "full",
  "toolNaming": "short"
}
```

| Field | Default | Description |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Auto-detect from OS. Supports 47 languages |
| `toolMode` | `full` | `full` (all tools) or `minimal` (shell only for search) |
| `toolNaming` | `short` | `short` (read, write) or `legacy` (read_file, write_file) |

No environment variables needed — everything is in the portable config file.

---

## Tunnel (Remote Access)

For ChatGPT web version (HTTPS required), MCP WebCoder auto-starts a tunnel:

| Tunnel | URL type | Setup |
|---|---|---|
| **Cloudflare** | Random (auto) | `cloudflared.exe` included in `tools/` |
| **Pinggy** | Stable | Needs SSH key (`ssh-keygen`) |

Server auto-detects which one is available. Restart the server for a new Cloudflare URL, or use Pinggy for a permanent URL.

---

## Shell Support

| OS | Default | Alternatives |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / any shell |
| **macOS** | bash | `sh` / `zsh` |

Set `"shell"` in config.json or choose in the GUI.

---

## Security

- **OAuth 2.0 with PKCE** — if owner password is set
- **Password-less mode** — if no password configured, runs without auth
- **Path containment** — all file ops validated against allowed roots
- **Optional tunnel** — Cloudflare Tunnel protects from direct exposure
- **No third-party uploads** — your code never leaves your machine

---

## Building from Source

```bash
git clone https://github.com/snakex21/mcp-webcoder
cd mcp-webcoder

# Build everything (all platforms)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Build just for current platform
go build -o mcp-webcoder ./cmd/devspace/
go build -o mcp-webcoder-gui ./cmd/devspace-gui/
```

---

## Platform Support

| Platform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compile natively) |
| **macOS Intel** | ✅ | 🔧 (compile natively) |
| **macOS M-chip** | ✅ | 🔧 (compile natively) |

GUI requires Fyne (OpenGL) — cannot cross-compile. Server compiles everywhere.

---

## Project Structure

```
mcp-webcoder/
├── cmd/
│   ├── devspace/           ← CLI + MCP server
│   └── devspace-gui/       ← Desktop GUI configurator (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE provider
│   ├── config/             ← Portable config system
│   ├── locales/            ← 47 language translations
│   ├── logger/             ← Structured logging (zerolog)
│   ├── server/             ← HTTP + MCP + tunnel orchestration
│   ├── store/              ← SQLite workspace sessions
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Workspace & path validation
├── scripts/
│   ├── windows/            ← PowerShell build script
│   ├── unix/               ← Bash + Makefile build scripts
│   └── userscripts/        ← Tampermonkey auto-approve script
├── readme/                 ← Translations of this file (47 languages)
├── tools/                  ← cloudflared.exe
├── go.mod / go.sum
└── README.md
```

---

Built in Go. Zero npm. Zero Node.js. One binary.
