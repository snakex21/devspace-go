# DevSpace (Go பதிப்பு)

**ChatGPT மற்றும் Claude-க்கு உங்கள் உள்ளூர் கணினிக்கு பாதுகாப்பான அணுகலை வழங்குங்கள். எந்த MCP ஹோஸ்டையும் உங்கள் குறியீட்டு பங்காளியாக மாற்றுங்கள்.**

DevSpace என்பது சுய-ஹோஸ்ட் செய்யப்பட்ட MCP சேவையகமாகும், இது AI உதவியாளர்கள் உங்கள் உண்மையான உள்ளூர் திட்டங்களில் — உங்கள் கோப்புகள், உங்கள் கருவிகள், உங்கள் முனையம் — எதையும் மூன்றாம் தரப்பினருக்கு பதிவேற்றாமல் குறியீட்டைப் படிக்க, திருத்த, தேட மற்றும் இயக்க அனுமதிக்கிறது. நீங்கள் அதை உங்கள் கணினியில் இயக்கி, நீங்கள் கட்டுப்படுத்தும் ஒரு சுரங்கப்பாதை வழியாக வெளிப்படுத்தி, விருப்பமாக கடவுச்சொல்லுடன் பாதுகாக்கலாம்.

---

## பொருளடக்கம்

- [விரைவான தொடக்கம்](#விரைவான-தொடக்கம்)
- [நிறுவல்](#நிறுவல்)
- [AI என்ன செய்ய முடியும்](#ai-என்ன-செய்ய-முடியும்)
- [உள்ளமைவு](#உள்ளமைவு)
- [சுரங்கப்பாதை (தொலைநிலை அணுகல்)](#சுரங்கப்பாதை-தொலைநிலை-அணுகல்)
- [ஷெல் ஆதரவு](#ஷெல்-ஆதரவு)
- [பாதுகாப்பு](#பாதுகாப்பு)
- [மூலத்திலிருந்து உருவாக்குதல்](#மூலத்திலிருந்து-உருவாக்குதல்)
- [தள ஆதரவு](#தள-ஆதரவு)
- [திட்ட அமைப்பு](#திட்ட-அமைப்பு)

### 🌍 மொழிபெயர்ப்புகள்

| மொழி | | மொழி | | மொழி | |
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

## விரைவான தொடக்கம்

### 1. பதிவிறக்கம்
[Releases](../../releases) இலிருந்து உங்கள் தளத்தைத் தேர்ந்தெடுக்கவும் அல்லது மூலத்திலிருந்து உருவாக்கவும்:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. உள்ளமைவு (GUI அல்லது உரை)
```bash
devspace-gui                  # டெஸ்க்டாப் உள்ளமைப்பான் (GUI)
devspace init                 # உரை-அடிப்படையிலான உள்ளமைப்பான்
```

### 3. இயக்கு
```bash
devspace                      # சேவையகத்தைத் தொடங்குகிறது. உள்ளமைவைத் தானாகக் கண்டறியும்.
```

`tools/` இல் `cloudflared` காணப்பட்டால், இது Cloudflare சுரங்கப்பாதையையும் தானாகத் தொடங்கும்.

### 4. உங்கள் MCP கிளையண்டை இணைக்கவும்
```
https://உங்கள்-சுரங்கப்பாதை.trycloudflare.com/mcp
```
அல்லது உள்ளூரில்: `http://127.0.0.1:7676/mcp`

---

## நிறுவல்

Node.js இல்லை, npm இல்லை, Python இல்லை. ஒரே ஒரு பைனரி.

| தளம் | பதிவிறக்கம் |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: நேட்டிவ்வாக தொகுக்கவும்) |
| **macOS Intel** | `devspace` (GUI: நேட்டிவ்வாக தொகுக்கவும்) |
| **macOS M-chip** | `devspace` (GUI: நேட்டிவ்வாக தொகுக்கவும்) |

மூலத்திலிருந்து உருவாக்கும்போது மட்டும் **Go 1.23+** தேவை.

---

## AI என்ன செய்ய முடியும்

இணைக்கப்பட்டவுடன், AI உங்கள் அங்கீகரிக்கப்பட்ட திட்ட கோப்புறைகளில் ஒன்றை பணியிடமாகத் திறக்க முடியும்:

- பணியிடத்திற்குள் கோப்புகளை **படிக்க, எழுத மற்றும் திருத்த**
- regex மூலம் **குறியீட்டைத் தேட** மற்றும் கோப்பகங்களை ஆய்வு செய்ய
- **ஷெல் கட்டளைகளை இயக்க** (Windows இல் PowerShell, Unix இல் bash)
- `AGENTS.md` / `CLAUDE.md` இலிருந்து **திட்ட வழிமுறைகளைக் கண்டறிய**
- கையடக்க `.devspace/config.json` மூலம் **தானாக உள்ளமைக்க**

8 MCP கருவிகள்: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## உள்ளமைவு

அனைத்து உள்ளமைவும் **இயங்கக்கூடிய கோப்பின் அதே கோப்புறையில்** உள்ளது (கையடக்கமானது):

```
.devspace/
├── config.json       ← அனுமதிக்கப்பட்ட ரூட்கள், போர்ட், ஷெல், மொழி, அங்கீகாரம்
└── auth.json         ← உரிமையாளர் கடவுச்சொல் (விரும்பினால்)
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

| புலம் | இயல்புநிலை | விளக்கம் |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | OS இலிருந்து தானாகக் கண்டறியும். 47 மொழிகளை ஆதரிக்கிறது |
| `toolMode` | `full` | `full` (அனைத்து கருவிகள்) அல்லது `minimal` (தேடலுக்கு ஷெல் மட்டும்) |
| `toolNaming` | `short` | `short` (read, write) அல்லது `legacy` (read_file, write_file) |

சூழல் மாறிகள் தேவையில்லை — அனைத்தும் கையடக்க உள்ளமைவு கோப்பில் உள்ளது.

---

## சுரங்கப்பாதை (தொலைநிலை அணுகல்)

ChatGPT வலை பதிப்பிற்கு (HTTPS தேவை), DevSpace தானாக ஒரு சுரங்கப்பாதையைத் தொடங்கும்:

| சுரங்கப்பாதை | URL வகை | அமைப்பு |
|---|---|---|
| **Cloudflare** | சீரற்ற (தானியங்கு) | `cloudflared.exe` ஐ `tools/` இல் வைக்கவும் |
| **Pinggy** | நிலையானது | SSH விசை தேவை (`ssh-keygen`) |

எது கிடைக்கிறது என்பதை சேவையகம் தானாகக் கண்டறியும். புதிய Cloudflare URL-க்கு சேவையகத்தை மறுதொடக்கம் செய்யவும் அல்லது நிரந்தர URL-க்கு Pinggy ஐப் பயன்படுத்தவும்.

---

## ஷெல் ஆதரவு

| OS | இயல்புநிலை | மாற்றுகள் |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / எந்த ஷெல்லும் |
| **macOS** | bash | `sh` / `zsh` |

config.json இல் `"shell"` ஐ அமைக்கவும் அல்லது GUI இல் தேர்வு செய்யவும்.

---

## பாதுகாப்பு

- **PKCE உடன் OAuth 2.0** — உரிமையாளர் கடவுச்சொல் அமைக்கப்பட்டால்
- **கடவுச்சொல் இல்லாத முறை** — கடவுச்சொல் உள்ளமைக்கப்படவில்லை எனில், அங்கீகாரம் இல்லாமல் இயங்கும்
- **பாதை கட்டுப்பாடு** — அனைத்து கோப்பு செயல்பாடுகளும் அனுமதிக்கப்பட்ட ரூட்களுக்கு எதிராக சரிபார்க்கப்படும்
- **விருப்ப சுரங்கப்பாதை** — Cloudflare சுரங்கப்பாதை நேரடி வெளிப்பாட்டிலிருந்து பாதுகாக்கிறது
- **மூன்றாம் தரப்பு பதிவேற்றங்கள் இல்லை** — உங்கள் குறியீடு உங்கள் கணினியை விட்டு வெளியேறாது

---

## மூலத்திலிருந்து உருவாக்குதல்

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# அனைத்தையும் உருவாக்கு (அனைத்து தளங்களும்)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# தற்போதைய தளத்திற்கு மட்டும் உருவாக்கு
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## தள ஆதரவு

| தளம் | சேவையகம் | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (நேட்டிவ்வாக தொகுக்கவும்) |
| **macOS Intel** | ✅ | 🔧 (நேட்டிவ்வாக தொகுக்கவும்) |
| **macOS M-chip** | ✅ | 🔧 (நேட்டிவ்வாக தொகுக்கவும்) |

GUI க்கு Fyne (OpenGL) தேவை — குறுக்கு-தொகுப்பு செய்ய முடியாது. சேவையகம் எல்லா இடங்களிலும் தொகுக்கப்படும்.

---

## திட்ட அமைப்பு

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP சேவையகம்
│   └── devspace-gui/       ← டெஸ்க்டாப் GUI உள்ளமைப்பான் (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE வழங்குநர்
│   ├── config/             ← கையடக்க உள்ளமைவு அமைப்பு
│   ├── locales/            ← 47 மொழி மொழிபெயர்ப்புகள்
│   ├── logger/             ← கட்டமைக்கப்பட்ட பதிவு (zerolog)
│   ├── server/             ← HTTP + MCP + சுரங்கப்பாதை ஒருங்கிணைப்பு
│   ├── skills/             ← AGENTS.md / திறன் கண்டுபிடிப்பு
│   ├── store/              ← SQLite பணியிட அமர்வுகள்
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← பணியிடம் & பாதை சரிபார்ப்பு
├── scripts/
│   ├── windows/            ← PowerShell உருவாக்க ஸ்கிரிப்ட்
│   └── unix/               ← Bash + Makefile உருவாக்க ஸ்கிரிப்ட்கள்
├── readme/                 ← இந்தக் கோப்பின் மொழிபெயர்ப்புகள் (47 மொழிகள்)
├── build/                  ← தொகுக்கப்பட்ட பைனரிகள் (அனைத்து தளங்களும்)
├── tools/                  ← cloudflared.exe, முதலியன.
├── go.mod / go.sum
└── README.md
```

---

Go-வில் கட்டப்பட்டது. பூஜ்யம் npm. பூஜ்யம் Node.js. ஒரு பைனரி.
