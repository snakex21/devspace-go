# DevSpace (Go版)

**ChatGPTとClaudeにあなたのローカルマシンへの安全なアクセスを提供します。任意のMCPホストをコーディングパートナーに変えましょう。**

DevSpaceは自己ホスト型のMCPサーバーで、AIアシスタントが実際のローカルプロジェクト内でファイルの読み取り、編集、検索、コードの実行を行えるようにします — あなたのファイル、あなたのツール、あなたのターミナル — サードパーティに何もアップロードすることなく。あなたのマシン上で実行し、あなたが制御するトンネルを通じて公開し、オプションでパスワードで保護します。

---

## 目次

- [クイックスタート](#クイックスタート)
- [インストール](#インストール)
- [AIができること](#aiができること)
- [設定](#設定)
- [トンネル（リモートアクセス）](#トンネルリモートアクセス)
- [シェルサポート](#シェルサポート)
- [セキュリティ](#セキュリティ)
- [ソースからのビルド](#ソースからのビルド)
- [プラットフォームサポート](#プラットフォームサポート)
- [プロジェクト構造](#プロジェクト構造)

### 🌍 翻訳

| 言語 | | 言語 | | 言語 | |
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

## クイックスタート

### 1. ダウンロード
[Releases](../../releases)からお使いのプラットフォームを選択するか、ソースからビルドします：
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. 設定（GUIまたはテキスト）
```bash
devspace-gui                  # デスクトップ設定ツール（GUI）
devspace init                 # テキストベース設定ツール
```

### 3. 実行
```bash
devspace                      # サーバーを起動。設定を自動検出します。
```

`tools/`に`cloudflared`が見つかった場合、Cloudflareトンネルも自動的に開始します。

### 4. MCPクライアントを接続
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
またはローカル：`http://127.0.0.1:7676/mcp`

---

## インストール

Node.js不要、npm不要、Python不要。シングルバイナリです。

| プラットフォーム | ダウンロード |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace`（GUI：ネイティブコンパイル） |
| **macOS Intel** | `devspace`（GUI：ネイティブコンパイル） |
| **macOS M-chip** | `devspace`（GUI：ネイティブコンパイル） |

ソースからビルドする場合のみ **Go 1.23以上** が必要です。

---

## AIができること

接続後、AIは承認されたプロジェクトフォルダの1つをワークスペースとして開くことができます：

- ワークスペース内のファイルの **読み取り、書き込み、編集**
- 正規表現による **コード検索** とディレクトリの検査
- **シェルコマンドの実行**（WindowsではPowerShell、Unixではbash）
- `AGENTS.md` / `CLAUDE.md`からの **プロジェクト指示の発見**
- ポータブルな `.devspace/config.json` による **自動設定**

8つのMCPツール：`open_workspace`、`read`、`write`、`edit`、`grep`、`glob`、`ls`、`bash`

---

## 設定

すべての設定は **実行ファイルと同じフォルダ** に保存されます（ポータブル）：

```
.devspace/
├── config.json       ← 許可されたルート、ポート、シェル、言語、認証
└── auth.json         ← 所有者パスワード（オプション）
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

| フィールド | デフォルト | 説明 |
|---|---|---|
| `shell` | `auto` | `auto`、`powershell`、`cmd`、`bash`、`sh` |
| `lang` | `auto` | OSから自動検出。47言語をサポート |
| `toolMode` | `full` | `full`（すべてのツール）または `minimal`（検索用シェルのみ） |
| `toolNaming` | `short` | `short`（read, write）または `legacy`（read_file, write_file） |

環境変数は不要 — すべてがポータブルな設定ファイルに含まれています。

---

## トンネル（リモートアクセス）

ChatGPT Web版（HTTPS必須）向けに、DevSpaceは自動的にトンネルを開始します：

| トンネル | URLタイプ | セットアップ |
|---|---|---|
| **Cloudflare** | ランダム（自動） | `cloudflared.exe` を `tools/` に配置 |
| **Pinggy** | 安定 | SSHキーが必要（`ssh-keygen`） |

サーバーは利用可能なものを自動検出します。新しいCloudflare URLが必要な場合はサーバーを再起動するか、固定URLにはPinggyを使用してください。

---

## シェルサポート

| OS | デフォルト | 代替 |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / 任意のシェル |
| **macOS** | bash | `sh` / `zsh` |

config.jsonで `"shell"` を設定するか、GUIで選択してください。

---

## セキュリティ

- **OAuth 2.0 with PKCE** — 所有者パスワードが設定されている場合
- **パスワードレスモード** — パスワードが設定されていない場合、認証なしで実行
- **パス制限** — すべてのファイル操作が許可されたルートに対して検証されます
- **オプショナルトンネル** — Cloudflareトンネルが直接の露出から保護
- **サードパーティへのアップロードなし** — コードがマシンから出ることはありません

---

## ソースからのビルド

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# すべてをビルド（全プラットフォーム）
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# 現在のプラットフォーム向けのみビルド
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## プラットフォームサポート

| プラットフォーム | サーバー | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧（ネイティブコンパイル） |
| **macOS Intel** | ✅ | 🔧（ネイティブコンパイル） |
| **macOS M-chip** | ✅ | 🔧（ネイティブコンパイル） |

GUIはFyne（OpenGL）が必要 — クロスコンパイル不可。サーバーはどこでもコンパイル可能です。

---

## プロジェクト構造

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCPサーバー
│   └── devspace-gui/       ← デスクトップGUI設定ツール（Fyne）
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE プロバイダー
│   ├── config/             ← ポータブル設定システム
│   ├── locales/            ← 47言語の翻訳
│   ├── logger/             ← 構造化ログ（zerolog）
│   ├── server/             ← HTTP + MCP + トンネルオーケストレーション
│   ├── skills/             ← AGENTS.md / スキル発見
│   ├── store/              ← SQLite ワークスペースセッション
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← ワークスペースとパス検証
├── scripts/
│   ├── windows/            ← PowerShellビルドスクリプト
│   └── unix/               ← Bash + Makefile ビルドスクリプト
├── readme/                 ← このファイルの翻訳（47言語）
├── build/                  ← コンパイル済みバイナリ（全プラットフォーム）
├── tools/                  ← cloudflared.exe など
├── go.mod / go.sum
└── README.md
```

---

Goで構築。npmゼロ。Node.jsゼロ。1つのバイナリ。
