# glogger

A minimal blog engine for Go web apps. Drop into an existing `net/http` server and get a fully functional blog without a build step, separate process, or database. Very few dependencies and super lightweight.

> Built for fun to add to my personal site [joeburgess.dev](https://joeburgess.dev), but usable for anyone running a go HTTP server that wants a blog but doesn't want to run a separate service or doesn't need something as feature-full as [Hugo](https://github.com/gohugoio/hugo).

<p align="center">
<img src="https://github.com/user-attachments/assets/061f15c9-c55b-4a82-a004-d0fc095074ee" width="800">
</p>

## Features

- Markdown posts with YAML frontmatter
- 4 built-in themes (default, light, dark, rose pine)
- Client side syntax highlighting w/ [highlight.js](https://highlightjs.org) — no extra Go deps
- RSS 2.0 feed at `/feed.xml`
- Tag filtering
- No database needed, posts are plain `.md` files on disk

## Installation

```bash
go get github.com/josephburgess/glogger
```

## Quick Start

```go
blog, err := glogger.New(glogger.Config{
    ContentDir:  "content/posts",
    URLPrefix:   "/blog",
    Theme:       glogger.ThemeRosePine,
    Title:       "My Blog",
    BaseURL:     "https://example.com",
})
if err != nil {
    log.Fatal(err)
}

// mount on your existing stdlib mux:
mux.Handle("/blog/", http.StripPrefix("/blog", blog.Handler()))
```

This registers these routes under your mux:

| Route | Description |
|---|---|
| `GET /blog/` | Post list |
| `GET /blog/{slug}` | Individual post |
| `GET /blog/feed.xml` | RSS 2.0 feed |
| `GET /blog/_tags/{tag}` | Posts filtered by tag |
| `GET /blog/_themes/{theme}.css` | Theme CSS |

## Post Format

Posts are `.md` files with required YAML frontmatter:

```markdown
---
title: "Hello world"
date: 2026-01-01
description: "Optional — shown in post list and RSS feed"
tags: [go, blogging]
draft: false
---

Content goes here.
```

The filename (without `.md`) becomes the URL slug. Draft posts are hidden from the listing and not served.

## Configuration

```go
type Config struct {
    ContentDir  string // directory containing markdown files (default: "content/posts")
    URLPrefix   string // URL prefix for the blog (default: "/blog")
    Theme       string // theme: "default", "dark", "light", "rosepine"
    SyntaxTheme string // highlight.js theme (sensible default will be set depending on Theme)
    Title       string // blog title, used for RSS and page header (default: "Blog")
    Description string // blog description for RSS (optional)
    BaseURL     string // used for absolute links in RSS
}
```

### Themes

`Theme` controls the overall page styling:

| Constant | Value | Style |
|---|---|---|
| `glogger.ThemeDefault` | `"default"` | Clean white |
| `glogger.ThemeLight` | `"light"` | Light grey |
| `glogger.ThemeDark` | `"dark"` | Dark |
| `glogger.ThemeRosePine` | `"rosepine"` | Rose Pine |

### Syntax Highlighting

`SyntaxTheme` controls the [highlight.js theme](https://highlightjs.org/examples) for code blocks. Sensible defaults are set per theme:

| Site theme | Default syntax theme |
|---|---|
| `default` / `light` | `github` |
| `dark` | `github-dark` |
| `rosepine` | `rose-pine` |

Override with any highlight.js theme name here: https://highlightjs.org/examples

```go
glogger.Config{
    Theme:       glogger.ThemeRosePine,
    SyntaxTheme: "tokyo-night-dark",
}
```

## Standalone markdown handler

Serve a single markdown file outside the blog structure. Useful for changelogs, about pages, etc:

```go
mux.HandleFunc("/changelog", glogger.PostHandler("content/changelog.md", glogger.ThemeDark))
```

## go dependencies

- [goldmark](https://github.com/yuin/goldmark) (markdown parsing)
- [yaml.v3](https://github.com/go-yaml/yaml) (frontmatter parsing)

## Contributing

Contributions welcome, themes especially!! Open a PR.

## License

MIT
