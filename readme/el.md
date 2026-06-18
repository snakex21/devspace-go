# DevSpace (Go Edition) - Ελληνικά

**Δώστε στους ChatGPT & Claude ασφαλή πρόσβαση στον τοπικό σας υπολογιστή. Μετατρέψτε οποιονδήποτε MCP host σε συνεργάτη προγραμματισμού.**

Το DevSpace είναι ένας αυτο-φιλοξενούμενος MCP διακομιστής που επιτρέπει στους βοηθούς AI να διαβάζουν, να επεξεργάζονται, να αναζητούν και να εκτελούν κώδικα στα πραγματικά τοπικά σας έργα — τα αρχεία σας, τα εργαλεία σας, το τερματικό σας — χωρίς να ανεβάζουν τίποτα σε τρίτο μέρος. Το εκτελείτε στον υπολογιστή σας, το εκθέτετε μέσω μιας σήραγγας που ελέγχετε και προαιρετικά το ασφαλίζετε με κωδικό πρόσβασης.

---

## Πίνακας Περιεχομένων

- [Γρήγορη Έναρξη](#γρήγορη-έναρξη)
- [Εγκατάσταση](#εγκατάσταση)
- [Τι Μπορεί να Κάνει η ΤΝ](#τι-μπορεί-να-κάνει-η-τν)
- [Διαμόρφωση](#διαμόρφωση)
- [Σήραγγα (Απομακρυσμένη Πρόσβαση)](#σήραγγα-απομακρυσμένη-πρόσβαση)
- [Υποστήριξη Κελύφους](#υποστήριξη-κελύφους)
- [Ασφάλεια](#ασφάλεια)
- [Δημιουργία από Πηγαίο Κώδικα](#δημιουργία-από-πηγαίο-κώδικα)
- [Υποστήριξη Πλατφορμών](#υποστήριξη-πλατφορμών)
- [Δομή Έργου](#δομή-έργου)

### 🌍 Μεταφράσεις

| Γλώσσα | | Γλώσσα | | Γλώσσα | |
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

## Γρήγορη Έναρξη

### 1. Λήψη
Επιλέξτε την πλατφόρμα σας από τις [Εκδόσεις](../../releases) ή δημιουργήστε από πηγαίο κώδικα:
```bash
./scripts/unix/build.sh      # Linux / Mac
.\scripts\windows\build.ps1   # Windows
```

### 2. Διαμόρφωση (GUI ή κείμενο)
```bash
devspace-gui                  # Διαμορφωτής επιφάνειας εργασίας (GUI)
devspace init                 # Διαμορφωτής βασισμένος σε κείμενο
```

### 3. Εκτέλεση
```bash
devspace                      # Εκκινεί τον διακομιστή. Ανιχνεύει αυτόματα τη διαμόρφωση.
```

Αυτό εκκινεί επίσης αυτόματα μια σήραγγα Cloudflare αν βρεθεί το `cloudflared` στο `tools/`.

### 4. Συνδέστε τον πελάτη MCP σας
```
https://Η-ΣΗΡΑΓΓΑ-ΣΑΣ.trycloudflare.com/mcp
```
Ή τοπικά: `http://127.0.0.1:7676/mcp`

---

## Εγκατάσταση

Χωρίς Node.js, χωρίς npm, χωρίς Python. Ένα μόνο δυαδικό αρχείο.

| Πλατφόρμα | Λήψη |
|---|---|
| **Windows** | `devspace.exe` + `devspace-gui.exe` |
| **Linux** | `devspace` (GUI: μεταγλωττίστε εγγενώς) |
| **macOS Intel** | `devspace` (GUI: μεταγλωττίστε εγγενώς) |
| **macOS M-chip** | `devspace` (GUI: μεταγλωττίστε εγγενώς) |

Απαιτεί **Go 1.23+** μόνο αν δημιουργείτε από πηγαίο κώδικα.

---

## Τι Μπορεί να Κάνει η ΤΝ

Μόλις συνδεθεί, η ΤΝ μπορεί να ανοίξει έναν από τους εγκεκριμένους φακέλους έργου σας ως χώρο εργασίας:

- **Ανάγνωση, εγγραφή και επεξεργασία** αρχείων εντός του χώρου εργασίας
- **Αναζήτηση κώδικα** με regex και επιθεώρηση καταλόγων
- **Εκτέλεση εντολών κελύφους** (PowerShell σε Windows, bash σε Unix)
- **Ανακάλυψη οδηγιών έργου** από `AGENTS.md` / `CLAUDE.md`
- **Αυτόματη διαμόρφωση** με φορητό `.devspace/config.json`

8 εργαλεία MCP: `open_workspace`, `read`, `write`, `edit`, `grep`, `glob`, `ls`, `bash`

---

## Διαμόρφωση

Όλη η διαμόρφωση βρίσκεται **στον ίδιο φάκελο με το εκτελέσιμο** (φορητή):

```
.devspace/
├── config.json       ← επιτρεπόμενες ρίζες, θύρα, κέλυφος, γλώσσα, αυθεντικοποίηση
└── auth.json         ← κωδικός πρόσβασης κατόχου (προαιρετικό)
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

| Πεδίο | Προεπιλογή | Περιγραφή |
|---|---|---|
| `shell` | `auto` | `auto`, `powershell`, `cmd`, `bash`, `sh` |
| `lang` | `auto` | Αυτόματη ανίχνευση από το ΛΣ. Υποστηρίζει 47 γλώσσες |
| `toolMode` | `full` | `full` (όλα τα εργαλεία) ή `minimal` (μόνο κέλυφος για αναζήτηση) |
| `toolNaming` | `short` | `short` (read, write) ή `legacy` (read_file, write_file) |

Δεν απαιτούνται μεταβλητές περιβάλλοντος — όλα βρίσκονται στο φορητό αρχείο διαμόρφωσης.

---

## Σήραγγα (Απομακρυσμένη Πρόσβαση)

Για την έκδοση web του ChatGPT (απαιτεί HTTPS), το DevSpace εκκινεί αυτόματα μια σήραγγα:

| Σήραγγα | Τύπος URL | Ρύθμιση |
|---|---|---|
| **Cloudflare** | Τυχαίο (auto) | Τοποθετήστε το `cloudflared.exe` στο `tools/` |
| **Pinggy** | Σταθερό | Χρειάζεται κλειδί SSH (`ssh-keygen`) |

Ο διακομιστής ανιχνεύει αυτόματα ποια είναι διαθέσιμη. Επανεκκινήστε τον διακομιστή για νέο Cloudflare URL ή χρησιμοποιήστε το Pinggy για μόνιμο URL.

---

## Υποστήριξη Κελύφους

| ΛΣ | Προεπιλογή | Εναλλακτικές |
|---|---|---|
| **Windows** | PowerShell | `cmd` / `pwsh` |
| **Linux** | bash | `sh` / οποιοδήποτε κέλυφος |
| **macOS** | bash | `sh` / `zsh` |

Ορίστε το `"shell"` στο config.json ή επιλέξτε στο GUI.

---

## Ασφάλεια

- **OAuth 2.0 με PKCE** — αν έχει οριστεί κωδικός πρόσβασης κατόχου
- **Λειτουργία χωρίς κωδικό** — αν δεν έχει διαμορφωθεί κωδικός, εκτελείται χωρίς αυθεντικοποίηση
- **Περιορισμός διαδρομής** — όλες οι λειτουργίες αρχείων επικυρώνονται έναντι των επιτρεπόμενων ριζών
- **Προαιρετική σήραγγα** — Η σήραγγα Cloudflare προστατεύει από άμεση έκθεση
- **Χωρίς μεταφορτώσεις σε τρίτους** — ο κώδικάς σας δεν εγκαταλείπει ποτέ τον υπολογιστή σας

---

## Δημιουργία από Πηγαίο Κώδικα

```bash
git clone https://github.com/waishnav/devspace-go
cd devspace-go

# Δημιουργία όλων (όλες οι πλατφόρμες)
.\scripts\windows\build.ps1     # Windows
./scripts/unix/build.sh          # Linux / Mac
make -f scripts/unix/Makefile    # Linux / Mac (make)

# Δημιουργία μόνο για την τρέχουσα πλατφόρμα
go build -o devspace ./cmd/devspace/
go build -o devspace-gui ./cmd/devspace-gui/
```

---

## Υποστήριξη Πλατφορμών

| Πλατφόρμα | Διακομιστής | GUI |
|---|---|---|
| **Windows** | ✅ | ✅ |
| **Linux** | ✅ | 🔧 (μεταγλωττίστε εγγενώς) |
| **macOS Intel** | ✅ | 🔧 (μεταγλωττίστε εγγενώς) |
| **macOS M-chip** | ✅ | 🔧 (μεταγλωττίστε εγγενώς) |

Το GUI απαιτεί Fyne (OpenGL) — δεν μπορεί να γίνει cross-compile. Ο διακομιστής μεταγλωττίζεται παντού.

---

## Δομή Έργου

```
devspace-go/
├── cmd/
│   ├── devspace/           ← CLI + διακομιστής MCP
│   └── devspace-gui/       ← Διαμορφωτής GUI επιφάνειας εργασίας (Fyne)
├── internal/
│   ├── auth/               ← Πάροχος OAuth 2.0 + PKCE
│   ├── config/             ← Φορητό σύστημα διαμόρφωσης
│   ├── locales/            ← Μεταφράσεις 47 γλωσσών
│   ├── logger/             ← Δομημένη καταγραφή (zerolog)
│   ├── server/             ← HTTP + MCP + ενορχήστρωση σήραγγας
│   ├── skills/             ← AGENTS.md / ανακάλυψη δεξιοτήτων
│   ├── store/              ← Συνεδρίες χώρου εργασίας SQLite
│   ├── tools/              ← read, write, edit, grep, glob, ls, bash
│   └── workspace/          ← Χώρος εργασίας & επικύρωση διαδρομής
├── scripts/
│   ├── windows/            ← Σενάριο δημιουργίας PowerShell
│   └── unix/               ← Σενάρια δημιουργίας Bash + Makefile
├── readme/                 ← Μεταφράσεις αυτού του αρχείου (47 γλώσσες)
├── build/                  ← Μεταγλωττισμένα δυαδικά αρχεία (όλες οι πλατφόρμες)
├── tools/                  ← cloudflared.exe, κλπ.
├── go.mod / go.sum
└── README.md
```

---

Χτισμένο σε Go. Μηδέν npm. Μηδέν Node.js. Ένα δυαδικό αρχείο.
