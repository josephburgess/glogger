package glogger

import (
	"html/template"
	"time"
)

type Post struct {
	Title       string
	Content     template.HTML
	RawContent  string
	PublishDate time.Time
	Slug        string
}

type Config struct {
	ContentDir    string // markdown files stored here
	URLPrefix     string // url prefix for the blog
	DefaultAuthor string // default author for posts
	PageSize      int    // post per page
	Theme         string // theme to use (default, dark, light, etc.)
}

type PostTemplateData struct {
	Post
	BlogPrefix string
	ThemeCSS   string
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

// default conf
func (c *Config) setDefaults() {
	if c.ContentDir == "" {
		c.ContentDir = "content/posts"
	}

	if c.URLPrefix == "" {
		c.URLPrefix = "/blog"
	}

	if c.PageSize == 0 {
		c.PageSize = 10
	}

	if c.Theme == "" {
		c.Theme = "default"
	}
}
