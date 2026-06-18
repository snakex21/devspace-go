# DevSpace (รุ่น Go)

**ให้ ChatGPT และ Claude เข้าถึงเครื่องท้องถิ่นของคุณอย่างปลอดภัย เปลี่ยนโฮสต์ MCP ใดๆ ให้เป็นคู่หูเขียนโค้ดของคุณ**

DevSpace คือเซิร์ฟเวอร์ MCP แบบโฮสต์เองที่ให้ผู้ช่วย AI อ่าน แก้ไข ค้นหา และรันโค้ดในโปรเจกต์ท้องถิ่นจริงของคุณ — ไฟล์ของคุณ เครื่องมือของคุณ เทอร์มินัลของคุณ — โดยไม่ต้องอัปโหลดอะไรไปยังบุคคลที่สาม คุณรันมันบนเครื่องของคุณ เปิดเผยผ่านทันเนลที่คุณควบคุม และป้องกันด้วยรหัสผ่านได้ตามต้องการ

---

## สารบัญ

- [เริ่มต้นอย่างรวดเร็ว](#เริ่มต้นอย่างรวดเร็ว)
- [การติดตั้ง](#การติดตั้ง)
- [สิ่งที่ AI ทำได้](#สิ่งที่-ai-ทำได้)
- [การกำหนดค่า](#การกำหนดค่า)
- [ทันเนล (การเข้าถึงระยะไกล)](#ทันเนล-การเข้าถึงระยะไกล)
- [การรองรับเชลล์](#การรองรับเชลล์)
- [ความปลอดภัย](#ความปลอดภัย)
- [สร้างจากซอร์สโค้ด](#สร้างจากซอร์สโค้ด)
- [การรองรับแพลตฟอร์ม](#การรองรับแพลตฟอร์ม)
- [โครงสร้างโปรเจกต์](#โครงสร้างโปรเจกต์)

### 🌍 การแปลภาษา

| ภาษา | | ภาษา | | ภาษา | |
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

## เริ่มต้นอย่างรวดเร็ว

### 1. ดาวน์โหลด
เลือกแพลตฟอร์มของคุณจาก [Releases](../../releases) หรือสร้างจากซอร์สโค้ด:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. กำหนดค่า (GUI หรือข้อความ)
```bash
devspace-gui                  # ตัวกำหนดค่าเดสก์ท็อป (GUI)
devspace init                 # ตัวกำหนดค่าแบบข้อความ
```

### 3. รัน
```bash
devspace                      # เริ่มเซิร์ฟเวอร์ ตรวจจับการกำหนดค่าโดยอัตโนมัติ
```

นอกจากนี้ยังเริ่มทันเนล Cloudflare โดยอัตโนมัติหากพบ `cloudflared` ใน `tools/`

### 4. เชื่อมต่อไคลเอนต์ MCP ของคุณ
```
https://ทันเนลของคุณ.trycloudflare.com/mcp
```
หรือในเครื่อง: `http://127.0.0.1:7676/mcp`

---

## การติดตั้ง

ไม่มี Node.js ไม่มี npm ไม่มี Python ไบนารีไฟล์เดียว

| แพลตฟอร์ม | ดาวน์โหลด |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: คอมไพล์แบบเนทีฟ) |
| **macOS Intel** | `devspace` (GUI: คอมไพล์แบบเนทีฟ) |
| **macOS M-chip** | `devspace` (GUI: คอมไพล์แบบเนทีฟ) |

ต้องการ **Go 1.23+** เฉพาะเมื่อสร้างจากซอร์สโค้ด

---

## สิ่งที่ AI ทำได้

เมื่อเชื่อมต่อแล้ว AI สามารถเปิดโฟลเดอร์โปรเจกต์ที่ได้รับอนุมัติของคุณเป็นพื้นที่ทำงาน:

- **อ่าน เขียน และแก้ไข** ไฟล์ภายในพื้นที่ทำงาน
- **ค้นหาโค้ด** ด้วย regex และตรวจสอบไดเรกทอรี
- **รันคำสั่งเชลล์** (PowerShell บน Windows, bash บน Unix)
- **ค้นพบคำแนะนำโปรเจกต์** จาก `AGENTS.md` / `CLAUDE.md`
- **กำหนดค่าอัตโนมัติ** ด้วย `.devspace/config.json` แบบพกพา

8 เครื่องมือ MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## การกำหนดค่า

การกำหนดค่าทั้งหมดอยู่ใน **โฟลเดอร์เดียวกับไฟล์ปฏิบัติการ** (พกพา):

```
.devspace/
├── config.json       ← ไดเรกทอรีที่อนุญาต, พอร์ต, เชลล์, ภาษา, การยืนยันตัวตน
└── auth.json         ← รหัสผ่านเจ้าของ (ไม่บังคับ)
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

| ฟิลด์ | ค่าเริ่มต้น | คำอธิบาย |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | ตรวจจับอัตโนมัติจาก OS รองรับ 47 ภาษา |
| `toolMode` | `full` | `full` (เครื่องมือทั้งหมด) หรือ `minimal` (เฉพาะเชลล์สำหรับการค้นหา) |
| `toolNaming` | `short` | `short` (read, write) หรือ `legacy` (read_file, write_file) |

ไม่จำเป็นต้องใช้ตัวแปรสภาพแวดล้อม — ทุกอย่างอยู่ในไฟล์กำหนดค่าแบบพกพา

---

## ทันเนล (การเข้าถึงระยะไกล)

สำหรับ ChatGPT เวอร์ชันเว็บ (ต้องใช้ HTTPS) DevSpace จะเริ่มทันเนลโดยอัตโนมัติ:

| ทันเนล | ประเภท URL | การตั้งค่า |
|---|---|---|
| **Cloudflare** | สุ่ม (อัตโนมัติ) | วาง `cloudflared.exe` ใน `tools/` |
| **Pinggy** | เสถียร | ต้องการคีย์ SSH (`ssh-keygen`) |

เซิร์ฟเวอร์ตรวจจับโดยอัตโนมัติว่ามีตัวไหนพร้อมใช้งาน รีสตาร์ทเซิร์ฟเวอร์เพื่อรับ URL Cloudflare ใหม่ หรือใช้ Pinggy สำหรับ URL ถาวร

---

## การรองรับเชลล์

| OS | ค่าเริ่มต้น | ทางเลือก |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / เชลล์ใดๆ |
| **macOS** | bash | `sh` / `zsh` |

ตั้งค่า `"shell"` ใน config.json หรือเลือกใน GUI

---

## ความปลอดภัย

- **OAuth 2.0 พร้อม PKCE** — หากตั้งรหัสผ่านเจ้าของ
- **โหมดไม่ใช้รหัสผ่าน** — หากไม่ได้กำหนดค่ารหัสผ่าน จะทำงานโดยไม่มีการยืนยันตัวตน
- **การจำกัดเส้นทาง** — การดำเนินการไฟล์ทั้งหมดถูกตรวจสอบกับไดเรกทอรีที่อนุญาต
- **ทันเนลเสริม** — ทันเนล Cloudflare ป้องกันการเปิดเผยโดยตรง
- **ไม่มีการอัปโหลดไปยังบุคคลที่สาม** — โค้ดของคุณไม่เคยออกจากเครื่องของคุณ

---

## สร้างจากซอร์สโค้ด

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# สร้างทุกอย่าง (ทุกแพลตฟอร์ม)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# สร้างเฉพาะแพลตฟอร์มปัจจุบัน
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## การรองรับแพลตฟอร์ม

| แพลตฟอร์ม | เซิร์ฟเวอร์ | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (คอมไพล์แบบเนทีฟ) |
| **macOS Intel** | ✅ | 🔧 (คอมไพล์แบบเนทีฟ) |
| **macOS M-chip** | ✅ | 🔧 (คอมไพล์แบบเนทีฟ) |

GUI ต้องการ Fyne (OpenGL) — ไม่สามารถคอมไพล์ข้ามแพลตฟอร์มได้ เซิร์ฟเวอร์คอมไพล์ได้ทุกที่

---

## โครงสร้างโปรเจกต์

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + เซิร์ฟเวอร์ MCP
│   └── devspace-gui/       ← ตัวกำหนดค่า GUI เดสก์ท็อป (Fyne)
├── internal/
│   ├── auth/               ← ผู้ให้บริการ OAuth 2.0 + PKCE
│   ├── config/             ← ระบบกำหนดค่าแบบพกพา
│   ├── locales/            ← การแปล 47 ภาษา
│   ├── logger/             ← การบันทึกแบบมีโครงสร้าง (zerolog)
│   ├── server/             ← HTTP + MCP + การจัดการทันเนล
│   ├── skills/             ← การค้นพบ AGENTS.md / ทักษะ
│   ├── store/              ← เซสชันพื้นที่ทำงาน SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← การตรวจสอบพื้นที่ทำงานและเส้นทาง
├── scripts/
│   ├── windows/            ← สคริปต์สร้าง PowerShell
│   └── unix/               ← สคริปต์สร้าง Bash + Makefile
├── readme/                 ← การแปลไฟล์นี้ (47 ภาษา)
├── build/                  ← ไบนารีที่คอมไพล์แล้ว (ทุกแพลตฟอร์ม)
├── tools/                  ← cloudflared.exe ฯลฯ
├── go.mod / go.sum
└── README.md
```

---

สร้างด้วย Go ไม่มี npm ไม่มี Node.js ไบนารีไฟล์เดียว
