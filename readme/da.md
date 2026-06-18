# DevSpace (Go Edition) - Dansk

**Giv ChatGPT & Claude sikker adgang til din lokale maskine. Gør enhver MCP-vært til din kodningspartner.**

DevSpace er en selv-hostet MCP-server, der lader AI-assistenter læse, redigere, søge og køre kode i dine rigtige lokale projekter — dine filer, dine værktøjer, din terminal — uden at uploade noget til en tredjepart. Du kører det på din maskine, eksponerer det gennem en tunnel, du kontrollerer, og sikrer det valgfrit med en adgangskode.

---

## Indholdsfortegnelse

- [Hurtig Start](#hurtig-start)
- [Installation](#installation)
- [Hvad AI Kan Gøre](#hvad-ai-kan-gøre)
- [Konfiguration](#konfiguration)
- [Tunnel (Fjernadgang)](#tunnel-fjernadgang)
- [Shell-understøttelse](#shell-understøttelse)
- [Sikkerhed](#sikkerhed)
- [Bygning fra Kildekode](#bygning-fra-kildekode)
- [Platformunderstøttelse](#platformunderstøttelse)
- [Projektstruktur](#projektstruktur)

### 🌍 Oversættelser

| Sprog | | Sprog | | Sprog | |
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

## Hurtig Start

### 1. Download
Vælg din platform fra [Udgivelser](../../releases) eller byg fra kildekode:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurer (GUI eller tekst)
```bash
devspace-gui                  # Skrivebordskonfigurator (GUI)
devspace init                 # Tekstbaseret konfigurator
```

### 3. Kør
```bash
devspace                      # Starter server. Auto-registrerer konfiguration.
```

Dette starter også automatisk en Cloudflare-tunnel, hvis `cloudflared` findes i `tools/`.

### 4. Tilslut din MCP-klient
```
https://DIN-TUNNEL.trycloudflare.com/mcp
```
Eller lokalt: `http://127.0.0.1:7676/mcp`

---

## Installation

Ingen Node.js, ingen npm, ingen Python. Enkelt binær fil.

| Platform | Download |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilér native) |
| **macOS Intel** | `devspace` (GUI: kompilér native) |
| **macOS M-chip** | `devspace` (GUI: kompilér native) |

Kræver **Go 1.23+** kun ved bygning fra kildekode.

---

## Hvad AI Kan Gøre

Når den er tilsluttet, kan AI'en åbne en af dine godkendte projektmapper som et arbejdsområde:

- **Læse, skrive og redigere** filer inden for arbejdsområdet
- **Søge i kode** med regex og inspicere mapper
- **Køre shell-kommandoer** (PowerShell på Windows, bash på Unix)
- **Opdage projektinstruktioner** fra `AGENTS.md` / `CLAUDE.md`
- **Auto-konfigurere** med bærbar `.devspace/config.json`

8 MCP-værktøjer: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguration

Al konfiguration ligger **i samme mappe som den eksekverbare fil** (bærbar):

```
.devspace/
├── config.json       ← tilladte rødder, port, shell, sprog, godkendelse
└── auth.json         ← ejerens adgangskode (valgfri)
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

| Felt | Standard | Beskrivelse |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Auto-registrer fra OS. Understøtter 47 sprog |
| `toolMode` | `full` | `full` (alle værktøjer) eller `minimal` (kun shell til søgning) |
| `toolNaming` | `short` | `short` (read, write) eller `legacy` (read_file, write_file) |

Ingen miljøvariabler nødvendige — alt er i den bærbare konfigurationsfil.

---

## Tunnel (Fjernadgang)

Til ChatGPT-webversion (kræver HTTPS) starter DevSpace automatisk en tunnel:

| Tunnel | URL-type | Opsætning |
|---|---|---|
| **Cloudflare** | Tilfældig (auto) | Placer `cloudflared.exe` i `tools/` |
| **Pinggy** | Stabil | Kræver SSH-nøgle (`ssh-keygen`) |

Serveren auto-registrerer, hvilken der er tilgængelig. Genstart serveren for en ny Cloudflare-URL, eller brug Pinggy for en permanent URL.

---

## Shell-understøttelse

| OS | Standard | Alternativer |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / enhver shell |
| **macOS** | bash | `sh` / `zsh` |

Indstil `"shell"` i config.json eller vælg i GUI'en.

---

## Sikkerhed

- **OAuth 2.0 med PKCE** — hvis ejerens adgangskode er indstillet
- **Adgangskodefri tilstand** — hvis ingen adgangskode er konfigureret, kører uden godkendelse
- **Stiindeslutning** — alle filoperationer valideres mod tilladte rødder
- **Valgfri tunnel** — Cloudflare-tunnel beskytter mod direkte eksponering
- **Ingen tredjepartsuploads** — din kode forlader aldrig din maskine

---

## Bygning fra Kildekode

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Byg alt (alle platforme)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Byg kun til nuværende platform
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformunderstøttelse

| Platform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilér native) |
| **macOS Intel** | ✅ | 🔧 (kompilér native) |
| **macOS M-chip** | ✅ | 🔧 (kompilér native) |

GUI kræver Fyne (OpenGL) — kan ikke krydskompileres. Server kompileres overalt.

---

## Projektstruktur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-server
│   └── devspace-gui/       ← Skrivebords-GUI-konfigurator (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-udbyder
│   ├── config/             ← Bærbart konfigurationssystem
│   ├── locales/            ← 47 sprogoversættelser
│   ├── logger/             ← Struktureret logning (zerolog)
│   ├── server/             ← HTTP + MCP + tunnelorkestrering
│   ├── skills/             ← AGENTS.md / færdighedsopdagelse
│   ├── store/              ← SQLite arbejdsområdesessioner
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Arbejdsområde- & stivalidering
├── scripts/
│   ├── windows/            ← PowerShell build-script
│   └── unix/               ← Bash + Makefile build-scripts
├── readme/                 ← Oversættelser af denne fil (47 sprog)
├── build/                  ← Kompilerede binære filer (alle platforme)
├── tools/                  ← cloudflared.exe, osv.
├── go.mod / go.sum
└── README.md
```

---

Bygget i Go. Nul npm. Nul Node.js. Én binær fil.
