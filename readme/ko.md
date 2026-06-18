# DevSpace (Go 에디션)

**ChatGPT와 Claude에게 로컬 머신에 대한 안전한 액세스를 제공하세요. 모든 MCP 호스트를 코딩 파트너로 만드세요.**

DevSpace는 AI 어시스턴트가 실제 로컬 프로젝트에서 파일을 읽고, 편집하고, 검색하고, 코드를 실행할 수 있게 해주는 자체 호스팅 MCP 서버입니다 — 여러분의 파일, 여러분의 도구, 여러분의 터미널 — 서드파티에 아무것도 업로드하지 않고요. 여러분의 머신에서 실행하고, 여러분이 제어하는 터널을 통해 노출하며, 선택적으로 비밀번호로 보호합니다.

---

## 목차

- [빠른 시작](#빠른-시작)
- [설치](#설치)
- [AI가 할 수 있는 일](#ai가-할-수-있는-일)
- [구성](#구성)
- [터널 (원격 액세스)](#터널-원격-액세스)
- [셸 지원](#셸-지원)
- [보안](#보안)
- [소스에서 빌드하기](#소스에서-빌드하기)
- [플랫폼 지원](#플랫폼-지원)
- [프로젝트 구조](#프로젝트-구조)

### 🌍 번역

| 언어 | | 언어 | | 언어 | |
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

## 빠른 시작

### 1. 다운로드
[Releases](../../releases)에서 플랫폼을 선택하거나 소스에서 빌드하세요:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. 구성 (GUI 또는 텍스트)
```bash
devspace-gui                  # 데스크톱 구성 도구 (GUI)
devspace init                 # 텍스트 기반 구성 도구
```

### 3. 실행
```bash
devspace                      # 서버 시작. 구성을 자동 감지합니다.
```

`tools/`에 `cloudflared`가 있으면 Cloudflare 터널도 자동으로 시작됩니다.

### 4. MCP 클라이언트 연결
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
또는 로컬에서: `http://127.0.0.1:7676/mcp`

---

## 설치

Node.js 불필요, npm 불필요, Python 불필요. 단일 바이너리.

| 플랫폼 | 다운로드 |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: 네이티브 컴파일) |
| **macOS Intel** | `devspace` (GUI: 네이티브 컴파일) |
| **macOS M-chip** | `devspace` (GUI: 네이티브 컴파일) |

소스에서 빌드할 때만 **Go 1.23+** 가 필요합니다.

---

## AI가 할 수 있는 일

연결되면 AI가 승인된 프로젝트 폴더 중 하나를 작업 공간으로 열 수 있습니다:

- 작업 공간 내 파일 **읽기, 쓰기 및 편집**
- 정규 표현식으로 **코드 검색** 및 디렉터리 검사
- **셸 명령 실행** (Windows에서는 PowerShell, Unix에서는 bash)
- `AGENTS.md` / `CLAUDE.md`에서 **프로젝트 지침 발견**
- 휴대용 `.devspace/config.json`으로 **자동 구성**

8개의 MCP 도구: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## 구성

모든 구성은 **실행 파일과 동일한 폴더**에 저장됩니다 (휴대용):

```
.devspace/
├── config.json       ← 허용된 루트, 포트, 셸, 언어, 인증
└── auth.json         ← 소유자 비밀번호 (선택 사항)
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

| 필드 | 기본값 | 설명 |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | OS에서 자동 감지. 47개 언어 지원 |
| `toolMode` | `full` | `full` (모든 도구) 또는 `minimal` (검색용 셸만) |
| `toolNaming` | `short` | `short` (read, write) 또는 `legacy` (read_file, write_file) |

환경 변수가 필요하지 않습니다 — 모든 것이 휴대용 구성 파일에 있습니다.

---

## 터널 (원격 액세스)

ChatGPT 웹 버전(HTTPS 필요)을 위해 DevSpace는 자동으로 터널을 시작합니다:

| 터널 | URL 유형 | 설정 |
|---|---|---|
| **Cloudflare** | 무작위 (자동) | `cloudflared.exe`를 `tools/`에 배치 |
| **Pinggy** | 안정적 | SSH 키 필요 (`ssh-keygen`) |

서버가 사용 가능한 것을 자동으로 감지합니다. 새로운 Cloudflare URL을 위해 서버를 다시 시작하거나, 영구 URL을 위해 Pinggy를 사용하세요.

---

## 셸 지원

| OS | 기본값 | 대안 |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / 모든 셸 |
| **macOS** | bash | `sh` / `zsh` |

config.json에서 `"shell"`을 설정하거나 GUI에서 선택하세요.

---

## 보안

- **OAuth 2.0 with PKCE** — 소유자 비밀번호가 설정된 경우
- **비밀번호 없는 모드** — 비밀번호가 구성되지 않은 경우 인증 없이 실행
- **경로 제한** — 모든 파일 작업이 허용된 루트에 대해 검증됨
- **선택적 터널** — Cloudflare 터널이 직접 노출로부터 보호
- **서드파티 업로드 없음** — 코드가 머신을 떠나지 않음

---

## 소스에서 빌드하기

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# 모두 빌드 (모든 플랫폼)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# 현재 플랫폼만 빌드
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## 플랫폼 지원

| 플랫폼 | 서버 | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (네이티브 컴파일) |
| **macOS Intel** | ✅ | 🔧 (네이티브 컴파일) |
| **macOS M-chip** | ✅ | 🔧 (네이티브 컴파일) |

GUI는 Fyne(OpenGL)이 필요 — 크로스 컴파일 불가. 서버는 어디서나 컴파일 가능합니다.

---

## 프로젝트 구조

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP 서버
│   └── devspace-gui/       ← 데스크톱 GUI 구성 도구 (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE 제공자
│   ├── config/             ← 휴대용 구성 시스템
│   ├── locales/            ← 47개 언어 번역
│   ├── logger/             ← 구조화된 로깅 (zerolog)
│   ├── server/             ← HTTP + MCP + 터널 오케스트레이션
│   ├── skills/             ← AGENTS.md / 스킬 발견
│   ├── store/              ← SQLite 작업 공간 세션
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← 작업 공간 및 경로 검증
├── scripts/
│   ├── windows/            ← PowerShell 빌드 스크립트
│   └── unix/               ← Bash + Makefile 빌드 스크립트
├── readme/                 ← 이 파일의 번역 (47개 언어)
├── build/                  ← 컴파일된 바이너리 (모든 플랫폼)
├── tools/                  ← cloudflared.exe 등
├── go.mod / go.sum
└── README.md
```

---

Go로 제작. npm 제로. Node.js 제로. 하나의 바이너리.
