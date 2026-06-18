# DevSpace (Edycja Go)

**Daj ChatGPT i Claude bezpieczny dostęp do swojego lokalnego komputera. Zamień dowolny host MCP w swojego partnera do kodowania.**

DevSpace to samodzielnie hostowany serwer MCP, który pozwala asystentom AI czytać, edytować, wyszukiwać i uruchamiać kod w Twoich rzeczywistych lokalnych projektach — Twoje pliki, Twoje narzędzia, Twój terminal — bez przesyłania czegokolwiek do stron trzecich. Uruchamiasz go na swoim komputerze, udostępniasz przez kontrolowany przez siebie tunel i opcjonalnie zabezpieczasz hasłem.

---

## Spis treści

- [Szybki start](#szybki-start)
- [Instalacja](#instalacja)
- [Co AI może zrobić](#co-ai-może-zrobić)
- [Konfiguracja](#konfiguracja)
- [Tunel (dostęp zdalny)](#tunel-dostęp-zdalny)
- [Obsługa powłoki](#obsługa-powłoki)
- [Bezpieczeństwo](#bezpieczeństwo)
- [Budowanie ze źródeł](#budowanie-ze-źródeł)
- [Wsparcie platform](#wsparcie-platform)
- [Struktura projektu](#struktura-projektu)

### 🌍 Tłumaczenia

| Język | | Język | | Język | |
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

## Szybki start

### 1. Pobierz
Wybierz swoją platformę ze strony [Releases](../../releases) lub zbuduj ze źródeł:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Skonfiguruj (GUI lub tekst)
```bash
devspace-gui                  # Konfigurator pulpitu (GUI)
devspace init                 # Konfigurator tekstowy
```

### 3. Uruchom
```bash
devspace                      # Uruchamia serwer. Automatycznie wykrywa konfigurację.
```

To również automatycznie uruchamia tunel Cloudflare, jeśli `cloudflared` zostanie znaleziony w `tools/`.

### 4. Połącz swojego klienta MCP
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Lub lokalnie: `http://127.0.0.1:7676/mcp`

---

## Instalacja

Bez Node.js, bez npm, bez Pythona. Pojedynczy plik binarny.

| Platforma | Pobieranie |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilować natywnie) |
| **macOS Intel** | `devspace` (GUI: kompilować natywnie) |
| **macOS M-chip** | `devspace` (GUI: kompilować natywnie) |

Wymaga **Go 1.23+** tylko przy budowaniu ze źródeł.

---

## Co AI może zrobić

Po połączeniu AI może otworzyć jeden z zatwierdzonych folderów projektu jako obszar roboczy:

- **Czytać, zapisywać i edytować** pliki w obszarze roboczym
- **Wyszukiwać kod** za pomocą wyrażeń regularnych i przeglądać katalogi
- **Uruchamiać polecenia powłoki** (PowerShell na Windows, bash na Unix)
- **Odkrywać instrukcje projektu** z `AGENTS.md` / `CLAUDE.md`
- **Automatycznie konfigurować** za pomocą przenośnego `.devspace/config.json`

8 narzędzi MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguracja

Cała konfiguracja znajduje się **w tym samym folderze co plik wykonywalny** (przenośna):

```
.devspace/
├── config.json       ← dozwolone katalogi główne, port, powłoka, język, uwierzytelnianie
└── auth.json         ← hasło właściciela (opcjonalnie)
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

| Pole | Domyślnie | Opis |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatyczne wykrywanie z systemu operacyjnego. Obsługuje 47 języków |
| `toolMode` | `full` | `full` (wszystkie narzędzia) lub `minimal` (tylko powłoka do wyszukiwania) |
| `toolNaming` | `short` | `short` (read, write) lub `legacy` (read_file, write_file) |

Nie są potrzebne zmienne środowiskowe — wszystko znajduje się w przenośnym pliku konfiguracyjnym.

---

## Tunel (dostęp zdalny)

Dla wersji webowej ChatGPT (wymagane HTTPS), DevSpace automatycznie uruchamia tunel:

| Tunel | Typ URL | Konfiguracja |
|---|---|---|
| **Cloudflare** | Losowy (auto) | Umieść `cloudflared.exe` w `tools/` |
| **Pinggy** | Stały | Wymaga klucza SSH (`ssh-keygen`) |

Serwer automatycznie wykrywa, który jest dostępny. Uruchom ponownie serwer, aby uzyskać nowy URL Cloudflare, lub użyj Pinggy dla stałego URL.

---

## Obsługa powłoki

| System | Domyślnie | Alternatywy |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / dowolna powłoka |
| **macOS** | bash | `sh` / `zsh` |

Ustaw `"shell"` w config.json lub wybierz w GUI.

---

## Bezpieczeństwo

- **OAuth 2.0 z PKCE** — jeśli ustawiono hasło właściciela
- **Tryb bez hasła** — jeśli nie skonfigurowano hasła, działa bez uwierzytelniania
- **Ograniczenie ścieżek** — wszystkie operacje na plikach są sprawdzane względem dozwolonych katalogów głównych
- **Opcjonalny tunel** — tunel Cloudflare chroni przed bezpośrednim wystawieniem
- **Brak wysyłania do stron trzecich** — Twój kod nigdy nie opuszcza Twojego komputera

---

## Budowanie ze źródeł

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Zbuduj wszystko (wszystkie platformy)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Zbuduj tylko dla bieżącej platformy
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Wsparcie platform

| Platforma | Serwer | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilować natywnie) |
| **macOS Intel** | ✅ | 🔧 (kompilować natywnie) |
| **macOS M-chip** | ✅ | 🔧 (kompilować natywnie) |

GUI wymaga Fyne (OpenGL) — nie można kompilować krzyżowo. Serwer kompiluje się wszędzie.

---

## Struktura projektu

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + serwer MCP
│   └── devspace-gui/       ← Konfigurator pulpitu GUI (Fyne)
├── internal/
│   ├── auth/               ← Dostawca OAuth 2.0 + PKCE
│   ├── config/             ← Przenośny system konfiguracji
│   ├── locales/            ← Tłumaczenia na 47 języków
│   ├── logger/             ← Strukturalne logowanie (zerolog)
│   ├── server/             ← HTTP + MCP + orkiestracja tunelu
│   ├── skills/             ← AGENTS.md / wykrywanie umiejętności
│   ├── store/              ← Sesje obszaru roboczego SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Obszar roboczy i walidacja ścieżek
├── scripts/
│   ├── windows/            ← Skrypt budowania PowerShell
│   └── unix/               ← Skrypty budowania Bash + Makefile
├── readme/                 ← Tłumaczenia tego pliku (47 języków)
├── build/                  ← Skompilowane pliki binarne (wszystkie platformy)
├── tools/                  ← cloudflared.exe, itp.
├── go.mod / go.sum
└── README.md
```

---

Zbudowane w Go. Zero npm. Zero Node.js. Jeden plik binarny.
