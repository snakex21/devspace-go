# DevSpace (Go Edition) - Afrikaans

**Gee ChatGPT & Claude veilige toegang tot jou plaaslike masjien. Verander enige MCP-gasheer in jou kodeermaat.**

DevSpace is 'n self-aangebode MCP-bediener wat KI-assistente toelaat om kode in jou werklike plaaslike projekte te lees, redigeer, soek en uit te voer — jou lêers, jou gereedskap, jou terminaal — sonder om enigiets na 'n derde party op te laai. Jy hardloop dit op jou masjien, stel dit bloot deur 'n tonnel wat jy beheer, en beveilig dit opsioneel met 'n wagwoord.

---

## Inhoudsopgawe

- [Vinnige Begin](#vinnige-begin)
- [Installasie](#installasie)
- [Wat KI Kan Doen](#wat-ki-kan-doen)
- [Konfigurasie](#konfigurasie)
- [Tonnel (Afstandtoegang)](#tonnel-afstandtoegang)
- [Dopondersteuning](#dopondersteuning)
- [Sekuriteit](#sekuriteit)
- [Bou vanaf Bron](#bou-vanaf-bron)
- [Platformondersteuning](#platformondersteuning)
- [Projekstruktuur](#projekstruktuur)

### 🌍 Vertalings

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

## Vinnige Begin

### 1. Laai Af
Kies jou platform vanaf [Vrystellings](../../releases) of bou vanaf bron:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigureer (GUI of teks)
```bash
devspace-gui                  # Werkskerm-konfigureerder (GUI)
devspace init                 # Teksgebaseerde konfigureerder
```

### 3. Hardloop
```bash
devspace                      # Begin bediener. Bespeur outomaties konfig.
```

Dit begin ook outomaties 'n Cloudflare-tonnel as `cloudflared` in `tools/` gevind word.

### 4. Koppel jou MCP-kliënt
```
https://JOU-TONNEL.trycloudflare.com/mcp
```
Of plaaslik: `http://127.0.0.1:7676/mcp`

---

## Installasie

Geen Node.js, geen npm, geen Python. Enkel binêre.

| Platform | Aflaai |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompileer natuurlik) |
| **macOS Intel** | `devspace` (GUI: kompileer natuurlik) |
| **macOS M-skyfie** | `devspace` (GUI: kompileer natuurlik) |

Vereis **Go 1.23+** slegs indien jy vanaf bron bou.

---

## Wat KI Kan Doen

Sodra gekoppel, kan die KI een van jou goedgekeurde projekgidse as 'n werkspasie oopmaak:

- **Lees, skryf en redigeer** lêers binne die werkspasie
- **Soek kode** met regex en inspekteer gidse
- **Voer dopopdragte uit** (PowerShell op Windows, bash op Unix)
- **Ontdek projekinstruksies** vanaf `AGENTS.md` / `CLAUDE.md`
- **Outo-konfigureer** met draagbare `.devspace/config.json`

8 MCP-gereedskap: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurasie

Alle konfigurasie woon **in dieselfde gids as die uitvoerbare lêer** (draagbaar):

```
.devspace/
├── config.json       ← toegelate wortels, poort, dop, taal, magtiging
└── auth.json         ← eienaar wagwoord (opsioneel)
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

| Veld | Verstek | Beskrywing |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Outo-bespeur vanaf bedryfstelsel. Ondersteun 47 tale |
| `toolMode` | `full` | `full` (alle gereedskap) of `minimal` (slegs dop vir soek) |
| `toolNaming` | `short` | `short` (read, write) of `legacy` (read_file, write_file) |

Geen omgewingsveranderlikes nodig nie — alles is in die draagbare konfigurasielêer.

---

## Tonnel (Afstandtoegang)

Vir ChatGPT-webweergawe (HTTPS vereis), begin DevSpace outomaties 'n tonnel:

| Tonnel | URL-tipe | Opstelling |
|---|---|---|
| **Cloudflare** | Ewekansig (outo) | Plaas `cloudflared.exe` in `tools/` |
| **Pinggy** | Stabiel | Benodig SSH-sleutel (`ssh-keygen`) |

Bediener bespeur outomaties watter een beskikbaar is. Herbegin die bediener vir 'n nuwe Cloudflare-URL, of gebruik Pinggy vir 'n permanente URL.

---

## Dopondersteuning

| Bedryfstelsel | Verstek | Alternatiewe |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / enige dop |
| **macOS** | bash | `sh` / `zsh` |

Stel `"shell"` in config.json of kies in die GUI.

---

## Sekuriteit

- **OAuth 2.0 met PKCE** — indien eienaarwagwoord gestel is
- **Wagwoordlose modus** — indien geen wagwoord gekonfigureer is, hardloop sonder magtiging
- **Padbeperking** — alle lêerbewerkings gevalideer teen toegelate wortels
- **Opsionele tonnel** — Cloudflare-tonnel beskerm teen direkte blootstelling
- **Geen derdeparty-oplaai** — jou kode verlaat nooit jou masjien nie

---

## Bou vanaf Bron

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bou alles (alle platforms)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bou slegs vir huidige platform
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformondersteuning

| Platform | Bediener | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompileer natuurlik) |
| **macOS Intel** | ✅ | 🔧 (kompileer natuurlik) |
| **macOS M-skyfie** | ✅ | 🔧 (kompileer natuurlik) |

GUI vereis Fyne (OpenGL) — kan nie kruiskompileer nie. Bediener kompileer oral.

---

## Projekstruktuur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-bediener
│   └── devspace-gui/       ← Werkskerm-GUI-konfigureerder (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE-verskaffer
│   ├── config/             ← Draagbare konfigurasiestelsel
│   ├── locales/            ← 47 taalvertalings
│   ├── logger/             ← Gestruktureerde logboek (zerolog)
│   ├── server/             ← HTTP + MCP + tonnelorkestrasie
│   ├── skills/             ← AGENTS.md / vaardigheidsontdekking
│   ├── store/              ← SQLite-werkspasiesessies
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Werkspasie- & padvalidering
├── scripts/
│   ├── windows/            ← PowerShell-bouskrip
│   └── unix/               ← Bash + Makefile-bouskripte
├── readme/                 ← Vertalings van hierdie lêer (47 tale)
├── build/                  ← Gekompileerde binêre (alle platforms)
├── tools/                  ← cloudflared.exe, ens.
├── go.mod / go.sum
└── README.md
```

---

Gebou in Go. Nul npm. Nul Node.js. Een binêre.
