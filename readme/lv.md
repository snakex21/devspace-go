# DevSpace (Go izdevums)

**Dodiet ChatGPT un Claude drošu piekļuvi savam lokālajam datoram. Pārvērtiet jebkuru MCP resursdatoru par savu programmēšanas partneri.**

DevSpace ir pašmitināts MCP serveris, kas ļauj AI asistentiem lasīt, rediģēt, meklēt un palaist kodu jūsu reālajos lokālajos projektos — jūsu faili, jūsu rīki, jūsu terminālis — neko neaugšupielādējot trešajām pusēm. Jūs to palaižat savā datorā, atklājat caur jūsu kontrolētu tuneli un pēc izvēles aizsargājat ar paroli.

---

## Saturs

- [Ātrā sākšana](#ātrā-sākšana)
- [Instalēšana](#instalēšana)
- [Ko AI var darīt](#ko-ai-var-darīt)
- [Konfigurācija](#konfigurācija)
- [Tunelis (attālā piekļuve)](#tunelis-attālā-piekļuve)
- [Čaulas atbalsts](#čaulas-atbalsts)
- [Drošība](#drošība)
- [Kompilēšana no pirmkoda](#kompilēšana-no-pirmkoda)
- [Platformu atbalsts](#platformu-atbalsts)
- [Projekta struktūra](#projekta-struktūra)

### 🌍 Tulkojumi

| Valoda | | Valoda | | Valoda | |
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

## Ātrā sākšana

### 1. Lejupielādēt
Izvēlieties savu platformu no [Releases](../../releases) vai kompilējiet no pirmkoda:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurēt (GUI vai teksts)
```bash
devspace-gui                  # Darbvirsmas konfigurators (GUI)
devspace init                 # Teksta konfigurators
```

### 3. Palaist
```bash
devspace                      # Palaiž serveri. Automātiski nosaka konfigurāciju.
```

Tas arī automātiski palaiž Cloudflare tuneli, ja `cloudflared` ir atrasts mapē `tools/`.

### 4. Pievienojiet savu MCP klientu
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Vai lokāli: `http://127.0.0.1:7676/mcp`

---

## Instalēšana

Nav Node.js, nav npm, nav Python. Viens binārs fails.

| Platforma | Lejupielāde |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilēt natīvi) |
| **macOS Intel** | `devspace` (GUI: kompilēt natīvi) |
| **macOS M-chip** | `devspace` (GUI: kompilēt natīvi) |

Nepieciešams **Go 1.23+** tikai tad, ja kompilējat no pirmkoda.

---

## Ko AI var darīt

Pēc savienojuma izveides AI var atvērt vienu no jūsu apstiprinātajām projektu mapēm kā darbvietu:

- **Lasīt, rakstīt un rediģēt** failus darbvietā
- **Meklēt kodu** ar regulārajām izteiksmēm un pārbaudīt direktorijas
- **Palaist čaulas komandas** (PowerShell uz Windows, bash uz Unix)
- **Atklāt projekta instrukcijas** no `AGENTS.md` / `CLAUDE.md`
- **Automātiski konfigurēt** ar pārnēsājamu `.devspace/config.json`

8 MCP rīki: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurācija

Visa konfigurācija atrodas **tajā pašā mapē kā izpildāmais fails** (pārnēsājama):

```
.devspace/
├── config.json       ← atļautās saknes, ports, čaula, valoda, autentifikācija
└── auth.json         ← īpašnieka parole (pēc izvēles)
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

| Lauks | Noklusējums | Apraksts |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automātiski noteikts no OS. Atbalsta 47 valodas |
| `toolMode` | `full` | `full` (visi rīki) vai `minimal` (tikai čaula meklēšanai) |
| `toolNaming` | `short` | `short` (read, write) vai `legacy` (read_file, write_file) |

Nav nepieciešami vides mainīgie — viss ir pārnēsājamā konfigurācijas failā.

---

## Tunelis (attālā piekļuve)

ChatGPT tīmekļa versijai (nepieciešams HTTPS), DevSpace automātiski palaiž tuneli:

| Tunelis | URL veids | Iestatīšana |
|---|---|---|
| **Cloudflare** | Nejaušs (auto) | Ievietojiet `cloudflared.exe` mapē `tools/` |
| **Pinggy** | Stabils | Nepieciešama SSH atslēga (`ssh-keygen`) |

Serveris automātiski nosaka, kurš ir pieejams. Pārstartējiet serveri, lai iegūtu jaunu Cloudflare URL, vai izmantojiet Pinggy pastāvīgam URL.

---

## Čaulas atbalsts

| OS | Noklusējums | Alternatīvas |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / jebkura čaula |
| **macOS** | bash | `sh` / `zsh` |

Iestatiet `"shell"` config.json vai izvēlieties GUI.

---

## Drošība

- **OAuth 2.0 ar PKCE** — ja ir iestatīta īpašnieka parole
- **Režīms bez paroles** — ja parole nav konfigurēta, darbojas bez autentifikācijas
- **Ceļa ierobežošana** — visas failu operācijas tiek validētas pret atļautajām saknēm
- **Izvēles tunelis** — Cloudflare tunelis aizsargā pret tiešu pakļaušanu
- **Nav augšupielādes trešajām pusēm** — jūsu kods nekad neatstāj jūsu datoru

---

## Kompilēšana no pirmkoda

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Kompilēt visu (visas platformas)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Kompilēt tikai pašreizējai platformai
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformu atbalsts

| Platforma | Serveris | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilēt natīvi) |
| **macOS Intel** | ✅ | 🔧 (kompilēt natīvi) |
| **macOS M-chip** | ✅ | 🔧 (kompilēt natīvi) |

GUI nepieciešams Fyne (OpenGL) — nevar krustkompilēt. Serveris kompilējas visur.

---

## Projekta struktūra

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP serveris
│   └── devspace-gui/       ← Darbvirsmas GUI konfigurators (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE nodrošinātājs
│   ├── config/             ← Pārnēsājama konfigurācijas sistēma
│   ├── locales/            ← Tulkojumi 47 valodās
│   ├── logger/             ← Strukturēta žurnālēšana (zerolog)
│   ├── server/             ← HTTP + MCP + tuneļa orķestrācija
│   ├── skills/             ← AGENTS.md / prasmju atklāšana
│   ├── store/              ← SQLite darbvietas sesijas
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Darbvieta un ceļa validācija
├── scripts/
│   ├── windows/            ← PowerShell kompilēšanas skripts
│   └── unix/               ← Bash + Makefile kompilēšanas skripti
├── readme/                 ← Šī faila tulkojumi (47 valodas)
├── build/                  ← Kompilēti binārie faili (visas platformas)
├── tools/                  ← cloudflared.exe utt.
├── go.mod / go.sum
└── README.md
```

---

Veidots Go valodā. Nulle npm. Nulle Node.js. Viens binārs fails.
