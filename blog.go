package glogger

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

type Blog struct {
	config   Config
	posts    []Post
	renderer *templateRenderer
	md       goldmark.Markdown
}

func New(config Config) (*Blog, error) {
	config.setDefaults()

	if !validateTheme(config.Theme) {
		config.Theme = "default"
	}

	renderer, err := newTemplateRenderer(config)
	if err != nil {
		return nil, err
	}

	b := &Blog{
		config:   config,
		posts:    []Post{},
		renderer: renderer,
		md:       newMarkdown(),
	}

	if err := b.Initialize(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Blog) Initialize() error {
	b.posts = []Post{}

	err := filepath.Walk(b.config.ContentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, err := parsePost(path, b.md)
		if err != nil {
			return err
		}

		if post.Draft {
			return nil
		}

		filename := filepath.Base(path)
		post.Slug = strings.TrimSuffix(filename, filepath.Ext(filename))

		b.posts = append(b.posts, post)
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(b.posts, func(i, j int) bool {
		return b.posts[i].PublishDate.After(b.posts[j].PublishDate)
	})

	return nil
}

func (b *Blog) GetPosts() []Post {
	result := make([]Post, len(b.posts))
	copy(result, b.posts)
	return result
}

func (b *Blog) URLPrefix() string {
	return b.config.URLPrefix
}

// Mount registers the blog's routes on mux using the configured URLPrefix.
// It also adds a redirect from URLPrefix to URLPrefix+"/" for clean URLs.
// This is the recommended way to attach a blog to an existing mux.
func (b *Blog) Mount(mux *http.ServeMux) {
	prefix := b.config.URLPrefix
	mux.HandleFunc("GET "+prefix, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, prefix+"/", http.StatusMovedPermanently)
	})
	mux.Handle(prefix+"/", http.StripPrefix(prefix, b.Handler()))
}

func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
}
