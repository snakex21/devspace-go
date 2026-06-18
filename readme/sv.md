# DevSpace (Go-utgåvan)

**Ge ChatGPT och Claude säker åtkomst till din lokala dator. Förvandla vilken MCP-värd som helst till din kodningspartner.**

DevSpace är en självhostad MCP-server som låter AI-assistenter läsa, redigera, söka och köra kod i dina riktiga lokala projekt — dina filer, dina verktyg, din terminal — utan att ladda upp något till tredje part. Du kör den på din dator, exponerar den genom en tunnel du kontrollerar och skyddar den eventuellt med ett lösenord.

---

## Innehållsförteckning

- [Snabbstart](#snabbstart)
- [Installation](#installation)
- [Vad AI Kan Göra](#vad-ai-kan-göra)
- [Konfiguration](#konfiguration)
- [Tunnel (Fjärråtkomst)](#tunnel-fjärråtkomst)
- [Skalstöd](#skalstöd)
- [Säkerhet](#säkerhet)
- [Bygga från Källkod](#bygga-från-källkod)
- [Plattformsstöd](#plattformsstöd)
- [Projektstruktur](#projektstruktur)

### 🌍 Översättningar

| Språk | | Språk | | Språk | |
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

## Snabbstart

### 1. Ladda ner
Välj din plattform från [Releases](../../releases) eller bygg från källkod:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurera (GUI eller text)
```bash
devspace-gui                  # Skrivbordskonfigurator (GUI)
devspace init                 # Textbaserad konfigurator
```

### 3. Kör
```bash
devspace                      # Startar servern. Autoupptäcker konfig.
```

Detta startar också automatiskt en Cloudflare-tunnel om `cloudflared` finns i `tools/`.

### 4. Anslut din MCP-klient
```
https://DIN-TUNNEL.trycloudflare.com/mcp
```
Eller lokalt: `http://127.0.0.1:7676/mcp`

---

## Installation

Ingen Node.js, ingen npm, ingen Python. En enda binär.

| Plattform | Nedladdning |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilera naturligt) |
| **macOS Intel** | `devspace` (GUI: kompilera naturligt) |
| **macOS M-chip** | `devspace` (GUI: kompilera naturligt) |

Kräver **Go 1.23+** endast vid bygge från källkod.

---

## Vad AI Kan Göra

När den är ansluten kan AI:n öppna en av dina godkända projektmappar som en arbetsyta:

- **Läsa, skriva och redigera** filer inom arbetsytan
- **Söka kod** med regex och inspektera kataloger
- **Köra skalkommandon** (PowerShell på Windows, bash på Unix)
- **Upptäcka projektinstruktioner** från `AGENTS.md` / `CLAUDE.md`
- **Autokonfigurera** med portabel `.devspace/config.json`

8 MCP-verktyg: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguration

All konfiguration finns **i samma mapp som den körbara filen** (portabel):

```
.devspace/
├── config.json       ← tillåtna rötter, port, skal, språk, autentisering
└── auth.json         ← ägarlösenord (valfritt)
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

| Fält | Standard | Beskrivning |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Autoupptäcker från OS. Stöder 47 språk |
| `toolMode` | `full` | `full` (alla verktyg) eller `minimal` (endast skal för sökning) |
| `toolNaming` | `short` | `short` (read, write) eller `legacy` (read_file, write_file) |

Inga miljövariabler behövs — allt finns i den portabla konfigurationsfilen.

---

## Tunnel (Fjärråtkomst)

För webbversionen av ChatGPT (HTTPS krävs) startar DevSpace automatiskt en tunnel:

| Tunnel | URL-typ | Installation |
|---|---|---|
| **Cloudflare** | Slumpmässig (auto) | Placera `cloudflared.exe` i `tools/` |
| **Pinggy** | Stabil | Kräver SSH-nyckel (`ssh-keygen`) |

Servern autoupptäcker vilken som är tillgänglig. Starta om servern för en ny Cloudflare-URL eller använd Pinggy för en permanent URL.

---

## Skalstöd

| OS | Standard | Alternativ |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / valfritt skal |
| **macOS** | bash | `sh` / `zsh` |

Ställ in `"shell"` i config.json eller välj i GUI:t.

---

## Säkerhet

- **OAuth 2.0 med PKCE** — om ägarlösenord är inställt
- **Lösenordslöst läge** — om inget lösenord är konfigurerat körs det utan autentisering
- **Sökvägsbegränsning** — alla filoperationer valideras mot tillåtna rötter
- **Valfri tunnel** — Cloudflare-tunnel skyddar mot direkt exponering
- **Inga uppladdningar till tredje part** — din kod lämnar aldrig din dator

---

## Bygga från Källkod

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bygg allt (alla plattformar)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bygg endast för aktuell plattform
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Plattformsstöd

| Plattform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilera naturligt) |
| **macOS Intel** | ✅ | 🔧 (kompilera naturligt) |
| **macOS M-chip** | ✅ | 🔧 (kompilera naturligt) |

GUI kräver Fyne (OpenGL) — kan inte korskompileras. Servern kompileras överallt.

---

## Projektstruktur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-server
│   └── devspace-gui/       ← Skrivbordskonfigurator GUI (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-leverantör
│   ├── config/             ← Portabelt konfigurationssystem
│   ├── locales/            ← Översättningar till 47 språk
│   ├── logger/             ← Strukturerad loggning (zerolog)
│   ├── server/             ← HTTP + MCP + tunnelorkestrering
│   ├── skills/             ← AGENTS.md / färdighetsupptäckt
│   ├── store/              ← SQLite-arbetsytesessioner
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Arbetsyte- och sökvägsvalidering
├── scripts/
│   ├── windows/            ← PowerShell-byggskript
│   └── unix/               ← Bash + Makefile-byggskript
├── readme/                 ← Översättningar av denna fil (47 språk)
├── build/                  ← Kompilerade binärer (alla plattformar)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Byggd i Go. Noll npm. Noll Node.js. En binär.
