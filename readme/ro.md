# DevSpace (Ediția Go)

**Oferă ChatGPT și Claude acces securizat la mașina ta locală. Transformă orice host MCP în partenerul tău de programare.**

DevSpace este un server MCP auto-găzduit care permite asistenților AI să citească, editeze, caute și ruleze cod în proiectele tale locale reale — fișierele tale, uneltele tale, terminalul tău — fără a încărca nimic către o terță parte. Îl rulezi pe mașina ta, îl expui printr-un tunel pe care îl controlezi și, opțional, îl securizezi cu o parolă.

---

## Cuprins

- [Pornire Rapidă](#pornire-rapidă)
- [Instalare](#instalare)
- [Ce Poate Face AI-ul](#ce-poate-face-ai-ul)
- [Configurare](#configurare)
- [Tunel (Acces la Distanță)](#tunel-acces-la-distanță)
- [Suport Shell](#suport-shell)
- [Securitate](#securitate)
- [Compilare din Surse](#compilare-din-surse)
- [Suport Platforme](#suport-platforme)
- [Structura Proiectului](#structura-proiectului)

### 🌍 Traduceri

| Limbă | | Limbă | | Limbă | |
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

## Pornire Rapidă

### 1. Descărcare
Alege platforma ta din [Releases](../../releases) sau compilează din surse:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configurare (GUI sau text)
```bash
devspace-gui                  # Configurator desktop (GUI)
devspace init                 # Configurator bazat pe text
```

### 3. Rulează
```bash
devspace                      # Pornește serverul. Detectează automat config.
```

Acest lucru pornește automat și un Tunel Cloudflare dacă `cloudflared` este găsit în `tools/`.

### 4. Conectează clientul tău MCP
```
https://TUNELUL-TĂU.trycloudflare.com/mcp
```
Sau local: `http://127.0.0.1:7676/mcp`

---

## Instalare

Fără Node.js, fără npm, fără Python. Un singur binar.

| Platformă | Descărcare |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compilează nativ) |
| **macOS Intel** | `devspace` (GUI: compilează nativ) |
| **macOS M-chip** | `devspace` (GUI: compilează nativ) |

Necesită **Go 1.23+** doar dacă compilezi din surse.

---

## Ce Poate Face AI-ul

Odată conectat, AI-ul poate deschide unul dintre dosarele tale de proiect aprobate ca spațiu de lucru:

- **Citește, scrie și editează** fișiere în interiorul spațiului de lucru
- **Caută cod** cu expresii regulate și inspectează directoare
- **Rulează comenzi shell** (PowerShell pe Windows, bash pe Unix)
- **Descoperă instrucțiunile proiectului** din `AGENTS.md` / `CLAUDE.md`
- **Auto-configurare** cu `.devspace/config.json` portabil

8 unelte MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configurare

Toată configurația se află **în același dosar cu executabilul** (portabil):

```
.devspace/
├── config.json       ← rădăcini permise, port, shell, limbă, autentificare
└── auth.json         ← parola proprietarului (opțional)
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

| Câmp | Implicit | Descriere |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Detectare automată din SO. Suportă 47 de limbi |
| `toolMode` | `full` | `full` (toate uneltele) sau `minimal` (doar shell pentru căutare) |
| `toolNaming` | `short` | `short` (read, write) sau `legacy` (read_file, write_file) |

Nu sunt necesare variabile de mediu — totul este în fișierul de configurare portabil.

---

## Tunel (Acces la Distanță)

Pentru versiunea web ChatGPT (necesită HTTPS), DevSpace pornește automat un tunel:

| Tunel | Tip URL | Configurare |
|---|---|---|
| **Cloudflare** | Aleatoriu (auto) | Pune `cloudflared.exe` în `tools/` |
| **Pinggy** | Stabil | Necesită cheie SSH (`ssh-keygen`) |

Serverul detectează automat care este disponibil. Repornește serverul pentru un nou URL Cloudflare sau folosește Pinggy pentru un URL permanent.

---

## Suport Shell

| SO | Implicit | Alternative |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / orice shell |
| **macOS** | bash | `sh` / `zsh` |

Setează `"shell"` în config.json sau alege în GUI.

---

## Securitate

- **OAuth 2.0 cu PKCE** — dacă este setată o parolă de proprietar
- **Mod fără parolă** — dacă nu este configurată nicio parolă, rulează fără autentificare
- **Izolare pe căi** — toate operațiile pe fișiere sunt validate față de rădăcinile permise
- **Tunel opțional** — Tunelul Cloudflare protejează împotriva expunerii directe
- **Fără încărcări către terți** — codul tău nu părăsește niciodată mașina ta

---

## Compilare din Surse

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compilează totul (toate platformele)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compilează doar pentru platforma curentă
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Suport Platforme

| Platformă | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compilează nativ) |
| **macOS Intel** | ✅ | 🔧 (compilează nativ) |
| **macOS M-chip** | ✅ | 🔧 (compilează nativ) |

GUI necesită Fyne (OpenGL) — nu se poate compila încrucișat. Serverul se compilează peste tot.

---

## Structura Proiectului

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + server MCP
│   └── devspace-gui/       ← Configurator GUI desktop (Fyne)
├── internal/
│   ├── auth/               ← Furnizor OAuth 2.0 + PKCE
│   ├── config/             ← Sistem de configurare portabil
│   ├── locales/            ← Traduceri în 47 de limbi
│   ├── logger/             ← Logging structurat (zerolog)
│   ├── server/             ← HTTP + MCP + orchestrare tunel
│   ├── skills/             ← Descoperire AGENTS.md / skill
│   ├── store/              ← Sesiuni spațiu de lucru SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Validare spațiu de lucru & căi
├── scripts/
│   ├── windows/            ← Script build PowerShell
│   └── unix/               ← Scripturi build Bash + Makefile
├── readme/                 ← Traduceri ale acestui fișier (47 de limbi)
├── build/                  ← Binare compilate (toate platformele)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Construit în Go. Zero npm. Zero Node.js. Un singur binar.
