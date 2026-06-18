# DevSpace (Go-editie)

**Geef ChatGPT en Claude veilige toegang tot je lokale machine. Maak van elke MCP-host je codeerpartner.**

DevSpace is een zelfgehoste MCP-server waarmee AI-assistenten bestanden kunnen lezen, bewerken, doorzoeken en code kunnen uitvoeren in jouw echte lokale projecten — jouw bestanden, jouw tools, jouw terminal — zonder iets naar een derde partij te uploaden. Je draait het op je eigen machine, stelt het bloot via een tunnel die jij beheert en beveiligt het optioneel met een wachtwoord.

---

## Inhoudsopgave

- [Snelle Start](#snelle-start)
- [Installatie](#installatie)
- [Wat AI Kan Doen](#wat-ai-kan-doen)
- [Configuratie](#configuratie)
- [Tunnel (Externe Toegang)](#tunnel-externe-toegang)
- [Shell-ondersteuning](#shell-ondersteuning)
- [Beveiliging](#beveiliging)
- [Bouwen vanaf Bron](#bouwen-vanaf-bron)
- [Platformondersteuning](#platformondersteuning)
- [Projectstructuur](#projectstructuur)

### 🌍 Vertalingen

| Taal | | Taal | | Taal | |
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

## Snelle Start

### 1. Downloaden
Kies je platform op de [Releases](../../releases)-pagina of bouw vanaf bron:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configureren (GUI of tekst)
```bash
devspace-gui                  # Desktop configurator (GUI)
devspace init                 # Tekstgebaseerde configurator
```

### 3. Uitvoeren
```bash
devspace                      # Start de server. Detecteert configuratie automatisch.
```

Dit start ook automatisch een Cloudflare-tunnel als `cloudflared` wordt gevonden in `tools/`.

### 4. Verbind je MCP-client
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Of lokaal: `http://127.0.0.1:7676/mcp`

---

## Installatie

Geen Node.js, geen npm, geen Python. Eén enkele binary.

| Platform | Download |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: native compileren) |
| **macOS Intel** | `devspace` (GUI: native compileren) |
| **macOS M-chip** | `devspace` (GUI: native compileren) |

Vereist **Go 1.23+** alleen als je vanaf bron bouwt.

---

## Wat AI Kan Doen

Eenmaal verbonden kan de AI een van je goedgekeurde projectmappen openen als werkruimte:

- Bestanden **lezen, schrijven en bewerken** binnen de werkruimte
- **Code doorzoeken** met regex en mappen inspecteren
- **Shell-commando's uitvoeren** (PowerShell op Windows, bash op Unix)
- **Projectinstructies ontdekken** uit `AGENTS.md` / `CLAUDE.md`
- **Automatisch configureren** met draagbare `.devspace/config.json`

8 MCP-tools: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuratie

Alle configuratie bevindt zich **in dezelfde map als het uitvoerbare bestand** (draagbaar):

```
.devspace/
├── config.json       ← toegestane roots, poort, shell, taal, authenticatie
└── auth.json         ← eigenaarswachtwoord (optioneel)
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

| Veld | Standaard | Beschrijving |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatisch gedetecteerd van besturingssysteem. Ondersteunt 47 talen |
| `toolMode` | `full` | `full` (alle tools) of `minimal` (alleen shell voor zoeken) |
| `toolNaming` | `short` | `short` (read, write) of `legacy` (read_file, write_file) |

Geen omgevingsvariabelen nodig — alles staat in het draagbare configuratiebestand.

---

## Tunnel (Externe Toegang)

Voor de ChatGPT-webversie (HTTPS vereist) start DevSpace automatisch een tunnel:

| Tunnel | URL-type | Setup |
|---|---|---|
| **Cloudflare** | Willekeurig (auto) | Plaats `cloudflared.exe` in `tools/` |
| **Pinggy** | Stabiel | Vereist SSH-sleutel (`ssh-keygen`) |

De server detecteert automatisch welke beschikbaar is. Herstart de server voor een nieuwe Cloudflare-URL of gebruik Pinggy voor een permanente URL.

---

## Shell-ondersteuning

| OS | Standaard | Alternatieven |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / elke shell |
| **macOS** | bash | `sh` / `zsh` |

Stel `"shell"` in config.json in of kies in de GUI.

---

## Beveiliging

- **OAuth 2.0 met PKCE** — als het eigenaarswachtwoord is ingesteld
- **Wachtwoordloze modus** — als er geen wachtwoord is geconfigureerd, draait het zonder authenticatie
- **Padbeperking** — alle bestandsoperaties worden gevalideerd tegen toegestane roots
- **Optionele tunnel** — Cloudflare-tunnel beschermt tegen directe blootstelling
- **Geen uploads naar derden** — je code verlaat nooit je machine

---

## Bouwen vanaf Bron

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bouw alles (alle platforms)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bouw alleen voor het huidige platform
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformondersteuning

| Platform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (native compileren) |
| **macOS Intel** | ✅ | 🔧 (native compileren) |
| **macOS M-chip** | ✅ | 🔧 (native compileren) |

De GUI vereist Fyne (OpenGL) — kan niet cross-compilen. De server compileert overal.

---

## Projectstructuur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-server
│   └── devspace-gui/       ← Desktop GUI-configurator (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-provider
│   ├── config/             ← Draagbaar configuratiesysteem
│   ├── locales/            ← Vertalingen in 47 talen
│   ├── logger/             ← Gestructureerde logging (zerolog)
│   ├── server/             ← HTTP + MCP + tunnelorkestratie
│   ├── skills/             ← AGENTS.md / vaardigheidsdetectie
│   ├── store/              ← SQLite werkruimtesessies
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Werkruimte en padvalidatie
├── scripts/
│   ├── windows/            ← PowerShell build-script
│   └── unix/               ← Bash + Makefile build-scripts
├── readme/                 ← Vertalingen van dit bestand (47 talen)
├── build/                  ← Gecompileerde binaries (alle platforms)
├── tools/                  ← cloudflared.exe, enz.
├── go.mod / go.sum
└── README.md
```

---

Gebouwd in Go. Nul npm. Nul Node.js. Eén binary.
