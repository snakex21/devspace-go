# DevSpace (Edisi Go)

**Berikan ChatGPT & Claude akses aman ke mesin lokal Anda. Ubah host MCP mana pun menjadi mitra koding Anda.**

DevSpace adalah server MCP yang di-host sendiri yang memungkinkan asisten AI membaca, mengedit, mencari, dan menjalankan kode di proyek lokal nyata Anda — file Anda, alat Anda, terminal Anda — tanpa mengunggah apa pun ke pihak ketiga. Anda menjalankannya di mesin Anda, mengeksposnya melalui terowongan yang Anda kendalikan, dan secara opsional mengamankannya dengan kata sandi.

---

## Daftar Isi

- [Mulai Cepat](#mulai-cepat)
- [Instalasi](#instalasi)
- [Apa yang Dapat Dilakukan AI](#apa-yang-dapat-dilakukan-ai)
- [Konfigurasi](#konfigurasi)
- [Terowongan (Akses Jarak Jauh)](#terowongan-akses-jarak-jauh)
- [Dukungan Shell](#dukungan-shell)
- [Keamanan](#keamanan)
- [Membangun dari Sumber](#membangun-dari-sumber)
- [Dukungan Platform](#dukungan-platform)
- [Struktur Proyek](#struktur-proyek)

### 🌍 Terjemahan

| Bahasa | | Bahasa | | Bahasa | |
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

## Mulai Cepat

### 1. Unduh
Pilih platform Anda dari [Releases](../../releases) atau bangun dari sumber:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurasi (GUI atau teks)
```bash
devspace-gui                  # Konfigurator desktop (GUI)
devspace init                 # Konfigurator berbasis teks
```

### 3. Jalankan
```bash
devspace                      # Memulai server. Otomatis mendeteksi konfigurasi.
```

Ini juga otomatis memulai Terowongan Cloudflare jika `cloudflared` ditemukan di `tools/`.

### 4. Hubungkan klien MCP Anda
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Atau secara lokal: `http://127.0.0.1:7676/mcp`

---

## Instalasi

Tanpa Node.js, tanpa npm, tanpa Python. Biner tunggal.

| Platform | Unduhan |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompilasi native) |
| **macOS Intel** | `devspace` (GUI: kompilasi native) |
| **macOS M-chip** | `devspace` (GUI: kompilasi native) |

Memerlukan **Go 1.23+** hanya jika membangun dari sumber.

---

## Apa yang Dapat Dilakukan AI

Setelah terhubung, AI dapat membuka salah satu folder proyek yang disetujui sebagai ruang kerja:

- **Membaca, menulis, dan mengedit** file di dalam ruang kerja
- **Mencari kode** dengan regex dan memeriksa direktori
- **Menjalankan perintah shell** (PowerShell di Windows, bash di Unix)
- **Menemukan instruksi proyek** dari `AGENTS.md` / `CLAUDE.md`
- **Konfigurasi otomatis** dengan `.devspace/config.json` portabel

8 alat MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurasi

Semua konfigurasi berada **di folder yang sama dengan file executable** (portabel):

```
.devspace/
├── config.json       ← root yang diizinkan, port, shell, bahasa, auth
└── auth.json         ← kata sandi pemilik (opsional)
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

| Kolom | Default | Deskripsi |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Deteksi otomatis dari OS. Mendukung 47 bahasa |
| `toolMode` | `full` | `full` (semua alat) atau `minimal` (hanya shell untuk pencarian) |
| `toolNaming` | `short` | `short` (read, write) atau `legacy` (read_file, write_file) |

Tidak perlu variabel lingkungan — semuanya ada dalam file konfigurasi portabel.

---

## Terowongan (Akses Jarak Jauh)

Untuk ChatGPT versi web (memerlukan HTTPS), DevSpace otomatis memulai terowongan:

| Terowongan | Tipe URL | Pengaturan |
|---|---|---|
| **Cloudflare** | Acak (otomatis) | Letakkan `cloudflared.exe` di `tools/` |
| **Pinggy** | Stabil | Memerlukan kunci SSH (`ssh-keygen`) |

Server otomatis mendeteksi mana yang tersedia. Mulai ulang server untuk URL Cloudflare baru, atau gunakan Pinggy untuk URL permanen.

---

## Dukungan Shell

| OS | Default | Alternatif |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / shell apa pun |
| **macOS** | bash | `sh` / `zsh` |

Atur `"shell"` di config.json atau pilih di GUI.

---

## Keamanan

- **OAuth 2.0 dengan PKCE** — jika kata sandi pemilik diatur
- **Mode tanpa kata sandi** — jika tidak ada kata sandi yang dikonfigurasi, berjalan tanpa auth
- **Pembatasan jalur** — semua operasi file divalidasi terhadap root yang diizinkan
- **Terowongan opsional** — Terowongan Cloudflare melindungi dari paparan langsung
- **Tanpa unggahan pihak ketiga** — kode Anda tidak pernah meninggalkan mesin Anda

---

## Membangun dari Sumber

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bangun semuanya (semua platform)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bangun hanya untuk platform saat ini
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Dukungan Platform

| Platform | Server | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompilasi native) |
| **macOS Intel** | ✅ | 🔧 (kompilasi native) |
| **macOS M-chip** | ✅ | 🔧 (kompilasi native) |

GUI memerlukan Fyne (OpenGL) — tidak dapat kompilasi silang. Server dapat dikompilasi di mana saja.

---

## Struktur Proyek

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + server MCP
│   └── devspace-gui/       ← Konfigurator GUI desktop (Fyne)
├── internal/
│   ├── auth/               ← Penyedia OAuth 2.0 + PKCE
│   ├── config/             ← Sistem konfigurasi portabel
│   ├── locales/            ← Terjemahan 47 bahasa
│   ├── logger/             ← Pencatatan terstruktur (zerolog)
│   ├── server/             ← HTTP + MCP + orkestrasi terowongan
│   ├── skills/             ← AGENTS.md / penemuan skill
│   ├── store/              ← Sesi ruang kerja SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Ruang kerja & validasi jalur
├── scripts/
│   ├── windows/            ← Skrip build PowerShell
│   └── unix/               ← Skrip build Bash + Makefile
├── readme/                 ← Terjemahan file ini (47 bahasa)
├── build/                  ← Biner yang dikompilasi (semua platform)
├── tools/                  ← cloudflared.exe, dll.
├── go.mod / go.sum
└── README.md
```

---

Dibangun dengan Go. Nol npm. Nol Node.js. Satu biner.
