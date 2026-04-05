# glogger

A lightweight, minimal blogging engine written in Go. Primarily built for my own site at [https://joeburgess.dev](https://joeburgess.dev) but thought it would be nice to have as a package.

<p align="center">
<img src="https://github.com/user-attachments/assets/061f15c9-c55b-4a82-a004-d0fc095074ee" width="800">
</p>

## Features

- Markdown posts with YAML frontmatter
- 4 built in themes
- Syntax highlighting in code blocks uses [highlight.js](https://highlightjs.org) (means no extra go deps)
- No database — posts are plain `.md` files

## Installation

```bash
go get github.com/josephburgess/glogger
```

## Quick Start

```go
blog, err := glogger.New(glogger.Config{
    ContentDir: "content/posts",
    URLPrefix:  "/blog",
    Theme:      glogger.ThemeRosePine,
})
if err != nil {
    log.Fatal(err)
}

// Mount on your router (stdlib mux):
mux.Handle("/blog/", http.StripPrefix("/blog", blog.Handler()))
```

This sets up three routes:

- `GET /blog/` — post list
- `GET /blog/{slug}` — individual post
- `GET /blog/_themes/{theme}.css` — theme CSS

## Post Format

Posts are `.md` files and you can use yaml frontmatter:

```markdown
---
title: "Hello world"
date: 2026-01-01
description: "optional description"
tags: [go, blogging]
draft: false
---

content goes here..
```

The filename (without `.md`) becomes the url slug. Draft posts are hidden from the listing and not served.

## Configuration

```go
type Config struct {
    ContentDir  string // directory containing markdown files (default: "content/posts")
    URLPrefix   string // URL prefix for the blog (default: "/blog")
    Theme       string // site theme: "default", "dark", "light", "rosepine"
    SyntaxTheme string // highlight.js theme name (default: best match for Theme)
}
```

### Site Themes

`Theme` controls the overall page styling (background, text, links):

| Constant | Value | Style |
|---|---|---|
| `glogger.ThemeDefault` | `"default"` | Clean white |
| `glogger.ThemeLight` | `"light"` | Light grey |
| `glogger.ThemeDark` | `"dark"` | Dark |
| `glogger.ThemeRosePine` | `"rosepine"` | Rose Pine |

### Syntax Highlighting

`SyntaxTheme` controls which [highlight.js theme](https://highlightjs.org/examples) is used for code blocks. If not set, sensible defaults are set for the site theme:

| Site theme | Default syntax theme |
|---|---|
| `default` / `light` | `github` |
| `dark` | `github-dark` |
| `rosepine` | `rose-pine` |

To use any other highlight.js theme, pass its name directly:

```go
glogger.Config{
    Theme:       glogger.ThemeRosePine,
    SyntaxTheme: "tokyo-night-dark", // any theme from highlightjs.org/examples
}
```

## Standalone post handler

To serve a single markdown file outside the blog structure:

```go
mux.HandleFunc("/changelog", glogger.PostHandler("content/changelog.md", glogger.ThemeDark))
```

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Markdown parsing
- [yaml.v3](https://github.com/go-yaml/yaml) - frontmatter parsing
- [highlight.js](https://highlightjs.org) - syntax highlighting (not a Go dependency)

## Roadmap

The following are just a few bits I'd like to add at some point!

### Short-term

- Pagination
- Authors

### Medium-term

- Search
- RSS feed generation
- Table of contents for longer posts
- Better image support


## Contributing

Contributions are welcome, in particular give me your themes. Please feel free to submit a Pull Request.

## License

MIT License

