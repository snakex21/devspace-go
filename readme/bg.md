# DevSpace (Go Edition) - Български

**Дайте на ChatGPT & Claude сигурен достъп до вашата локална машина. Превърнете всеки MCP хост във ваш партньор за програмиране.**

DevSpace е самостоятелно хостван MCP сървър, който позволява на AI асистенти да четат, редактират, търсят и изпълняват код във вашите реални локални проекти — вашите файлове, вашите инструменти, вашият терминал — без да качват нищо на трета страна. Стартирате го на вашата машина, излагате го чрез тунел, който контролирате, и по желание го защитавате с парола.

---

## Съдържание

- [Бърз старт](#бърз-старт)
- [Инсталация](#инсталация)
- [Какво може AI](#какво-може-ai)
- [Конфигурация](#конфигурация)
- [Тунел (Отдалечен достъп)](#тунел-отдалечен-достъп)
- [Поддръжка на обвивка](#поддръжка-на-обвивка)
- [Сигурност](#сигурност)
- [Изграждане от изходен код](#изграждане-от-изходен-код)
- [Поддръжка на платформи](#поддръжка-на-платформи)
- [Структура на проекта](#структура-на-проекта)

### 🌍 Преводи

| Език | | Език | | Език | |
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

## Бърз старт

### 1. Изтегляне
Изберете вашата платформа от [Издания](../../releases) или изградете от изходен код:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Конфигуриране (GUI или текст)
```bash
devspace-gui                  # Десктоп конфигуратор (GUI)
devspace init                 # Текстов конфигуратор
```

### 3. Стартиране
```bash
devspace                      # Стартира сървъра. Автоматично открива конфигурацията.
```

Това също автоматично стартира Cloudflare тунел, ако `cloudflared` бъде намерен в `tools/`.

### 4. Свържете вашия MCP клиент
```
https://ВАШИЯ-ТУНЕЛ.trycloudflare.com/mcp
```
Или локално: `http://127.0.0.1:7676/mcp`

---

## Инсталация

Без Node.js, без npm, без Python. Единствен бинарен файл.

| Платформа | Изтегляне |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: компилирайте нативно) |
| **macOS Intel** | `devspace` (GUI: компилирайте нативно) |
| **macOS M-чип** | `devspace` (GUI: компилирайте нативно) |

Изисква **Go 1.23+** само при изграждане от изходен код.

---

## Какво може AI

Веднъж свързан, AI може да отвори една от вашите одобрени проектни папки като работно пространство:

- **Четене, писане и редактиране** на файлове в работното пространство
- **Търсене на код** с regex и инспектиране на директории
- **Изпълнение на команди в обвивка** (PowerShell на Windows, bash на Unix)
- **Откриване на инструкции за проекта** от `AGENTS.md` / `CLAUDE.md`
- **Автоматично конфигуриране** с преносимо `.devspace/config.json`

8 MCP инструмента: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Конфигурация

Цялата конфигурация се намира **в същата папка като изпълнимия файл** (преносима):

```
.devspace/
├── config.json       ← позволени корени, порт, обвивка, език, удостоверяване
└── auth.json         ← парола на собственика (по желание)
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

| Поле | По подразбиране | Описание |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Автоматично откриване от ОС. Поддържа 47 езика |
| `toolMode` | `full` | `full` (всички инструменти) или `minimal` (само обвивка за търсене) |
| `toolNaming` | `short` | `short` (read, write) или `legacy` (read_file, write_file) |

Не са необходими променливи на средата — всичко е в преносимия конфигурационен файл.

---

## Тунел (Отдалечен достъп)

За уеб версията на ChatGPT (изисква HTTPS), DevSpace автоматично стартира тунел:

| Тунел | Тип URL | Настройка |
|---|---|---|
| **Cloudflare** | Случаен (авто) | Поставете `cloudflared.exe` в `tools/` |
| **Pinggy** | Стабилен | Нуждае се от SSH ключ (`ssh-keygen`) |

Сървърът автоматично открива кой е наличен. Рестартирайте сървъра за нов Cloudflare URL или използвайте Pinggy за постоянен URL.

---

## Поддръжка на обвивка

| ОС | По подразбиране | Алтернативи |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / всяка обвивка |
| **macOS** | bash | `sh` / `zsh` |

Задайте `"shell"` в config.json или изберете в GUI.

---

## Сигурност

- **OAuth 2.0 с PKCE** — ако е зададена парола на собственика
- **Режим без парола** — ако не е конфигурирана парола, работи без удостоверяване
- **Ограничаване на пътя** — всички файлови операции се валидират спрямо позволените корени
- **Опционален тунел** — Cloudflare тунел предпазва от директно излагане
- **Без качвания към трети страни** — вашият код никога не напуска вашата машина

---

## Изграждане от изходен код

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Изграждане на всичко (всички платформи)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Изграждане само за текущата платформа
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Поддръжка на платформи

| Платформа | Сървър | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (компилирайте нативно) |
| **macOS Intel** | ✅ | 🔧 (компилирайте нативно) |
| **macOS M-чип** | ✅ | 🔧 (компилирайте нативно) |

GUI изисква Fyne (OpenGL) — не може кръстосано компилиране. Сървърът се компилира навсякъде.

---

## Структура на проекта

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP сървър
│   └── devspace-gui/       ← Десктоп GUI конфигуратор (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE доставчик
│   ├── config/             ← Преносима конфигурационна система
│   ├── locales/            ← Преводи на 47 езика
│   ├── logger/             ← Структурирано логване (zerolog)
│   ├── server/             ← HTTP + MCP + тунелна оркестрация
│   ├── skills/             ← AGENTS.md / откриване на умения
│   ├── store/              ← SQLite сесии на работно пространство
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Работно пространство и валидация на път
├── scripts/
│   ├── windows/            ← PowerShell скрипт за изграждане
│   └── unix/               ← Bash + Makefile скриптове за изграждане
├── readme/                 ← Преводи на този файл (47 езика)
├── build/                  ← Компилирани бинарни файлове (всички платформи)
├── tools/                  ← cloudflared.exe, и др.
├── go.mod / go.sum
└── README.md
```

---

Изграден на Go. Нула npm. Нула Node.js. Един бинарен файл.
