# DevSpace (Go संस्करण)

**ChatGPT और Claude को अपनी स्थानीय मशीन तक सुरक्षित पहुँच दें। किसी भी MCP होस्ट को अपना कोडिंग पार्टनर बनाएँ।**

DevSpace एक सेल्फ-होस्टेड MCP सर्वर है जो AI असिस्टेंट को आपके वास्तविक स्थानीय प्रोजेक्ट्स में फ़ाइलें पढ़ने, संपादित करने, खोजने और कोड चलाने देता है — आपकी फ़ाइलें, आपके टूल, आपका टर्मिनल — बिना किसी तीसरे पक्ष को कुछ भी अपलोड किए। आप इसे अपनी मशीन पर चलाते हैं, अपने नियंत्रण वाली टनल के माध्यम से एक्सपोज़ करते हैं, और वैकल्पिक रूप से पासवर्ड से सुरक्षित करते हैं।

---

## विषय सूची

- [त्वरित शुरुआत](#त्वरित-शुरुआत)
- [इंस्टॉलेशन](#इंस्टॉलेशन)
- [AI क्या कर सकता है](#ai-क्या-कर-सकता-है)
- [कॉन्फ़िगरेशन](#कॉन्फ़िगरेशन)
- [टनल (रिमोट एक्सेस)](#टनल-रिमोट-एक्सेस)
- [शेल समर्थन](#शेल-समर्थन)
- [सुरक्षा](#सुरक्षा)
- [सोर्स से बिल्ड करना](#सोर्स-से-बिल्ड-करना)
- [प्लेटफ़ॉर्म समर्थन](#प्लेटफ़ॉर्म-समर्थन)
- [प्रोजेक्ट संरचना](#प्रोजेक्ट-संरचना)

### 🌍 अनुवाद

| भाषा | | भाषा | | भाषा | |
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

## त्वरित शुरुआत

### 1. डाउनलोड करें
[Releases](../../releases) से अपना प्लेटफ़ॉर्म चुनें या सोर्स से बिल्ड करें:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. कॉन्फ़िगर करें (GUI या टेक्स्ट)
```bash
devspace-gui                  # डेस्कटॉप कॉन्फ़िगरेटर (GUI)
devspace init                 # टेक्स्ट-आधारित कॉन्फ़िगरेटर
```

### 3. चलाएँ
```bash
devspace                      # सर्वर शुरू करता है। कॉन्फ़िग स्वतः पहचानता है।
```

यदि `tools/` में `cloudflared` मौजूद है तो यह स्वतः Cloudflare टनल भी शुरू करता है।

### 4. अपना MCP क्लाइंट कनेक्ट करें
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
या स्थानीय रूप से: `http://127.0.0.1:7676/mcp`

---

## इंस्टॉलेशन

कोई Node.js, कोई npm, कोई Python नहीं। एकल बाइनरी।

| प्लेटफ़ॉर्म | डाउनलोड |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: मूल रूप से कंपाइल करें) |
| **macOS Intel** | `devspace` (GUI: मूल रूप से कंपाइल करें) |
| **macOS M-chip** | `devspace` (GUI: मूल रूप से कंपाइल करें) |

केवल सोर्स से बिल्ड करने पर **Go 1.23+** आवश्यक है।

---

## AI क्या कर सकता है

कनेक्ट होने के बाद, AI आपके स्वीकृत प्रोजेक्ट फ़ोल्डरों में से एक को वर्कस्पेस के रूप में खोल सकता है:

- वर्कस्पेस के भीतर फ़ाइलें **पढ़ना, लिखना और संपादित करना**
- रेगेक्स के साथ **कोड खोजना** और निर्देशिकाएँ निरीक्षण करना
- **शेल कमांड चलाना** (Windows पर PowerShell, Unix पर bash)
- `AGENTS.md` / `CLAUDE.md` से **प्रोजेक्ट निर्देश खोजना**
- पोर्टेबल `.devspace/config.json` के साथ **स्वतः-कॉन्फ़िगर**

8 MCP टूल: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## कॉन्फ़िगरेशन

सारी कॉन्फ़िग **निष्पादन योग्य फ़ाइल वाले फ़ोल्डर में** रहती है (पोर्टेबल):

```
.devspace/
├── config.json       ← अनुमत रूट, पोर्ट, शेल, भाषा, प्रमाणीकरण
└── auth.json         ← मालिक का पासवर्ड (वैकल्पिक)
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

| फ़ील्ड | डिफ़ॉल्ट | विवरण |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | OS से स्वतः पहचान। 47 भाषाओं का समर्थन |
| `toolMode` | `full` | `full` (सभी टूल) या `minimal` (केवल खोज के लिए शेल) |
| `toolNaming` | `short` | `short` (read, write) या `legacy` (read_file, write_file) |

किसी एनवायरनमेंट वेरिएबल की आवश्यकता नहीं — सब कुछ पोर्टेबल कॉन्फ़िग फ़ाइल में है।

---

## टनल (रिमोट एक्सेस)

ChatGPT वेब संस्करण (HTTPS आवश्यक) के लिए, DevSpace स्वतः टनल शुरू करता है:

| टनल | URL प्रकार | सेटअप |
|---|---|---|
| **Cloudflare** | यादृच्छिक (स्वतः) | `cloudflared.exe` को `tools/` में रखें |
| **Pinggy** | स्थिर | SSH की आवश्यकता (`ssh-keygen`) |

सर्वर स्वतः पहचानता है कि कौन सी उपलब्ध है। नए Cloudflare URL के लिए सर्वर पुनः आरंभ करें, या स्थायी URL के लिए Pinggy का उपयोग करें।

---

## शेल समर्थन

| OS | डिफ़ॉल्ट | विकल्प |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / कोई भी शेल |
| **macOS** | bash | `sh` / `zsh` |

config.json में `"shell"` सेट करें या GUI में चुनें।

---

## सुरक्षा

- **OAuth 2.0 with PKCE** — यदि मालिक का पासवर्ड सेट है
- **पासवर्ड-रहित मोड** — यदि कोई पासवर्ड कॉन्फ़िगर नहीं, तो बिना प्रमाणीकरण के चलता है
- **पाथ कंटेनमेंट** — सभी फ़ाइल ऑपरेशन अनुमत रूट के विरुद्ध मान्य होते हैं
- **वैकल्पिक टनल** — Cloudflare टनल सीधे एक्सपोज़र से बचाती है
- **कोई तृतीय-पक्ष अपलोड नहीं** — आपका कोड कभी आपकी मशीन नहीं छोड़ता

---

## सोर्स से बिल्ड करना

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# सब कुछ बिल्ड करें (सभी प्लेटफ़ॉर्म)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# केवल वर्तमान प्लेटफ़ॉर्म के लिए बिल्ड करें
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## प्लेटफ़ॉर्म समर्थन

| प्लेटफ़ॉर्म | सर्वर | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (मूल रूप से कंपाइल करें) |
| **macOS Intel** | ✅ | 🔧 (मूल रूप से कंपाइल करें) |
| **macOS M-chip** | ✅ | 🔧 (मूल रूप से कंपाइल करें) |

GUI को Fyne (OpenGL) की आवश्यकता है — क्रॉस-कंपाइल नहीं कर सकते। सर्वर हर जगह कंपाइल होता है।

---

## प्रोजेक्ट संरचना

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP सर्वर
│   └── devspace-gui/       ← डेस्कटॉप GUI कॉन्फ़िगरेटर (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE प्रदाता
│   ├── config/             ← पोर्टेबल कॉन्फ़िग सिस्टम
│   ├── locales/            ← 47 भाषा अनुवाद
│   ├── logger/             ← संरचित लॉगिंग (zerolog)
│   ├── server/             ← HTTP + MCP + टनल ऑर्केस्ट्रेशन
│   ├── skills/             ← AGENTS.md / स्किल खोज
│   ├── store/              ← SQLite वर्कस्पेस सेशन
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← वर्कस्पेस और पाथ सत्यापन
├── scripts/
│   ├── windows/            ← PowerShell बिल्ड स्क्रिप्ट
│   └── unix/               ← Bash + Makefile बिल्ड स्क्रिप्ट
├── readme/                 ← इस फ़ाइल के अनुवाद (47 भाषाएँ)
├── build/                  ← कंपाइल्ड बाइनरीज़ (सभी प्लेटफ़ॉर्म)
├── tools/                  ← cloudflared.exe, आदि
├── go.mod / go.sum
└── README.md
```

---

Go में निर्मित। शून्य npm। शून्य Node.js। एक बाइनरी।
