package glogger

import (
	"embed"
	"fmt"
	"slices"
)

//go:embed assets/themes/*.css
var themeFS embed.FS

var availableThemes = []string{
	"default",
	"dark",
	"light",
	"rosepine",
}

// ValidateTheme checks if the provided theme name is valid
func ValidateTheme(theme string) bool {
	return slices.Contains(availableThemes, theme)
}

// GetThemePath returns the path to the theme CSS file
func GetThemePath(urlPrefix, theme string) string {
	return fmt.Sprintf("%s/_themes/%s.css", urlPrefix, theme)
}

// ListAvailableThemes returns a slice of all available theme names
func ListAvailableThemes() []string {
	return availableThemes
}
