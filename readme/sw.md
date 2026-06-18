# DevSpace (Toleo la Go)

**Wape ChatGPT na Claude ufikiaji salama kwa mashine yako ya ndani. Geuza seva yoyote ya MCP kuwa mshirika wako wa usimbaji.**

DevSpace ni seva ya MCP inayojiendesha yenyewe inayoruhusu wasaidizi wa AI kusoma, kuhariri, kutafuta, na kuendesha msimbo katika miradi yako halisi ya ndani — faili zako, zana zako, terminal yako — bila kupakia chochote kwa mhusika wa tatu. Unaiendesha kwenye mashine yako, unaitoa kupitia handaki unayodhibiti, na kwa hiari unailinda kwa nenosiri.

---

## Yaliyomo

- [Kuanza Haraka](#kuanza-haraka)
- [Usakinishaji](#usakinishaji)
- [Nini AI Inaweza Kufanya](#nini-ai-inaweza-kufanya)
- [Usanidi](#usanidi)
- [Handaki (Ufikiaji wa Mbali)](#handaki-ufikiaji-wa-mbali)
- [Usaidizi wa Shell](#usaidizi-wa-shell)
- [Usalama](#usalama)
- [Kujenga kutoka Chanzo](#kujenga-kutoka-chanzo)
- [Usaidizi wa Jukwaa](#usaidizi-wa-jukwaa)
- [Muundo wa Mradi](#muundo-wa-mradi)

### 🌍 Tafsiri

| Lugha | | Lugha | | Lugha | |
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

## Kuanza Haraka

### 1. Pakua
Chagua jukwaa lako kutoka [Releases](../../releases) au jenga kutoka chanzo:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Sanidi (GUI au maandishi)
```bash
devspace-gui                  # Kisanidi cha desktop (GUI)
devspace init                 # Kisanidi cha maandishi
```

### 3. Endesha
```bash
devspace                      # Huanzisha seva. Hugundua usanidi kiotomatiki.
```

Hii pia huanzisha Handaki ya Cloudflare kiotomatiki ikiwa `cloudflared` inapatikana kwenye `tools/`.

### 4. Unganisha kiteja chako cha MCP
```
https://HANDAKI-YAKO.trycloudflare.com/mcp
```
Au kwa ndani: `http://127.0.0.1:7676/mcp`

---

## Usakinishaji

Hakuna Node.js, hakuna npm, hakuna Python. Faili moja ya binary.

| Jukwaa | Upakuaji |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: unganisha kiasili) |
| **macOS Intel** | `devspace` (GUI: unganisha kiasili) |
| **macOS M-chip** | `devspace` (GUI: unganisha kiasili) |

Inahitaji **Go 1.23+** tu ikiwa unajenga kutoka chanzo.

---

## Nini AI Inaweza Kufanya

Mara ikiwa imeunganishwa, AI inaweza kufungua moja ya folda zako za mradi zilizoidhinishwa kama eneo la kazi:

- **Kusoma, kuandika, na kuhariri** faili ndani ya eneo la kazi
- **Kutafuta msimbo** kwa regex na kukagua saraka
- **Kuendesha amri za shell** (PowerShell kwenye Windows, bash kwenye Unix)
- **Kugundua maagizo ya mradi** kutoka `AGENTS.md` / `CLAUDE.md`
- **Kujisanidi kiotomatiki** na `.devspace/config.json` inayobebeka

Zana 8 za MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Usanidi

Usanidi wote upo **kwenye folda sawa na faili inayotekelezeka** (inayobebeka):

```
.devspace/
├── config.json       ← mizizi inayoruhusiwa, mlango, shell, lugha, uthibitisho
└── auth.json         ← nenosiri la mmiliki (si lazima)
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

| Sehemu | Chaguomsingi | Maelezo |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Hugundua kiotomatiki kutoka OS. Inasaidia lugha 47 |
| `toolMode` | `full` | `full` (zana zote) au `minimal` (shell tu kwa utafutaji) |
| `toolNaming` | `short` | `short` (read, write) au `legacy` (read_file, write_file) |

Vigezo vya mazingira havihitajiki — kila kitu kipo kwenye faili ya usanidi inayobebeka.

---

## Handaki (Ufikiaji wa Mbali)

Kwa toleo la wavuti la ChatGPT (HTTPS inahitajika), DevSpace huanzisha handaki kiotomatiki:

| Handaki | Aina ya URL | Usanidi |
|---|---|---|
| **Cloudflare** | Nasibu (kiotomatiki) | Weka `cloudflared.exe` kwenye `tools/` |
| **Pinggy** | Imara | Inahitaji ufunguo wa SSH (`ssh-keygen`) |

Seva hugundua kiotomatiki ni ipi inayopatikana. Anzisha upya seva kwa URL mpya ya Cloudflare, au tumia Pinggy kwa URL ya kudumu.

---

## Usaidizi wa Shell

| OS | Chaguomsingi | Njia Mbadala |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / shell yoyote |
| **macOS** | bash | `sh` / `zsh` |

Weka `"shell"` kwenye config.json au chagua kwenye GUI.

---

## Usalama

- **OAuth 2.0 na PKCE** — ikiwa nenosiri la mmiliki limewekwa
- **Hali bila nenosiri** — ikiwa hakuna nenosiri lililosanidiwa, inaendesha bila uthibitisho
- **Uzuiaji wa njia** — shughuli zote za faili zinathibitishwa dhidi ya mizizi inayoruhusiwa
- **Handaki ya hiari** — Handaki ya Cloudflare inalinda dhidi ya mfichuo wa moja kwa moja
- **Hakuna upakiaji kwa wahusika wengine** — msimbo wako hauachi mashine yako kamwe

---

## Kujenga kutoka Chanzo

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Jenga kila kitu (majukwaa yote)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Jenga kwa jukwaa la sasa tu
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Usaidizi wa Jukwaa

| Jukwaa | Seva | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (unganisha kiasili) |
| **macOS Intel** | ✅ | 🔧 (unganisha kiasili) |
| **macOS M-chip** | ✅ | 🔧 (unganisha kiasili) |

GUI inahitaji Fyne (OpenGL) — haiwezi kuunganishwa kwa majukwaa tofauti. Seva inaunganishwa kila mahali.

---

## Muundo wa Mradi

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + seva ya MCP
│   └── devspace-gui/       ← Kisanidi cha desktop GUI (Fyne)
├── internal/
│   ├── auth/               ← Mtoa huduma wa OAuth 2.0 + PKCE
│   ├── config/             ← Mfumo wa usanidi unaobebeka
│   ├── locales/            ← Tafsiri za lugha 47
│   ├── logger/             ← Uwekaji kumbukumbu uliopangwa (zerolog)
│   ├── server/             ← HTTP + MCP + usimamizi wa handaki
│   ├── skills/             ← Ugunduzi wa AGENTS.md / ujuzi
│   ├── store/              ← Vipindi vya eneo la kazi SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Uthibitishaji wa eneo la kazi na njia
├── scripts/
│   ├── windows/            ← Hati ya ujenzi ya PowerShell
│   └── unix/               ← Hati za ujenzi za Bash + Makefile
├── readme/                 ← Tafsiri za faili hii (lugha 47)
├── build/                  ← Faili za binary zilizounganishwa (majukwaa yote)
├── tools/                  ← cloudflared.exe, n.k.
├── go.mod / go.sum
└── README.md
```

---

Imejengwa kwa Go. Zero npm. Zero Node.js. Binary moja.
