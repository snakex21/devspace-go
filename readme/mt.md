# DevSpace (Edizzjoni Go)

**Agħti lil ChatGPT u Claude aċċess sigur għall-magna lokali tiegħek. Ibbiddel kwalunkwe host MCP fis-sieħeb tal-kodifikazzjoni tiegħek.**

DevSpace huwa server MCP li jospita lilek innifsek li jippermetti lill-assistenti tal-AI jaqraw, jeditjaw, ifittxu u jmexxu kodiċi fil-proġetti lokali reali tiegħek — il-fajls tiegħek, l-għodod tiegħek, it-terminal tiegħek — mingħajr ma ttella' xejn lil parti terza. Tħaddmu fuq il-magna tiegħek, tesponih permezz ta' mina li tikkontrolla, u b'mod fakultattiv tassigurah b'password.

---

## Werrej

- [Bidu Mgħaġġel](#bidu-mgħaġġel)
- [Installazzjoni](#installazzjoni)
- [X'Jista' Jagħmel l-AI](#x-jista-jagħmel-l-ai)
- [Konfigurazzjoni](#konfigurazzjoni)
- [Mina (Aċċess Remot)](#mina-aċċess-remot)
- [Appoġġ għall-Shell](#appoġġ-għall-shell)
- [Sigurtà](#sigurtà)
- [Bini mis-Sors](#bini-mis-sors)
- [Appoġġ għall-Pjattaformi](#appoġġ-għall-pjattaformi)
- [Struttura tal-Proġett](#struttura-tal-proġett)

### 🌍 Traduzzjonijiet

| Lingwa | | Lingwa | | Lingwa | |
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

## Bidu Mgħaġġel

### 1. Niżżel
Agħżel il-pjattaforma tiegħek minn [Releases](../../releases) jew ibni mis-sors:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Ikkonfigura (GUI jew test)
```bash
devspace-gui                  # Konfiguratur tad-desktop (GUI)
devspace init                 # Konfiguratur ibbażat fuq test
```

### 3. Mexxi
```bash
devspace                      # Jibda s-server. Awtomatikament jiskopri l-konfigurazzjoni.
```

Dan jibda wkoll awtomatikament Mina Cloudflare jekk `cloudflared` jinstab f'`tools/`.

### 4. Qabbad il-klijent MCP tiegħek
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Jew lokalment: `http://127.0.0.1:7676/mcp`

---

## Installazzjoni

Ebda Node.js, ebda npm, ebda Python. Fajl wieħed binarju.

| Pjattaforma | Tniżżil |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: ikkompilja b'mod nattiv) |
| **macOS Intel** | `devspace` (GUI: ikkompilja b'mod nattiv) |
| **macOS M-chip** | `devspace` (GUI: ikkompilja b'mod nattiv) |

Jeħtieġ **Go 1.23+** biss jekk tibni mis-sors.

---

## X'Jista' Jagħmel l-AI

Ladarba tkun konnessa, l-AI tista' tiftaħ waħda mill-fowlders tal-proġett approvati tiegħek bħala spazju tax-xogħol:

- **Taqra, tikteb u teditja** fajls ġewwa l-ispazju tax-xogħol
- **Tfittex kodiċi** b'regex u tispezzjona direttorji
- **Tmexxi kmandi shell** (PowerShell fuq Windows, bash fuq Unix)
- **Tiskopri struzzjonijiet tal-proġett** minn `AGENTS.md` / `CLAUDE.md`
- **Awtokonfigurazzjoni** b'`.devspace/config.json` portabbli

8 għodod MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurazzjoni

Il-konfigurazzjoni kollha tinsab **fl-istess folder bħall-eżegwibbli** (portabbli):

```
.devspace/
├── config.json       ← għeruq permessi, port, shell, lingwa, awtentikazzjoni
└── auth.json         ← password tas-sid (fakultattiv)
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

| Qasam | Default | Deskrizzjoni |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Awtoskoperta mill-OS. Jappoġġja 47 lingwa |
| `toolMode` | `full` | `full` (l-għodod kollha) jew `minimal` (shell biss għat-tiftix) |
| `toolNaming` | `short` | `short` (read, write) jew `legacy` (read_file, write_file) |

L-ebda varjabbli tal-ambjent meħtieġa — kollox jinsab fil-fajl tal-konfigurazzjoni portabbli.

---

## Mina (Aċċess Remot)

Għall-verżjoni web ta' ChatGPT (HTTPS meħtieġ), DevSpace jibda awtomatikament mina:

| Mina | Tip ta' URL | Setup |
|---|---|---|
| **Cloudflare** | Każwali (awto) | Poġġi `cloudflared.exe` f'`tools/` |
| **Pinggy** | Stabbli | Jeħtieġ ċavetta SSH (`ssh-keygen`) |

Is-server awtomatikament jiskopri liema waħda hija disponibbli. Ibda mill-ġdid is-server għal URL ġdid ta' Cloudflare, jew uża Pinggy għal URL permanenti.

---

## Appoġġ għall-Shell

| OS | Default | Alternattivi |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / kwalunkwe shell |
| **macOS** | bash | `sh` / `zsh` |

Issettja `"shell"` f'config.json jew agħżel fil-GUI.

---

## Sigurtà

- **OAuth 2.0 ma' PKCE** — jekk il-password tas-sid hija ssettjata
- **Modalità mingħajr password** — jekk l-ebda password mhija kkonfigurata, jaħdem mingħajr awtentikazzjoni
- **Konteniment tal-passaġġ** — l-operazzjonijiet kollha tal-fajls huma vvalidati kontra l-għeruq permessi
- **Mina fakultattiva** — Il-Mina Cloudflare tipproteġi minn espożizzjoni diretta
- **Ebda uploads lil partijiet terzi** — il-kodiċi tiegħek qatt ma jitlaq il-magna tiegħek

---

## Bini mis-Sors

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Ibni kollox (il-pjattaformi kollha)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Ibni biss għall-pjattaforma kurrenti
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Appoġġ għall-Pjattaformi

| Pjattaforma | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (ikkompilja b'mod nattiv) |
| **macOS Intel** | ✅ | 🔧 (ikkompilja b'mod nattiv) |
| **macOS M-chip** | ✅ | 🔧 (ikkompilja b'mod nattiv) |

Il-GUI teħtieġ Fyne (OpenGL) — ma tistax tikkompilja bejn il-pjattaformi. Is-server jikkompila kullimkien.

---

## Struttura tal-Proġett

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + server MCP
│   └── devspace-gui/       ← Konfiguratur GUI tad-desktop (Fyne)
├── internal/
│   ├── auth/               ← Fornitur OAuth 2.0 + PKCE
│   ├── config/             ← Sistema ta' konfigurazzjoni portabbli
│   ├── locales/            ← Traduzzjonijiet f'47 lingwa
│   ├── logger/             ← Logging strutturat (zerolog)
│   ├── server/             ← HTTP + MCP + orkestrazzjoni tal-mini
│   ├── skills/             ← AGENTS.md / skoperta tal-ħiliet
│   ├── store/              ← Sessjonijiet tal-ispazju tax-xogħol SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Spazju tax-xogħol u validazzjoni tal-passaġġ
├── scripts/
│   ├── windows/            ← Skript tal-bini PowerShell
│   └── unix/               ← Skripts tal-bini Bash + Makefile
├── readme/                 ← Traduzzjonijiet ta' dan il-fajl (47 lingwa)
├── build/                  ← Fajls binarji kkompilati (il-pjattaformi kollha)
├── tools/                  ← cloudflared.exe, eċċ.
├── go.mod / go.sum
└── README.md
```

---

Mibni f'Go. Żero npm. Żero Node.js. Fajl wieħed binarju.
