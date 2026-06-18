# DevSpace (Go Видання)

**Надайте ChatGPT і Claude безпечний доступ до вашого локального комп'ютера. Перетворіть будь-який MCP-хост на вашого партнера з програмування.**

DevSpace — це самостійно розміщений MCP-сервер, який дозволяє AI-асистентам читати, редагувати, шукати та запускати код у ваших реальних локальних проєктах — ваші файли, ваші інструменти, ваш термінал — без завантаження чого-небудь третім сторонам. Ви запускаєте його на своєму комп'ютері, надаєте доступ через контрольований вами тунель і за бажанням захищаєте паролем.

---

## Зміст

- [Швидкий Старт](#швидкий-старт)
- [Встановлення](#встановлення)
- [Що Може Робити AI](#що-може-робити-ai)
- [Конфігурація](#конфігурація)
- [Тунель (Віддалений Доступ)](#тунель-віддалений-доступ)
- [Підтримка Оболонок](#підтримка-оболонок)
- [Безпека](#безпека)
- [Збірка з Вихідного Коду](#збірка-з-вихідного-коду)
- [Підтримка Платформ](#підтримка-платформ)
- [Структура Проєкту](#структура-проєкту)

### 🌍 Переклади

| Мова | | Мова | | Мова | |
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

## Швидкий Старт

### 1. Завантажити
Виберіть свою платформу в [Releases](../../releases) або зберіть з вихідного коду:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Налаштувати (GUI або текст)
```bash
devspace-gui                  # Настільний конфігуратор (GUI)
devspace init                 # Текстовий конфігуратор
```

### 3. Запустити
```bash
devspace                      # Запускає сервер. Автовизначення конфігурації.
```

Це також автоматично запускає тунель Cloudflare, якщо `cloudflared` знайдено в `tools/`.

### 4. Підключіть ваш MCP-клієнт
```
https://ВАШ-ТУНЕЛЬ.trycloudflare.com/mcp
```
Або локально: `http://127.0.0.1:7676/mcp`

---

## Встановлення

Без Node.js, без npm, без Python. Один бінарний файл.

| Платформа | Завантаження |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: скомпілюйте нативно) |
| **macOS Intel** | `devspace` (GUI: скомпілюйте нативно) |
| **macOS M-chip** | `devspace` (GUI: скомпілюйте нативно) |

Потребує **Go 1.23+** лише при збірці з вихідного коду.

---

## Що Може Робити AI

Після підключення AI може відкрити одну з ваших схвалених папок проєкту як робочий простір:

- **Читати, записувати та редагувати** файли в межах робочого простору
- **Шукати код** за допомогою регулярних виразів і переглядати директорії
- **Запускати команди оболонки** (PowerShell на Windows, bash на Unix)
- **Виявляти інструкції проєкту** з `AGENTS.md` / `CLAUDE.md`
- **Автоматично налаштовувати** за допомогою портативного `.devspace/config.json`

8 інструментів MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Конфігурація

Вся конфігурація знаходиться **в тій же папці, що й виконуваний файл** (портативна):

```
.devspace/
├── config.json       ← дозволені корені, порт, оболонка, мова, автентифікація
└── auth.json         ← пароль власника (необов'язково)
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

| Поле | За замовчуванням | Опис |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Автовизначення з ОС. Підтримує 47 мов |
| `toolMode` | `full` | `full` (всі інструменти) або `minimal` (лише оболонка для пошуку) |
| `toolNaming` | `short` | `short` (read, write) або `legacy` (read_file, write_file) |

Змінні середовища не потрібні — все в портативному файлі конфігурації.

---

## Тунель (Віддалений Доступ)

Для веб-версії ChatGPT (потрібен HTTPS) DevSpace автоматично запускає тунель:

| Тунель | Тип URL | Налаштування |
|---|---|---|
| **Cloudflare** | Випадковий (авто) | Помістіть `cloudflared.exe` у `tools/` |
| **Pinggy** | Стабільний | Потрібен SSH-ключ (`ssh-keygen`) |

Сервер автоматично визначає, який із них доступний. Перезапустіть сервер для нового URL Cloudflare або використовуйте Pinggy для постійного URL.

---

## Підтримка Оболонок

| ОС | За замовчуванням | Альтернативи |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / будь-яка оболонка |
| **macOS** | bash | `sh` / `zsh` |

Встановіть `"shell"` у config.json або виберіть у GUI.

---

## Безпека

- **OAuth 2.0 з PKCE** — якщо встановлено пароль власника
- **Режим без пароля** — якщо пароль не налаштовано, працює без автентифікації
- **Обмеження шляхів** — всі файлові операції перевіряються на відповідність дозволеним кореням
- **Необов'язковий тунель** — тунель Cloudflare захищає від прямого доступу
- **Без завантажень третім сторонам** — ваш код ніколи не залишає ваш комп'ютер

---

## Збірка з Вихідного Коду

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Зібрати все (всі платформи)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Зібрати лише для поточної платформи
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Підтримка Платформ

| Платформа | Сервер | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (скомпілюйте нативно) |
| **macOS Intel** | ✅ | 🔧 (скомпілюйте нативно) |
| **macOS M-chip** | ✅ | 🔧 (скомпілюйте нативно) |

GUI потребує Fyne (OpenGL) — крос-компіляція неможлива. Сервер компілюється всюди.

---

## Структура Проєкту

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-сервер
│   └── devspace-gui/       ← Настільний GUI-конфігуратор (Fyne)
├── internal/
│   ├── auth/               ← Провайдер OAuth 2.0 + PKCE
│   ├── config/             ← Система портативної конфігурації
│   ├── locales/            ← Переклади на 47 мов
│   ├── logger/             ← Структуроване журналювання (zerolog)
│   ├── server/             ← HTTP + MCP + оркестрація тунелю
│   ├── skills/             ← Виявлення AGENTS.md / навичок
│   ├── store/              ← Сесії робочого простору SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Перевірка робочого простору та шляхів
├── scripts/
│   ├── windows/            ← PowerShell скрипт збірки
│   └── unix/               ← Bash + Makefile скрипти збірки
├── readme/                 ← Переклади цього файлу (47 мов)
├── build/                  ← Скомпільовані бінарні файли (всі платформи)
├── tools/                  ← cloudflared.exe тощо
├── go.mod / go.sum
└── README.md
```

---

Створено на Go. Нуль npm. Нуль Node.js. Один бінарний файл.
