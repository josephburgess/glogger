// simple blog engine for go websites. mainly built this for my personal site http://joeburgess.dev
// but thought it would be worth packaging for others to use!
package glogger

const Version = "0.1.0"

// Usage:
//
//   // Create a new blog with default settings
//   blog, err := glogger.New(glogger.Config{
//     ContentDir: "content/posts",  // Where your markdown files are stored
//     URLPrefix: "/blog",           // URL prefix for the blog routes
//   })
//
//   // Register the blog handlers with your router
//   blog.RegisterHandlers(router)
//
// Unless set otherwise, this will set up routes for "/blog" (list of posts) and "/blog/{slug}" (individual posts)
