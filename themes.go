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

// ValidateTheme checks whether the given theme name is valid.
func ValidateTheme(theme string) bool {
	return slices.Contains(availableThemes, theme)
}

// GetThemePath returns the URL path to the theme's CSS file.
func GetThemePath(urlPrefix, theme string) string {
	return fmt.Sprintf("%s/_themes/%s.css", urlPrefix, theme)
}

const highlightJSBase = "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1"

// HighlightJSStyleURL returns the CDN URL for the highlight.js stylesheet matching the theme.
func HighlightJSStyleURL(theme string) string {
	styles := map[string]string{
		"default":  "github",
		"light":    "github",
		"dark":     "github-dark",
		"rosepine": "rose-pine",
	}
	style, ok := styles[theme]
	if !ok {
		style = "github"
	}
	return fmt.Sprintf("%s/styles/%s.min.css", highlightJSBase, style)
}

