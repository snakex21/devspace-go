# DevSpace (Go versija)

**Suteikite ChatGPT ir Claude saugią prieigą prie savo vietinio kompiuterio. Paverskite bet kurį MCP serverį savo programavimo partneriu.**

DevSpace yra savarankiškai talpinamas MCP serveris, leidžiantis AI asistentams skaityti, redaguoti, ieškoti ir vykdyti kodą jūsų tikruose vietiniuose projektuose — jūsų failai, jūsų įrankiai, jūsų terminalas — nieko neįkeliant į trečiųjų šalių serverius. Jūs paleidžiate jį savo kompiuteryje, atveriate per jūsų valdomą tunelį ir pasirinktinai apsaugote slaptažodžiu.

---

## Turinys

- [Greita pradžia](#greita-pradžia)
- [Diegimas](#diegimas)
- [Ką gali AI](#ką-gali-ai)
- [Konfigūracija](#konfigūracija)
- [Tunelis (nuotolinė prieiga)](#tunelis-nuotolinė-prieiga)
- [Apvalkalo palaikymas](#apvalkalo-palaikymas)
- [Saugumas](#saugumas)
- [Kompiliavimas iš šaltinio](#kompiliavimas-iš-šaltinio)
- [Platformų palaikymas](#platformų-palaikymas)
- [Projekto struktūra](#projekto-struktūra)

### 🌍 Vertimai

| Kalba | | Kalba | | Kalba | |
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

## Greita pradžia

### 1. Atsisiųskite
Pasirinkite savo platformą iš [Releases](../../releases) arba kompiliuokite iš šaltinio:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigūruokite (GUI arba tekstu)
```bash
devspace-gui                  # Darbalaukio konfigūratorius (GUI)
devspace init                 # Tekstinis konfigūratorius
```

### 3. Paleiskite
```bash
devspace                      # Paleidžia serverį. Automatiškai aptinka konfigūraciją.
```

Tai taip pat automatiškai paleidžia Cloudflare tunelį, jei `cloudflared` randamas kataloge `tools/`.

### 4. Prijunkite savo MCP klientą
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Arba vietiškai: `http://127.0.0.1:7676/mcp`

---

## Diegimas

Jokio Node.js, jokio npm, jokio Python. Vienas vykdomasis failas.

| Platforma | Atsisiuntimas |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompiliuoti natūraliai) |
| **macOS Intel** | `devspace` (GUI: kompiliuoti natūraliai) |
| **macOS M-chip** | `devspace` (GUI: kompiliuoti natūraliai) |

Reikalauja **Go 1.23+** tik jei kompiliuojate iš šaltinio.

---

## Ką gali AI

Prisijungęs AI gali atidaryti vieną iš jūsų patvirtintų projektų katalogų kaip darbo sritį:

- **Skaityti, rašyti ir redaguoti** failus darbo srityje
- **Ieškoti kodo** su reguliariosiomis išraiškomis ir tikrinti katalogus
- **Vykdyti apvalkalo komandas** (PowerShell Windows, bash Unix)
- **Atrasti projekto instrukcijas** iš `AGENTS.md` / `CLAUDE.md`
- **Automatiškai konfigūruoti** su nešiojama `.devspace/config.json`

8 MCP įrankiai: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigūracija

Visa konfigūracija yra **tame pačiame kataloge kaip ir vykdomasis failas** (nešiojama):

```
.devspace/
├── config.json       ← leidžiami šakniniai katalogai, prievadas, apvalkalas, kalba, autentifikacija
└── auth.json         ← savininko slaptažodis (neprivaloma)
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

| Laukas | Numatyta | Aprašymas |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatiškai aptinkama iš OS. Palaiko 47 kalbas |
| `toolMode` | `full` | `full` (visi įrankiai) arba `minimal` (tik apvalkalas paieškai) |
| `toolNaming` | `short` | `short` (read, write) arba `legacy` (read_file, write_file) |

Nereikia aplinkos kintamųjų — viskas yra nešiojamame konfigūracijos faile.

---

## Tunelis (nuotolinė prieiga)

ChatGPT žiniatinklio versijai (reikalingas HTTPS), DevSpace automatiškai paleidžia tunelį:

| Tunelis | URL tipas | Nustatymas |
|---|---|---|
| **Cloudflare** | Atsitiktinis (auto) | Įdėkite `cloudflared.exe` į `tools/` |
| **Pinggy** | Stabilus | Reikia SSH rakto (`ssh-keygen`) |

Serveris automatiškai aptinka, kuris iš jų prieinamas. Perkraukite serverį naujam Cloudflare URL arba naudokite Pinggy nuolatiniam URL.

---

## Apvalkalo palaikymas

| OS | Numatyta | Alternatyvos |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / bet koks apvalkalas |
| **macOS** | bash | `sh` / `zsh` |

Nustatykite `"shell"` config.json arba pasirinkite GUI.

---

## Saugumas

- **OAuth 2.0 su PKCE** — jei nustatytas savininko slaptažodis
- **Režimas be slaptažodžio** — jei slaptažodis nesukonfigūruotas, veikia be autentifikacijos
- **Kelio apribojimas** — visos failų operacijos tikrinamos pagal leidžiamus šakninius katalogus
- **Pasirenkamas tunelis** — Cloudflare tunelis apsaugo nuo tiesioginio poveikio
- **Jokių įkėlimų į trečiųjų šalių serverius** — jūsų kodas niekada nepalieka jūsų kompiuterio

---

## Kompiliavimas iš šaltinio

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Kompiliuoti viską (visos platformos)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Kompiliuoti tik dabartinei platformai
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformų palaikymas

| Platforma | Serveris | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompiliuoti natūraliai) |
| **macOS Intel** | ✅ | 🔧 (kompiliuoti natūraliai) |
| **macOS M-chip** | ✅ | 🔧 (kompiliuoti natūraliai) |

GUI reikalauja Fyne (OpenGL) — negalima kryžminio kompiliavimo. Serveris kompiliuojasi visur.

---

## Projekto struktūra

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP serveris
│   └── devspace-gui/       ← Darbalaukio GUI konfigūratorius (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE teikėjas
│   ├── config/             ← Nešiojama konfigūracijos sistema
│   ├── locales/            ← Vertimai į 47 kalbas
│   ├── logger/             ← Struktūrizuotas registravimas (zerolog)
│   ├── server/             ← HTTP + MCP + tunelio orkestravimas
│   ├── skills/             ← AGENTS.md / įgūdžių atradimas
│   ├── store/              ← SQLite darbo srities sesijos
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Darbo sritis ir kelio tikrinimas
├── scripts/
│   ├── windows/            ← PowerShell kompiliavimo scenarijus
│   └── unix/               ← Bash + Makefile kompiliavimo scenarijai
├── readme/                 ← Šio failo vertimai (47 kalbos)
├── build/                  ← Sukompiliuoti vykdomieji failai (visos platformos)
├── tools/                  ← cloudflared.exe ir kt.
├── go.mod / go.sum
└── README.md
```

---

Sukurta Go kalba. Nulis npm. Nulis Node.js. Vienas vykdomasis failas.
