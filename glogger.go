// simple blog engine for go websites. mainly built this for my personal site http://joeburgess.dev
// but thought it would be worth packaging for others to use!
package glogger

const Version = "0.2.0"

// themes available
const (
	ThemeDefault  = "default"
	ThemeDark     = "dark"
	ThemeLight    = "light"
	ThemeRosePine = "rosepine"
)

// Usage:
//
//   // Create a new blog with default settings
//   blog, err := glogger.New(glogger.Config{
//     ContentDir: "content/posts",  // Where your markdown files are stored
//     URLPrefix: "/blog",           // URL prefix for the blog routes
//     Theme: glogger.ThemeRosePine, // Optional theme selection
//   })
//
//   // Register the blog handlers with your router
//   blog.RegisterHandlers(router)
//
// This will set up routes for:
// - "/blog" (list of posts)
// - "/blog/{slug}" (individual posts)
// - "/blog/_themes/{theme}.css" (theme CSS files)
//
// Available themes:
// - glogger.ThemeDefault: A clean, minimal light theme
// - glogger.ThemeDark: A dark theme for reduced eye strain
// - glogger.ThemeLight: A light theme with subtle colors
// - glogger.ThemeRosePine: A soothing dark theme inspired by Rose Pine
