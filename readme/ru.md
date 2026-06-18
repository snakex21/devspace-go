# DevSpace (Go-версия)

**Предоставьте ChatGPT и Claude безопасный доступ к вашему локальному компьютеру. Превратите любой MCP-хост в вашего партнёра по программированию.**

DevSpace — это самостоятельно размещаемый MCP-сервер, который позволяет ИИ-ассистентам читать, редактировать, искать и запускать код в ваших реальных локальных проектах — ваши файлы, ваши инструменты, ваш терминал — без загрузки чего-либо третьим лицам. Вы запускаете его на своём компьютере, предоставляете доступ через контролируемый вами туннель и при необходимости защищаете паролем.

---

## Содержание

- [Быстрый Старт](#быстрый-старт)
- [Установка](#установка)
- [Что Может Делать ИИ](#что-может-делать-ии)
- [Конфигурация](#конфигурация)
- [Туннель (Удалённый Доступ)](#туннель-удалённый-доступ)
- [Поддержка Оболочек](#поддержка-оболочек)
- [Безопасность](#безопасность)
- [Сборка из Исходников](#сборка-из-исходников)
- [Поддержка Платформ](#поддержка-платформ)
- [Структура Проекта](#структура-проекта)

### 🌍 Переводы

| Язык | | Язык | | Язык | |
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

## Быстрый Старт

### 1. Скачать
Выберите вашу платформу в [Releases](../../releases) или соберите из исходников:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Настроить (GUI или текст)
```bash
devspace-gui                  # Настольный конфигуратор (GUI)
devspace init                 # Текстовый конфигуратор
```

### 3. Запустить
```bash
devspace                      # Запускает сервер. Автоопределение конфига.
```

Это также автоматически запускает туннель Cloudflare, если `cloudflared` найден в `tools/`.

### 4. Подключить ваш MCP-клиент
```
https://ВАШ-ТУННЕЛЬ.trycloudflare.com/mcp
```
Или локально: `http://127.0.0.1:7676/mcp`

---

## Установка

Без Node.js, без npm, без Python. Один бинарный файл.

| Платформа | Загрузка |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: скомпилируйте нативно) |
| **macOS Intel** | `devspace` (GUI: скомпилируйте нативно) |
| **macOS M-chip** | `devspace` (GUI: скомпилируйте нативно) |

Требуется **Go 1.23+** только при сборке из исходников.

---

## Что Может Делать ИИ

После подключения ИИ может открыть одну из ваших разрешённых папок проекта как рабочее пространство:

- **Читать, записывать и редактировать** файлы внутри рабочего пространства
- **Искать код** с помощью регулярных выражений и просматривать директории
- **Запускать команды оболочки** (PowerShell на Windows, bash на Unix)
- **Обнаруживать инструкции проекта** из `AGENTS.md` / `CLAUDE.md`
- **Автонастройка** с помощью переносимого `.devspace/config.json`

8 инструментов MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Конфигурация

Вся конфигурация находится **в той же папке, что и исполняемый файл** (переносимая):

```
.devspace/
├── config.json       ← разрешённые корни, порт, оболочка, язык, аутентификация
└── auth.json         ← пароль владельца (необязательно)
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

| Поле | По умолчанию | Описание |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Автоопределение из ОС. Поддерживает 47 языков |
| `toolMode` | `full` | `full` (все инструменты) или `minimal` (только оболочка для поиска) |
| `toolNaming` | `short` | `short` (read, write) или `legacy` (read_file, write_file) |

Переменные окружения не требуются — всё в переносимом файле конфигурации.

---

## Туннель (Удалённый Доступ)

Для веб-версии ChatGPT (требуется HTTPS) DevSpace автоматически запускает туннель:

| Туннель | Тип URL | Настройка |
|---|---|---|
| **Cloudflare** | Случайный (авто) | Поместите `cloudflared.exe` в `tools/` |
| **Pinggy** | Стабильный | Требуется SSH-ключ (`ssh-keygen`) |

Сервер автоматически определяет, какой из них доступен. Перезапустите сервер для нового URL Cloudflare или используйте Pinggy для постоянного URL.

---

## Поддержка Оболочек

| ОС | По умолчанию | Альтернативы |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / любая оболочка |
| **macOS** | bash | `sh` / `zsh` |

Установите `"shell"` в config.json или выберите в GUI.

---

## Безопасность

- **OAuth 2.0 с PKCE** — если установлен пароль владельца
- **Режим без пароля** — если пароль не настроен, работает без аутентификации
- **Изоляция по путям** — все файловые операции проверяются на соответствие разрешённым корням
- **Опциональный туннель** — туннель Cloudflare защищает от прямого доступа
- **Без загрузок третьим лицам** — ваш код никогда не покидает ваш компьютер

---

## Сборка из Исходников

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Собрать всё (все платформы)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Собрать только для текущей платформы
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Поддержка Платформ

| Платформа | Сервер | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (скомпилируйте нативно) |
| **macOS Intel** | ✅ | 🔧 (скомпилируйте нативно) |
| **macOS M-chip** | ✅ | 🔧 (скомпилируйте нативно) |

GUI требует Fyne (OpenGL) — кросс-компиляция невозможна. Сервер компилируется везде.

---

## Структура Проекта

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-сервер
│   └── devspace-gui/       ← Настольный GUI-конфигуратор (Fyne)
├── internal/
│   ├── auth/               ← Провайдер OAuth 2.0 + PKCE
│   ├── config/             ← Система переносимой конфигурации
│   ├── locales/            ← Переводы на 47 языков
│   ├── logger/             ← Структурированное логирование (zerolog)
│   ├── server/             ← HTTP + MCP + оркестрация туннеля
│   ├── skills/             ← Обнаружение AGENTS.md / навыков
│   ├── store/              ← Сессии рабочих пространств SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Проверка рабочего пространства и путей
├── scripts/
│   ├── windows/            ← Сценарий сборки PowerShell
│   └── unix/               ← Сценарии сборки Bash + Makefile
├── readme/                 ← Переводы этого файла (47 языков)
├── build/                  ← Скомпилированные бинарные файлы (все платформы)
├── tools/                  ← cloudflared.exe и т.д.
├── go.mod / go.sum
└── README.md
```

---

Создано на Go. Ноль npm. Ноль Node.js. Один бинарный файл.
