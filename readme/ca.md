# DevSpace (Go Edition) - Català

**Doneu a ChatGPT i Claude accés segur al vostre ordinador local. Convertiu qualsevol amfitrió MCP en el vostre company de codificació.**

DevSpace és un servidor MCP autoallotjat que permet als assistents d'IA llegir, editar, buscar i executar codi als vostres projectes locals reals — els vostres fitxers, les vostres eines, el vostre terminal — sense penjar res a un tercer. L'executeu a la vostra màquina, l'exposeu a través d'un túnel que controleu i, opcionalment, el protegiu amb una contrasenya.

---

## Taula de Continguts

- [Inici Ràpid](#inici-ràpid)
- [Instal·lació](#instal·lació)
- [Què Pot Fer la IA](#què-pot-fer-la-ia)
- [Configuració](#configuració)
- [Túnel (Accés Remot)](#túnel-accés-remot)
- [Suport d'Intèrpret d'Ordres](#suport-dintèrpret-dordres)
- [Seguretat](#seguretat)
- [Compilació des del Codi Font](#compilació-des-del-codi-font)
- [Suport de Plataformes](#suport-de-plataformes)
- [Estructura del Projecte](#estructura-del-projecte)

### 🌍 Traduccions

| Idioma | | Idioma | | Idioma | |
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

## Inici Ràpid

### 1. Descarregar
Trieu la vostra plataforma des de [Llançaments](../../releases) o compileu des del codi font:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configurar (GUI o text)
```bash
devspace-gui                  # Configurador d'escriptori (GUI)
devspace init                 # Configurador basat en text
```

### 3. Executar
```bash
devspace                      # Inicia el servidor. Detecta la configuració automàticament.
```

Això també inicia automàticament un túnel de Cloudflare si es troba `cloudflared` a `tools/`.

### 4. Connecteu el vostre client MCP
```
https://EL-VOSTRE-TÚNEL.trycloudflare.com/mcp
```
O localment: `http://127.0.0.1:7676/mcp`

---

## Instal·lació

Sense Node.js, sense npm, sense Python. Un sol binari.

| Plataforma | Descàrrega |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compileu de manera nativa) |
| **macOS Intel** | `devspace` (GUI: compileu de manera nativa) |
| **macOS M-chip** | `devspace` (GUI: compileu de manera nativa) |

Requereix **Go 1.23+** només si es compila des del codi font.

---

## Què Pot Fer la IA

Un cop connectat, la IA pot obrir una de les vostres carpetes de projecte aprovades com a espai de treball:

- **Llegir, escriure i editar** fitxers dins de l'espai de treball
- **Buscar codi** amb regex i inspeccionar directoris
- **Executar ordres d'intèrpret** (PowerShell a Windows, bash a Unix)
- **Descobrir instruccions del projecte** des de `AGENTS.md` / `CLAUDE.md`
- **Autoconfigurar** amb `.devspace/config.json` portàtil

8 eines MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuració

Tota la configuració resideix **a la mateixa carpeta que l'executable** (portàtil):

```
.devspace/
├── config.json       ← arrels permeses, port, intèrpret, idioma, autenticació
└── auth.json         ← contrasenya del propietari (opcional)
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

| Camp | Per defecte | Descripció |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Detecció automàtica des del SO. Admet 47 idiomes |
| `toolMode` | `full` | `full` (totes les eines) o `minimal` (només intèrpret per cercar) |
| `toolNaming` | `short` | `short` (read, write) o `legacy` (read_file, write_file) |

No calen variables d'entorn — tot és al fitxer de configuració portàtil.

---

## Túnel (Accés Remot)

Per a la versió web de ChatGPT (requereix HTTPS), DevSpace inicia automàticament un túnel:

| Túnel | Tipus d'URL | Configuració |
|---|---|---|
| **Cloudflare** | Aleatori (auto) | Poseu `cloudflared.exe` a `tools/` |
| **Pinggy** | Estable | Necessita clau SSH (`ssh-keygen`) |

El servidor detecta automàticament quin està disponible. Reinicieu el servidor per a una nova URL de Cloudflare o utilitzeu Pinggy per a una URL permanent.

---

## Suport d'Intèrpret d'Ordres

| SO | Per defecte | Alternatives |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / qualsevol intèrpret |
| **macOS** | bash | `sh` / `zsh` |

Establiu `"shell"` a config.json o trieu a la GUI.

---

## Seguretat

- **OAuth 2.0 amb PKCE** — si s'estableix la contrasenya del propietari
- **Mode sense contrasenya** — si no es configura cap contrasenya, s'executa sense autenticació
- **Contenció de ruta** — totes les operacions de fitxers es validen contra les arrels permeses
- **Túnel opcional** — El túnel de Cloudflare protegeix de l'exposició directa
- **Sense càrregues a tercers** — el vostre codi mai no surt de la vostra màquina

---

## Compilació des del Codi Font

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compilar tot (totes les plataformes)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compilar només per a la plataforma actual
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Suport de Plataformes

| Plataforma | Servidor | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compileu de manera nativa) |
| **macOS Intel** | ✅ | 🔧 (compileu de manera nativa) |
| **macOS M-chip** | ✅ | 🔧 (compileu de manera nativa) |

La GUI requereix Fyne (OpenGL) — no es pot compilar de manera creuada. El servidor es compila a tot arreu.

---

## Estructura del Projecte

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + servidor MCP
│   └── devspace-gui/       ← Configurador d'escriptori GUI (Fyne)
├── internal/
│   ├── auth/               ← Proveïdor OAuth 2.0 + PKCE
│   ├── config/             ← Sistema de configuració portàtil
│   ├── locales/            ← Traduccions a 47 idiomes
│   ├── logger/             ← Registre estructurat (zerolog)
│   ├── server/             ← HTTP + MCP + orquestració de túnel
│   ├── skills/             ← AGENTS.md / descobriment d'habilitats
│   ├── store/              ← Sessions d'espai de treball SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Espai de treball i validació de ruta
├── scripts/
│   ├── windows/            ← Script de compilació PowerShell
│   └── unix/               ← Scripts de compilació Bash + Makefile
├── readme/                 ← Traduccions d'aquest fitxer (47 idiomes)
├── build/                  ← Binaris compilats (totes les plataformes)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Construït en Go. Zero npm. Zero Node.js. Un binari.
