# DevSpace (Go Edition) - Suomi

**Anna ChatGPT:lle & Claudelle turvallinen pääsy paikalliselle koneellesi. Muuta mikä tahansa MCP-isäntä koodauskumppaniksesi.**

DevSpace on itseisännöity MCP-palvelin, joka antaa tekoälyavustajien lukea, muokata, etsiä ja suorittaa koodia oikeissa paikallisissa projekteissasi — tiedostosi, työkalusi, päätteesi — lataamatta mitään kolmannelle osapuolelle. Suoritat sen koneellasi, paljastat sen hallitsemasi tunnelin kautta ja valinnaisesti suojaat sen salasanalla.

---

## Sisällysluettelo

- [Pika-aloitus](#pika-aloitus)
- [Asennus](#asennus)
- [Mitä Tekoäly Voi Tehdä](#mitä-tekoäly-voi-tehdä)
- [Määritykset](#määritykset)
- [Tunneli (Etäkäyttö)](#tunneli-etäkäyttö)
- [Komentotulkin Tuki](#komentotulkin-tuki)
- [Tietoturva](#tietoturva)
- [Kääntäminen Lähdekoodista](#kääntäminen-lähdekoodista)
- [Alustatuki](#alustatuki)
- [Projektin Rakenne](#projektin-rakenne)

### 🌍 Käännökset

| Kieli | | Kieli | | Kieli | |
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

## Pika-aloitus

### 1. Lataa
Valitse alustasi [Julkaisuista](../../releases) tai käännä lähdekoodista:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Määritä (GUI tai teksti)
```bash
devspace-gui                  # Työpöytämääritin (GUI)
devspace init                 # Tekstipohjainen määritin
```

### 3. Suorita
```bash
devspace                      # Käynnistää palvelimen. Havaitsee määritykset automaattisesti.
```

Tämä käynnistää myös automaattisesti Cloudflare-tunnelin, jos `cloudflared` löytyy kansiosta `tools/`.

### 4. Yhdistä MCP-asiakkaasi
```
https://SINUN-TUNNELISI.trycloudflare.com/mcp
```
Tai paikallisesti: `http://127.0.0.1:7676/mcp`

---

## Asennus

Ei Node.js:ää, ei npm:ää, ei Pythonia. Yksi binääritiedosto.

| Alusta | Lataus |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: käännä natiivisti) |
| **macOS Intel** | `devspace` (GUI: käännä natiivisti) |
| **macOS M-siru** | `devspace` (GUI: käännä natiivisti) |

Vaatii **Go 1.23+** vain lähdekoodista käännettäessä.

---

## Mitä Tekoäly Voi Tehdä

Yhdistettynä tekoäly voi avata yhden hyväksytyistä projektikansioistasi työtilana:

- **Lukea, kirjoittaa ja muokata** tiedostoja työtilan sisällä
- **Etsiä koodia** regexillä ja tarkastella hakemistoja
- **Suorittaa komentotulkkikomentoja** (PowerShell Windowsissa, bash Unixissa)
- **Löytää projektiohjeita** tiedostoista `AGENTS.md` / `CLAUDE.md`
- **Määrittää automaattisesti** kannettavalla `.devspace/config.json`-tiedostolla

8 MCP-työkalua: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Määritykset

Kaikki määritykset sijaitsevat **samassa kansiossa kuin suoritettava tiedosto** (kannettava):

```
.devspace/
├── config.json       ← sallitut juuret, portti, komentotulkki, kieli, todennus
└── auth.json         ← omistajan salasana (valinnainen)
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

| Kenttä | Oletus | Kuvaus |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Automaattinen tunnistus käyttöjärjestelmästä. Tukee 47 kieltä |
| `toolMode` | `full` | `full` (kaikki työkalut) tai `minimal` (vain komentotulkki hakuun) |
| `toolNaming` | `short` | `short` (read, write) tai `legacy` (read_file, write_file) |

Ympäristömuuttujia ei tarvita — kaikki on kannettavassa määritystiedostossa.

---

## Tunneli (Etäkäyttö)

ChatGPT-verkkoversiolle (vaatii HTTPS) DevSpace käynnistää tunnelin automaattisesti:

| Tunneli | URL-tyyppi | Asennus |
|---|---|---|
| **Cloudflare** | Satunnainen (auto) | Laita `cloudflared.exe` kansioon `tools/` |
| **Pinggy** | Vakaa | Tarvitsee SSH-avaimen (`ssh-keygen`) |

Palvelin havaitsee automaattisesti, kumpi on käytettävissä. Käynnistä palvelin uudelleen saadaksesi uuden Cloudflare-URL-osoitteen tai käytä Pinggytä pysyvään URL-osoitteeseen.

---

## Komentotulkin Tuki

| Käyttöjärjestelmä | Oletus | Vaihtoehdot |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / mikä tahansa komentotulkki |
| **macOS** | bash | `sh` / `zsh` |

Aseta `"shell"` config.json-tiedostossa tai valitse GUI:ssa.

---

## Tietoturva

- **OAuth 2.0 PKCE:lla** — jos omistajan salasana on asetettu
- **Salasanaton tila** — jos salasanaa ei ole määritetty, toimii ilman todennusta
- **Polun rajoittaminen** — kaikki tiedosto-operaatiot validoidaan sallittuja juuria vasten
- **Valinnainen tunneli** — Cloudflare-tunneli suojaa suoralta altistumiselta
- **Ei kolmannen osapuolen latauksia** — koodisi ei koskaan poistu koneeltasi

---

## Kääntäminen Lähdekoodista

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Käännä kaikki (kaikki alustat)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Käännä vain nykyiselle alustalle
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Alustatuki

| Alusta | Palvelin | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (käännä natiivisti) |
| **macOS Intel** | ✅ | 🔧 (käännä natiivisti) |
| **macOS M-siru** | ✅ | 🔧 (käännä natiivisti) |

GUI vaatii Fynen (OpenGL) — ei voi ristikääntää. Palvelin kääntyy kaikkialla.

---

## Projektin Rakenne

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + MCP-palvelin
│   └── devspace-gui/       ← Työpöydän GUI-määritin (Fyne)
├── internal/
│   ├── auth/               ← OAuth 2.0 + PKCE -tarjoaja
│   ├── config/             ← Kannettava määritysjärjestelmä
│   ├── locales/            ← 47 kielen käännökset
│   ├── logger/             ← Jäsennelty lokitus (zerolog)
│   ├── server/             ← HTTP + MCP + tunneliorkestrointi
│   ├── skills/             ← AGENTS.md / taitojen löytäminen
│   ├── store/              ← SQLite-työtilasessiot
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Työtila ja polun validointi
├── scripts/
│   ├── windows/            ← PowerShell-käännösskripti
│   └── unix/               ← Bash + Makefile -käännösskriptit
├── readme/                 ← Tämän tiedoston käännökset (47 kieltä)
├── build/                  ← Käännettyt binääritiedostot (kaikki alustat)
├── tools/                  ← cloudflared.exe, jne.
├── go.mod / go.sum
└── README.md
```

---

Rakennettu Go:lla. Nolla npm. Nolla Node.js. Yksi binääritiedosto.
