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

func validateTheme(theme string) bool {
	return slices.Contains(availableThemes, theme)
}

func getThemePath(urlPrefix, theme string) string {
	return fmt.Sprintf("%s/_themes/%s.css", urlPrefix, theme)
}

const highlightJSBase = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1"

func highlightJSStyleURL(syntaxTheme string) string {
	return fmt.Sprintf("%s/styles/%s.min.css", highlightJSBase, syntaxTheme)
}

func defaultSyntaxTheme(theme string) string {
	switch theme {
	case "rosepine":
		return "rose-pine"
	case "dark":
		return "github-dark"
	default:
		return "github"
	}
}
