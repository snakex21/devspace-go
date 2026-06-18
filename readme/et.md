# DevSpace (Go Edition) - Eesti

**Andke ChatGPT-le ja Claude'ile turvaline juurdepääs oma kohalikule masinale. Muutke iga MCP host oma kodeerimispartneriks.**

DevSpace on isemajutatud MCP-server, mis võimaldab AI-assistentidel lugeda, redigeerida, otsida ja käivitada koodi teie päris kohalikes projektides — teie failid, teie tööriistad, teie terminal — ilma midagi kolmandale osapoolele üles laadimata. Käivitate selle oma masinas, eksponeerite selle enda kontrollitava tunneli kaudu ja valikuliselt kaitsete parooliga.

---

## Sisukord

- [Kiirkäivitus](#kiirkäivitus)
- [Paigaldamine](#paigaldamine)
- [Mida AI Suudab Teha](#mida-ai-suudab-teha)
- [Konfiguratsioon](#konfiguratsioon)
- [Tunnel (Kaugjuurdepääs)](#tunnel-kaugjuurdepääs)
- [Kesta Tugi](#kesta-tugi)
- [Turvalisus](#turvalisus)
- [Lähtekoodist Ehitamine](#lähtekoodist-ehitamine)
- [Platvormi Tugi](#platvormi-tugi)
- [Projekti Struktuur](#projekti-struktuur)

### 🌍 Tõlked

| Keel | | Keel | | Keel | |
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

## Kiirkäivitus

### 1. Lae Alla
Valige oma platvorm [Väljalasetest](../../releases) või ehitage lähtekoodist:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Seadistage (GUI või tekst)
```bash
devspace-gui                  # Töölaua seadistaja (GUI)
devspace init                 # Tekstipõhine seadistaja
```

### 3. Käivitage
```bash
devspace                      # Käivitab serveri. Tuvastab seadistuse automaatselt.
```

See käivitab automaatselt ka Cloudflare'i tunneli, kui `cloudflared` leitakse kaustast `tools/`.

### 4. Ühendage oma MCP klient
```
https://TEIE-TUNNEL.trycloudflare.com/mcp
```
Või lokaalselt: `http://127.0.0.1:7676/mcp`

---

## Paigaldamine

Ei Node.js, ei npm, ei Python. Üks binaarfail.

| Platvorm | Allalaadimine |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompileerige natiivselt) |
| **macOS Intel** | `devspace` (GUI: kompileerige natiivselt) |
| **macOS M-kiip** | `devspace` (GUI: kompileerige natiivselt) |

Nõuab **Go 1.23+** ainult lähtekoodist ehitamisel.

---

## Mida AI Suudab Teha

Kui ühendus on loodud, saab AI avada ühe teie kinnitatud projektikaustadest tööruumina:

- Tööruumis faile **lugeda, kirjutada ja redigeerida**
- **Koodi otsida** regexiga ja katalooge inspekteerida
- **Kesta käske käivitada** (PowerShell Windowsis, bash Unixis)
- **Projekti juhiseid avastada** failidest `AGENTS.md` / `CLAUDE.md`
- **Automaatselt seadistada** kaasaskantava `.devspace/config.json` abil

8 MCP tööriista: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguratsioon

Kogu seadistus asub **käivitatava failiga samas kaustas** (kaasaskantav):

```
.devspace/
├── config.json       ← lubatud juured, port, kest, keel, autentimine
└── auth.json         ← omaniku parool (valikuline)
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

| Väli | Vaikimisi | Kirjeldus |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automaatne tuvastus OS-ist. Toetab 47 keelt |
| `toolMode` | `full` | `full` (kõik tööriistad) või `minimal` (ainult kest otsinguks) |
| `toolNaming` | `short` | `short` (read, write) või `legacy` (read_file, write_file) |

Keskkonnamuutujaid pole vaja — kõik on kaasaskantavas seadistusfailis.

---

## Tunnel (Kaugjuurdepääs)

ChatGPT veebiversiooni jaoks (nõuab HTTPS) käivitab DevSpace automaatselt tunneli:

| Tunnel | URL-i tüüp | Seadistamine |
|---|---|---|
| **Cloudflare** | Juhuslik (auto) | Pange `cloudflared.exe` kausta `tools/` |
| **Pinggy** | Stabiilne | Vajab SSH võtit (`ssh-keygen`) |

Server tuvastab automaatselt, kumb on saadaval. Taaskäivitage server uue Cloudflare'i URL-i saamiseks või kasutage Pinggyt püsiva URL-i jaoks.

---

## Kesta Tugi

| OS | Vaikimisi | Alternatiivid |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / mis tahes kest |
| **macOS** | bash | `sh` / `zsh` |

Määrake `"shell"` failis config.json või valige GUI-s.

---

## Turvalisus

- **OAuth 2.0 koos PKCE-ga** — kui omaniku parool on määratud
- **Paroolita režiim** — kui parooli pole seadistatud, töötab ilma autentimiseta
- **Tee piiramine** — kõik failitoimingud valideeritakse lubatud juurte suhtes
- **Valikuline tunnel** — Cloudflare'i tunnel kaitseb otsese eksponeerimise eest
- **Pole kolmanda osapoole üleslaadimisi** — teie kood ei lahku kunagi teie masinast

---

## Lähtekoodist Ehitamine

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Ehita kõik (kõik platvormid)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Ehita ainult praeguse platvormi jaoks
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platvormi Tugi

| Platvorm | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompileerige natiivselt) |
| **macOS Intel** | ✅ | 🔧 (kompileerige natiivselt) |
| **macOS M-kiip** | ✅ | 🔧 (kompileerige natiivselt) |

GUI nõuab Fyne'i (OpenGL) — ei saa ristkompileerida. Server kompileerub igal pool.

---

## Projekti Struktuur

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP server
│   └── devspace-gui/       ← Töölaua GUI seadistaja (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE pakkuja
│   ├── config/             ← Kaasaskantav seadistussüsteem
│   ├── locales/            ← 47 keeletõlget
│   ├── logger/             ← Struktureeritud logimine (zerolog)
│   ├── server/             ← HTTP + MCP + tunneli orkestreerimine
│   ├── skills/             ← AGENTS.md / oskuste avastamine
│   ├── store/              ← SQLite tööruumi seansid
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Tööruum ja tee valideerimine
├── scripts/
│   ├── windows/            ← PowerShell ehitusskript
│   └── unix/               ← Bash + Makefile ehitusskriptid
├── readme/                 ← Selle faili tõlked (47 keeles)
├── build/                  ← Kompileeritud binaarfailid (kõik platvormid)
├── tools/                  ← cloudflared.exe, jne.
├── go.mod / go.sum
└── README.md
```

---

Ehitatud Go-s. Null npm. Null Node.js. Üks binaarfail.
