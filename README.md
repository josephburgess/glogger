# glogger

A lightweight, minimal blogging engine written in Go. Primarily built for my own site at [https://joeburgess.dev](https://joeburgess.dev) but thought it would be nice to have as a package.

<p align="center">
<img src="https://github.com/user-attachments/assets/061f15c9-c55b-4a82-a004-d0fc095074ee" width="800">
</p>

## Features

- **Simple setup**: Add to an existing Go site with a few lines of code
- **Routing setup**: Registers routes using [gorilla/mux](https://github.com/gorilla/mux)
- **Themes**: Includes 4 themes
- **Responsive**: Mobile-friendly
- **No database**: Posts are stored as Markdown files in your project

## Installation

```bash
go get github.com/josephburgess/glogger
```

## Quick Start

1. Create a directory for your blog posts (default: `content/posts/`)
2. Add .md files for your blog posts in this directory
3. Example initialising glogger using gorilla/mux:

```go
package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/josephburgess/glogger"
)

func main() {
    router := mux.NewRouter()

    // Initialize the blog
    blog, err := glogger.New(glogger.Config{
        ContentDir: "content/posts",   // where you are storing your .md blog posts
        URLPrefix:  "/blog",           // url prefix for blog routes
        Theme:      glogger.ThemeLight, // theme
    })
    if err != nil {
        panic(err)
    }

    // register routes
    blog.RegisterHandlers(router)

    // other routes
    router.HandleFunc("/", homeHandler).Methods("GET")

    http.ListenAndServe(":8080", router)
}
```

## Post Format

Posts should be written in `.md` in your content directory. The first line starting with `#` will be used as the post title, and the first line containing a date in `YYYY-MM-DD` format will be used as the publish date.

Example post (`hello-world.md`):

```markdown
2024-03-21

# Hello world

First blog post eyyyyyyyy!

## Section 1

Interesting stuff...

### Section 2

Thanks for making it this far...
```

The filename will be used as the post slug in urls.

## Configuration

The config struct accepts the following:

```go
type Config struct {
    ContentDir    string // directory where .md posts are stored
    URLPrefix     string // URL prefix for the blog
    Theme         string // theme
}
```

### Available Themes

- `glogger.ThemeDefault` - clean, white theme
- `glogger.ThemeDark` - standard dark theme
- `glogger.ThemeLight` - standard light grey theme
- `glogger.ThemeRosePine` - rose pine because its what I use

## API

### Creating a Blog Instance

```go
blog, err := glogger.New(glogger.Config{
    ContentDir: "content/posts",
    URLPrefix:  "/blog",
    Theme:      glogger.ThemeRosePine,
})
```

### Registering Routes

```go
blog.RegisterHandlers(router)
```

This will register the following routes:

- `/blog` - List of all blog posts
- `/blog/{slug}` - Individual blog post
- `/blog/_themes/{theme}.css` - Theme CSS files

### Standalone post handler

If you want to use glogger to add a specific post outside of the default route structure:

```go
router.HandleFunc("/special-post", glogger.PostHandler("path/to/post.md", glogger.ThemeDark))
```

## Customization

### Templates

The built-in templates are embedded in the glogger package. If you want to customize them, you'll need to fork the repository and modify the files in the `assets/templates` directory.

### Themes

Similarly, the CSS themes are embedded in the package. You can modify them by forking the repository and updating the files in the `assets/themes` directory. The themes are really standardised so I'd love theme contributions!

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) for routing
  - This choice was simply because its what I use for my personal site - hope to expand the compatibility with other routers in future
- [goldmark](https://github.com/yuin/goldmark) for Markdown parsing

## Roadmap

The following are just a few bits I'd like to add at some point!

### Short-term

- Make routing easier outside of gorilla/mux
- Pagination
- Authors
- Custom user themes set in config
- Syntax highlighting in code blocks

### Medium-term

- Tag/category support
- Search
- RSS feed generation
- Table of contents for longer posts
- Better image support

### Long-term

- If I somehow get through all the above I'll think about this more!

## License

MIT License

## Contributing

Contributions are welcome, in particular give me your themes. Please feel free to submit a Pull Request.
