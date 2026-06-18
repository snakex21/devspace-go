# DevSpace (Uhlelo lwe-Go)

**Nikeza u-ChatGPT no-Claude ukufinyelela okuphephile emshinini wakho wendawo. Guqula noma iyiphi i-host ye-MCP ibe umlingani wakho wokubhala amakhodi.**

I-DevSpace iyiseva ye-MCP ezisingathayo evumela abasizi be-AI ukuthi bafunde, bahlele, bacinge, futhi basebenzise ikhodi kumaphrojekthi akho angempela endawo — amafayela akho, amathuluzi akho, itheminali yakho — ngaphandle kokulayisha noma yini kumuntu wesithathu. Uyisebenzisa emshinini wakho, uyiveze ngomhubhe owulawulayo, futhi ngokuzikhethela uyivikele ngephasiwedi.

---

## Okuqukethwe

- [Ukuqala Ngokushesha](#ukuqala-ngokushesha)
- [Ukufakwa](#ukufakwa)
- [I-AI Engakwenza](#i-ai-engakwenza)
- [Ukumisa](#ukumisa)
- [Umhubhe (Ukufinyelela Okukude)](#umhubhe-ukufinyelela-okukude)
- [Ukusekelwa Kwe-Shell](#ukusekelwa-kwe-shell)
- [Ukuvikeleka](#ukuvikeleka)
- [Ukwakha Kusuka Emthonjeni](#ukwakha-kusuka-emthonjeni)
- [Ukusekelwa Kwepulatifomu](#ukusekelwa-kwepulatifomu)
- [Isakhiwo Sephrojekthi](#isakhiwo-sephrojekthi)

### 🌍 Ukuhunyushwa

| Ulimi | | Ulimi | | Ulimi | |
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

## Ukuqala Ngokushesha

### 1. Landa
Khetha ipulatifomu yakho ku-[Releases](../../releases) noma wakhe kusuka emthonjeni:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Misela (i-GUI noma umbhalo)
```bash
devspace-gui                  # Isimisi sedeskithophu (i-GUI)
devspace init                 # Isimisi esisekelwe embhalweni
```

### 3. Sebenzisa
```bash
devspace                      # Iqala iseva. Ithola ukumisa ngokuzenzakalelayo.
```

Lokhu futhi kuqala ngokuzenzakalelayo Umhubhe we-Cloudflare uma `cloudflared` itholakala ku-`tools/`.

### 4. Xhuma iklayenti lakho le-MCP
```
https://UMHUBHE-WAKHO.trycloudflare.com/mcp
```
Noma endaweni: `http://127.0.0.1:7676/mcp`

---

## Ukufakwa

Akukho i-Node.js, akukho i-npm, akukho i-Python. Ibhinari eyodwa.

| Ipulatifomu | Ukulanda |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (i-GUI: hlanganisa ngokwendabuko) |
| **macOS Intel** | `devspace` (i-GUI: hlanganisa ngokwendabuko) |
| **macOS M-chip** | `devspace` (i-GUI: hlanganisa ngokwendabuko) |

Idinga **i-Go 1.23+** kuphela uma wakha kusuka emthonjeni.

---

## I-AI Engakwenza

Uma isixhunyiwe, i-AI ingavula enye yamafolda akho ephrojekthi agunyaziwe njengendawo yokusebenza:

- **Ukufunda, ukubhala, nokuhlela** amafayela ngaphakathi kwendawo yokusebenza
- **Ukucinga ikhodi** nge-regex nokuhlola izikhombisi
- **Ukusebenzisa imiyalo ye-shell** (i-PowerShell ku-Windows, i-bash ku-Unix)
- **Ukuthola imiyalo yephrojekthi** ku-`AGENTS.md` / `CLAUDE.md`
- **Ukuzimisela ngokuzenzakalelayo** nge-`.devspace/config.json` ephathekayo

Amathuluzi e-MCP ayi-8: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Ukumisa

Konke ukumisa kutholakala **kufolda efanayo nefayela elisebenzisekayo** (okuphathekayo):

```
.devspace/
├── config.json       ← izimpande ezigunyaziwe, imbobo, i-shell, ulimi, ukuqinisekiswa
└── auth.json         ← iphasiwedi yomnikazi (ngokuzikhethela)
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

| Inkambu | Okuzenzakalelayo | Incazelo |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Ukuthola ngokuzenzakalelayo kwi-OS. Kweseka izilimi ezingama-47 |
| `toolMode` | `full` | `full` (wonke amathuluzi) noma `minimal` (i-shell kuphela yokucinga) |
| `toolNaming` | `short` | `short` (read, write) noma `legacy` (read_file, write_file) |

Akukho ziguquguquko zemvelo ezidingekayo — konke kusefayeleni lokumisa eliphathekayo.

---

## Umhubhe (Ukufinyelela Okukude)

Kunguqulo yewebhu ye-ChatGPT (i-HTTPS iyadingeka), i-DevSpace iqala umhubhe ngokuzenzakalelayo:

| Umhubhe | Uhlobo lwe-URL | Ukusetha |
|---|---|---|
| **Cloudflare** | Okungahleliwe (okuzenzakalelayo) | Beka `cloudflared.exe` ku-`tools/` |
| **Pinggy** | Okuzinzile | Idinga ukhiye we-SSH (`ssh-keygen`) |

Iseva ithola ngokuzenzakalelayo ukuthi iyiphi etholakalayo. Qala kabusha iseva ukuze uthole i-URL entsha ye-Cloudflare, noma usebenzise i-Pinggy ukuthola i-URL ehlala njalo.

---

## Ukusekelwa Kwe-Shell

| I-OS | Okuzenzakalelayo | Okunye |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / noma iyiphi i-shell |
| **macOS** | bash | `sh` / `zsh` |

Setha `"shell"` ku-config.json noma ukhethe ku-GUI.

---

## Ukuvikeleka

- **I-OAuth 2.0 ene-PKCE** — uma iphasiwedi yomnikazi isethiwe
- **Imodi engenaphasiwedi** — uma kungekho phasiwedi emisiwe, isebenza ngaphandle kokuqinisekiswa
- **Ukuvinjelwa kwendlela** — yonke imisebenzi yamafayela iqinisekiswa ngokumelene nezimpande ezigunyaziwe
- **Umhubhe ongakhethwa** — Umhubhe we-Cloudflare uvikela ekuvezweni okuqondile
- **Akukho ukulayishwa kumuntu wesithathu** — ikhodi yakho ayilokothi ishiye umshini wakho

---

## Ukwakha Kusuka Emthonjeni

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Yakha konke (wonke amapulatifomu)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Yakha kuphela ipulatifomu yamanje
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Ukusekelwa Kwepulatifomu

| Ipulatifomu | Iseva | I-GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (hlanganisa ngokwendabuko) |
| **macOS Intel** | ✅ | 🔧 (hlanganisa ngokwendabuko) |
| **macOS M-chip** | ✅ | 🔧 (hlanganisa ngokwendabuko) |

I-GUI idinga i-Fyne (OpenGL) — ayikwazi ukuhlanganiswa ngezindlela ezahlukene. Iseva ihlanganiswa yonke indawo.

---

## Isakhiwo Sephrojekthi

```
devspace-go/
├── cmd/
│   ├── devspace/           ← I-CLI + iseva ye-MCP
│   └── devspace-gui/       ← Isimisi se-GUI sedeskithophu (Fyne)
├── internal/
│   ├── auth/               ← Umhlinzeki we-OAuth 2.0 + PKCE
│   ├── config/             ← Isistimu yokumisa ephathekayo
│   ├── locales/            ← Ukuhunyushwa kwezilimi ezingama-47
│   ├── logger/             ← Ukuloga okuhlelekile (zerolog)
│   ├── server/             ← I-HTTP + MCP + ukuphathwa komhubhe
│   ├── skills/             ← Ukutholwa kwe-AGENTS.md / amakhono
│   ├── store/              ← Izikhathi zendawo yokusebenza ze-SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Ukuqinisekiswa kwendawo yokusebenza nendlela
├── scripts/
│   ├── windows/            ← Iskripthi sokwakha se-PowerShell
│   └── unix/               ← Izikripthi zokwakha ze-Bash + Makefile
├── readme/                 ← Ukuhunyushwa kwaleli fayela (izilimi ezingama-47)
├── build/                  ← Amabhinari ahlanganisiwe (wonke amapulatifomu)
├── tools/                  ← cloudflared.exe, njll.
├── go.mod / go.sum
└── README.md
```

---

Yakhelwe nge-Go. Zero npm. Zero Node.js. Ibhinari eyodwa.
