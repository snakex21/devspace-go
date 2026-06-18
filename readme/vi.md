# DevSpace (Phiên bản Go)

**Cấp cho ChatGPT và Claude quyền truy cập an toàn vào máy tính cục bộ của bạn. Biến bất kỳ máy chủ MCP nào thành đối tác lập trình của bạn.**

DevSpace là một máy chủ MCP tự lưu trữ cho phép trợ lý AI đọc, chỉnh sửa, tìm kiếm và chạy mã trong các dự án cục bộ thực tế của bạn — tệp của bạn, công cụ của bạn, thiết bị đầu cuối của bạn — mà không tải lên bất cứ thứ gì cho bên thứ ba. Bạn chạy nó trên máy của mình, phơi bày nó qua một đường hầm mà bạn kiểm soát và tùy chọn bảo mật bằng mật khẩu.

---

## Mục Lục

- [Bắt Đầu Nhanh](#bắt-đầu-nhanh)
- [Cài Đặt](#cài-đặt)
- [AI Có Thể Làm Gì](#ai-có-thể-làm-gì)
- [Cấu Hình](#cấu-hình)
- [Đường Hầm (Truy Cập Từ Xa)](#đường-hầm-truy-cập-từ-xa)
- [Hỗ Trợ Shell](#hỗ-trợ-shell)
- [Bảo Mật](#bảo-mật)
- [Xây Dựng từ Mã Nguồn](#xây-dựng-từ-mã-nguồn)
- [Hỗ Trợ Nền Tảng](#hỗ-trợ-nền-tảng)
- [Cấu Trúc Dự Án](#cấu-trúc-dự-án)

### 🌍 Bản Dịch

| Ngôn ngữ | | Ngôn ngữ | | Ngôn ngữ | |
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

## Bắt Đầu Nhanh

### 1. Tải Xuống
Chọn nền tảng của bạn từ [Releases](../../releases) hoặc xây dựng từ mã nguồn:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Cấu Hình (GUI hoặc văn bản)
```bash
devspace-gui                  # Trình cấu hình máy tính để bàn (GUI)
devspace init                 # Trình cấu hình dựa trên văn bản
```

### 3. Chạy
```bash
devspace                      # Khởi động máy chủ. Tự động phát hiện cấu hình.
```

Điều này cũng tự động khởi động Đường hầm Cloudflare nếu `cloudflared` được tìm thấy trong `tools/`.

### 4. Kết nối máy khách MCP của bạn
```
https://ĐƯỜNG-HẦM-CỦA-BẠN.trycloudflare.com/mcp
```
Hoặc cục bộ: `http://127.0.0.1:7676/mcp`

---

## Cài Đặt

Không Node.js, không npm, không Python. Một tệp nhị phân duy nhất.

| Nền tảng | Tải xuống |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: biên dịch tự nhiên) |
| **macOS Intel** | `devspace` (GUI: biên dịch tự nhiên) |
| **macOS M-chip** | `devspace` (GUI: biên dịch tự nhiên) |

Yêu cầu **Go 1.23+** chỉ khi xây dựng từ mã nguồn.

---

## AI Có Thể Làm Gì

Sau khi kết nối, AI có thể mở một trong các thư mục dự án đã được phê duyệt của bạn làm không gian làm việc:

- **Đọc, ghi và chỉnh sửa** tệp trong không gian làm việc
- **Tìm kiếm mã** bằng regex và kiểm tra thư mục
- **Chạy lệnh shell** (PowerShell trên Windows, bash trên Unix)
- **Khám phá hướng dẫn dự án** từ `AGENTS.md` / `CLAUDE.md`
- **Tự động cấu hình** với `.devspace/config.json` di động

8 công cụ MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Cấu Hình

Tất cả cấu hình nằm **trong cùng thư mục với tệp thực thi** (di động):

```
.devspace/
├── config.json       ← thư mục gốc được phép, cổng, shell, ngôn ngữ, xác thực
└── auth.json         ← mật khẩu chủ sở hữu (tùy chọn)
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

| Trường | Mặc định | Mô tả |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Tự động phát hiện từ HĐH. Hỗ trợ 47 ngôn ngữ |
| `toolMode` | `full` | `full` (tất cả công cụ) hoặc `minimal` (chỉ shell để tìm kiếm) |
| `toolNaming` | `short` | `short` (read, write) hoặc `legacy` (read_file, write_file) |

Không cần biến môi trường — mọi thứ đều nằm trong tệp cấu hình di động.

---

## Đường Hầm (Truy Cập Từ Xa)

Đối với phiên bản web ChatGPT (yêu cầu HTTPS), DevSpace tự động khởi động một đường hầm:

| Đường hầm | Loại URL | Thiết lập |
|---|---|---|
| **Cloudflare** | Ngẫu nhiên (tự động) | Đặt `cloudflared.exe` vào `tools/` |
| **Pinggy** | Ổn định | Cần khóa SSH (`ssh-keygen`) |

Máy chủ tự động phát hiện cái nào khả dụng. Khởi động lại máy chủ để có URL Cloudflare mới hoặc sử dụng Pinggy để có URL vĩnh viễn.

---

## Hỗ Trợ Shell

| HĐH | Mặc định | Lựa chọn thay thế |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / bất kỳ shell nào |
| **macOS** | bash | `sh` / `zsh` |

Đặt `"shell"` trong config.json hoặc chọn trong GUI.

---

## Bảo Mật

- **OAuth 2.0 với PKCE** — nếu mật khẩu chủ sở hữu được đặt
- **Chế độ không mật khẩu** — nếu không có mật khẩu được cấu hình, chạy mà không cần xác thực
- **Giới hạn đường dẫn** — tất cả thao tác tệp được xác thực dựa trên thư mục gốc được phép
- **Đường hầm tùy chọn** — Đường hầm Cloudflare bảo vệ khỏi phơi nhiễm trực tiếp
- **Không tải lên bên thứ ba** — mã của bạn không bao giờ rời khỏi máy của bạn

---

## Xây Dựng từ Mã Nguồn

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Xây dựng mọi thứ (tất cả nền tảng)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Chỉ xây dựng cho nền tảng hiện tại
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Hỗ Trợ Nền Tảng

| Nền tảng | Máy chủ | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (biên dịch tự nhiên) |
| **macOS Intel** | ✅ | 🔧 (biên dịch tự nhiên) |
| **macOS M-chip** | ✅ | 🔧 (biên dịch tự nhiên) |

GUI yêu cầu Fyne (OpenGL) — không thể biên dịch chéo. Máy chủ biên dịch ở mọi nơi.

---

## Cấu Trúc Dự Án

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + máy chủ MCP
│   └── devspace-gui/       ← Trình cấu hình GUI máy tính để bàn (Fyne)
├── internal/
│   ├── auth/               ← Nhà cung cấp OAuth 2.0 + PKCE
│   ├── config/             ← Hệ thống cấu hình di động
│   ├── locales/            ← Bản dịch 47 ngôn ngữ
│   ├── logger/             ← Ghi nhật ký có cấu trúc (zerolog)
│   ├── server/             ← HTTP + MCP + điều phối đường hầm
│   ├── skills/             ← Khám phá AGENTS.md / kỹ năng
│   ├── store/              ← Phiên không gian làm việc SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Xác thực không gian làm việc và đường dẫn
├── scripts/
│   ├── windows/            ← Tập lệnh xây dựng PowerShell
│   └── unix/               ← Tập lệnh xây dựng Bash + Makefile
├── readme/                 ← Bản dịch của tệp này (47 ngôn ngữ)
├── build/                  ← Tệp nhị phân đã biên dịch (tất cả nền tảng)
├── tools/                  ← cloudflared.exe, v.v.
├── go.mod / go.sum
└── README.md
```

---

Được xây dựng bằng Go. Không npm. Không Node.js. Một tệp nhị phân.
