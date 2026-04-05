package glogger

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Draft       bool     `yaml:"draft"`
}

func parsePost(filename string, md goldmark.Markdown) (Post, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Post{}, err
	}

	raw := string(content)

	if !strings.HasPrefix(raw, "---\n") {
		return Post{}, fmt.Errorf("missing frontmatter in %s", filename)
	}

	parts := strings.SplitN(raw, "---\n", 3)
	if len(parts) < 3 {
		return Post{}, fmt.Errorf("invalid frontmatter in %s", filename)
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return Post{}, err
	}

	title := fm.Title
	if title == "" {
		title = "Untitled Post"
	}

	var publishDate time.Time
	if fm.Date != "" {
		if d, err := time.Parse("2006-01-02", fm.Date); err == nil {
			publishDate = d
		}
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(strings.TrimSpace(parts[2])), &buf); err != nil {
		return Post{}, err
	}

	return Post{
		Title:       title,
		Content:     template.HTML(buf.String()),
		PublishDate: publishDate,
		Description: fm.Description,
		Tags:        fm.Tags,
		Draft:       fm.Draft,
	}, nil
}
