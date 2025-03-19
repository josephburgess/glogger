package glogger

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

type Blog struct {
	config   Config
	posts    []Post
	renderer *templateRenderer
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

	return &Blog{
		config:   config,
		posts:    []Post{},
		renderer: renderer,
	}, nil
}

func (b *Blog) Initialize() error {
	return b.initialize()
}

func (b *Blog) initialize() error {
	b.posts = []Post{}

	err := filepath.Walk(b.config.ContentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, err := parsePost(path)
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
