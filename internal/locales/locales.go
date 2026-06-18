package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed *.json
var localeFS embed.FS

var (
	currentLocale = "en"
	mu            sync.RWMutex
	stringsMap    = map[string]map[string]string{}
)

type localeFile struct {
	Lang    string            `json:"lang"`
	Name    string            `json:"name"`
	Strings map[string]string `json:"strings"`
}

// Init loads all locale files and sets the active locale.
func Init(lang string) {
	entries, err := localeFS.ReadDir(".")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := localeFS.ReadFile(entry.Name())
		if err != nil {
			continue
		}
		data = stripBOM(data)
		var lf localeFile
		if err := json.Unmarshal(data, &lf); err != nil {
			continue
		}
		mu.Lock()
		stringsMap[lf.Lang] = lf.Strings
		mu.Unlock()
	}

	SetLocale(lang)
}

// SetLocale changes the active locale.
// Falls back to "en" if the requested locale is not available.
func SetLocale(lang string) {
	mu.RLock()
	_, ok := stringsMap[lang]
	mu.RUnlock()

	if !ok {
		// Try to match just the first 2 chars
		if len(lang) > 2 {
			short := lang[:2]
			mu.RLock()
			_, ok = stringsMap[short]
			mu.RUnlock()
			if ok {
				lang = short
			}
		}
	}

	if !ok {
		// Fallback: try "en"
		mu.RLock()
		_, enOk := stringsMap["en"]
		mu.RUnlock()
		if enOk {
			lang = "en"
		} else {
			return // nothing to fallback to
		}
	}

	mu.Lock()
	currentLocale = lang
	mu.Unlock()
}

// T returns a translated string for the given key.
// Optional args are passed to fmt.Sprintf for formatting.
func T(key string, args ...interface{}) string {
	mu.RLock()
	locale := stringsMap[currentLocale]
	mu.RUnlock()

	if locale == nil {
		// fallback to English
		mu.RLock()
		locale = stringsMap["en"]
		mu.RUnlock()
	}
	if locale == nil {
		if len(args) > 0 {
			return fmt.Sprintf(key, args...)
		}
		return key
	}

	str, ok := locale[key]
	if !ok {
		// fallback to key itself
		str = key
	}

	if len(args) > 0 {
		return fmt.Sprintf(str, args...)
	}
	return str
}

// LocaleName returns the human-readable name of the current locale.
func LocaleName() string {
	mu.RLock()
	locale := stringsMap[currentLocale]
	mu.RUnlock()
	if locale != nil {
		if name, ok := locale["lang.name"]; ok {
			return name
		}
	}
	return currentLocale
}

// AvailableLocales returns all available locale codes.
func AvailableLocales() []string {
	mu.RLock()
	defer mu.RUnlock()
	var codes []string
	for code := range stringsMap {
		codes = append(codes, code)
	}
	return codes
}

// stripBOM removes UTF-8 BOM from data.
func stripBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}
