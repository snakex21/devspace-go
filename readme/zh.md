# DevSpace（Go 版）

**让 ChatGPT 和 Claude 安全访问您的本地计算机。将任何 MCP 主机变成您的编程伙伴。**

DevSpace 是一个自托管的 MCP 服务器，让 AI 助手能够在您真实的本地项目中读取、编辑、搜索和运行代码——您的文件、您的工具、您的终端——无需上传任何内容给第三方。您在机器上运行它，通过您控制的隧道公开它，并可选择用密码保护它。

---

## 目录

- [快速入门](#快速入门)
- [安装](#安装)
- [AI 能做什么](#ai-能做什么)
- [配置](#配置)
- [隧道（远程访问）](#隧道远程访问)
- [Shell 支持](#shell-支持)
- [安全性](#安全性)
- [从源代码构建](#从源代码构建)
- [平台支持](#平台支持)
- [项目结构](#项目结构)

### 🌍 翻译

| 语言 | | 语言 | | 语言 | |
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

## 快速入门

### 1. 下载
从 [Releases](../../releases) 选择您的平台或从源代码构建：
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. 配置（GUI 或文本）
```bash
devspace-gui                  # 桌面配置器（GUI）
devspace init                 # 基于文本的配置器
```

### 3. 运行
```bash
devspace                      # 启动服务器。自动检测配置。
```

如果在 `tools/` 中找到 `cloudflared`，这也会自动启动 Cloudflare 隧道。

### 4. 连接您的 MCP 客户端
```
https://您的隧道.trycloudflare.com/mcp
```
或本地：`http://127.0.0.1:7676/mcp`

---

## 安装

无需 Node.js、无需 npm、无需 Python。单一二进制文件。

| 平台 | 下载 |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace`（GUI：请原生编译） |
| **macOS Intel** | `devspace`（GUI：请原生编译） |
| **macOS M-chip** | `devspace`（GUI：请原生编译） |

仅在从源代码构建时需要 **Go 1.23+**。

---

## AI 能做什么

连接后，AI 可以打开您已批准的项目文件夹之一作为工作区：

- 在工作区中**读取、写入和编辑**文件
- 使用正则表达式**搜索代码**并检查目录
- **运行 Shell 命令**（Windows 上使用 PowerShell，Unix 上使用 bash）
- 从 `AGENTS.md` / `CLAUDE.md` **发现项目指令**
- 使用便携式 `.devspace/config.json` **自动配置**

8 个 MCP 工具：`open_workspace`、`read`、`write`、`edit`、`grep`、`glob`、`ls`、`bash`

---

## 配置

所有配置都位于**可执行文件所在的同一文件夹中**（便携式）：

```
.devspace/
├── config.json       ← 允许的根目录、端口、Shell、语言、认证
└── auth.json         ← 所有者密码（可选）
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

| 字段 | 默认值 | 描述 |
|---|---|---|
| `shell` | `auto` | `auto`、`powershell`、`cmd`、`bash`、`sh` |
| `lang` | `auto` | 从操作系统自动检测。支持 47 种语言 |
| `toolMode` | `full` | `full`（所有工具）或 `minimal`（仅 Shell 用于搜索） |
| `toolNaming` | `short` | `short`（read、write）或 `legacy`（read_file、write_file） |

无需环境变量——一切都在便携式配置文件中。

---

## 隧道（远程访问）

对于 ChatGPT 网页版（需要 HTTPS），DevSpace 会自动启动隧道：

| 隧道 | URL 类型 | 设置 |
|---|---|---|
| **Cloudflare** | 随机（自动） | 将 `cloudflared.exe` 放入 `tools/` |
| **Pinggy** | 稳定 | 需要 SSH 密钥（`ssh-keygen`） |

服务器会自动检测哪一个可用。重启服务器以获取新的 Cloudflare URL，或使用 Pinggy 获取永久 URL。

---

## Shell 支持

| 操作系统 | 默认 | 替代方案 |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / 任何 Shell |
| **macOS** | bash | `sh` / `zsh` |

在 config.json 中设置 `"shell"` 或在 GUI 中选择。

---

## 安全性

- **带 PKCE 的 OAuth 2.0**——如果设置了所有者密码
- **无密码模式**——如果未配置密码，则无需认证即可运行
- **路径限制**——所有文件操作都根据允许的根目录进行验证
- **可选隧道**——Cloudflare 隧道可防止直接暴露
- **无第三方上传**——您的代码永远不会离开您的机器

---

## 从源代码构建

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# 构建所有内容（所有平台）
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# 仅为当前平台构建
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## 平台支持

| 平台 | 服务器 | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧（请原生编译） |
| **macOS Intel** | ✅ | 🔧（请原生编译） |
| **macOS M-chip** | ✅ | 🔧（请原生编译） |

GUI 需要 Fyne（OpenGL）——无法交叉编译。服务器可在任何地方编译。

---

## 项目结构

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP 服务器
│   └── devspace-gui/       ← 桌面 GUI 配置器（Fyne）
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE 提供程序
│   ├── config/             ← 便携式配置系统
│   ├── locales/            ← 47 种语言翻译
│   ├── logger/             ← 结构化日志记录（zerolog）
│   ├── server/             ← HTTP + MCP + 隧道编排
│   ├── skills/             ← AGENTS.md / 技能发现
│   ├── store/              ← SQLite 工作区会话
│   ├── tools/              ← read、write、edit、grep、glob、ls、bash
│   └── workspace/          ← 工作区和路径验证
├── scripts/
│   ├── windows/            ← PowerShell 构建脚本
│   └── unix/               ← Bash + Makefile 构建脚本
├── readme/                 ← 此文件的翻译（47 种语言）
├── build/                  ← 已编译的二进制文件（所有平台）
├── tools/                  ← cloudflared.exe 等
├── go.mod / go.sum
└── README.md
```

---

使用 Go 构建。零 npm。零 Node.js。一个二进制文件。
