//go:build !windows

package config

func detectOSSystemLang() string {
	return ""
}
