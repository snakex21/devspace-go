# DevSpace (Go izdanje)

**Dajte ChatGPT-u i Claudeu siguran pristup vašem lokalnom računalu. Pretvorite bilo koji MCP domaćin u svog partnera za programiranje.**

DevSpace je samostalno hostirani MCP poslužitelj koji AI asistentima omogućuje čitanje, uređivanje, pretraživanje i pokretanje koda u vašim stvarnim lokalnim projektima — vaše datoteke, vaši alati, vaš terminal — bez slanja ičega trećoj strani. Pokrećete ga na svom računalu, izlažete ga kroz tunel koji kontrolirate i opcionalno osiguravate lozinkom.

---

## Sadržaj

- [Brzi početak](#brzi-početak)
- [Instalacija](#instalacija)
- [Što AI može učiniti](#što-ai-može-učiniti)
- [Konfiguracija](#konfiguracija)
- [Tunel (udaljeni pristup)](#tunel-udaljeni-pristup)
- [Podrška za ljusku](#podrška-za-ljusku)
- [Sigurnost](#sigurnost)
- [Izgradnja iz izvornog koda](#izgradnja-iz-izvornog-koda)
- [Podrška za platforme](#podrška-za-platforme)
- [Struktura projekta](#struktura-projekta)

### 🌍 Prijevodi

| Jezik | | Jezik | | Jezik | |
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

## Brzi početak

### 1. Preuzimanje
Odaberite svoju platformu sa stranice [Releases](../../releases) ili izgradite iz izvornog koda:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Konfigurirajte (GUI ili tekst)
```bash
devspace-gui                  # Konfigurator za radnu površinu (GUI)
devspace init                 # Tekstualni konfigurator
```

### 3. Pokretanje
```bash
devspace                      # Pokreće poslužitelj. Automatski otkriva konfiguraciju.
```

Ovo također automatski pokreće Cloudflare tunel ako je `cloudflared` pronađen u `tools/`.

### 4. Povežite svoj MCP klijent
```
https://YOUR-TUNNEL.trycloudflare.com/mcp
```
Ili lokalno: `http://127.0.0.1:7676/mcp`

---

## Instalacija

Bez Node.js-a, bez npm-a, bez Pythona. Jedna binarna datoteka.

| Platforma | Preuzimanje |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: kompajlirati nativno) |
| **macOS Intel** | `devspace` (GUI: kompajlirati nativno) |
| **macOS M-chip** | `devspace` (GUI: kompajlirati nativno) |

Zahtijeva **Go 1.23+** samo ako gradite iz izvornog koda.

---

## Što AI može učiniti

Nakon povezivanja, AI može otvoriti jednu od vaših odobrenih mapa projekta kao radni prostor:

- **Čitati, pisati i uređivati** datoteke unutar radnog prostora
- **Pretraživati kod** s regularnim izrazima i pregledavati direktorije
- **Pokretati naredbe ljuske** (PowerShell na Windowsu, bash na Unixu)
- **Otkrivati upute projekta** iz `AGENTS.md` / `CLAUDE.md`
- **Automatski konfigurirati** s prijenosnim `.devspace/config.json`

8 MCP alata: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Konfiguracija

Sva konfiguracija nalazi se **u istoj mapi kao izvršna datoteka** (prijenosno):

```
.devspace/
├── config.json       ← dopušteni korijeni, port, ljuska, jezik, autentifikacija
└── auth.json         ← lozinka vlasnika (opcionalno)
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

| Polje | Zadano | Opis |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automatsko otkrivanje iz OS-a. Podržava 47 jezika |
| `toolMode` | `full` | `full` (svi alati) ili `minimal` (samo ljuska za pretragu) |
| `toolNaming` | `short` | `short` (read, write) ili `legacy` (read_file, write_file) |

Nisu potrebne varijable okoline — sve je u prijenosnoj konfiguracijskoj datoteci.

---

## Tunel (udaljeni pristup)

Za web verziju ChatGPT-a (potreban HTTPS), DevSpace automatski pokreće tunel:

| Tunel | Vrsta URL-a | Postavljanje |
|---|---|---|
| **Cloudflare** | Nasumični (auto) | Stavite `cloudflared.exe` u `tools/` |
| **Pinggy** | Stabilan | Potreban SSH ključ (`ssh-keygen`) |

Poslužitelj automatski otkriva koji je dostupan. Ponovno pokrenite poslužitelj za novi Cloudflare URL ili koristite Pinggy za trajni URL.

---

## Podrška za ljusku

| OS | Zadano | Alternative |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / bilo koja ljuska |
| **macOS** | bash | `sh` / `zsh` |

Postavite `"shell"` u config.json ili odaberite u GUI-u.

---

## Sigurnost

- **OAuth 2.0 s PKCE** — ako je postavljena lozinka vlasnika
- **Način bez lozinke** — ako nije konfigurirana lozinka, radi bez autentifikacije
- **Ograničenje putanje** — sve operacije s datotekama provjeravaju se prema dopuštenim korijenima
- **Opcionalni tunel** — Cloudflare tunel štiti od izravnog izlaganja
- **Bez slanja trećim stranama** — vaš kod nikada ne napušta vaše računalo

---

## Izgradnja iz izvornog koda

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Izgradi sve (sve platforme)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Izgradi samo za trenutnu platformu
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Podrška za platforme

| Platforma | Poslužitelj | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (kompajlirati nativno) |
| **macOS Intel** | ✅ | 🔧 (kompajlirati nativno) |
| **macOS M-chip** | ✅ | 🔧 (kompajlirati nativno) |

GUI zahtijeva Fyne (OpenGL) — ne može se unakrsno kompajlirati. Poslužitelj se kompajlira svugdje.

---

## Struktura projekta

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP poslužitelj
│   └── devspace-gui/       ← Konfigurator radne površine GUI (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE pružatelj
│   ├── config/             ← Prijenosni konfiguracijski sustav
│   ├── locales/            ← Prijevodi na 47 jezika
│   ├── logger/             ← Strukturirano zapisivanje (zerolog)
│   ├── server/             ← HTTP + MCP + orkestracija tunela
│   ├── skills/             ← AGENTS.md / otkrivanje vještina
│   ├── store/              ← SQLite sesije radnih prostora
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Radni prostor i provjera putanje
├── scripts/
│   ├── windows/            ← PowerShell skripta za izgradnju
│   └── unix/               ← Bash + Makefile skripte za izgradnju
├── readme/                 ← Prijevodi ove datoteke (47 jezika)
├── build/                  ← Kompajlirane binarne datoteke (sve platforme)
├── tools/                  ← cloudflared.exe, itd.
├── go.mod / go.sum
└── README.md
```

---

Izgrađeno u Go-u. Nula npm. Nula Node.js. Jedna binarna datoteka.
