# DevSpace (Go Edition) - فارسی

**به ChatGPT و Claude دسترسی امن به ماشین محلی خود بدهید. هر میزبان MCP را به شریک کدنویسی خود تبدیل کنید.**

DevSpace یک سرور MCP خودمیزبان است که به دستیارهای هوش مصنوعی اجازه می‌دهد کد را در پروژه‌های محلی واقعی شما بخوانند، ویرایش کنند، جستجو کنند و اجرا کنند — فایل‌های شما، ابزارهای شما، ترمینال شما — بدون آپلود هیچ چیزی به شخص ثالث. شما آن را روی ماشین خود اجرا می‌کنید، از طریق تونلی که کنترل می‌کنید در معرض نمایش قرار می‌دهید و به صورت اختیاری با یک رمز عبور ایمن می‌کنید.

---

## فهرست مطالب

- [شروع سریع](#شروع-سریع)
- [نصب](#نصب)
- [هوش مصنوعی چه کارهایی می‌تواند انجام دهد](#هوش-مصنوعی-چه-کارهایی-می‌تواند-انجام-دهد)
- [پیکربندی](#پیکربندی)
- [تونل (دسترسی از راه دور)](#تونل-دسترسی-از-راه-دور)
- [پشتیبانی از شل](#پشتیبانی-از-شل)
- [امنیت](#امنیت)
- [ساخت از سورس](#ساخت-از-سورس)
- [پشتیبانی از پلتفرم‌ها](#پشتیبانی-از-پلتفرم‌ها)
- [ساختار پروژه](#ساختار-پروژه)

### 🌍 ترجمه‌ها

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

## شروع سریع

### ۱. دانلود
پلتفرم خود را از [انتشارات](../../releases) انتخاب کنید یا از سورس بسازید:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### ۲. پیکربندی (رابط کاربری گرافیکی یا متنی)
```bash
devspace-gui                  # پیکربند دسکتاپ (رابط گرافیکی)
devspace init                 # پیکربند مبتنی بر متن
```

### ۳. اجرا
```bash
devspace                      # سرور را شروع می‌کند. پیکربندی را خودکار تشخیص می‌دهد.
```

اگر `cloudflared` در `tools/` یافت شود، این کار همچنین یک تونل Cloudflare را به طور خودکار شروع می‌کند.

### ۴. کلاینت MCP خود را متصل کنید
```
https://تونل-شما.trycloudflare.com/mcp
```
یا به صورت محلی: `http://127.0.0.1:7676/mcp`

---

## نصب

بدون Node.js، بدون npm، بدون Python. یک فایل باینری واحد.

| پلتفرم | دانلود |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (رابط گرافیکی: به صورت بومی کامپایل کنید) |
| **macOS Intel** | `devspace` (رابط گرافیکی: به صورت بومی کامپایل کنید) |
| **macOS M-chip** | `devspace` (رابط گرافیکی: به صورت بومی کامپایل کنید) |

فقط در صورت ساخت از سورس به **Go 1.23+** نیاز دارد.

---

## هوش مصنوعی چه کارهایی می‌تواند انجام دهد

پس از اتصال، هوش مصنوعی می‌تواند یکی از پوشه‌های پروژه تأیید شده شما را به عنوان فضای کاری باز کند:

- **خواندن، نوشتن و ویرایش** فایل‌ها در داخل فضای کاری
- **جستجوی کد** با regex و بازرسی دایرکتوری‌ها
- **اجرای دستورات شل** (PowerShell در ویندوز، bash در یونیکس)
- **کشف دستورالعمل‌های پروژه** از `AGENTS.md` / `CLAUDE.md`
- **پیکربندی خودکار** با `.devspace/config.json` قابل حمل

۸ ابزار MCP: `open_workspace`، `read`، `write`، `edit`، `grep`، `glob`، `ls`، `bash`

---

## پیکربندی

تمام پیکربندی **در همان پوشه فایل اجرایی** قرار دارد (قابل حمل):

```
.devspace/
├── config.json       ← ریشه‌های مجاز، پورت، شل، زبان، احراز هویت
└── auth.json         ← رمز عبور مالک (اختیاری)
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

| فیلد | پیش‌فرض | توضیحات |
|---|---|---|
| `shell` | `auto` | `auto`، `powershell`، `cmd`، `bash`، `sh` |
| `lang` | `auto` | تشخیص خودکار از سیستم عامل. پشتیبانی از ۴۷ زبان |
| `toolMode` | `full` | `full` (همه ابزارها) یا `minimal` (فقط شل برای جستجو) |
| `toolNaming` | `short` | `short` (read, write) یا `legacy` (read_file, write_file) |

نیازی به متغیرهای محیطی نیست — همه چیز در فایل پیکربندی قابل حمل است.

---

## تونل (دسترسی از راه دور)

برای نسخه وب ChatGPT (نیاز به HTTPS)، DevSpace به طور خودکار یک تونل را شروع می‌کند:

| تونل | نوع URL | راه‌اندازی |
|---|---|---|
| **Cloudflare** | تصادفی (خودکار) | `cloudflared.exe` را در `tools/` قرار دهید |
| **Pinggy** | پایدار | نیاز به کلید SSH دارد (`ssh-keygen`) |

سرور به طور خودکار تشخیص می‌دهد کدام یک در دسترس است. برای URL جدید Cloudflare سرور را مجدداً راه‌اندازی کنید یا از Pinggy برای URL دائمی استفاده کنید.

---

## پشتیبانی از شل

| سیستم عامل | پیش‌فرض | جایگزین‌ها |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / هر شل |
| **macOS** | bash | `sh` / `zsh` |

`"shell"` را در config.json تنظیم کنید یا در رابط گرافیکی انتخاب کنید.

---

## امنیت

- **OAuth 2.0 با PKCE** — اگر رمز عبور مالک تنظیم شده باشد
- **حالت بدون رمز عبور** — اگر رمز عبوری پیکربندی نشده باشد، بدون احراز هویت اجرا می‌شود
- **محدودیت مسیر** — تمام عملیات فایل در برابر ریشه‌های مجاز اعتبارسنجی می‌شوند
- **تونل اختیاری** — تونل Cloudflare از افشای مستقیم محافظت می‌کند
- **بدون آپلود به شخص ثالث** — کد شما هرگز ماشین شما را ترک نمی‌کند

---

## ساخت از سورس

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# ساخت همه چیز (همه پلتفرم‌ها)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# ساخت فقط برای پلتفرم فعلی
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## پشتیبانی از پلتفرم‌ها

| پلتفرم | سرور | رابط گرافیکی |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (به صورت بومی کامپایل کنید) |
| **macOS Intel** | ✅ | 🔧 (به صورت بومی کامپایل کنید) |
| **macOS M-chip** | ✅ | 🔧 (به صورت بومی کامپایل کنید) |

رابط گرافیکی نیاز به Fyne (OpenGL) دارد — نمی‌توان به صورت متقابل کامپایل کرد. سرور در همه جا کامپایل می‌شود.

---

## ساختار پروژه

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + سرور MCP
│   └── devspace-gui/       ← پیکربند دسکتاپ (Fyne)
├── internal/
│   ├── auth/               ← ارائه‌دهنده OAuth 2.0 + PKCE
│   ├── config/             ← سیستم پیکربندی قابل حمل
│   ├── locales/            ← ترجمه‌های ۴۷ زبان
│   ├── logger/             ← ثبت ساختاریافته (zerolog)
│   ├── server/             ← HTTP + MCP + هماهنگی تونل
│   ├── skills/             ← AGENTS.md / کشف مهارت
│   ├── store/              ← نشست‌های فضای کاری SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← فضای کاری و اعتبارسنجی مسیر
├── scripts/
│   ├── windows/            ← اسکریپت ساخت PowerShell
│   └── unix/               ← اسکریپت‌های ساخت Bash + Makefile
├── readme/                 ← ترجمه‌های این فایل (۴۷ زبان)
├── build/                  ← باینری‌های کامپایل شده (همه پلتفرم‌ها)
├── tools/                  ← cloudflared.exe و غیره
├── go.mod / go.sum
└── README.md
```

---

ساخته شده با Go. صفر npm. صفر Node.js. یک باینری.
