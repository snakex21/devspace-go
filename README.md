# Dev Space Go

**Give ChatGPT & Claude secure access to your local machine. Turn any MCP host into your coding partner.**

Dev Space Go is a self-hosted MCP server that lets AI assistants read, edit, search, and run code in your real local projects — your files, your tools, your terminal — without uploading anything to a third party. You run it on your machine and expose it through a tunnel you control.

---

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [What AI Can Do](#what-ai-can-do)
- [Configuration](#configuration)
- [Tunnel (Remote Access)](#tunnel-remote-access)
- [Shell Support](#shell-support)
- [Security](#security)
- [Robustness and Limits](#robustness-and-limits)
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
devspace-gui                  # Desktop configurator (GUI)
devspace init                 # Text-based configurator
```

### 3. Run
```bash
devspace                      # Starts server. Auto-detects config.
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
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compile natively) |
| **macOS Intel** | `devspace` (GUI: compile natively) |
| **macOS M-chip** | `devspace` (GUI: compile natively) |

Requires **Go 1.26+** (see `go.mod`) only if building from source.

---

## What AI Can Do

Once connected, the AI can open one of your approved project folders as a workspace:

- **Read, write, and edit** files inside the workspace
- **Create directories and move/rename files** safely inside the workspace
- **Search code** with regex and inspect directories
- **Run shell commands** (PowerShell on Windows, bash on Unix)
- **Discover project instructions** from `AGENTS.md` / `CLAUDE.md`
- **Auto-configure** with portable `.devspace/config.json`

10 MCP tools: `open_workspace`, `read`, `write`, `mkdir`, `move`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuration

All config lives **in the same folder as the executable** (portable):

```
.devspace/
└── config.json       ← allowed roots, port, shell, language
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

For ChatGPT web version (HTTPS required), Dev Space Go auto-starts a tunnel:

| Tunnel | URL type | Setup |
|---|---|---|
| **Cloudflare** | Random (auto) | `cloudflared.exe` included in `tools/` |
| **Pinggy** | Stable | Needs SSH key (`ssh-keygen`) |

Server auto-detects which one is available. Restart the server for a new Cloudflare URL, or use Pinggy for a permanent URL.

### Connection readiness and troubleshooting

The local server starts listening before the tunnel is launched. A URL printed by
`cloudflared` is only an allocated address: wait for Dev Space Go's **TUNNEL ACTIVE**
banner, which appears after the public `/healthz` endpoint reaches this server.
Startup checks retry transient DNS and gateway errors within a bounded timeout.
Ctrl+C also cancels an in-progress tunnel startup and reaps the child process.

Use **`/mcp`** with the default Cloudflare Quick Tunnel. Quick Tunnels
[do not support SSE](https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/),
so switching to `/sse` will not fix them. The legacy `/sse` endpoint remains
available locally or through a provider that supports SSE.

If the tunnel process exits later (network loss, sleep/hibernation, a crash of
`cloudflared`), Dev Space Go restarts it automatically with exponential backoff
and jitter (5 s doubling up to 5 min). A restarted Quick Tunnel gets a **new
URL**, announced in a new TUNNEL ACTIVE banner. If neither `cloudflared` nor
`ssh` is installed, the server simply stays local and does not retry.

If the public connection fails, check `http://127.0.0.1:7676/healthz` locally
(adjust the configured host/port) and keep the tunnel log visible. The local
server remains usable if tunnel startup fails. A tunnel process that exits is
reported explicitly; its old URL is no longer active. Restarting a Quick Tunnel
changes the public hostname, so update the MCP client's URL as well.

The readiness check needs outbound HTTPS access to the generated hostname from
the server machine. A local proxy, DNS filter, or firewall blocking that check
can prevent the active banner even if another machine could reach the tunnel.
Quick Tunnels are for development and have no uptime guarantee.

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

- **No built-in authentication** — anyone with the public tunnel URL can use the server
- **Path containment** — every tool path (read, write, edit, mkdir, move, ls, grep, glob and the shell working directory) is resolved inside the workspace root; `..`, absolute paths elsewhere and symlinks/junctions pointing outside are rejected
- **The shell tool is not a sandbox** — `bash` runs arbitrary commands with your user's permissions; path containment only applies to the file tools
- **Tunnel access** — treat the public tunnel URL as a secret and stop the server when it is not in use
- **No third-party uploads** — your code never leaves your machine

---

## Robustness and Limits

**Shutdown.** Ctrl+C (or SIGTERM / closing the console window) cancels running
requests first — open `/sse` streams end and running shell commands are killed
together with all processes they started — then shuts the HTTP server down
(at most 10 s) and stops the tunnel. Press Ctrl+C a second time to force-quit
if shutdown hangs.

**Crash / hard kill.** The tunnel process (`cloudflared`/`ssh`) runs in a
Windows job object (Linux: parent-death signal) and dies with Dev Space Go even
after a crash or a Task Manager kill, so no orphaned tunnel keeps publishing
the port. The workspace database is SQLite with a 5 s busy timeout, so a
restart that overlaps the previous shutdown waits instead of failing. If the
port is still taken at startup, the error says whether another Dev Space Go
instance or a different program holds it.

**Configuration problems** (invalid JSON in `config.json`, a non-numeric
`PORT`, an allowed root that does not exist, an out-of-range port) are printed
as warnings at startup and by `devspace doctor`. Config files are written
atomically, so a crash while saving never leaves a truncated file.

**Limits and timeouts**

| What | Limit |
|---|---|
| HTTP request headers | 10 s to arrive, 1 MiB |
| HTTP request body | 2 min to arrive, 64 MiB |
| Idle keep-alive connection | 2 min |
| Shell command | 30 s default, 300 s max; the whole process tree is killed on timeout |
| Shell output kept | 256 KiB (the rest is counted and discarded) |
| Shell background processes | left running; output collection stops 2 s after the shell exits |
| `read` response | 2000 lines by default, max 5000 lines / 4 MiB |
| `edit` file size | 64 MiB |
| `grep` / `glob` | 500 matches, files up to 2 MiB, 20 s |
| `open_workspace` scan for `AGENTS.md` / `CLAUDE.md` | 5000 files / 3 s (Windows, Linux, macOS) |
| Graceful shutdown | 10 s |

A panic inside a tool is logged and returned to the client as a tool error; it
no longer terminates the server.

**Known limitations**

- macOS has no parent-death signal: after a hard kill (`kill -9`) of Dev Space
  Go, a running tunnel child can survive there. Normal exits and Ctrl+C are fine.
- A process the shell tool starts within the first milliseconds, before it is
  attached to its Windows job object, can escape the tree kill.
- Link checks happen before each file operation; a link swapped in by another
  local process between the check and the operation (a race) is not prevented.
- `open_workspace` stops looking for `AGENTS.md` / `CLAUDE.md` after 5000
  files or 3 s, so on a very large root (for example a whole drive) instruction
  files deep in the tree may not be listed; open the project folder instead.
- Restarting a Quick Tunnel changes its URL; use a provider with a stable URL
  if clients must not be reconfigured.

---

## Building from Source

```bash
git clone https://github.com/snakex21/devspace-go
cd devspace-go

# Build everything (all platforms)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Build just for current platform (release flags: smaller binaries)
go build -trimpath -ldflags="-s -w" -o devspace ./cmd/devspace/
go build -trimpath -ldflags="-s -w" -o devspace-gui ./cmd/devspace-gui/

# Tests (add -race when a C compiler is available)
go vet ./...
go test ./...
```

The build scripts use `-trimpath -ldflags="-s -w"`, which strips symbol and
debug tables and local paths. Panics still print stack traces with function
names; only debuggers lose symbol information.

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
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP server
│   └── devspace-gui/       ← Desktop GUI configurator (Fyne)
├── internal/
│   ├── config/             ← Portable config system
│   ├── locales/            ← 47 language translations
│   ├── logger/             ← Structured logging (zerolog)
│   ├── procutil/           ← Child process trees (job objects / process groups)
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
