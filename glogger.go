// Package glogger provides a lightweight blog engine for go web apps.
// It supports markdown content, multiple themes, and easy integration with existing sites.
package glogger

const Version = "0.3.0"

// themes available
const (
	ThemeDefault  = "default"
	ThemeDark     = "dark"
	ThemeLight    = "light"
	ThemeRosePine = "rosepine"
)

// Usage:
//
// blog, err := glogger.New(glogger.Config{
//     ContentDir:  "content/posts",
//     URLPrefix:   "/blog",
//     Theme:       glogger.ThemeRosePine,
//     Title:       "My Blog",
//     BaseURL:     "https://example.com",
// })
//
// blog.Mount(mux)  // registers all routes under URLPrefix
//
// assuming default conf, this will set up these routes (relative to prefix)
//   - GET /                    — post list
//   - GET /feed.xml            — RSS 2.0 feed
//   - GET /{slug}              — individual post
//   - GET /_tags/{tag}         — posts filtered by tag
//   - GET /_themes/{theme}.css — theme CSS
