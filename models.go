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
}
