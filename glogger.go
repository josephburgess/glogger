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
//   // Register with yyour router
//   blog.RegisterHandlers(router)
//
// assuming default conf, this will set up these routes:
// - "/blog" (list of posts)
// - "/blog/{slug}" (individual posts)
// - "/blog/_themes/{theme}.css" (theme CSS files)
//
