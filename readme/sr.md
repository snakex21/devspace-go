# DevSpace (Go издање)

**Омогућите ChatGPT и Claude сигуран приступ вашем локалном рачунару. Претворите било ког MCP домаћина у свог програмског партнера.**

DevSpace је само-хостовани MCP сервер који омогућава AI асистентима да читају, уређују, претражују и покрећу код у вашим стварним локалним пројектима — ваше датотеке, ваши алати, ваш терминал — без отпремања ичега трећим странама. Покрећете га на свом рачунару, излажете га кроз тунел који контролишете и опционо га обезбеђујете лозинком.

---

## Садржај

- [Брзи Почетак](#брзи-почетак)
- [Инсталација](#инсталација)
- [Шта AI Може](#шта-ai-може)
- [Конфигурација](#конфигурација)
- [Тунел (Удаљени Приступ)](#тунел-удаљени-приступ)
- [Подршка Шкољке](#подршка-шкољке)
- [Безбедност](#безбедност)
- [Изградња из Извора](#изградња-из-извора)
- [Подршка Платформи](#подршка-платформи)
- [Структура Пројекта](#структура-пројекта)

### 🌍 Преводи

| Језик | | Језик | | Језик | |
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

## Брзи Почетак

### 1. Преузимање
Изаберите своју платформу са [Releases](../../releases) или изградите из извора:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Конфигурација (GUI или текст)
```bash
devspace-gui                  # Десктоп конфигуратор (GUI)
devspace init                 # Текстуални конфигуратор
```

### 3. Покретање
```bash
devspace                      # Покреће сервер. Аутоматски детектује конфиг.
```

Ово такође аутоматски покреће Cloudflare тунел ако је `cloudflared` пронађен у `tools/`.

### 4. Повежите свог MCP клијента
```
https://ВАШ-ТУНЕЛ.trycloudflare.com/mcp
```
Или локално: `http://127.0.0.1:7676/mcp`

---

## Инсталација

Без Node.js, без npm, без Python-а. Један бинарни фајл.

| Платформа | Преузимање |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: компајлирајте нативно) |
| **macOS Intel** | `devspace` (GUI: компајлирајте нативно) |
| **macOS M-chip** | `devspace` (GUI: компајлирајте нативно) |

Захтева **Go 1.23+** само при изградњи из извора.

---

## Шта AI Може

Када се повеже, AI може отворити једну од ваших одобрених фасцикли пројекта као радни простор:

- **Чита, пише и уређује** датотеке унутар радног простора
- **Претражује код** регуларним изразима и прегледа директоријуме
- **Покреће команде шкољке** (PowerShell на Windows-у, bash на Unix-у)
- **Открива упутства пројекта** из `AGENTS.md` / `CLAUDE.md`
- **Аутоматски конфигурише** са преносивим `.devspace/config.json`

8 MCP алата: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Конфигурација

Сва конфигурација се налази **у истој фасцикли као и извршна датотека** (преносива):

```
.devspace/
├── config.json       ← дозвољени корени, порт, шкољка, језик, аутентификација
└── auth.json         ← лозинка власника (опционо)
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

| Поље | Подразумевано | Опис |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Аутоматско откривање из ОС-а. Подржава 47 језика |
| `toolMode` | `full` | `full` (сви алати) или `minimal` (само шкољка за претрагу) |
| `toolNaming` | `short` | `short` (read, write) или `legacy` (read_file, write_file) |

Нису потребне променљиве окружења — све је у преносивој конфигурационој датотеци.

---

## Тунел (Удаљени Приступ)

За веб верзију ChatGPT-а (потребан HTTPS), DevSpace аутоматски покреће тунел:

| Тунел | Тип URL-а | Подешавање |
|---|---|---|
| **Cloudflare** | Насумичан (ауто) | Ставите `cloudflared.exe` у `tools/` |
| **Pinggy** | Стабилан | Потребан SSH кључ (`ssh-keygen`) |

Сервер аутоматски открива који је доступан. Поново покрените сервер за нови Cloudflare URL или користите Pinggy за трајни URL.

---

## Подршка Шкољке

| ОС | Подразумевано | Алтернативе |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / било која шкољка |
| **macOS** | bash | `sh` / `zsh` |

Подесите `"shell"` у config.json или изаберите у GUI-ју.

---

## Безбедност

- **OAuth 2.0 са PKCE** — ако је постављена лозинка власника
- **Режим без лозинке** — ако лозинка није конфигурисана, ради без аутентификације
- **Ограничење путања** — све датотечне операције су потврђене у односу на дозвољене корене
- **Опциони тунел** — Cloudflare тунел штити од директног излагања
- **Без отпремања трећим странама** — ваш код никада не напушта ваш рачунар

---

## Изградња из Извора

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Изгради све (све платформе)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Изгради само за тренутну платформу
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Подршка Платформи

| Платформа | Сервер | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (компајлирајте нативно) |
| **macOS Intel** | ✅ | 🔧 (компајлирајте нативно) |
| **macOS M-chip** | ✅ | 🔧 (компајлирајте нативно) |

GUI захтева Fyne (OpenGL) — унакрсно компајлирање није могуће. Сервер се компајлира свуда.

---

## Структура Пројекта

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP сервер
│   └── devspace-gui/       ← Десктоп GUI конфигуратор (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE провајдер
│   ├── config/             ← Систем преносиве конфигурације
│   ├── locales/            ← Преводи на 47 језика
│   ├── logger/             ← Структурирано логовање (zerolog)
│   ├── server/             ← HTTP + MCP + оркестрација тунела
│   ├── skills/             ← Откривање AGENTS.md / вештина
│   ├── store/              ← SQLite сесије радног простора
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Валидација радног простора и путања
├── scripts/
│   ├── windows/            ← PowerShell скрипта за изградњу
│   └── unix/               ← Bash + Makefile скрипте за изградњу
├── readme/                 ← Преводи ове датотеке (47 језика)
├── build/                  ← Компајлирани бинарни фајлови (све платформе)
├── tools/                  ← cloudflared.exe, итд.
├── go.mod / go.sum
└── README.md
```

---

Изграђено у Go-у. Нула npm. Нула Node.js. Један бинарни фајл.
