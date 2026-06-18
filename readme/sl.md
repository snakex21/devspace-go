# DevSpace (Go izdaja)

**Omogočite ChatGPT in Claude varen dostop do vašega lokalnega računalnika. Spremenite katerega koli gostitelja MCP v svojega programerskega partnerja.**

DevSpace je samo-gostujoč strežnik MCP, ki omogoča pomočnikom AI branje, urejanje, iskanje in zagon kode v vaših resničnih lokalnih projektih — vaše datoteke, vaša orodja, vaš terminal — brez nalaganja česarkoli tretjim osebam. Zaženete ga na svojem računalniku, izpostavite ga prek tunela, ki ga nadzorujete, in ga po želji zaščitite z geslom.

---

## Kazalo

- [Hitri Začetek](#hitri-začetek)
- [Namestitev](#namestitev)
- [Kaj Zmore AI](#kaj-zmore-ai)
- [Konfiguracija](#konfiguracija)
- [Tunel (Oddaljeni Dostop)](#tunel-oddaljeni-dostop)
- [Podpora Lupine](#podpora-lupine)
- [Varnost](#varnost)
- [Gradnja iz Izvorne Kode](#gradnja-iz-izvorne-kode)
- [Podpora Platform](#podpora-platform)
- [Struktura Projekta](#struktura-projekta)

### 🌍 Prevodi

| Jezik | | Jezik | | Jezik | |
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

## Hitri Začetek

### 1. Prenos
Izberite svojo platformo v [Releases](../../releases) ali zgradite iz izvorne kode:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfiguracija (GUI ali besedilo)
```bash
devspace-gui                  # Namizni konfigurator (GUI)
devspace init                 # Besedilni konfigurator
```

### 3. Zagon
```bash
devspace                      # Zažene strežnik. Samodejno zazna konfig.
```

To tudi samodejno zažene tunel Cloudflare, če je `cloudflared` najden v `tools/`.

### 4. Povežite svojega odjemalca MCP
```
https://VAŠ-TUNEL.trycloudflare.com/mcp
```
Ali lokalno: `http://127.0.0.1:7676/mcp`

---

## Namestitev

Brez Node.js, brez npm, brez Pythona. Ena binarna datoteka.

| Platforma | Prenos |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: prevedite domorodno) |
| **macOS Intel** | `devspace` (GUI: prevedite domorodno) |
| **macOS M-chip** | `devspace` (GUI: prevedite domorodno) |

Zahteva **Go 1.23+** samo pri gradnji iz izvorne kode.

---

## Kaj Zmore AI

Ko je povezan, lahko AI odpre eno od vaših odobrenih map projekta kot delovni prostor:

- **Bere, piše in ureja** datoteke znotraj delovnega prostora
- **Išče kodo** z regularnimi izrazi in pregleduje imenike
- **Zaganja ukaze lupine** (PowerShell na Windows, bash na Unix)
- **Odkriva navodila projekta** iz `AGENTS.md` / `CLAUDE.md`
- **Samodejno konfigurira** s prenosnim `.devspace/config.json`

8 orodij MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguracija

Vsa konfiguracija je **v isti mapi kot izvedljiva datoteka** (prenosna):

```
.devspace/
├── config.json       ← dovoljeni koreni, vrata, lupina, jezik, overjanje
└── auth.json         ← geslo lastnika (neobvezno)
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

| Polje | Privzeto | Opis |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Samodejno zaznavanje iz OS. Podpira 47 jezikov |
| `toolMode` | `full` | `full` (vsa orodja) ali `minimal` (samo lupina za iskanje) |
| `toolNaming` | `short` | `short` (read, write) ali `legacy` (read_file, write_file) |

Spremenljivke okolja niso potrebne — vse je v prenosni konfiguracijski datoteki.

---

## Tunel (Oddaljeni Dostop)

Za spletno različico ChatGPT (potreben HTTPS) DevSpace samodejno zažene tunel:

| Tunel | Vrsta URL | Nastavitev |
|---|---|---|
| **Cloudflare** | Naključni (avto) | Postavite `cloudflared.exe` v `tools/` |
| **Pinggy** | Stabilen | Potreben SSH ključ (`ssh-keygen`) |

Strežnik samodejno zazna, kateri je na voljo. Znova zaženite strežnik za nov URL Cloudflare ali uporabite Pinggy za trajni URL.

---

## Podpora Lupine

| OS | Privzeto | Alternative |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / katera koli lupina |
| **macOS** | bash | `sh` / `zsh` |

Nastavite `"shell"` v config.json ali izberite v GUI.

---

## Varnost

- **OAuth 2.0 s PKCE** — če je nastavljeno geslo lastnika
- **Način brez gesla** — če geslo ni konfigurirano, deluje brez overjanja
- **Omejitev poti** — vse datotečne operacije so preverjene glede na dovoljene korene
- **Neobvezen tunel** — tunel Cloudflare ščiti pred neposredno izpostavljenostjo
- **Brez nalaganj tretjim osebam** — vaša koda nikoli ne zapusti vašega računalnika

---

## Gradnja iz Izvorne Kode

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Zgradi vse (vse platforme)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Zgradi samo za trenutno platformo
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Podpora Platform

| Platforma | Strežnik | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (prevedite domorodno) |
| **macOS Intel** | ✅ | 🔧 (prevedite domorodno) |
| **macOS M-chip** | ✅ | 🔧 (prevedite domorodno) |

GUI zahteva Fyne (OpenGL) — navzkrižno prevajanje ni mogoče. Strežnik se prevede povsod.

---

## Struktura Projekta

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + strežnik MCP
│   └── devspace-gui/       ← Namizni konfigurator GUI (Fyne)
├── internal/
│   ├── auth/               ← Ponudnik OAuth 2.0 + PKCE
│   ├── config/             ← Sistem prenosne konfiguracije
│   ├── locales/            ← Prevodi v 47 jezikov
│   ├── logger/             ← Strukturirano beleženje (zerolog)
│   ├── server/             ← HTTP + MCP + orkestracija tunela
│   ├── skills/             ← Odkrivanje AGENTS.md / veščin
│   ├── store/              ← Seje delovnega prostora SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Preverjanje delovnega prostora in poti
├── scripts/
│   ├── windows/            ← Skripta za gradnjo PowerShell
│   └── unix/               ← Skripte za gradnjo Bash + Makefile
├── readme/                 ← Prevodi te datoteke (47 jezikov)
├── build/                  ← Prevedene binarne datoteke (vse platforme)
├── tools/                  ← cloudflared.exe, itd.
├── go.mod / go.sum
└── README.md
```

---

Zgrajeno v Go. Nič npm. Nič Node.js. Ena binarna datoteka.
