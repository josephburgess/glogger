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

// HighlightJSStyleURL returns the CDN URL for a named highlight.js theme.
// Any theme available at https://highlightjs.org/examples can be used.
func HighlightJSStyleURL(syntaxTheme string) string {
	return fmt.Sprintf("%s/styles/%s.min.css", highlightJSBase, syntaxTheme)
}

// defaultSyntaxTheme returns the highlight.js theme that best matches a glogger theme.
func defaultSyntaxTheme(theme string) string {
	switch theme {
	case "rosepine":
		return "rose-pine"
	case "dark":
		return "github-dark"
	case "light", "default":
		return "github"
	default:
		return "github"
	}
}

