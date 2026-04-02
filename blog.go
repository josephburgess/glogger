package glogger

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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

	if !ValidateTheme(config.Theme) {
		config.Theme = "default"
	}

	renderer, err := newTemplateRenderer(config.Theme, config.URLPrefix)
	if err != nil {
		return nil, err
	}

	md := newMarkdown(config.HighlightStyle)

	return &Blog{
		config:   config,
		posts:    []Post{},
		renderer: renderer,
		md:       md,
	}, nil
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

// URLPrefix returns the configured URL prefix for the blog.
func (b *Blog) URLPrefix() string {
	return b.config.URLPrefix
}

func newMarkdown(highlightStyle string) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(highlightStyle),
			),
		),
	)
}
