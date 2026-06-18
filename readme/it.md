# DevSpace (Edizione Go)

**Dai a ChatGPT e Claude accesso sicuro alla tua macchina locale. Trasforma qualsiasi host MCP nel tuo partner di codifica.**

DevSpace è un server MCP auto-ospitato che consente agli assistenti AI di leggere, modificare, cercare ed eseguire codice nei tuoi veri progetti locali — i tuoi file, i tuoi strumenti, il tuo terminale — senza caricare nulla su terze parti. Lo esegui sulla tua macchina, lo esponi attraverso un tunnel che controlli e opzionalmente lo proteggi con una password.

---

## Indice

- [Avvio Rapido](#avvio-rapido)
- [Installazione](#installazione)
- [Cosa Può Fare l'AI](#cosa-può-fare-lai)
- [Configurazione](#configurazione)
- [Tunnel (Accesso Remoto)](#tunnel-accesso-remoto)
- [Supporto Shell](#supporto-shell)
- [Sicurezza](#sicurezza)
- [Compilazione dai Sorgenti](#compilazione-dai-sorgenti)
- [Supporto Piattaforme](#supporto-piattaforme)
- [Struttura del Progetto](#struttura-del-progetto)

### 🌍 Traduzioni

| Lingua | | Lingua | | Lingua | |
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

## Avvio Rapido

### 1. Scarica
Scegli la tua piattaforma dalla pagina [Releases](../../releases) o compila dai sorgenti:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configura (GUI o testo)
```bash
devspace-gui                  # Configuratore desktop (GUI)
devspace init                 # Configuratore testuale
```

### 3. Esegui
```bash
devspace                      # Avvia il server. Rileva automaticamente la configurazione.
```

Questo avvia anche automaticamente un Tunnel Cloudflare se `cloudflared` è presente in `tools/`.

### 4. Connetti il tuo client MCP
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
O localmente: `http://127.0.0.1:7676/mcp`

---

## Installazione

Niente Node.js, niente npm, niente Python. Un singolo binario.

| Piattaforma | Download |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compilare nativamente) |
| **macOS Intel** | `devspace` (GUI: compilare nativamente) |
| **macOS M-chip** | `devspace` (GUI: compilare nativamente) |

Richiede **Go 1.23+** solo se compili dai sorgenti.

---

## Cosa Può Fare l'AI

Una volta connessa, l'AI può aprire una delle tue cartelle di progetto approvate come area di lavoro:

- **Leggere, scrivere e modificare** file all'interno dell'area di lavoro
- **Cercare codice** con regex e ispezionare directory
- **Eseguire comandi shell** (PowerShell su Windows, bash su Unix)
- **Scoprire istruzioni di progetto** da `AGENTS.md` / `CLAUDE.md`
- **Auto-configurazione** con `.devspace/config.json` portatile

8 strumenti MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configurazione

Tutta la configurazione risiede **nella stessa cartella dell'eseguibile** (portatile):

```
.devspace/
├── config.json       ← root consentite, porta, shell, lingua, autenticazione
└── auth.json         ← password proprietario (opzionale)
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

| Campo | Predefinito | Descrizione |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Rilevamento automatico dal sistema operativo. Supporta 47 lingue |
| `toolMode` | `full` | `full` (tutti gli strumenti) o `minimal` (solo shell per la ricerca) |
| `toolNaming` | `short` | `short` (read, write) o `legacy` (read_file, write_file) |

Nessuna variabile d'ambiente necessaria — tutto è nel file di configurazione portatile.

---

## Tunnel (Accesso Remoto)

Per la versione web di ChatGPT (richiede HTTPS), DevSpace avvia automaticamente un tunnel:

| Tunnel | Tipo URL | Configurazione |
|---|---|---|
| **Cloudflare** | Casuale (auto) | Inserisci `cloudflared.exe` in `tools/` |
| **Pinggy** | Stabile | Richiede chiave SSH (`ssh-keygen`) |

Il server rileva automaticamente quale è disponibile. Riavvia il server per un nuovo URL Cloudflare, o usa Pinggy per un URL permanente.

---

## Supporto Shell

| SO | Predefinita | Alternative |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / qualsiasi shell |
| **macOS** | bash | `sh` / `zsh` |

Imposta `"shell"` in config.json o scegli nella GUI.

---

## Sicurezza

- **OAuth 2.0 con PKCE** — se è impostata la password del proprietario
- **Modalità senza password** — se nessuna password è configurata, funziona senza autenticazione
- **Contenimento del percorso** — tutte le operazioni sui file sono convalidate rispetto alle root consentite
- **Tunnel opzionale** — il Tunnel Cloudflare protegge dall'esposizione diretta
- **Nessun caricamento a terzi** — il tuo codice non lascia mai la tua macchina

---

## Compilazione dai Sorgenti

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compila tutto (tutte le piattaforme)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compila solo per la piattaforma corrente
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Supporto Piattaforme

| Piattaforma | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compilare nativamente) |
| **macOS Intel** | ✅ | 🔧 (compilare nativamente) |
| **macOS M-chip** | ✅ | 🔧 (compilare nativamente) |

La GUI richiede Fyne (OpenGL) — non è possibile la compilazione incrociata. Il server si compila ovunque.

---

## Struttura del Progetto

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + server MCP
│   └── devspace-gui/       ← Configuratore GUI desktop (Fyne)
├── internal/
│   ├── auth/               ← Provider OAuth 2.0 + PKCE
│   ├── config/             ← Sistema di configurazione portatile
│   ├── locales/            ← Traduzioni in 47 lingue
│   ├── logger/             ← Logging strutturato (zerolog)
│   ├── server/             ← HTTP + MCP + orchestrazione tunnel
│   ├── skills/             ← AGENTS.md / scoperta competenze
│   ├── store/              ← Sessioni area di lavoro SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Area di lavoro e validazione percorso
├── scripts/
│   ├── windows/            ← Script di build PowerShell
│   └── unix/               ← Script di build Bash + Makefile
├── readme/                 ← Traduzioni di questo file (47 lingue)
├── build/                  ← Binari compilati (tutte le piattaforme)
├── tools/                  ← cloudflared.exe, ecc.
├── go.mod / go.sum
└── README.md
```

---

Realizzato in Go. Zero npm. Zero Node.js. Un binario.
