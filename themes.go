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

// check the theme in cfg is valid
func ValidateTheme(theme string) bool {
	return slices.Contains(availableThemes, theme)
}

// grab the path to css file
func GetThemePath(urlPrefix, theme string) string {
	return fmt.Sprintf("%s/_themes/%s.css", urlPrefix, theme)
}
