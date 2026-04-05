package glogger

import (
	"html/template"
	"time"
)

type Post struct {
	Title       string
	Content     template.HTML
	PublishDate time.Time
	Slug        string
	Description string
	Tags        []string
	Draft       bool
}

type Config struct {
	ContentDir  string // directory containing markdown files
	URLPrefix   string // URL prefix for the blog (e.g. "/blog")
	Theme       string // theme name: "default", "dark", "light", "rosepine"
	SyntaxTheme string // highlight.js theme name (e.g. "rose-pine", "github-dark"); defaults to best match for Theme
	Title       string // blog title used in RSS feed channel (default: "Blog")
	Description string // blog description used in RSS feed channel (optional)
	BaseURL     string // base URL of the site (e.g. "https://example.com") — used to build absolute links in RSS feed
}

type PostTemplateData struct {
	Post
	BlogPrefix   string
	ThemeCSS     string
	HighlightCSS string
}

type ListTemplateData struct {
	Posts      []Post
	BlogPrefix string
	ThemeCSS   string
	Title      string
}

type templateRenderer struct {
	postTemplate    *template.Template
	listTemplate    *template.Template
	theme           string
	urlPrefix       string
	highlightCSSURL string
}

func (c *Config) setDefaults() {
	if c.ContentDir == "" {
		c.ContentDir = "content/posts"
	}
	if c.URLPrefix == "" {
		c.URLPrefix = "/blog"
	}
	if c.Theme == "" {
		c.Theme = "default"
	}
	if c.SyntaxTheme == "" {
		c.SyntaxTheme = defaultSyntaxTheme(c.Theme)
	}
	if c.Title == "" {
		c.Title = "Blog"
	}
}
