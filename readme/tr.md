# DevSpace (Go Sürümü)

**ChatGPT ve Claude'a yerel makinenize güvenli erişim verin. Herhangi bir MCP sunucusunu kodlama ortağınıza dönüştürün.**

DevSpace, AI asistanlarının gerçek yerel projelerinizde — dosyalarınız, araçlarınız, terminaliniz — hiçbir şeyi üçüncü taraflara yüklemeden kod okumasına, düzenlemesine, aramasına ve çalıştırmasına olanak tanıyan, kendi kendine barındırılan bir MCP sunucusudur. Kendi makinenizde çalıştırır, kontrol ettiğiniz bir tünel aracılığıyla erişime açar ve isteğe bağlı olarak bir parola ile güvence altına alırsınız.

---

## İçindekiler

- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Kurulum](#kurulum)
- [AI Neler Yapabilir](#ai-neler-yapabilir)
- [Yapılandırma](#yapılandırma)
- [Tünel (Uzaktan Erişim)](#tünel-uzaktan-erişim)
- [Kabuk Desteği](#kabuk-desteği)
- [Güvenlik](#güvenlik)
- [Kaynaktan Derleme](#kaynaktan-derleme)
- [Platform Desteği](#platform-desteği)
- [Proje Yapısı](#proje-yapısı)

### 🌍 Çeviriler

| Dil | | Dil | | Dil | |
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

## Hızlı Başlangıç

### 1. İndir
[Releases](../../releases) sayfasından platformunuzu seçin veya kaynaktan derleyin:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Yapılandır (GUI veya metin)
```bash
devspace-gui                  # Masaüstü yapılandırıcı (GUI)
devspace init                 # Metin tabanlı yapılandırıcı
```

### 3. Çalıştır
```bash
devspace                      # Sunucuyu başlatır. Yapılandırmayı otomatik algılar.
```

Bu ayrıca, `tools/` içinde `cloudflared` bulunursa otomatik olarak bir Cloudflare Tüneli başlatır.

### 4. MCP istemcinizi bağlayın
```
https://SİZİN-TÜNELİNİZ.trycloudflare.com/mcp
```
Veya yerel olarak: `http://127.0.0.1:7676/mcp`

---

## Kurulum

Node.js yok, npm yok, Python yok. Tek bir binary.

| Platform | İndirme |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: yerel olarak derleyin) |
| **macOS Intel** | `devspace` (GUI: yerel olarak derleyin) |
| **macOS M-chip** | `devspace` (GUI: yerel olarak derleyin) |

Yalnızca kaynaktan derlerken **Go 1.23+** gerektirir.

---

## AI Neler Yapabilir

Bağlandıktan sonra AI, onaylanmış proje klasörlerinizden birini çalışma alanı olarak açabilir:

- Çalışma alanı içindeki dosyaları **okuma, yazma ve düzenleme**
- Regex ile **kod arama** ve dizinleri inceleme
- **Kabuk komutları çalıştırma** (Windows'ta PowerShell, Unix'te bash)
- `AGENTS.md` / `CLAUDE.md` dosyalarından **proje talimatlarını keşfetme**
- Taşınabilir `.devspace/config.json` ile **otomatik yapılandırma**

8 MCP aracı: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Yapılandırma

Tüm yapılandırma **çalıştırılabilir dosya ile aynı klasördedir** (taşınabilir):

```
.devspace/
├── config.json       ← izin verilen kökler, port, kabuk, dil, kimlik doğrulama
└── auth.json         ← sahip parolası (isteğe bağlı)
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

| Alan | Varsayılan | Açıklama |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | İşletim sisteminden otomatik algılar. 47 dili destekler |
| `toolMode` | `full` | `full` (tüm araçlar) veya `minimal` (yalnızca arama için kabuk) |
| `toolNaming` | `short` | `short` (read, write) veya `legacy` (read_file, write_file) |

Ortam değişkenlerine gerek yok — her şey taşınabilir yapılandırma dosyasındadır.

---

## Tünel (Uzaktan Erişim)

ChatGPT web sürümü için (HTTPS gereklidir), DevSpace otomatik olarak bir tünel başlatır:

| Tünel | URL türü | Kurulum |
|---|---|---|
| **Cloudflare** | Rastgele (otomatik) | `cloudflared.exe` dosyasını `tools/` içine koyun |
| **Pinggy** | Kararlı | SSH anahtarı gerektirir (`ssh-keygen`) |

Sunucu hangisinin kullanılabilir olduğunu otomatik olarak algılar. Yeni bir Cloudflare URL'si için sunucuyu yeniden başlatın veya kalıcı bir URL için Pinggy kullanın.

---

## Kabuk Desteği

| İS | Varsayılan | Alternatifler |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / herhangi bir kabuk |
| **macOS** | bash | `sh` / `zsh` |

config.json içinde `"shell"` ayarını yapın veya GUI'den seçin.

---

## Güvenlik

- **PKCE ile OAuth 2.0** — sahip parolası ayarlanmışsa
- **Parolasız mod** — parola yapılandırılmamışsa, kimlik doğrulama olmadan çalışır
- **Yol kısıtlaması** — tüm dosya işlemleri izin verilen köklere göre doğrulanır
- **İsteğe bağlı tünel** — Cloudflare Tüneli doğrudan maruz kalmaya karşı korur
- **Üçüncü taraf yüklemeleri yok** — kodunuz asla makinenizden ayrılmaz

---

## Kaynaktan Derleme

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Her şeyi derle (tüm platformlar)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Yalnızca mevcut platform için derle
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Platform Desteği

| Platform | Sunucu | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (yerel olarak derleyin) |
| **macOS Intel** | ✅ | 🔧 (yerel olarak derleyin) |
| **macOS M-chip** | ✅ | 🔧 (yerel olarak derleyin) |

GUI, Fyne (OpenGL) gerektirir — çapraz derlenemez. Sunucu her yerde derlenir.

---

## Proje Yapısı

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP sunucusu
│   └── devspace-gui/       ← Masaüstü GUI yapılandırıcı (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE sağlayıcı
│   ├── config/             ← Taşınabilir yapılandırma sistemi
│   ├── locales/            ← 47 dil çevirisi
│   ├── logger/             ← Yapılandırılmış günlükleme (zerolog)
│   ├── server/             ← HTTP + MCP + tünel orkestrasyonu
│   ├── skills/             ← AGENTS.md / beceri keşfi
│   ├── store/              ← SQLite çalışma alanı oturumları
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Çalışma alanı ve yol doğrulama
├── scripts/
│   ├── windows/            ← PowerShell derleme betiği
│   └── unix/               ← Bash + Makefile derleme betikleri
├── readme/                 ← Bu dosyanın çevirileri (47 dil)
├── build/                  ← Derlenmiş binary'ler (tüm platformlar)
├── tools/                  ← cloudflared.exe, vb.
├── go.mod / go.sum
└── README.md
```

---

Go ile geliştirildi. Sıfır npm. Sıfır Node.js. Tek binary.
