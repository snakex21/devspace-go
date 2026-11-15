# Dev Space Go (Edycja Go)

**Daj ChatGPT i Claude bezpieczny dostęp do swojego lokalnego komputera. Zamień dowolny host MCP w swojego partnera do kodowania.**

Dev Space Go to samodzielnie hostowany serwer MCP, który pozwala asystentom AI czytać, edytować, wyszukiwać i uruchamiać kod w Twoich rzeczywistych lokalnych projektach — Twoje pliki, Twoje narzędzia, Twój terminal — bez przesyłania czegokolwiek do stron trzecich. Uruchamiasz go na swoim komputerze i udostępniasz przez kontrolowany przez siebie tunel.

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
└── config.json       ← dozwolone katalogi główne, port, powłoka, język
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

Dla wersji webowej ChatGPT (wymagane HTTPS), Dev Space Go automatycznie uruchamia tunel:

| Tunel | Typ URL | Konfiguracja |
|---|---|---|
| **Cloudflare** | Losowy (auto) | Umieść `cloudflared.exe` w `tools/` |
| **Pinggy** | Stały | Wymaga klucza SSH (`ssh-keygen`) |

Serwer automatycznie wykrywa, który jest dostępny. Uruchom ponownie serwer, aby uzyskać nowy URL Cloudflare, lub użyj Pinggy dla stałego URL.

### Gotowość połączenia i diagnostyka

Serwer lokalny zaczyna nasłuchiwać przed uruchomieniem tunelu. Sam adres wypisany
przez `cloudflared` nie oznacza gotowego połączenia: poczekaj na ramkę
**TUNEL AKTYWNY** w Dev Space Go. Pojawi się dopiero po sprawdzeniu publicznego
`/healthz`. Przejściowe błędy DNS i bramy są ponawiane w ograniczonym czasie.
Ctrl+C przerywa również uruchamianie tunelu i kończy jego proces potomny.

Z domyślnym Cloudflare Quick Tunnel używaj **`/mcp`**. Quick Tunnels
[nie obsługują SSE](https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/),
więc przełączenie na `/sse` nie naprawi tego połączenia. Starszy endpoint `/sse`
pozostaje dostępny lokalnie lub przez usługę obsługującą SSE.

Przy problemach sprawdź lokalnie `http://127.0.0.1:7676/healthz` (uwzględnij swój
host i port) oraz komunikaty procesu tunelu. Serwer lokalny działa także wtedy,
gdy uruchomienie tunelu się nie powiedzie. Zakończenie procesu tunelu jest
zgłaszane jawnie; jego poprzedni adres przestaje działać. Restart Quick Tunnel
zmienia hostname, dlatego trzeba zaktualizować URL w kliencie MCP.

Kontrola gotowości wymaga wychodzącego HTTPS do wygenerowanej domeny z komputera
serwera. Lokalny proxy, filtr DNS lub zapora blokująca tę kontrolę może uniemożliwić
potwierdzenie gotowości, nawet jeśli tunel jest osiągalny z innego komputera.
Quick Tunnels służą do testów i nie mają gwarancji dostępności.

---

### Dostęp poza projekt (`outsideAccess`)

Domyślnie (`"mode": "full"`) AI **nie jest ograniczone do folderu projektu**:
narzędzia plikowe (`read`, `ls`, `grep`, `glob`, `write`, `edit`, `mkdir`, `move`),
katalog roboczy `bash` i `open_workspace` przyjmują dowolną ścieżkę bezwzględną
na lokalnych dyskach komputera, z odczytem i zapisem, lokalnie i przez tunel.
Ścieżki względne nadal liczą się od otwartego workspace.

```json
"outsideAccess": {
  "mode": "full",
  "remote": true,
  "denyDefaults": true,
  "deny": ["D:/prywatne"],
  "audit": true
}
```

