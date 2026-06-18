# DevSpace (Go Edition) - Čeština

**Dejte ChatGPT a Claude bezpečný přístup k vašemu lokálnímu počítači. Proměňte jakéhokoli MCP hostitele ve svého parťáka pro programování.**

DevSpace je samo-hostovaný MCP server, který umožňuje AI asistentům číst, upravovat, vyhledávat a spouštět kód ve vašich skutečných lokálních projektech — vaše soubory, vaše nástroje, váš terminál — bez nahrávání čehokoli třetí straně. Spustíte ho na svém počítači, vystavíte ho přes tunel, který ovládáte, a volitelně ho zabezpečíte heslem.

---

## Obsah

- [Rychlý Start](#rychlý-start)
- [Instalace](#instalace)
- [Co AI Umí](#co-ai-umí)
- [Konfigurace](#konfigurace)
- [Tunel (Vzdálený Přístup)](#tunel-vzdálený-přístup)
- [Podpora Shellu](#podpora-shellu)
- [Zabezpečení](#zabezpečení)
- [Sestavení ze Zdrojového Kódu](#sestavení-ze-zdrojového-kódu)
- [Podpora Platforem](#podpora-platforem)
- [Struktura Projektu](#struktura-projektu)

### 🌍 Překlady

| Jazyk | | Jazyk | | Jazyk | |
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

## Rychlý Start

### 1. Stáhnout
Vyberte svou platformu z [Vydání](../../releases) nebo sestavte ze zdrojového kódu:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurovat (GUI nebo text)
```bash
devspace-gui                  # Desktopový konfigurátor (GUI)
devspace init                 # Textový konfigurátor
```

### 3. Spustit
```bash
devspace                      # Spustí server. Automaticky detekuje konfiguraci.
```

Tím se také automaticky spustí Cloudflare tunel, pokud je `cloudflared` nalezen v `tools/`.

### 4. Připojte svého MCP klienta
```
https://VÁŠ-TUNEL.trycloudflare.com/mcp
```
Nebo lokálně: `http://127.0.0.1:7676/mcp`

---

## Instalace

Žádný Node.js, žádný npm, žádný Python. Jediný binární soubor.

| Platforma | Stažení |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: zkompilujte nativně) |
| **macOS Intel** | `devspace` (GUI: zkompilujte nativně) |
| **macOS M-čip** | `devspace` (GUI: zkompilujte nativně) |

Vyžaduje **Go 1.23+** pouze při sestavování ze zdrojového kódu.

---

## Co AI Umí

Po připojení může AI otevřít jednu z vašich schválených složek projektu jako pracovní prostor:

- **Číst, zapisovat a upravovat** soubory uvnitř pracovního prostoru
- **Vyhledávat kód** pomocí regexu a prohlížet adresáře
- **Spouštět příkazy shellu** (PowerShell na Windows, bash na Unixu)
- **Objevovat instrukce projektu** z `AGENTS.md` / `CLAUDE.md`
- **Automaticky konfigurovat** pomocí přenosného `.devspace/config.json`

8 MCP nástrojů: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurace

Veškerá konfigurace je **ve stejné složce jako spustitelný soubor** (přenosná):

```
.devspace/
├── config.json       ← povolené kořeny, port, shell, jazyk, autentizace
└── auth.json         ← heslo vlastníka (volitelné)
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

| Pole | Výchozí | Popis |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatická detekce z OS. Podporuje 47 jazyků |
| `toolMode` | `full` | `full` (všechny nástroje) nebo `minimal` (pouze shell pro vyhledávání) |
| `toolNaming` | `short` | `short` (read, write) nebo `legacy` (read_file, write_file) |

Nejsou potřeba žádné proměnné prostředí — vše je v přenosném konfiguračním souboru.

---

## Tunel (Vzdálený Přístup)

Pro webovou verzi ChatGPT (vyžaduje HTTPS), DevSpace automaticky spouští tunel:

| Tunel | Typ URL | Nastavení |
|---|---|---|
| **Cloudflare** | Náhodné (auto) | Vložte `cloudflared.exe` do `tools/` |
| **Pinggy** | Stabilní | Potřebuje SSH klíč (`ssh-keygen`) |

Server automaticky detekuje, který je k dispozici. Restartujte server pro novou Cloudflare URL nebo použijte Pinggy pro trvalou URL.

---

## Podpora Shellu

| OS | Výchozí | Alternativy |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / jakýkoli shell |
| **macOS** | bash | `sh` / `zsh` |

Nastavte `"shell"` v config.json nebo vyberte v GUI.

---

## Zabezpečení

- **OAuth 2.0 s PKCE** — pokud je nastaveno heslo vlastníka
- **Režim bez hesla** — pokud není heslo nakonfigurováno, běží bez autentizace
- **Omezení cesty** — všechny souborové operace jsou validovány proti povoleným kořenům
- **Volitelný tunel** — Cloudflare tunel chrání před přímým vystavením
- **Žádná nahrávání třetím stranám** — váš kód nikdy neopouští váš počítač

---

## Sestavení ze Zdrojového Kódu

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Sestavit vše (všechny platformy)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Sestavit pouze pro aktuální platformu
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Podpora Platforem

| Platforma | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (zkompilujte nativně) |
| **macOS Intel** | ✅ | 🔧 (zkompilujte nativně) |
| **macOS M-čip** | ✅ | 🔧 (zkompilujte nativně) |

GUI vyžaduje Fyne (OpenGL) — nelze křížově kompilovat. Server se zkompiluje všude.

---

## Struktura Projektu

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP server
│   └── devspace-gui/       ← Desktopový GUI konfigurátor (Fyne)
├── internal/
│   ├── auth/               ← Poskytovatel OAuth 2.0 + PKCE
│   ├── config/             ← Přenosný konfigurační systém
│   ├── locales/            ← Překlady do 47 jazyků
│   ├── logger/             ← Strukturované logování (zerolog)
│   ├── server/             ← HTTP + MCP + orchestrace tunelu
│   ├── skills/             ← AGENTS.md / objevování dovedností
│   ├── store/              ← SQLite relace pracovního prostoru
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Pracovní prostor a validace cesty
├── scripts/
│   ├── windows/            ← PowerShell skript pro sestavení
│   └── unix/               ← Bash + Makefile skripty pro sestavení
├── readme/                 ← Překlady tohoto souboru (47 jazyků)
├── build/                  ← Zkompilované binární soubory (všechny platformy)
├── tools/                  ← cloudflared.exe, atd.
├── go.mod / go.sum
└── README.md
```

---

Postaveno v Go. Nula npm. Nula Node.js. Jeden binární soubor.
