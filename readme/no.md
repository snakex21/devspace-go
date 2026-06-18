# DevSpace (Go-utgave)

**Gi ChatGPT og Claude sikker tilgang til din lokale maskin. Gjør enhver MCP-vert til din kodepartner.**

DevSpace er en selvhostet MCP-server som lar AI-assistenter lese, redigere, søke og kjøre kode i dine virkelige lokale prosjekter — dine filer, dine verktøy, din terminal — uten å laste opp noe til en tredjepart. Du kjører den på maskinen din, eksponerer den gjennom en tunnel du kontrollerer, og sikrer den eventuelt med et passord.

---

## Innholdsfortegnelse

- [Hurtigstart](#hurtigstart)
- [Installasjon](#installasjon)
- [Hva AI Kan Gjøre](#hva-ai-kan-gjøre)
- [Konfigurasjon](#konfigurasjon)
- [Tunnel (Fjerntilgang)](#tunnel-fjerntilgang)
- [Skallstøtte](#skallstøtte)
- [Sikkerhet](#sikkerhet)
- [Bygging fra Kildekode](#bygging-fra-kildekode)
- [Plattformstøtte](#plattformstøtte)
- [Prosjektstruktur](#prosjektstruktur)

### 🌍 Oversettelser

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

## Hurtigstart

### 1. Last ned
Velg din plattform fra [Releases](../../releases) eller bygg fra kildekode:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurer (GUI eller tekst)
```bash
devspace-gui                  # Skrivebordskonfigurator (GUI)
devspace init                 # Tekstbasert konfigurator
```

### 3. Kjør
```bash
devspace                      # Starter serveren. Oppdager konfigurasjon automatisk.
```

Dette starter også automatisk en Cloudflare-tunnel hvis `cloudflared` finnes i `tools/`.

### 4. Koble til din MCP-klient
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Eller lokalt: `http://127.0.0.1:7676/mcp`

---

## Installasjon

Ingen Node.js, ingen npm, ingen Python. Én enkelt binærfil.

| Plattform | Nedlasting |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompileres naturlig) |
| **macOS Intel** | `devspace` (GUI: kompileres naturlig) |
| **macOS M-chip** | `devspace` (GUI: kompileres naturlig) |

Krever **Go 1.23+** bare hvis du bygger fra kildekode.

---

## Hva AI Kan Gjøre

Når den er tilkoblet, kan AI-en åpne en av dine godkjente prosjektmapper som et arbeidsområde:

- **Lese, skrive og redigere** filer innenfor arbeidsområdet
- **Søke i kode** med regex og inspisere mapper
- **Kjøre skallkommandoer** (PowerShell på Windows, bash på Unix)
- **Oppdage prosjektinstruksjoner** fra `AGENTS.md` / `CLAUDE.md`
- **Automatisk konfigurere** med bærbar `.devspace/config.json`

8 MCP-verktøy: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurasjon

All konfigurasjon ligger **i samme mappe som den kjørbare filen** (bærbar):

```
.devspace/
├── config.json       ← tillatte rotmappr, port, skall, språk, autentisering
└── auth.json         ← eierpassord (valgfritt)
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
| `lang` | `auto` | Oppdages automatisk fra OS. Støtter 47 språk |
| `toolMode` | `full` | `full` (alle verktøy) eller `minimal` (bare skall for søk) |
| `toolNaming` | `short` | `short` (read, write) eller `legacy` (read_file, write_file) |

Ingen miljøvariabler nødvendig — alt er i den bærbare konfigurasjonsfilen.

---

## Tunnel (Fjerntilgang)

For ChatGPT-nettversjonen (HTTPS påkrevd) starter DevSpace automatisk en tunnel:

| Tunnel | URL-type | Oppsett |
|---|---|---|
| **Cloudflare** | Tilfeldig (auto) | Plasser `cloudflared.exe` i `tools/` |
| **Pinggy** | Stabil | Krever SSH-nøkkel (`ssh-keygen`) |

Serveren oppdager automatisk hvilken som er tilgjengelig. Start serveren på nytt for en ny Cloudflare-URL, eller bruk Pinggy for en permanent URL.

---

## Skallstøtte

| OS | Standard | Alternativer |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / hvilket som helst skall |
| **macOS** | bash | `sh` / `zsh` |

Sett `"shell"` i config.json eller velg i GUI.

---

## Sikkerhet

- **OAuth 2.0 med PKCE** — hvis eierpassord er satt
- **Passordløs modus** — hvis ingen passord er konfigurert, kjører uten autentisering
- **Stiinneslutning** — alle filoperasjoner valideres mot tillatte rotmapper
- **Valgfri tunnel** — Cloudflare-tunnel beskytter mot direkte eksponering
- **Ingen tredjepartsopplastinger** — koden din forlater aldri maskinen din

---

## Bygging fra Kildekode

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bygg alt (alle plattformer)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bygg bare for gjeldende plattform
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Plattformstøtte

| Plattform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompileres naturlig) |
| **macOS Intel** | ✅ | 🔧 (kompileres naturlig) |
| **macOS M-chip** | ✅ | 🔧 (kompileres naturlig) |

GUI krever Fyne (OpenGL) — kan ikke krysskompileres. Serveren kompileres overalt.

---

## Prosjektstruktur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-server
│   └── devspace-gui/       ← Skrivebords-GUI-konfigurator (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-leverandør
│   ├── config/             ← Bærbart konfigurasjonssystem
│   ├── locales/            ← Oversettelser til 47 språk
│   ├── logger/             ← Strukturert logging (zerolog)
│   ├── server/             ← HTTP + MCP + tunnelorkestrering
│   ├── skills/             ← AGENTS.md / ferdighetsoppdagelse
│   ├── store/              ← SQLite arbeidsområdeøkter
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Arbeidsområde og stivalidering
├── scripts/
│   ├── windows/            ← PowerShell-byggeskript
│   └── unix/               ← Bash + Makefile-byggeskript
├── readme/                 ← Oversettelser av denne filen (47 språk)
├── build/                  ← Kompilerte binærfiler (alle plattformer)
├── tools/                  ← cloudflared.exe, osv.
├── go.mod / go.sum
└── README.md
```

---

Bygget i Go. Null npm. Null Node.js. Én binærfil.
