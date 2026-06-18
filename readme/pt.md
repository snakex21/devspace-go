# DevSpace (Edição Go)

**Dê ao ChatGPT e ao Claude acesso seguro à sua máquina local. Transforme qualquer host MCP no seu parceiro de codificação.**

O DevSpace é um servidor MCP auto-hospedado que permite que assistentes de IA leiam, editem, pesquisem e executem código nos seus projetos locais reais — seus arquivos, suas ferramentas, seu terminal — sem enviar nada para terceiros. Você o executa na sua máquina, expõe através de um túnel que você controla e, opcionalmente, protege com uma senha.

---

## Índice

- [Início Rápido](#início-rápido)
- [Instalação](#instalação)
- [O Que a IA Pode Fazer](#o-que-a-ia-pode-fazer)
- [Configuração](#configuração)
- [Túnel (Acesso Remoto)](#túnel-acesso-remoto)
- [Suporte a Shell](#suporte-a-shell)
- [Segurança](#segurança)
- [Compilando do Código Fonte](#compilando-do-código-fonte)
- [Suporte a Plataformas](#suporte-a-plataformas)
- [Estrutura do Projeto](#estrutura-do-projeto)

### 🌍 Traduções

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

## Início Rápido

### 1. Baixar
Escolha sua plataforma na página de [Releases](../../releases) ou compile do código fonte:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configurar (GUI ou texto)
```bash
devspace-gui                  # Configurador de desktop (GUI)
devspace init                 # Configurador baseado em texto
```

### 3. Executar
```bash
devspace                      # Inicia o servidor. Detecta automaticamente a configuração.
```

Isso também inicia automaticamente um Túnel Cloudflare se `cloudflared` for encontrado em `tools/`.

### 4. Conectar seu cliente MCP
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Ou localmente: `http://127.0.0.1:7676/mcp`

---

## Instalação

Sem Node.js, sem npm, sem Python. Binário único.

| Plataforma | Download |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: compilar nativamente) |
| **macOS Intel** | `devspace` (GUI: compilar nativamente) |
| **macOS M-chip** | `devspace` (GUI: compilar nativamente) |

Requer **Go 1.23+** apenas se compilar do código fonte.

---

## O Que a IA Pode Fazer

Uma vez conectada, a IA pode abrir uma das suas pastas de projeto aprovadas como espaço de trabalho:

- **Ler, escrever e editar** arquivos dentro do espaço de trabalho
- **Pesquisar código** com regex e inspecionar diretórios
- **Executar comandos shell** (PowerShell no Windows, bash no Unix)
- **Descobrir instruções do projeto** a partir de `AGENTS.md` / `CLAUDE.md`
- **Auto-configurar** com `.devspace/config.json` portátil

8 ferramentas MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuração

Toda a configuração fica **na mesma pasta do executável** (portátil):

```
.devspace/
├── config.json       ← raízes permitidas, porta, shell, idioma, autenticação
└── auth.json         ← senha do proprietário (opcional)
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

| Campo | Padrão | Descrição |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Detectado automaticamente do SO. Suporta 47 idiomas |
| `toolMode` | `full` | `full` (todas as ferramentas) ou `minimal` (apenas shell para pesquisa) |
| `toolNaming` | `short` | `short` (read, write) ou `legacy` (read_file, write_file) |

Nenhuma variável de ambiente necessária — tudo está no arquivo de configuração portátil.

---

## Túnel (Acesso Remoto)

Para a versão web do ChatGPT (HTTPS necessário), o DevSpace inicia automaticamente um túnel:

| Túnel | Tipo de URL | Configuração |
|---|---|---|
| **Cloudflare** | Aleatório (auto) | Coloque `cloudflared.exe` em `tools/` |
| **Pinggy** | Estável | Requer chave SSH (`ssh-keygen`) |

O servidor detecta automaticamente qual está disponível. Reinicie o servidor para um novo URL do Cloudflare ou use o Pinggy para um URL permanente.

---

## Suporte a Shell

| SO | Padrão | Alternativas |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / qualquer shell |
| **macOS** | bash | `sh` / `zsh` |

Defina `"shell"` no config.json ou escolha na GUI.

---

## Segurança

- **OAuth 2.0 com PKCE** — se a senha do proprietário estiver definida
- **Modo sem senha** — se nenhuma senha for configurada, executa sem autenticação
- **Contenção de caminho** — todas as operações de arquivo são validadas contra as raízes permitidas
- **Túnel opcional** — o Túnel Cloudflare protege contra exposição direta
- **Sem uploads para terceiros** — seu código nunca sai da sua máquina

---

## Compilando do Código Fonte

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compilar tudo (todas as plataformas)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compilar apenas para a plataforma atual
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Suporte a Plataformas

| Plataforma | Servidor | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compilar nativamente) |
| **macOS Intel** | ✅ | 🔧 (compilar nativamente) |
| **macOS M-chip** | ✅ | 🔧 (compilar nativamente) |

A GUI requer Fyne (OpenGL) — não pode ser compilada cruzadamente. O servidor compila em qualquer lugar.

---

## Estrutura do Projeto

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + servidor MCP
│   └── devspace-gui/       ← Configurador GUI de desktop (Fyne)
├── internal/
│   ├── auth/               ← Provedor OAuth 2.0 + PKCE
│   ├── config/             ← Sistema de configuração portátil
│   ├── locales/            ← Traduções para 47 idiomas
│   ├── logger/             ← Logging estruturado (zerolog)
│   ├── server/             ← HTTP + MCP + orquestração de túnel
│   ├── skills/             ← AGENTS.md / descoberta de habilidades
│   ├── store/              ← Sessões de espaço de trabalho SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Espaço de trabalho e validação de caminho
├── scripts/
│   ├── windows/            ← Script de build PowerShell
│   └── unix/               ← Scripts de build Bash + Makefile
├── readme/                 ← Traduções deste arquivo (47 idiomas)
├── build/                  ← Binários compilados (todas as plataformas)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Construído em Go. Zero npm. Zero Node.js. Um binário.
