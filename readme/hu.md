# DevSpace (Go kiadás)

**Adjon biztonságos hozzáférést a ChatGPT-nek és a Claude-nak a helyi gépéhez. Alakítson bármilyen MCP hosztot a kódolási partnerévé.**

A DevSpace egy önállóan üzemeltetett MCP szerver, amely lehetővé teszi az AI asszisztensek számára, hogy fájlokat olvassanak, szerkesszenek, keressenek és kódot futtassanak a valódi helyi projektjeiben — a fájljai, az eszközei, a terminálja — anélkül, hogy bármit feltöltene egy harmadik félhez. Ön futtatja a saját gépén, egy ön által vezérelt alagúton keresztül teszi elérhetővé, és opcionálisan jelszóval védi.

---

## Tartalomjegyzék

- [Gyors kezdés](#gyors-kezdés)
- [Telepítés](#telepítés)
- [Mit tud az AI](#mit-tud-az-ai)
- [Konfiguráció](#konfiguráció)
- [Alagút (távoli hozzáférés)](#alagút-távoli-hozzáférés)
- [Shell támogatás](#shell-támogatás)
- [Biztonság](#biztonság)
- [Fordítás forrásból](#fordítás-forrásból)
- [Platformtámogatás](#platformtámogatás)
- [Projektstruktúra](#projektstruktúra)

### 🌍 Fordítások

| Nyelv | | Nyelv | | Nyelv | |
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

## Gyors kezdés

### 1. Letöltés
Válassza ki a platformját a [Releases](../../releases) oldalról, vagy fordítsa forrásból:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurálás (GUI vagy szöveges)
```bash
devspace-gui                  # Asztali konfigurátor (GUI)
devspace init                 # Szöveges konfigurátor
```

### 3. Futtatás
```bash
devspace                      # Elindítja a szervert. Automatikusan felismeri a konfigot.
```

Ez automatikusan elindít egy Cloudflare alagutat is, ha a `cloudflared` megtalálható a `tools/` mappában.

### 4. Csatlakoztassa az MCP kliensét
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Vagy helyben: `http://127.0.0.1:7676/mcp`

---

## Telepítés

Nincs Node.js, nincs npm, nincs Python. Egyetlen bináris.

| Platform | Letöltés |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: natívan fordítható) |
| **macOS Intel** | `devspace` (GUI: natívan fordítható) |
| **macOS M-chip** | `devspace` (GUI: natívan fordítható) |

**Go 1.23+** csak akkor szükséges, ha forrásból épít.

---

## Mit tud az AI

Csatlakozás után az AI megnyithatja az egyik jóváhagyott projektmappáját munkaterületként:

- Fájlok **olvasása, írása és szerkesztése** a munkaterületen belül
- **Kód keresése** reguláris kifejezésekkel és könyvtárak vizsgálata
- **Shell parancsok futtatása** (PowerShell Windowson, bash Unixon)
- **Projektutasítások felfedezése** az `AGENTS.md` / `CLAUDE.md` fájlokból
- **Automatikus konfigurálás** a hordozható `.devspace/config.json` segítségével

8 MCP eszköz: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguráció

Minden konfiguráció **a futtatható fájllal azonos mappában** található (hordozható):

```
.devspace/
├── config.json       ← engedélyezett gyökerek, port, shell, nyelv, hitelesítés
└── auth.json         ← tulajdonosi jelszó (opcionális)
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

| Mező | Alapértelmezett | Leírás |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatikus felismerés az OS-ből. 47 nyelvet támogat |
| `toolMode` | `full` | `full` (minden eszköz) vagy `minimal` (csak shell kereséshez) |
| `toolNaming` | `short` | `short` (read, write) vagy `legacy` (read_file, write_file) |

Nincs szükség környezeti változókra — minden a hordozható konfigurációs fájlban van.

---

## Alagút (távoli hozzáférés)

A ChatGPT webes verziójához (HTTPS szükséges) a DevSpace automatikusan alagutat indít:

| Alagút | URL típusa | Beállítás |
|---|---|---|
| **Cloudflare** | Véletlenszerű (auto) | Helyezze a `cloudflared.exe` fájlt a `tools/` mappába |
| **Pinggy** | Stabil | SSH kulcs szükséges (`ssh-keygen`) |

A szerver automatikusan felismeri, melyik érhető el. Indítsa újra a szervert új Cloudflare URL-ért, vagy használja a Pinggy-t állandó URL-ért.

---

## Shell támogatás

| OS | Alapértelmezett | Alternatívák |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / bármely shell |
| **macOS** | bash | `sh` / `zsh` |

Állítsa be a `"shell"` értéket a config.json-ban, vagy válassza ki a GUI-ban.

---

## Biztonság

- **OAuth 2.0 PKCE-vel** — ha a tulajdonosi jelszó be van állítva
- **Jelszó nélküli mód** — ha nincs jelszó konfigurálva, hitelesítés nélkül fut
- **Útvonal-korlátozás** — minden fájlművelet ellenőrzése az engedélyezett gyökerek alapján
- **Opcionális alagút** — a Cloudflare alagút véd a közvetlen kitettségtől
- **Nincs harmadik félnek történő feltöltés** — a kódja soha nem hagyja el a gépét

---

## Fordítás forrásból

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Minden platformra fordítás
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Csak az aktuális platformra fordítás
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platformtámogatás

| Platform | Szerver | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (natívan fordítható) |
| **macOS Intel** | ✅ | 🔧 (natívan fordítható) |
| **macOS M-chip** | ✅ | 🔧 (natívan fordítható) |

A GUI Fyne-t (OpenGL) igényel — nem fordítható keresztbe. A szerver mindenhol fordítható.

---

## Projektstruktúra

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP szerver
│   └── devspace-gui/       ← Asztali GUI konfigurátor (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE szolgáltató
│   ├── config/             ← Hordozható konfigurációs rendszer
│   ├── locales/            ← 47 nyelvű fordítások
│   ├── logger/             ← Strukturált naplózás (zerolog)
│   ├── server/             ← HTTP + MCP + alagút vezénylés
│   ├── skills/             ← AGENTS.md / készségfelderítés
│   ├── store/              ← SQLite munkaterület munkamenetek
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Munkaterület és útvonal-ellenőrzés
├── scripts/
│   ├── windows/            ← PowerShell fordítási szkript
│   └── unix/               ← Bash + Makefile fordítási szkriptek
├── readme/                 ← E fájl fordításai (47 nyelv)
├── build/                  ← Lefordított binárisok (minden platform)
├── tools/                  ← cloudflared.exe, stb.
├── go.mod / go.sum
└── README.md
```

---

Go-ban készült. Nulla npm. Nulla Node.js. Egy bináris.
