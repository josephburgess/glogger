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

// New creates a new Blog instance with the given configuration
func New(config Config) (*Blog, error) {
	// Set default values for config
	config.setDefaults()

	// Create template renderer
	renderer, err := newTemplateRenderer()
	if err != nil {
		return nil, err
	}

	return &Blog{
		config:   config,
		posts:    []Post{},
		renderer: renderer,
	}, nil
}

// Initialize loads all posts from the content directory
func (b *Blog) Initialize() error {
	return b.initialize()
}

// initialize is an internal method to load all posts
func (b *Blog) initialize() error {
	b.posts = []Post{} // Reset posts

	// Walk through the content directory
	err := filepath.Walk(b.config.ContentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Parse the post
		post, err := parsePost(path)
		if err != nil {
			return err
		}

		// Generate slug from filename
		filename := filepath.Base(path)
		post.Slug = strings.TrimSuffix(filename, filepath.Ext(filename))

		// Add post to the blog
		b.posts = append(b.posts, post)
		return nil
	})
	if err != nil {
		return err
	}

	// Sort posts by publish date (newest first)
	sort.Slice(b.posts, func(i, j int) bool {
		return b.posts[i].PublishDate.After(b.posts[j].PublishDate)
	})

	return nil
}

// GetPosts returns a copy of the current posts
func (b *Blog) GetPosts() []Post {
	result := make([]Post, len(b.posts))
	copy(result, b.posts)
	return result
}
