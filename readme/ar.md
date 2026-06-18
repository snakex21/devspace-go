# DevSpace (Go Edition) - العربية

**امنح ChatGPT و Claude وصولاً آمناً إلى جهازك المحلي. حوّل أي مضيف MCP إلى شريك برمجة خاص بك.**

DevSpace هو خادم MCP ذاتي الاستضافة يتيح للمساعدين الذكاء الاصطناعي قراءة وتحرير والبحث وتشغيل الكود في مشاريعك المحلية الحقيقية — ملفاتك، أدواتك، الطرفية الخاصة بك — دون تحميل أي شيء إلى طرف ثالث. تقوم بتشغيله على جهازك، وتعرضه عبر نفق تتحكم به، وتؤمنه اختيارياً بكلمة مرور.

---

## جدول المحتويات

- [البدء السريع](#البدء-السريع)
- [التثبيت](#التثبيت)
- [ما يمكن للذكاء الاصطناعي فعله](#ما-يمكن-للذكاء-الاصطناعي-فعله)
- [الإعدادات](#الإعدادات)
- [النفق (الوصول عن بعد)](#النفق-الوصول-عن-بعد)
- [دعم الصدفة](#دعم-الصدفة)
- [الأمان](#الأمان)
- [البناء من المصدر](#البناء-من-المصدر)
- [دعم المنصات](#دعم-المنصات)
- [هيكل المشروع](#هيكل-المشروع)

### 🌍 الترجمات

| اللغة | | اللغة | | اللغة | |
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

## البدء السريع

### ١. التحميل
اختر منصتك من [الإصدارات](../../releases) أو قم بالبناء من المصدر:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### ٢. الإعداد (واجهة رسومية أو نصية)
```bash
devspace-gui                  # مكوّن سطح المكتب (واجهة رسومية)
devspace init                 # مكوّن نصي
```

### ٣. التشغيل
```bash
devspace                      # يبدأ الخادم. يكتشف الإعدادات تلقائياً.
```

يبدأ هذا أيضاً نفق Cloudflare تلقائياً إذا تم العثور على `cloudflared` في `tools/`.

### ٤. توصيل عميل MCP الخاص بك
```
https://نفقك.trycloudflare.com/mcp
```
أو محلياً: `http://127.0.0.1:7676/mcp`

---

## التثبيت

لا Node.js، لا npm، لا Python. ملف ثنائي واحد.

| المنصة | التحميل |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: تجميع محلي) |
| **macOS Intel** | `devspace` (GUI: تجميع محلي) |
| **macOS M-chip** | `devspace` (GUI: تجميع محلي) |

يتطلب **Go 1.23+** فقط عند البناء من المصدر.

---

## ما يمكن للذكاء الاصطناعي فعله

بمجرد الاتصال، يمكن للذكاء الاصطناعي فتح أحد مجلدات المشاريع المعتمدة كمساحة عمل:

- **قراءة وكتابة وتحرير** الملفات داخل مساحة العمل
- **البحث في الكود** باستخدام regex وفحص المجلدات
- **تشغيل أوامر الصدفة** (PowerShell على Windows، bash على Unix)
- **اكتشاف تعليمات المشروع** من `AGENTS.md` / `CLAUDE.md`
- **التكوين التلقائي** مع `.devspace/config.json` المحمول

٨ أدوات MCP: `open_workspace`، `read`، `write`، `edit`، `grep`، `glob`، `ls`، `bash`

---

## الإعدادات

جميع الإعدادات تعيش **في نفس مجلد الملف التنفيذي** (محمولة):

```
.devspace/
├── config.json       ← الجذور المسموحة، المنفذ، الصدفة، اللغة، المصادقة
└── auth.json         ← كلمة مرور المالك (اختياري)
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

| الحقل | الافتراضي | الوصف |
|---|---|---|
| `shell` | `auto` | `auto`، `powershell`، `cmd`، `bash`، `sh` |
| `lang` | `auto` | اكتشاف تلقائي من نظام التشغيل. يدعم ٤٧ لغة |
| `toolMode` | `full` | `full` (جميع الأدوات) أو `minimal` (الصدفة فقط للبحث) |
| `toolNaming` | `short` | `short` (read, write) أو `legacy` (read_file, write_file) |

لا حاجة لمتغيرات البيئة — كل شيء في ملف الإعدادات المحمول.

---

## النفق (الوصول عن بعد)

لإصدار ChatGPT على الويب (يتطلب HTTPS)، يبدأ DevSpace نفقاً تلقائياً:

| النفق | نوع الرابط | الإعداد |
|---|---|---|
| **Cloudflare** | عشوائي (تلقائي) | ضع `cloudflared.exe` في `tools/` |
| **Pinggy** | مستقر | يحتاج مفتاح SSH (`ssh-keygen`) |

يكتشف الخادم تلقائياً أيهما متاح. أعد تشغيل الخادم للحصول على رابط Cloudflare جديد، أو استخدم Pinggy للحصول على رابط دائم.

---

## دعم الصدفة

| نظام التشغيل | الافتراضي | البدائل |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / أي صدفة |
| **macOS** | bash | `sh` / `zsh` |

اضبط `"shell"` في config.json أو اختر في الواجهة الرسومية.

---

## الأمان

- **OAuth 2.0 مع PKCE** — إذا تم تعيين كلمة مرور المالك
- **وضع بدون كلمة مرور** — إذا لم يتم تكوين كلمة مرور، يعمل بدون مصادقة
- **احتواء المسار** — جميع عمليات الملفات تتحقق ضد الجذور المسموحة
- **نفق اختياري** — نفق Cloudflare يحمي من التعرض المباشر
- **لا تحميلات لطرف ثالث** — كودك لا يغادر جهازك أبداً

---

## البناء من المصدر

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# بناء الكل (جميع المنصات)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# بناء للمنصة الحالية فقط
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## دعم المنصات

| المنصة | الخادم | الواجهة الرسومية |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (تجميع محلي) |
| **macOS Intel** | ✅ | 🔧 (تجميع محلي) |
| **macOS M-chip** | ✅ | 🔧 (تجميع محلي) |

الواجهة الرسومية تتطلب Fyne (OpenGL) — لا يمكن التجميع المتقاطع. الخادم يتجمع في كل مكان.

---

## هيكل المشروع

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + خادم MCP
│   └── devspace-gui/       ← مكوّن سطح المكتب (Fyne)
├── internal/
│   ├── auth/               ← مزود OAuth 2.0 + PKCE
│   ├── config/             ← نظام إعدادات محمول
│   ├── locales/            ← ترجمات ٤٧ لغة
│   ├── logger/             ← تسجيل منظم (zerolog)
│   ├── server/             ← HTTP + MCP + تنسيق النفق
│   ├── skills/             ← AGENTS.md / اكتشاف المهارات
│   ├── store/              ← جلسات عمل SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← مساحة العمل والتحقق من المسار
├── scripts/
│   ├── windows/            ← سكربت بناء PowerShell
│   └── unix/               ← سكربتات بناء Bash + Makefile
├── readme/                 ← ترجمات هذا الملف (٤٧ لغة)
├── build/                  ← ملفات ثنائية مجمعة (جميع المنصات)
├── tools/                  ← cloudflared.exe، إلخ.
├── go.mod / go.sum
└── README.md
```

---

مبني بلغة Go. صفر npm. صفر Node.js. ملف ثنائي واحد.
