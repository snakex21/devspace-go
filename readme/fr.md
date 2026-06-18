# DevSpace (Go Edition) - Français

**Donnez à ChatGPT et Claude un accès sécurisé à votre machine locale. Transformez n'importe quel hôte MCP en votre partenaire de codage.**

DevSpace est un serveur MCP auto-hébergé qui permet aux assistants IA de lire, modifier, rechercher et exécuter du code dans vos vrais projets locaux — vos fichiers, vos outils, votre terminal — sans rien télécharger vers un tiers. Vous l'exécutez sur votre machine, l'exposez via un tunnel que vous contrôlez et le sécurisez optionnellement avec un mot de passe.

---

## Table des Matières

- [Démarrage Rapide](#démarrage-rapide)
- [Installation](#installation)
- [Ce Que l'IA Peut Faire](#ce-que-lia-peut-faire)
- [Configuration](#configuration)
- [Tunnel (Accès à Distance)](#tunnel-accès-à-distance)
- [Support Shell](#support-shell)
- [Sécurité](#sécurité)
- [Compilation depuis les Sources](#compilation-depuis-les-sources)
- [Support des Plateformes](#support-des-plateformes)
- [Structure du Projet](#structure-du-projet)

### 🌍 Traductions

| Langue | | Langue | | Langue | |
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

## Démarrage Rapide

### 1. Télécharger
Choisissez votre plateforme depuis les [Versions](../../releases) ou compilez depuis les sources :
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Configurer (GUI ou texte)
```bash
devspace-gui                  # Configurateur de bureau (GUI)
devspace init                 # Configurateur en mode texte
```

### 3. Exécuter
```bash
devspace                      # Démarre le serveur. Détecte automatiquement la config.
```

Cela démarre également automatiquement un tunnel Cloudflare si `cloudflared` est trouvé dans `tools/`.

### 4. Connectez votre client MCP
```
https://VOTRE-TUNNEL.trycloudflare.com/mcp
```
Ou localement : `http://127.0.0.1:7676/mcp`

---

## Installation

Pas de Node.js, pas de npm, pas de Python. Un seul binaire.

| Plateforme | Téléchargement |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI : compiler nativement) |
| **macOS Intel** | `devspace` (GUI : compiler nativement) |
| **macOS M-chip** | `devspace` (GUI : compiler nativement) |

Nécessite **Go 1.23+** uniquement en cas de compilation depuis les sources.

---

## Ce Que l'IA Peut Faire

Une fois connectée, l'IA peut ouvrir l'un de vos dossiers de projet approuvés comme espace de travail :

- **Lire, écrire et modifier** des fichiers dans l'espace de travail
- **Rechercher du code** avec des regex et inspecter des répertoires
- **Exécuter des commandes shell** (PowerShell sous Windows, bash sous Unix)
- **Découvrir les instructions du projet** depuis `AGENTS.md` / `CLAUDE.md`
- **Auto-configurer** avec `.devspace/config.json` portable

8 outils MCP : `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Configuration

Toute la configuration réside **dans le même dossier que l'exécutable** (portable) :

```
.devspace/
├── config.json       ← racines autorisées, port, shell, langue, authentification
└── auth.json         ← mot de passe du propriétaire (optionnel)
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

| Champ | Défaut | Description |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Détection automatique depuis l'OS. Prend en charge 47 langues |
| `toolMode` | `full` | `full` (tous les outils) ou `minimal` (shell uniquement pour la recherche) |
| `toolNaming` | `short` | `short` (read, write) ou `legacy` (read_file, write_file) |

Aucune variable d'environnement nécessaire — tout est dans le fichier de configuration portable.

---

## Tunnel (Accès à Distance)

Pour la version web de ChatGPT (HTTPS requis), DevSpace démarre automatiquement un tunnel :

| Tunnel | Type d'URL | Configuration |
|---|---|---|
| **Cloudflare** | Aléatoire (auto) | Placez `cloudflared.exe` dans `tools/` |
| **Pinggy** | Stable | Nécessite une clé SSH (`ssh-keygen`) |

Le serveur détecte automatiquement lequel est disponible. Redémarrez le serveur pour une nouvelle URL Cloudflare, ou utilisez Pinggy pour une URL permanente.

---

## Support Shell

| OS | Défaut | Alternatives |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / n'importe quel shell |
| **macOS** | bash | `sh` / `zsh` |

Définissez `"shell"` dans config.json ou choisissez dans l'interface graphique.

---

## Sécurité

- **OAuth 2.0 avec PKCE** — si le mot de passe du propriétaire est défini
- **Mode sans mot de passe** — si aucun mot de passe n'est configuré, fonctionne sans authentification
- **Confinement de chemin** — toutes les opérations sur les fichiers sont validées par rapport aux racines autorisées
- **Tunnel optionnel** — le tunnel Cloudflare protège contre l'exposition directe
- **Aucun téléchargement vers des tiers** — votre code ne quitte jamais votre machine

---

## Compilation depuis les Sources

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Compiler tout (toutes les plateformes)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Compiler uniquement pour la plateforme actuelle
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Support des Plateformes

| Plateforme | Serveur | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (compiler nativement) |
| **macOS Intel** | ✅ | 🔧 (compiler nativement) |
| **macOS M-chip** | ✅ | 🔧 (compiler nativement) |

L'interface graphique nécessite Fyne (OpenGL) — ne peut pas être compilée de manière croisée. Le serveur se compile partout.

---

## Structure du Projet

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + serveur MCP
│   └── devspace-gui/       ← Configurateur graphique de bureau (Fyne)
├── internal/
│   ├── auth/               ← Fournisseur OAuth 2.0 + PKCE
│   ├── config/             ← Système de configuration portable
│   ├── locales/            ← Traductions en 47 langues
│   ├── logger/             ← Journalisation structurée (zerolog)
│   ├── server/             ← HTTP + MCP + orchestration de tunnel
│   ├── skills/             ← AGENTS.md / découverte de compétences
│   ├── store/              ← Sessions d'espace de travail SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Espace de travail et validation de chemin
├── scripts/
│   ├── windows/            ← Script de compilation PowerShell
│   └── unix/               ← Scripts de compilation Bash + Makefile
├── readme/                 ← Traductions de ce fichier (47 langues)
├── build/                  ← Binaires compilés (toutes les plateformes)
├── tools/                  ← cloudflared.exe, etc.
├── go.mod / go.sum
└── README.md
```

---

Construit en Go. Zéro npm. Zéro Node.js. Un binaire.
