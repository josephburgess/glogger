package glogger

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

// parse md with goldmark
func parsePost(filename string) (Post, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Post{}, err
	}

	// default untitled post unless # found in a line
	lines := strings.Split(string(content), "\n")
	title := "Untitled Post"
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	info, err := os.Stat(filename)
	if err != nil {
		return Post{}, err
	}

	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return Post{}, err
	}

	return Post{
		Title:       title,
		Content:     template.HTML(buf.String()),
		RawContent:  string(content),
		PublishDate: info.ModTime(),
	}, nil
}
