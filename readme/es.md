# DevSpace (Go Edition) - Español

**Dele a ChatGPT y Claude acceso seguro a su máquina local. Convierta cualquier host MCP en su compañero de programación.**

DevSpace es un servidor MCP autoalojado que permite a los asistentes de IA leer, editar, buscar y ejecutar código en sus proyectos locales reales — sus archivos, sus herramientas, su terminal — sin subir nada a un tercero. Lo ejecuta en su máquina, lo expone a través de un túnel que usted controla y, opcionalmente, lo asegura con una contraseña.

---

## Tabla de Contenidos

- [Inicio Rápido](#inicio-rápido)
- [Instalación](#instalación)
- [Qué Puede Hacer la IA](#qué-puede-hacer-la-ia)
- [Configuración](#configuración)
- [Túnel (Acceso Remoto)](#túnel-acceso-remoto)
- [Soporte de Shell](#soporte-de-shell)
- [Seguridad](#seguridad)
- [Compilación desde el Código Fuente](#compilación-desde-el-código-fuente)
- [Soporte de Plataformas](#soporte-de-plataformas)
- [Estructura del Proyecto](#estructura-del-proyecto)

### 🌍 Traducciones

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

## Inicio Rápido

### 1. Descargar
Elija su plataforma desde [Versiones](../../releases) o compile desde el código fuente:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configurar (GUI o texto)
```bash
devspace-gui                  # Configurador de escritorio (GUI)
devspace init                 # Configurador basado en texto
```

### 3. Ejecutar
```bash
devspace                      # Inicia el servidor. Detecta automáticamente la configuración.
```

Esto también inicia automáticamente un túnel de Cloudflare si se encuentra `cloudflared` en `tools/`.

### 4. Conecte su cliente MCP
```
https://SU-TÚNEL.trycloudflare.com/mcp
```
O localmente: `http://127.0.0.1:7676/mcp`

---

## Instalación

Sin Node.js, sin npm, sin Python. Un solo binario.

| Plataforma | Descarga |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compilar de forma nativa) |
| **macOS Intel** | `devspace` (GUI: compilar de forma nativa) |
| **macOS M-chip** | `devspace` (GUI: compilar de forma nativa) |

Requiere **Go 1.23+** solo si compila desde el código fuente.

---

## Qué Puede Hacer la IA

Una vez conectado, la IA puede abrir una de sus carpetas de proyecto aprobadas como espacio de trabajo:

- **Leer, escribir y editar** archivos dentro del espacio de trabajo
- **Buscar código** con regex e inspeccionar directorios
- **Ejecutar comandos de shell** (PowerShell en Windows, bash en Unix)
- **Descubrir instrucciones del proyecto** desde `AGENTS.md` / `CLAUDE.md`
- **Autoconfigurar** con `.devspace/config.json` portátil

8 herramientas MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuración

Toda la configuración reside **en la misma carpeta que el ejecutable** (portátil):

```
.devspace/
├── config.json       ← raíces permitidas, puerto, shell, idioma, autenticación
└── auth.json         ← contraseña del propietario (opcional)
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

| Campo | Predeterminado | Descripción |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Detección automática desde el SO. Soporta 47 idiomas |
| `toolMode` | `full` | `full` (todas las herramientas) o `minimal` (solo shell para búsqueda) |
| `toolNaming` | `short` | `short` (read, write) o `legacy` (read_file, write_file) |

No se necesitan variables de entorno — todo está en el archivo de configuración portátil.

---

## Túnel (Acceso Remoto)

Para la versión web de ChatGPT (requiere HTTPS), DevSpace inicia automáticamente un túnel:

| Túnel | Tipo de URL | Configuración |
|---|---|---|
| **Cloudflare** | Aleatorio (auto) | Coloque `cloudflared.exe` en `tools/` |
| **Pinggy** | Estable | Necesita clave SSH (`ssh-keygen`) |

El servidor detecta automáticamente cuál está disponible. Reinicie el servidor para una nueva URL de Cloudflare o use Pinggy para una URL permanente.

---

## Soporte de Shell

| SO | Predeterminado | Alternativas |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / cualquier shell |
| **macOS** | bash | `sh` / `zsh` |

Establezca `"shell"` en config.json o elija en la GUI.

---

## Seguridad

- **OAuth 2.0 con PKCE** — si se establece la contraseña del propietario
- **Modo sin contraseña** — si no se configura contraseña, se ejecuta sin autenticación
- **Contención de ruta** — todas las operaciones de archivos se validan contra las raíces permitidas
- **Túnel opcional** — El túnel de Cloudflare protege de la exposición directa
- **Sin cargas a terceros** — su código nunca sale de su máquina

---

## Compilación desde el Código Fuente

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compilar todo (todas las plataformas)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compilar solo para la plataforma actual
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Soporte de Plataformas

| Plataforma | Servidor | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compilar de forma nativa) |
| **macOS Intel** | ✅ | 🔧 (compilar de forma nativa) |
| **macOS M-chip** | ✅ | 🔧 (compilar de forma nativa) |

La GUI requiere Fyne (OpenGL) — no se puede compilar de forma cruzada. El servidor se compila en todas partes.

---

## Estructura del Proyecto

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + servidor MCP
│   └── devspace-gui/       ← Configurador de escritorio GUI (Fyne)
├── internal/
│   ├── auth/               ← Proveedor OAuth 2.0 + PKCE
│   ├── config/             ← Sistema de configuración portátil
│   ├── locales/            ← Traducciones a 47 idiomas
│   ├── logger/             ← Registro estructurado (zerolog)
│   ├── server/             ← HTTP + MCP + orquestación de túnel
│   ├── skills/             ← AGENTS.md / descubrimiento de habilidades
│   ├── store/              ← Sesiones de espacio de trabajo SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Espacio de trabajo y validación de ruta
├── scripts/
│   ├── windows/            ← Script de compilación PowerShell
│   └── unix/               ← Scripts de compilación Bash + Makefile
├── readme/                 ← Traducciones de este archivo (47 idiomas)
├── build/                  ← Binarios compilados (todas las plataformas)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Construido en Go. Cero npm. Cero Node.js. Un binario.
