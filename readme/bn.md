# DevSpace (Go Edition) - বাংলা

**ChatGPT ও Claude-কে আপনার স্থানীয় মেশিনে নিরাপদ অ্যাক্সেস দিন। যেকোনো MCP হোস্টকে আপনার কোডিং পার্টনারে রূপান্তর করুন।**

DevSpace হল একটি স্ব-হোস্টেড MCP সার্ভার যা AI সহায়কদের আপনার বাস্তব স্থানীয় প্রকল্পগুলিতে কোড পড়তে, সম্পাদনা করতে, অনুসন্ধান করতে এবং চালাতে দেয় — আপনার ফাইল, আপনার সরঞ্জাম, আপনার টার্মিনাল — কোনো তৃতীয় পক্ষের কাছে কিছু আপলোড না করেই। আপনি এটি আপনার মেশিনে চালান, আপনার নিয়ন্ত্রিত একটি টানেলের মাধ্যমে প্রকাশ করেন এবং ঐচ্ছিকভাবে একটি পাসওয়ার্ড দিয়ে সুরক্ষিত করেন।

---

## সূচিপত্র

- [দ্রুত শুরু](#দ্রুত-শুরু)
- [ইনস্টলেশন](#ইনস্টলেশন)
- [AI কী করতে পারে](#ai-কী-করতে-পারে)
- [কনফিগারেশন](#কনফিগারেশন)
- [টানেল (দূরবর্তী অ্যাক্সেস)](#টানেল-দূরবর্তী-অ্যাক্সেস)
- [শেল সমর্থন](#শেল-সমর্থন)
- [নিরাপত্তা](#নিরাপত্তা)
- [সোর্স থেকে বিল্ড](#সোর্স-থেকে-বিল্ড)
- [প্ল্যাটফর্ম সমর্থন](#প্ল্যাটফর্ম-সমর্থন)
- [প্রকল্প কাঠামো](#প্রকল্প-কাঠামো)

### 🌍 অনুবাদসমূহ

| ভাষা | | ভাষা | | ভাষা | |
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

## দ্রুত শুরু

### ১. ডাউনলোড
[রিলিজ](../../releases) থেকে আপনার প্ল্যাটফর্ম বেছে নিন অথবা সোর্স থেকে বিল্ড করুন:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### ২. কনফিগার করুন (GUI বা টেক্সট)
```bash
devspace-gui                  # ডেস্কটপ কনফিগারেটর (GUI)
devspace init                 # টেক্সট-ভিত্তিক কনফিগারেটর
```

### ৩. চালান
```bash
devspace                      # সার্ভার শুরু করে। স্বয়ংক্রিয়ভাবে কনফিগ সনাক্ত করে।
```

এটি `tools/`-এ `cloudflared` পাওয়া গেলে স্বয়ংক্রিয়ভাবে একটি Cloudflare টানেলও শুরু করে।

### ৪. আপনার MCP ক্লায়েন্ট সংযুক্ত করুন
```
https://আপনার-টানেল.trycloudflare.com/mcp
```
অথবা স্থানীয়ভাবে: `http://127.0.0.1:7676/mcp`

---

## ইনস্টলেশন

কোনো Node.js নেই, কোনো npm নেই, কোনো Python নেই। একক বাইনারি।

| প্ল্যাটফর্ম | ডাউনলোড |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: নেটিভভাবে কম্পাইল করুন) |
| **macOS Intel** | `devspace` (GUI: নেটিভভাবে কম্পাইল করুন) |
| **macOS M-chip** | `devspace` (GUI: নেটিভভাবে কম্পাইল করুন) |

শুধুমাত্র সোর্স থেকে বিল্ড করলে **Go 1.23+** প্রয়োজন।

---

## AI কী করতে পারে

একবার সংযুক্ত হলে, AI আপনার অনুমোদিত প্রকল্প ফোল্ডারগুলির একটি ওয়ার্কস্পেস হিসেবে খুলতে পারে:

- ওয়ার্কস্পেসের ভিতরে ফাইল **পড়তে, লিখতে এবং সম্পাদনা করতে**
- regex দিয়ে **কোড অনুসন্ধান** এবং ডিরেক্টরি পরিদর্শন করতে
- **শেল কমান্ড চালাতে** (Windows-এ PowerShell, Unix-এ bash)
- `AGENTS.md` / `CLAUDE.md` থেকে **প্রকল্প নির্দেশাবলী আবিষ্কার** করতে
- পোর্টেবল `.devspace/config.json` দিয়ে **স্বয়ংক্রিয়-কনফিগার** করতে

৮টি MCP টুল: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## কনফিগারেশন

সমস্ত কনফিগ **এক্সিকিউটেবলের মতো একই ফোল্ডারে** থাকে (পোর্টেবল):

```
.devspace/
├── config.json       ← অনুমোদিত রুট, পোর্ট, শেল, ভাষা, অথেনটিকেশন
└── auth.json         ← মালিকের পাসওয়ার্ড (ঐচ্ছিক)
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

| ফিল্ড | ডিফল্ট | বিবরণ |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | OS থেকে স্বয়ংক্রিয়-সনাক্ত। ৪৭টি ভাষা সমর্থন করে |
| `toolMode` | `full` | `full` (সব টুল) অথবা `minimal` (শুধু অনুসন্ধানের জন্য শেল) |
| `toolNaming` | `short` | `short` (read, write) অথবা `legacy` (read_file, write_file) |

কোনো এনভায়রনমেন্ট ভেরিয়েবলের প্রয়োজন নেই — সবকিছু পোর্টেবল কনফিগ ফাইলে রয়েছে।

---

## টানেল (দূরবর্তী অ্যাক্সেস)

ChatGPT ওয়েব সংস্করণের জন্য (HTTPS প্রয়োজন), DevSpace স্বয়ংক্রিয়ভাবে একটি টানেল শুরু করে:

| টানেল | URL ধরণ | সেটআপ |
|---|---|---|
| **Cloudflare** | র্যান্ডম (স্বয়ংক্রিয়) | `tools/`-এ `cloudflared.exe` রাখুন |
| **Pinggy** | স্থিতিশীল | SSH কী প্রয়োজন (`ssh-keygen`) |

সার্ভার স্বয়ংক্রিয়ভাবে সনাক্ত করে কোনটি উপলব্ধ। নতুন Cloudflare URL-এর জন্য সার্ভার পুনরায় শুরু করুন, অথবা স্থায়ী URL-এর জন্য Pinggy ব্যবহার করুন।

---

## শেল সমর্থন

| OS | ডিফল্ট | বিকল্প |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / যেকোনো শেল |
| **macOS** | bash | `sh` / `zsh` |

config.json-এ `"shell"` সেট করুন অথবা GUI-তে বেছে নিন।

---

## নিরাপত্তা

- **PKCE সহ OAuth 2.0** — যদি মালিকের পাসওয়ার্ড সেট করা থাকে
- **পাসওয়ার্ড-বিহীন মোড** — যদি কোনো পাসওয়ার্ড কনফিগার না করা থাকে, অথেনটিকেশন ছাড়াই চলে
- **পাথ কন্টেইনমেন্ট** — সমস্ত ফাইল অপারেশন অনুমোদিত রুটের বিরুদ্ধে যাচাই করা হয়
- **ঐচ্ছিক টানেল** — Cloudflare টানেল সরাসরি প্রকাশ থেকে রক্ষা করে
- **কোনো তৃতীয়-পক্ষ আপলোড নেই** — আপনার কোড কখনো আপনার মেশিন ছেড়ে যায় না

---

## সোর্স থেকে বিল্ড

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# সবকিছু বিল্ড করুন (সব প্ল্যাটফর্ম)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# শুধু বর্তমান প্ল্যাটফর্মের জন্য বিল্ড করুন
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## প্ল্যাটফর্ম সমর্থন

| প্ল্যাটফর্ম | সার্ভার | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (নেটিভভাবে কম্পাইল করুন) |
| **macOS Intel** | ✅ | 🔧 (নেটিভভাবে কম্পাইল করুন) |
| **macOS M-chip** | ✅ | 🔧 (নেটিভভাবে কম্পাইল করুন) |

GUI-এর জন্য Fyne (OpenGL) প্রয়োজন — ক্রস-কম্পাইল করা যায় না। সার্ভার সর্বত্র কম্পাইল হয়।

---

## প্রকল্প কাঠামো

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP সার্ভার
│   └── devspace-gui/       ← ডেস্কটপ GUI কনফিগারেটর (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE প্রদানকারী
│   ├── config/             ← পোর্টেবল কনফিগ সিস্টেম
│   ├── locales/            ← ৪৭টি ভাষার অনুবাদ
│   ├── logger/             ← স্ট্রাকচার্ড লগিং (zerolog)
│   ├── server/             ← HTTP + MCP + টানেল অর্কেস্ট্রেশন
│   ├── skills/             ← AGENTS.md / দক্ষতা আবিষ্কার
│   ├── store/              ← SQLite ওয়ার্কস্পেস সেশন
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← ওয়ার্কস্পেস ও পাথ যাচাইকরণ
├── scripts/
│   ├── windows/            ← PowerShell বিল্ড স্ক্রিপ্ট
│   └── unix/               ← Bash + Makefile বিল্ড স্ক্রিপ্ট
├── readme/                 ← এই ফাইলের অনুবাদ (৪৭টি ভাষা)
├── build/                  ← কম্পাইল করা বাইনারি (সব প্ল্যাটফর্ম)
├── tools/                  ← cloudflared.exe, ইত্যাদি
├── go.mod / go.sum
└── README.md
```

---

Go-তে নির্মিত। শূন্য npm। শূন্য Node.js। একটি বাইনারি।