| Pole | Domyślnie | Opis |
|---|---|---|
| `mode` | `full` | `full`: dowolny lokalny folder, odczyt i zapis. `roots`: tylko foldery z `roots`. `off`: tylko projekt |
| `remote` | `true` | Ten sam dostęp dla żądań przez publiczny tunel. `false` (lub odznaczony checkbox w GUI): żądania z tunelu zostają w projekcie. Tunel nie ma uwierzytelniania: kto zna jego tymczasowy adres, ma ten dostęp |
| `denyDefaults` | `true` | Wbudowana lista wrażliwych miejsc (niżej) pozostaje zablokowana. `false` ją wyłącza |
| `deny` | `[]` | Twoje własne ścieżki bezwzględne do zablokowania. `[]` czyści listę |
| `roots` | `[]` | Tylko tryb `roots`: `{ "path": "...", "access": "read" \| "readwrite" }` |
| `audit` | `true` | Każdy dostęp poza projekt (dozwolony i odrzucony) trafia do `.devspace-state/outside-access.log` (tylko ścieżki, nigdy treść) |

Wszystko można też zmienić w GUI: tryb, checkbox tunelu, listę `roots`
(linie `ro <ścieżka>` / `rw <ścieżka>`), przełącznik listy wbudowanej i listę
własnych zablokowanych ścieżek (dodaj, usuń, przeglądaj). Zapis z GUI zachowuje
wszystkie inne wartości wpisane ręcznie w `config.json`.

**Wbudowana lista w trybie `full`** (`denyDefaults`), obowiązuje wszędzie,
także wewnątrz workspace:

- Prywatne klucze SSH: `id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa` (`.pub` pozostaje czytelny)
- Menedżer poświadczeń Windows, Vault i klucze DPAPI: `%APPDATA%\Microsoft\Credentials`, `Protect`, `Vault`, `%LOCALAPPDATA%\Microsoft\Credentials`, `Vault` (macOS: `~/Library/Keychains`)
- Hasła i ciasteczka przeglądarek: `Login Data`, `Cookies`, `logins.json`, `key4.db`, `cookies.sqlite`
- `.git-credentials` i token GitHub CLI (`gh/hosts.yml`)
- Własna konfiguracja i stan Dev Space Go, żeby AI nie mogło zmienić sobie uprawnień

`.env`, `.npmrc`, foldery systemowe i cała reszta pozostają dostępne. W trybie
`roots` lista jest szersza (także `.aws`, `.gnupg`, `.kube`, `.docker`,
`.git/config`, `.npmrc`, profile przeglądarek oraz foldery systemowe, chyba że
skonfigurujesz folder wewnątrz nich).

Higiena ścieżek działa w każdym trybie: przed sprawdzeniem system kanonikalizuje
ścieżkę (symlinki, junction, dyski subst, nazwy 8.3, wielkość liter), a ścieżki
UNC (`\\serwer\udział`, które ujawniłyby hash logowania Windows), `\\?\` / `\\.\`,
strumienie NTFS (`plik:strumień`), nazwy zarezerwowane (`CON`, `NUL`, ...) i ścieżki
względne względem dysku (`C:foo`) są odrzucane. Workspace i jego folderów
nadrzędnych nie da się przenieść. `devspace doctor` pokazuje aktywny tryb i listy.

pnpm na Windows: `node_modules` to junction do magazynu pnpm poza projektem.
W trybie `full` po prostu działa; w trybie `roots` dodaj magazyn
(`pnpm store path`) jako folder `read`.

## Obsługa powłoki

| System | Domyślnie | Alternatywy |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / dowolna powłoka |
| **macOS** | bash | `sh` / `zsh` |

Ustaw `"shell"` w config.json lub wybierz w GUI.

---

## Bezpieczeństwo

- **Brak wbudowanego uwierzytelniania** — każdy, kto zna publiczny adres tunelu, może używać serwera
- **Ograniczenie ścieżek** — każda ścieżka narzędzia (read, write, edit, mkdir, move, ls, grep, glob i katalog roboczy powłoki) jest rozwiązywana wewnątrz workspace; `..`, ścieżki bezwzględne gdzie indziej oraz symlinki/junction prowadzące na zewnątrz są odrzucane. Tak działa `outsideAccess.mode: off`. Domyślny tryb `full` celowo pozwala na ścieżki bezwzględne na wszystkich lokalnych dyskach (lokalnie i przez tunel) poza wbudowaną listą; `roots` pozwala tylko na wybrane foldery (patrz wyżej)
- **Dostęp przez tunel** — traktuj publiczny adres jak sekret i wyłączaj serwer, gdy go nie używasz
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
