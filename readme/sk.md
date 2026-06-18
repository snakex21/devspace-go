# DevSpace (Go edícia)

**Poskytnite ChatGPT a Claude bezpečný prístup k vášmu lokálnemu počítaču. Premeňte akéhokoľvek MCP hostiteľa na svojho programovacieho partnera.**

DevSpace je samo-hostovaný MCP server, ktorý umožňuje AI asistentom čítať, upravovať, vyhľadávať a spúšťať kód vo vašich skutočných lokálnych projektoch — vaše súbory, vaše nástroje, váš terminál — bez nahrávania čohokoľvek tretej strane. Spustíte ho na svojom počítači, sprístupníte ho cez tunel, ktorý ovládate, a voliteľne ho zabezpečíte heslom.

---

## Obsah

- [Rýchly Štart](#rýchly-štart)
- [Inštalácia](#inštalácia)
- [Čo Dokáže AI](#čo-dokáže-ai)
- [Konfigurácia](#konfigurácia)
- [Tunel (Vzdialený Prístup)](#tunel-vzdialený-prístup)
- [Podpora Shellu](#podpora-shellu)
- [Bezpečnosť](#bezpečnosť)
- [Kompilácia zo Zdrojov](#kompilácia-zo-zdrojov)
- [Podpora Platforiem](#podpora-platforiem)
- [Štruktúra Projektu](#štruktúra-projektu)

### 🌍 Preklady

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

## Rýchly Štart

### 1. Stiahnutie
Vyberte si platformu z [Releases](../../releases) alebo skompilujte zo zdrojov:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurácia (GUI alebo text)
```bash
devspace-gui                  # Desktopový konfigurátor (GUI)
devspace init                 # Textový konfigurátor
```

### 3. Spustenie
```bash
devspace                      # Spustí server. Automaticky deteguje konfig.
```

Toto tiež automaticky spustí Cloudflare tunel, ak sa `cloudflared` nachádza v `tools/`.

### 4. Pripojenie MCP klienta
```
https://VÁŠ-TUNEL.trycloudflare.com/mcp
```
Alebo lokálne: `http://127.0.0.1:7676/mcp`

---

## Inštalácia

Žiadny Node.js, žiadny npm, žiadny Python. Jeden binárny súbor.

| Platforma | Stiahnutie |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilujte natívne) |
| **macOS Intel** | `devspace` (GUI: kompilujte natívne) |
| **macOS M-chip** | `devspace` (GUI: kompilujte natívne) |

Vyžaduje **Go 1.23+** len pri kompilácii zo zdrojov.

---

## Čo Dokáže AI

Po pripojení môže AI otvoriť jeden z vašich schválených priečinkov projektu ako pracovný priestor:

- **Čítať, zapisovať a upravovať** súbory v rámci pracovného priestoru
- **Vyhľadávať kód** pomocou regulárnych výrazov a prezerať adresáre
- **Spúšťať príkazy shellu** (PowerShell na Windows, bash na Unix)
- **Objavovať inštrukcie projektu** z `AGENTS.md` / `CLAUDE.md`
- **Automaticky konfigurovať** pomocou prenosného `.devspace/config.json`

8 MCP nástrojov: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurácia

Celá konfigurácia sa nachádza **v rovnakom priečinku ako spustiteľný súbor** (prenosná):

```
.devspace/
├── config.json       ← povolené korene, port, shell, jazyk, overenie
└── auth.json         ← heslo vlastníka (voliteľné)
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

| Pole | Predvolené | Popis |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatická detekcia z OS. Podporuje 47 jazykov |
| `toolMode` | `full` | `full` (všetky nástroje) alebo `minimal` (len shell pre vyhľadávanie) |
| `toolNaming` | `short` | `short` (read, write) alebo `legacy` (read_file, write_file) |

Nie sú potrebné žiadne premenné prostredia — všetko je v prenosnom konfiguračnom súbore.

---

## Tunel (Vzdialený Prístup)

Pre webovú verziu ChatGPT (vyžaduje HTTPS) DevSpace automaticky spúšťa tunel:

| Tunel | Typ URL | Nastavenie |
|---|---|---|
| **Cloudflare** | Náhodný (auto) | Vložte `cloudflared.exe` do `tools/` |
| **Pinggy** | Stabilný | Vyžaduje SSH kľúč (`ssh-keygen`) |

Server automaticky deteguje, ktorý je k dispozícii. Reštartujte server pre novú Cloudflare URL alebo použite Pinggy pre trvalú URL.

---

## Podpora Shellu

| OS | Predvolený | Alternatívy |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / akýkoľvek shell |
| **macOS** | bash | `sh` / `zsh` |

Nastavte `"shell"` v config.json alebo vyberte v GUI.

---

## Bezpečnosť

- **OAuth 2.0 s PKCE** — ak je nastavené heslo vlastníka
- **Režim bez hesla** — ak nie je nakonfigurované žiadne heslo, beží bez overenia
- **Obmedzenie ciest** — všetky súborové operácie sú overené voči povoleným koreňom
- **Voliteľný tunel** — Cloudflare tunel chráni pred priamym vystavením
- **Žiadne nahrávanie tretím stranám** — váš kód nikdy neopustí váš počítač

---

## Kompilácia zo Zdrojov

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Kompilácia všetkého (všetky platformy)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Kompilácia len pre aktuálnu platformu
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Podpora Platforiem

| Platforma | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilujte natívne) |
| **macOS Intel** | ✅ | 🔧 (kompilujte natívne) |
| **macOS M-chip** | ✅ | 🔧 (kompilujte natívne) |

GUI vyžaduje Fyne (OpenGL) — nie je možné krížovo kompilovať. Server sa kompiluje všade.

---

## Štruktúra Projektu

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP server
│   └── devspace-gui/       ← Desktopový GUI konfigurátor (Fyne)
├── internal/
│   ├── auth/               ← Poskytovateľ OAuth 2.0 + PKCE
│   ├── config/             ← Systém prenosnej konfigurácie
│   ├── locales/            ← Preklady do 47 jazykov
│   ├── logger/             ← Štruktúrované logovanie (zerolog)
│   ├── server/             ← HTTP + MCP + orchestrácia tunela
│   ├── skills/             ← Objavovanie AGENTS.md / zručností
│   ├── store/              ← Relácie pracovného priestoru SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Overenie pracovného priestoru a ciest
├── scripts/
│   ├── windows/            ← PowerShell skript na kompiláciu
│   └── unix/               ← Bash + Makefile skripty na kompiláciu
├── readme/                 ← Preklady tohto súboru (47 jazykov)
├── build/                  ← Skompilované binárne súbory (všetky platformy)
├── tools/                  ← cloudflared.exe, atď.
├── go.mod / go.sum
└── README.md
```

---

Vytvorené v Go. Nula npm. Nula Node.js. Jeden binárny súbor.
