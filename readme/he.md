# DevSpace (מהדורת Go)

**תנו ל־ChatGPT ול־Claude גישה מאובטחת למחשב המקומי שלכם. הפכו כל מארח MCP לשותף התכנות שלכם.**

DevSpace הוא שרת MCP באירוח עצמי המאפשר לעוזרי בינה מלאכותית לקרוא, לערוך, לחפש ולהריץ קוד בפרויקטים המקומיים האמיתיים שלכם — הקבצים שלכם, הכלים שלכם, הטרמינל שלכם — מבלי להעלות דבר לצד שלישי. אתם מריצים אותו על המחשב שלכם, חושפים אותו דרך מנהרה שאתם שולטים בה, ובאופן אופציונלי מאבטחים אותו עם סיסמה.

---

## תוכן עניינים

- [התחלה מהירה](#התחלה-מהירה)
- [התקנה](#התקנה)
- [מה הבינה המלאכותית יכולה לעשות](#מה-הבינה-המלאכותית-יכולה-לעשות)
- [תצורה](#תצורה)
- [מנהרה (גישה מרחוק)](#מנהרה-גישה-מרחוק)
- [תמיכה במעטפת פקודה](#תמיכה-במעטפת-פקודה)
- [אבטחה](#אבטחה)
- [בנייה מקוד מקור](#בנייה-מקוד-מקור)
- [תמיכה בפלטפורמות](#תמיכה-בפלטפורמות)
- [מבנה הפרויקט](#מבנה-הפרויקט)

### 🌍 תרגומים

| שפה | | שפה | | שפה | |
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

## התחלה מהירה

### 1. הורדה
בחרו את הפלטפורמה שלכם מעמוד ה-[Releases](../../releases) או בנו מקוד מקור:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. תצורה (GUI או טקסט)
```bash
devspace-gui                  # כלי תצורה שולחני (GUI)
devspace init                 # כלי תצורה מבוסס טקסט
```

### 3. הרצה
```bash
devspace                      # מפעיל את השרת. מזהה תצורה אוטומטית.
```

זה גם מפעיל אוטומטית מנהרת Cloudflare אם `cloudflared` נמצא בתיקיית `tools/`.

### 4. חיבור לקוח ה-MCP שלכם
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
או מקומית: `http://127.0.0.1:7676/mcp`

---

## התקנה

ללא Node.js, ללא npm, ללא Python. קובץ בינארי יחיד.

| פלטפורמה | הורדה |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: הידור מקורי) |
| **macOS Intel** | `devspace` (GUI: הידור מקורי) |
| **macOS M-chip** | `devspace` (GUI: הידור מקורי) |

דורש **Go 1.23+** רק אם בונים מקוד מקור.

---

## מה הבינה המלאכותית יכולה לעשות

לאחר החיבור, הבינה המלאכותית יכולה לפתוח את אחת מתיקיות הפרויקט המאושרות שלכם כמרחב עבודה:

- **לקרוא, לכתוב ולערוך** קבצים בתוך מרחב העבודה
- **לחפש קוד** עם ביטויים רגולריים ולסרוק תיקיות
- **להריץ פקודות מעטפת** (PowerShell ב-Windows, bash ב-Unix)
- **לגלות הוראות פרויקט** מתוך `AGENTS.md` / `CLAUDE.md`
- **תצורה אוטומטית** עם `.devspace/config.json` נייד

8 כלי MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## תצורה

כל התצורה נמצאת **באותה תיקייה של קובץ ההרצה** (נייד):

```
.devspace/
├── config.json       ← שורשים מורשים, פורט, מעטפת, שפה, אימות
└── auth.json         ← סיסמת בעלים (אופציונלי)
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

| שדה | ברירת מחדל | תיאור |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | זיהוי אוטומטי ממערכת ההפעלה. תומך ב-47 שפות |
| `toolMode` | `full` | `full` (כל הכלים) או `minimal` (מעטפת בלבד לחיפוש) |
| `toolNaming` | `short` | `short` (read, write) או `legacy` (read_file, write_file) |

אין צורך במשתני סביבה — הכל נמצא בקובץ התצורה הנייד.

---

## מנהרה (גישה מרחוק)

עבור גרסת האינטרנט של ChatGPT (נדרש HTTPS), DevSpace מפעיל אוטומטית מנהרה:

| מנהרה | סוג כתובת | הגדרה |
|---|---|---|
| **Cloudflare** | אקראית (אוטומטית) | שימו את `cloudflared.exe` בתיקיית `tools/` |
| **Pinggy** | קבועה | דורש מפתח SSH (`ssh-keygen`) |

השרת מזהה אוטומטית איזו מהן זמינה. הפעילו מחדש את השרת לכתובת Cloudflare חדשה, או השתמשו ב-Pinggy לכתובת קבועה.

---

## תמיכה במעטפת פקודה

| מערכת הפעלה | ברירת מחדל | חלופות |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / כל מעטפת |
| **macOS** | bash | `sh` / `zsh` |

הגדירו `"shell"` ב-config.json או בחרו ב-GUI.

---

## אבטחה

- **OAuth 2.0 עם PKCE** — אם הוגדרה סיסמת בעלים
- **מצב ללא סיסמה** — אם לא הוגדרה סיסמה, פועל ללא אימות
- **הגבלת נתיבים** — כל פעולות הקבצים מאומתות מול השורשים המורשים
- **מנהרה אופציונלית** — מנהרת Cloudflare מגנה מפני חשיפה ישירה
- **ללא העלאות לצד שלישי** — הקוד שלכם לעולם לא עוזב את המחשב שלכם

---

## בנייה מקוד מקור

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# בניית הכל (כל הפלטפורמות)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# בנייה רק לפלטפורמה הנוכחית
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## תמיכה בפלטפורמות

| פלטפורמה | שרת | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (הידור מקורי) |
| **macOS Intel** | ✅ | 🔧 (הידור מקורי) |
| **macOS M-chip** | ✅ | 🔧 (הידור מקורי) |

ה-GUI דורש Fyne (OpenGL) — לא ניתן להדר צולב. השרת מתהדר בכל מקום.

---

## מבנה הפרויקט

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + שרת MCP
│   └── devspace-gui/       ← כלי תצורה שולחני GUI (Fyne)
├── internal/
│   ├── auth/               ← ספק OAuth 2.0 + PKCE
│   ├── config/             ← מערכת תצורה ניידת
│   ├── locales/            ← תרגומים ל-47 שפות
│   ├── logger/             ← תיעוד מובנה (zerolog)
│   ├── server/             ← HTTP + MCP + תזמור מנהרה
│   ├── skills/             ← AGENTS.md / גילוי מיומנויות
│   ├── store/              ← SQLite — סשנים של מרחבי עבודה
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← מרחב עבודה ואימות נתיבים
├── scripts/
│   ├── windows/            ← סקריפט בנייה PowerShell
│   └── unix/               ← סקריפטים Bash + Makefile
├── readme/                 ← תרגומים של קובץ זה (47 שפות)
├── build/                  ← קבצים בינאריים מהודרים (כל הפלטפורמות)
├── tools/                  ← cloudflared.exe, וכו'
├── go.mod / go.sum
└── README.md
```

---

נבנה ב-Go. אפס npm. אפס Node.js. קובץ בינארי אחד.
