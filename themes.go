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

// ValidateTheme checks the theme in cfg is valid
func ValidateTheme(theme string) bool {
	return slices.Contains(availableThemes, theme)
}

// GetThemePath grabs the path to the relevant css file
func GetThemePath(urlPrefix, theme string) string {
	return fmt.Sprintf("%s/_themes/%s.css", urlPrefix, theme)
}

func highlightStyleForTheme(theme string) string {
	switch theme {
	case "rosepine":
		return "rose-pine"
	case "dark":
		return "dracula"
	case "light":
		return "rose-pine-dawn"
	case "default":
		return "github"
	default:
		return "dracula"
	}
}
