# DevSpace (گو ایڈیشن)

**ChatGPT اور Claude کو اپنی مقامی مشین تک محفوظ رسائی دیں۔ کسی بھی MCP میزبان کو اپنا کوڈنگ پارٹنر بنائیں۔**

DevSpace ایک خود میزبان MCP سرور ہے جو AI معاونین کو آپ کے حقیقی مقامی منصوبوں میں — آپ کی فائلیں، آپ کے ٹولز، آپ کا ٹرمینل — بغیر کسی تیسرے فریق کو کچھ اپ لوڈ کیے کوڈ پڑھنے، ترمیم کرنے، تلاش کرنے اور چلانے کی اجازت دیتا ہے۔ آپ اسے اپنی مشین پر چلاتے ہیں، ایک سرنگ کے ذریعے ظاہر کرتے ہیں جسے آپ کنٹرول کرتے ہیں، اور اختیاری طور پر پاس ورڈ سے محفوظ کرتے ہیں۔

---

## فہرست مضامین

- [فوری آغاز](#فوری-آغاز)
- [تنصیب](#تنصیب)
- [AI کیا کر سکتا ہے](#ai-کیا-کر-سکتا-ہے)
- [تشکیل](#تشکیل)
- [سرنگ (دور دراز رسائی)](#سرنگ-دور-دراز-رسائی)
- [شیل سپورٹ](#شیل-سپورٹ)
- [سیکورٹی](#سیکورٹی)
- [ماخذ سے تعمیر](#ماخذ-سے-تعمیر)
- [پلیٹ فارم سپورٹ](#پلیٹ-فارم-سپورٹ)
- [منصوبے کا ڈھانچہ](#منصوبے-کا-ڈھانچہ)

### 🌍 تراجم

| زبان | | زبان | | زبان | |
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

## فوری آغاز

### 1. ڈاؤن لوڈ کریں
[Releases](../../releases) سے اپنا پلیٹ فارم منتخب کریں یا ماخذ سے تعمیر کریں:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. تشکیل دیں (GUI یا متن)
```bash
devspace-gui                  # ڈیسک ٹاپ تشکیل کار (GUI)
devspace init                 # متن پر مبنی تشکیل کار
```

### 3. چلائیں
```bash
devspace                      # سرور شروع کرتا ہے۔ تشکیل کو خودکار طور پر شناخت کرتا ہے۔
```

یہ خودکار طور پر Cloudflare سرنگ بھی شروع کرتا ہے اگر `cloudflared` کو `tools/` میں پایا جائے۔

### 4. اپنے MCP کلائنٹ کو جوڑیں
```
https://آپ-کی-سرنگ.trycloudflare.com/mcp
```
یا مقامی طور پر: `http://127.0.0.1:7676/mcp`

---

## تنصیب

کوئی Node.js نہیں، کوئی npm نہیں، کوئی Python نہیں۔ ایک واحد بائنری۔

| پلیٹ فارم | ڈاؤن لوڈ |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: مقامی طور پر مرتب کریں) |
| **macOS Intel** | `devspace` (GUI: مقامی طور پر مرتب کریں) |
| **macOS M-chip** | `devspace` (GUI: مقامی طور پر مرتب کریں) |

صرف ماخذ سے تعمیر کرتے وقت **Go 1.23+** درکار ہے۔

---

## AI کیا کر سکتا ہے

منسلک ہونے کے بعد، AI آپ کے منظور شدہ منصوبے کے فولڈرز میں سے ایک کو ورک اسپیس کے طور پر کھول سکتا ہے:

- ورک اسپیس کے اندر فائلوں کو **پڑھنا، لکھنا اور ترمیم کرنا**
- regex کے ساتھ **کوڈ تلاش کرنا** اور ڈائریکٹریوں کا معائنہ کرنا
- **شیل کمانڈز چلانا** (Windows پر PowerShell، Unix پر bash)
- `AGENTS.md` / `CLAUDE.md` سے **منصوبے کی ہدایات دریافت کرنا**
- پورٹیبل `.devspace/config.json` کے ساتھ **خودکار تشکیل**

8 MCP ٹولز: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## تشکیل

تمام تشکیل **قابل عمل فائل کے اسی فولڈر میں** ہوتی ہے (پورٹیبل):

```
.devspace/
├── config.json       ← اجازت شدہ جڑیں، پورٹ، شیل، زبان، تصدیق
└── auth.json         ← مالک کا پاس ورڈ (اختیاری)
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

| فیلڈ | طے شدہ | تفصیل |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | OS سے خودکار شناخت۔ 47 زبانوں کو سپورٹ کرتا ہے |
| `toolMode` | `full` | `full` (تمام ٹولز) یا `minimal` (صرف تلاش کے لیے شیل) |
| `toolNaming` | `short` | `short` (read, write) یا `legacy` (read_file, write_file) |

ماحولیاتی متغیرات کی ضرورت نہیں — سب کچھ پورٹیبل تشکیل فائل میں ہے۔

---

## سرنگ (دور دراز رسائی)

ChatGPT ویب ورژن کے لیے (HTTPS درکار ہے)، DevSpace خودکار طور پر ایک سرنگ شروع کرتا ہے:

| سرنگ | URL قسم | سیٹ اپ |
|---|---|---|
| **Cloudflare** | بے ترتیب (خودکار) | `cloudflared.exe` کو `tools/` میں رکھیں |
| **Pinggy** | مستحکم | SSH کلید درکار ہے (`ssh-keygen`) |

سرور خودکار طور پر شناخت کرتا ہے کہ کون سا دستیاب ہے۔ نئے Cloudflare URL کے لیے سرور کو دوبارہ شروع کریں، یا مستقل URL کے لیے Pinggy استعمال کریں۔

---

## شیل سپورٹ

| OS | طے شدہ | متبادل |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / کوئی بھی شیل |
| **macOS** | bash | `sh` / `zsh` |

config.json میں `"shell"` سیٹ کریں یا GUI میں منتخب کریں۔

---

## سیکورٹی

- **PKCE کے ساتھ OAuth 2.0** — اگر مالک کا پاس ورڈ سیٹ ہو
- **پاس ورڈ کے بغیر موڈ** — اگر کوئی پاس ورڈ تشکیل نہیں دیا گیا، تو تصدیق کے بغیر چلتا ہے
- **راہ کی پابندی** — تمام فائل آپریشنز اجازت شدہ جڑوں کے خلاف توثیق شدہ ہیں
- **اختیاری سرنگ** — Cloudflare سرنگ براہ راست نمائش سے بچاتی ہے
- **کوئی تیسرے فریق کو اپ لوڈ نہیں** — آپ کا کوڈ کبھی آپ کی مشین سے باہر نہیں جاتا

---

## ماخذ سے تعمیر

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# سب کچھ تعمیر کریں (تمام پلیٹ فارم)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# صرف موجودہ پلیٹ فارم کے لیے تعمیر کریں
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## پلیٹ فارم سپورٹ

| پلیٹ فارم | سرور | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (مقامی طور پر مرتب کریں) |
| **macOS Intel** | ✅ | 🔧 (مقامی طور پر مرتب کریں) |
| **macOS M-chip** | ✅ | 🔧 (مقامی طور پر مرتب کریں) |

GUI کو Fyne (OpenGL) درکار ہے — کراس مرتب نہیں کیا جا سکتا۔ سرور ہر جگہ مرتب ہوتا ہے۔

---

## منصوبے کا ڈھانچہ

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP سرور
│   └── devspace-gui/       ← ڈیسک ٹاپ GUI تشکیل کار (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE فراہم کنندہ
│   ├── config/             ← پورٹیبل تشکیل کا نظام
│   ├── locales/            ← 47 زبانوں کے تراجم
│   ├── logger/             ← منظم لاگنگ (zerolog)
│   ├── server/             ← HTTP + MCP + سرنگ کی آرکیسٹریشن
│   ├── skills/             ← AGENTS.md / مہارت کی دریافت
│   ├── store/              ← SQLite ورک اسپیس سیشنز
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← ورک اسپیس اور راہ کی توثیق
├── scripts/
│   ├── windows/            ← PowerShell تعمیر اسکرپٹ
│   └── unix/               ← Bash + Makefile تعمیر اسکرپٹس
├── readme/                 ← اس فائل کے تراجم (47 زبانیں)
├── build/                  ← مرتب شدہ بائنریاں (تمام پلیٹ فارم)
├── tools/                  ← cloudflared.exe، وغیرہ
├── go.mod / go.sum
└── README.md
```

---

Go میں بنایا گیا۔ صفر npm۔ صفر Node.js۔ ایک بائنری۔
