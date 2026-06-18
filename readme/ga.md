# DevSpace (Go Edition) - Gaeilge

**Tabhair rochtain shlán do ChatGPT & Claude ar do mheaisín áitiúil. Déan comhpháirtí códála d'aon óstach MCP.**

Is freastalaí MCP féin-óstáilte é DevSpace a ligeann do chúntóirí AI cód a léamh, a chur in eagar, a chuardach agus a rith i do fhíor-thionscadail áitiúla — do chomhaid, d'uirlisí, do theirminéal — gan aon rud a uaslódáil chuig tríú páirtí. Rithfidh tú é ar do mheaisín, nochtfaidh tú é trí thollán a rialaíonn tú, agus roghnachdaíonn tú é le pasfhocal.

---

## Clár na nÁbhar

- [Tús Tapa](#tús-tapa)
- [Suiteáil](#suiteáil)
- [Cad is Féidir le AI a Dhéanamh](#cad-is-féidir-le-ai-a-dhéanamh)
- [Cumraíocht](#cumraíocht)
- [Tollán (Rochtain Chianach)](#tollán-rochtain-chianach)
- [Tacaíocht Bhlaoise](#tacaíocht-bhlaoise)
- [Slándáil](#slándáil)
- [Tógáil ón bhFoinse](#tógáil-ón-bhfhoinse)
- [Tacaíocht Ardáin](#tacaíocht-ardáin)
- [Struchtúr an Tionscadail](#struchtúr-an-tionscadail)

### 🌍 Aistriúcháin

| Teanga | | Teanga | | Teanga | |
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

## Tús Tapa

### 1. Íoslódáil
Roghnaigh d'ardán ó [Eisiúintí](../../releases) nó tóg ón bhfoinse:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Cumraigh (GUI nó téacs)
```bash
devspace-gui                  # Cumraitheoir deisce (GUI)
devspace init                 # Cumraitheoir téacsbhunaithe
```

### 3. Rith
```bash
devspace                      # Tosaíonn an freastalaí. Aimsíonn cumraíocht go huathoibríoch.
```

Tosaíonn sé seo Tollán Cloudflare go huathoibríoch freisin má aimsítear `cloudflared` i `tools/`.

### 4. Ceangail do chliant MCP
```
https://DO-THOLLÁN.trycloudflare.com/mcp
```
Nó go háitiúil: `http://127.0.0.1:7676/mcp`

---

## Suiteáil

Gan Node.js, gan npm, gan Python. Comhad dénártha amháin.

| Ardán | Íoslódáil |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: tiomsaigh go dúchasach) |
| **macOS Intel** | `devspace` (GUI: tiomsaigh go dúchasach) |
| **macOS M-chip** | `devspace` (GUI: tiomsaigh go dúchasach) |

Teastaíonn **Go 1.23+** ach amháin má tá tú ag tógáil ón bhfoinse.

---

## Cad is Féidir le AI a Dhéanamh

Nuair atá sé ceangailte, is féidir leis an AI ceann de do fhillteáin tionscadail ceadaithe a oscailt mar spás oibre:

- Comhaid laistigh den spás oibre a **léamh, a scríobh agus a chur in eagar**
- **Cód a chuardach** le regex agus eolairí a iniúchadh
- **Orduithe blaosca a rith** (PowerShell ar Windows, bash ar Unix)
- **Treoracha tionscadail a aimsiú** ó `AGENTS.md` / `CLAUDE.md`
- **Uathchumrú** le `.devspace/config.json` iniompartha

8 n-uirlis MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Cumraíocht

Tá an chumraíocht ar fad **san fhillteán céanna leis an gcomhad inrite** (iniompartha):

```
.devspace/
├── config.json       ← fréamhacha ceadaithe, port, blaosc, teanga, fíordheimhniú
└── auth.json         ← pasfhocal úinéara (roghnach)
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

| Réimse | Réamhshocrú | Cur Síos |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Uathbhrath ón OS. Tacaíonn le 47 teanga |
| `toolMode` | `full` | `full` (gach uirlis) nó `minimal` (blaosc amháin le haghaidh cuardaigh) |
| `toolNaming` | `short` | `short` (read, write) nó `legacy` (read_file, write_file) |

Níl aon athróga timpeallachta de dhíth — tá gach rud sa chomhad cumraíochta iniompartha.

---

## Tollán (Rochtain Chianach)

Do leagan gréasáin ChatGPT (HTTPS de dhíth), tosaíonn DevSpace tollán go huathoibríoch:

| Tollán | Cineál URL | Socrú |
|---|---|---|
| **Cloudflare** | Randamach (uathoibríoch) | Cuir `cloudflared.exe` i `tools/` |
| **Pinggy** | Cobhsaí | Teastaíonn eochair SSH (`ssh-keygen`) |

Aimsíonn an freastalaí go huathoibríoch cé acu atá ar fáil. Atosaigh an freastalaí le haghaidh URL Cloudflare nua, nó úsáid Pinggy le haghaidh URL buan.

---

## Tacaíocht Bhlaoise

| OS | Réamhshocrú | Roghanna Eile |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / aon bhlaosc |
| **macOS** | bash | `sh` / `zsh` |

Socraigh `"shell"` i config.json nó roghnaigh sa GUI.

---

## Slándáil

- **OAuth 2.0 le PKCE** — má tá pasfhocal úinéara socraithe
- **Mód gan phasfhocal** — mura bhfuil pasfhocal cumraithe, ritheann gan fíordheimhniú
- **Coimeádán conaire** — déantar gach oibríocht comhaid a bhailíochtú i gcoinne fréamhacha ceadaithe
- **Tollán roghnach** — Cosnaíonn Tollán Cloudflare ar nochtadh díreach
- **Gan uaslódálacha tríú páirtí** — ní fhágann do chód do mheaisín choíche

---

## Tógáil ón bhFoinse

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Tóg gach rud (gach ardán)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Tóg don ardán reatha amháin
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Tacaíocht Ardáin

| Ardán | Freastalaí | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (tiomsaigh go dúchasach) |
| **macOS Intel** | ✅ | 🔧 (tiomsaigh go dúchasach) |
| **macOS M-chip** | ✅ | 🔧 (tiomsaigh go dúchasach) |

Teastaíonn Fyne (OpenGL) ón GUI — ní féidir cros-thiomsú. Tiomsaíonn an freastalaí i ngach áit.

---

## Struchtúr an Tionscadail

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + freastalaí MCP
│   └── devspace-gui/       ← Cumraitheoir deisce GUI (Fyne)
├── internal/
│   ├── auth/               ← Soláthraí OAuth 2.0 + PKCE
│   ├── config/             ← Córas cumraíochta iniompartha
│   ├── locales/            ← Aistriúcháin 47 teanga
│   ├── logger/             ← Logáil struchtúrtha (zerolog)
│   ├── server/             ← HTTP + MCP + ceolfhoireann tolláin
│   ├── skills/             ← AGENTS.md / aimsiú scileanna
│   ├── store/              ← Seisiúin spáis oibre SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Spás oibre & bailíochtú conaire
├── scripts/
│   ├── windows/            ← Script tógála PowerShell
│   └── unix/               ← Scripteanna tógála Bash + Makefile
├── readme/                 ← Aistriúcháin an chomhaid seo (47 teanga)
├── build/                  ← Comhaid dhéanártha tiomsaithe (gach ardán)
├── tools/                  ← cloudflared.exe, srl.
├── go.mod / go.sum
└── README.md
```

---

Tógtha i Go. Náid npm. Náid Node.js. Comhad dénártha amháin.
