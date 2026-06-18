# DevSpace (Go Edition) - Deutsch

**Geben Sie ChatGPT & Claude sicheren Zugriff auf Ihren lokalen Rechner. Machen Sie jeden MCP-Host zu Ihrem Coding-Partner.**

DevSpace ist ein selbst gehosteter MCP-Server, der KI-Assistenten erlaubt, Code in Ihren echten lokalen Projekten zu lesen, zu bearbeiten, zu durchsuchen und auszuführen — Ihre Dateien, Ihre Werkzeuge, Ihr Terminal — ohne etwas an Dritte hochzuladen. Sie führen ihn auf Ihrem Rechner aus, stellen ihn über einen von Ihnen kontrollierten Tunnel bereit und sichern ihn optional mit einem Passwort.

---

## Inhaltsverzeichnis

- [Schnellstart](#schnellstart)
- [Installation](#installation)
- [Was KI Kann](#was-ki-kann)
- [Konfiguration](#konfiguration)
- [Tunnel (Fernzugriff)](#tunnel-fer nzugriff)
- [Shell-Unterstützung](#shell-unterstützung)
- [Sicherheit](#sicherheit)
- [Aus dem Quellcode Bauen](#aus-dem-quellcode-bauen)
- [Plattform-Unterstützung](#plattform-unterstützung)
- [Projektstruktur](#projektstruktur)

### 🌍 Übersetzungen

| Sprache | | Sprache | | Sprache | |
|---|---|---|---|---|---|
| [Afrikaans](readme/af.md) | [العربية](readme/ar.md) | [Български](readme/bg.md) | [বাংলা](readme/bn.md) | [Català](readme/ca.md) |
| [Čeština](readme/cs.md) | [Dansk](readme/da.md) | [Deutsch](readme/de.md) | [Ελληνικά](readme/el.md) | [English](readme/en.md) |
| [Español](readme/es.md) | [Eesti](readme/et.md) | [فارسی](readme/fa.md) | [Suomi](readme/fi.md) | [Français](readme/fr.md) |
| [Gaeilge](readme/ga.md) | [עברית](readme/he.md) | [हिन्दी](readme/hi.md) | [Hrvatski](readme/hr.md) | [Magyar](readme/hu.md) |
| [Bahasa](readme/id.md) | [Italiano](readme/it.md) | [日本語](readme/ja.md) | [한국어](readme/ko.md) | [Lietuvių](readme/lt.md) |
| [Latviešu](readme/lv.md) | [Melayu](readme/ms.md) | [Malti](readme/mt.md) | [Nederlands](readme/nl.md) | [Norsk](readme/no.md) |
| [Polski](readme/pl.md) | [Português](readme/pt.md) | [Română](readme/ro.md) | [Русский](readme/ru.md) | [Slovenčina](readme/sk.md) |
| [Slovenščina](readme/sl.md) | [Српски](readme/sr.md) | [Svenska](readme/sv.md) | [Kiswahili](readme/sw.md) | [தமிழ்](readme/ta.md) |
| [ไทย](readme/th.md) | [Türkçe](readme/tr.md) | [Українська](readme/uk.md) | [اردو](readme/ur.md) | [Tiếng Việt](readme/vi.md) |
| [简体中文](readme/zh.md) | [isiZulu](readme/zu.md) |

---

## Schnellstart

### 1. Herunterladen
Wählen Sie Ihre Plattform von [Releases](../../releases) oder bauen Sie aus dem Quellcode:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurieren (GUI oder Text)
```bash
devspace-gui                  # Desktop-Konfigurator (GUI)
devspace init                 # Textbasierter Konfigurator
```

### 3. Ausführen
```bash
devspace                      # Startet Server. Erkennt Konfiguration automatisch.
```

Dies startet auch automatisch einen Cloudflare-Tunnel, wenn `cloudflared` in `tools/` gefunden wird.

### 4. Verbinden Sie Ihren MCP-Client
```
https://IHR-TUNNEL.trycloudflare.com/mcp
```
Oder lokal: `http://127.0.0.1:7676/mcp`

---

## Installation

Kein Node.js, kein npm, kein Python. Einzelne Binärdatei.

| Plattform | Download |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: nativ kompilieren) |
| **macOS Intel** | `devspace` (GUI: nativ kompilieren) |
| **macOS M-Chip** | `devspace` (GUI: nativ kompilieren) |

Erfordert **Go 1.23+** nur beim Bauen aus dem Quellcode.

---

## Was KI Kann

Einmal verbunden, kann die KI einen Ihrer genehmigten Projektordner als Arbeitsbereich öffnen:

- Dateien im Arbeitsbereich **lesen, schreiben und bearbeiten**
- **Code durchsuchen** mit Regex und Verzeichnisse inspizieren
- **Shell-Befehle ausführen** (PowerShell unter Windows, bash unter Unix)
- **Projektanweisungen entdecken** aus `AGENTS.md` / `CLAUDE.md`
- **Automatisch konfigurieren** mit portabler `.devspace/config.json`

8 MCP-Werkzeuge: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguration

Die gesamte Konfiguration befindet sich **im selben Ordner wie die ausführbare Datei** (portabel):

```
.devspace/
├── config.json       ← erlaubte Wurzeln, Port, Shell, Sprache, Authentifizierung
└── auth.json         ← Besitzer-Passwort (optional)
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

| Feld | Standard | Beschreibung |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatische Erkennung vom Betriebssystem. Unterstützt 47 Sprachen |
| `toolMode` | `full` | `full` (alle Werkzeuge) oder `minimal` (nur Shell für Suche) |
| `toolNaming` | `short` | `short` (read, write) oder `legacy` (read_file, write_file) |

Keine Umgebungsvariablen erforderlich — alles befindet sich in der portablen Konfigurationsdatei.

---

## Tunnel (Fernzugriff)

Für die ChatGPT-Webversion (HTTPS erforderlich) startet DevSpace automatisch einen Tunnel:

| Tunnel | URL-Typ | Einrichtung |
|---|---|---|
| **Cloudflare** | Zufällig (auto) | Legen Sie `cloudflared.exe` in `tools/` ab |
| **Pinggy** | Stabil | Benötigt SSH-Schlüssel (`ssh-keygen`) |

Der Server erkennt automatisch, welcher verfügbar ist. Starten Sie den Server für eine neue Cloudflare-URL neu oder verwenden Sie Pinggy für eine permanente URL.

---

## Shell-Unterstützung

| Betriebssystem | Standard | Alternativen |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / beliebige Shell |
| **macOS** | bash | `sh` / `zsh` |

Setzen Sie `"shell"` in config.json oder wählen Sie in der GUI.

---

## Sicherheit

- **OAuth 2.0 mit PKCE** — wenn ein Besitzer-Passwort gesetzt ist
- **Passwortloser Modus** — wenn kein Passwort konfiguriert ist, läuft ohne Authentifizierung
- **Pfadbegrenzung** — alle Dateioperationen werden gegen erlaubte Wurzeln validiert
- **Optionaler Tunnel** — Cloudflare-Tunnel schützt vor direkter Offenlegung
- **Keine Drittanbieter-Uploads** — Ihr Code verlässt niemals Ihren Rechner

---

## Aus dem Quellcode Bauen

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Alles bauen (alle Plattformen)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Nur für aktuelle Plattform bauen
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Plattform-Unterstützung

| Plattform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (nativ kompilieren) |
| **macOS Intel** | ✅ | 🔧 (nativ kompilieren) |
| **macOS M-Chip** | ✅ | 🔧 (nativ kompilieren) |

GUI erfordert Fyne (OpenGL) — kann nicht cross-kompiliert werden. Server kompiliert überall.

---

## Projektstruktur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-Server
│   └── devspace-gui/       ← Desktop-GUI-Konfigurator (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-Anbieter
│   ├── config/             ← Portables Konfigurationssystem
│   ├── locales/            ← 47 Sprachübersetzungen
│   ├── logger/             ← Strukturiertes Logging (zerolog)
│   ├── server/             ← HTTP + MCP + Tunnel-Orchestrierung
│   ├── skills/             ← AGENTS.md / Fähigkeitserkennung
│   ├── store/              ← SQLite-Arbeitsbereichssitzungen
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Arbeitsbereich- & Pfadvalidierung
├── scripts/
│   ├── windows/            ← PowerShell-Bauskript
│   └── unix/               ← Bash + Makefile-Bauskripte
├── readme/                 ← Übersetzungen dieser Datei (47 Sprachen)
├── build/                  ← Kompilierte Binärdateien (alle Plattformen)
├── tools/                  ← cloudflared.exe, usw.
├── go.mod / go.sum
└── README.md
```

---

In Go gebaut. Null npm. Null Node.js. Eine Binärdatei.
