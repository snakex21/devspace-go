# DevSpace (Edisi Go)

**Berikan ChatGPT & Claude akses selamat ke mesin tempatan anda. Jadikan mana-mana hos MCP sebagai rakan pengekodan anda.**

DevSpace ialah pelayan MCP yang dihos sendiri yang membolehkan pembantu AI membaca, mengedit, mencari dan menjalankan kod dalam projek tempatan sebenar anda — fail anda, alatan anda, terminal anda — tanpa memuat naik apa-apa kepada pihak ketiga. Anda jalankannya pada mesin anda, dedahkannya melalui terowong yang anda kawal, dan secara pilihan melindunginya dengan kata laluan.

---

## Isi Kandungan

- [Permulaan Pantas](#permulaan-pantas)
- [Pemasangan](#pemasangan)
- [Apa yang AI Boleh Lakukan](#apa-yang-ai-boleh-lakukan)
- [Konfigurasi](#konfigurasi)
- [Terowong (Akses Jauh)](#terowong-akses-jauh)
- [Sokongan Shell](#sokongan-shell)
- [Keselamatan](#keselamatan)
- [Membina dari Sumber](#membina-dari-sumber)
- [Sokongan Platform](#sokongan-platform)
- [Struktur Projek](#struktur-projek)

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

## Permulaan Pantas

### 1. Muat Turun
Pilih platform anda dari [Releases](../../releases) atau bina dari sumber:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurasi (GUI atau teks)
```bash
devspace-gui                  # Konfigurator desktop (GUI)
devspace init                 # Konfigurator berasaskan teks
```

### 3. Jalankan
```bash
devspace                      # Mulakan pelayan. Auto-kesan konfigurasi.
```

Ini juga auto-mulakan Terowong Cloudflare jika `cloudflared` dijumpai dalam `tools/`.

### 4. Sambungkan klien MCP anda
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Atau secara tempatan: `http://127.0.0.1:7676/mcp`

---

## Pemasangan

Tiada Node.js, tiada npm, tiada Python. Satu binari tunggal.

| Platform | Muat Turun |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompil secara natif) |
| **macOS Intel** | `devspace` (GUI: kompil secara natif) |
| **macOS M-chip** | `devspace` (GUI: kompil secara natif) |

Memerlukan **Go 1.23+** hanya jika membina dari sumber.

---

## Apa yang AI Boleh Lakukan

Setelah disambungkan, AI boleh membuka salah satu folder projek yang diluluskan sebagai ruang kerja:

- **Membaca, menulis dan mengedit** fail dalam ruang kerja
- **Mencari kod** dengan regex dan memeriksa direktori
- **Menjalankan arahan shell** (PowerShell pada Windows, bash pada Unix)
- **Menemui arahan projek** dari `AGENTS.md` / `CLAUDE.md`
- **Auto-konfigurasi** dengan `.devspace/config.json` mudah alih

8 alatan MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfigurasi

Semua konfigurasi berada **dalam folder yang sama dengan fail boleh laku** (mudah alih):

```
.devspace/
├── config.json       ← root dibenarkan, port, shell, bahasa, pengesahan
└── auth.json         ← kata laluan pemilik (pilihan)
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

| Medan | Lalai | Penerangan |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Auto-kesan dari OS. Menyokong 47 bahasa |
| `toolMode` | `full` | `full` (semua alatan) atau `minimal` (shell sahaja untuk carian) |
| `toolNaming` | `short` | `short` (read, write) atau `legacy` (read_file, write_file) |

Tiada pembolehubah persekitaran diperlukan — semuanya dalam fail konfigurasi mudah alih.

---

## Terowong (Akses Jauh)

Untuk ChatGPT versi web (HTTPS diperlukan), DevSpace auto-mulakan terowong:

| Terowong | Jenis URL | Persediaan |
|---|---|---|
| **Cloudflare** | Rawak (auto) | Letakkan `cloudflared.exe` dalam `tools/` |
| **Pinggy** | Stabil | Memerlukan kunci SSH (`ssh-keygen`) |

Pelayan auto-kesan yang mana satu tersedia. Mulakan semula pelayan untuk URL Cloudflare baharu, atau gunakan Pinggy untuk URL kekal.

---

## Sokongan Shell

| OS | Lalai | Alternatif |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / mana-mana shell |
| **macOS** | bash | `sh` / `zsh` |

Tetapkan `"shell"` dalam config.json atau pilih dalam GUI.

---

## Keselamatan

- **OAuth 2.0 dengan PKCE** — jika kata laluan pemilik ditetapkan
- **Mod tanpa kata laluan** — jika tiada kata laluan dikonfigurasi, berjalan tanpa pengesahan
- **Sekatan laluan** — semua operasi fail disahkan terhadap root yang dibenarkan
- **Terowong pilihan** — Terowong Cloudflare melindungi dari pendedahan langsung
- **Tiada muat naik pihak ketiga** — kod anda tidak pernah meninggalkan mesin anda

---

## Membina dari Sumber

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Bina semua (semua platform)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Bina hanya untuk platform semasa
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Sokongan Platform

| Platform | Pelayan | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompil secara natif) |
| **macOS Intel** | ✅ | 🔧 (kompil secara natif) |
| **macOS M-chip** | ✅ | 🔧 (kompil secara natif) |

GUI memerlukan Fyne (OpenGL) — tidak boleh kompil silang. Pelayan boleh dikompil di mana-mana.

---

## Struktur Projek

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + pelayan MCP
│   └── devspace-gui/       ← Konfigurator GUI desktop (Fyne)
├── internal/
│   ├── auth/               ← Pembekal OAuth 2.0 + PKCE
│   ├── config/             ← Sistem konfigurasi mudah alih
│   ├── locales/            ← Terjemahan 47 bahasa
│   ├── logger/             ← Pengelogan berstruktur (zerolog)
│   ├── server/             ← HTTP + MCP + orkestrasi terowong
│   ├── skills/             ← AGENTS.md / penemuan kemahiran
│   ├── store/              ← Sesi ruang kerja SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Ruang kerja & pengesahan laluan
├── scripts/
│   ├── windows/            ← Skrip bina PowerShell
│   └── unix/               ← Skrip bina Bash + Makefile
├── readme/                 ← Terjemahan fail ini (47 bahasa)
├── build/                  ← Binari terkompil (semua platform)
├── tools/                  ← cloudflared.exe, dll.
├── go.mod / go.sum
└── README.md
```

---

Dibina dalam Go. Sifar npm. Sifar Node.js. Satu binari.
