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
	ContentDir string // directory containing markdown files
	URLPrefix  string // URL prefix for the blog (e.g. "/blog")
	Theme      string // theme name: "default", "dark", "light", "rosepine"
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
	postTemplate *template.Template
	listTemplate *template.Template
	theme        string
	urlPrefix    string
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
}
