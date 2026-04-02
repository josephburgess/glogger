package glogger

import (
	"bytes"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

// parse md with goldmark
func parsePost(filename string, md goldmark.Markdown) (Post, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Post{}, err
	}

	rawContent := string(content)
	lines := strings.Split(rawContent, "\n")

	title := "Untitled Post"
	var publishDate time.Time

	// find title
	titleLine := -1
	for i, line := range lines {
		if after, ok := strings.CutPrefix(line, "# "); ok {
			title = after
			titleLine = i
			break
		}
	}

	// check for date in YYYY-MM-DD in first 3 lines
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	dateLine := -1
	for i, line := range lines[:min(3, len(lines))] {
		if datePattern.MatchString(line) {
			if date, err := time.Parse("2006-01-02", strings.TrimSpace(line)); err == nil {
				publishDate = date
				dateLine = i
				break
			}
		}
	}

	// else use file modified date
	if publishDate.IsZero() {
		info, err := os.Stat(filename)
		if err != nil {
			return Post{}, err
		}
		publishDate = info.ModTime()
	}

	// remove date/title
	contentLines := make([]string, 0, len(lines))
	for i, line := range lines {
		if i != dateLine && i != titleLine {
			contentLines = append(contentLines, line)
		}
	}

	filteredContent := strings.Join(contentLines, "\n")

	var buf bytes.Buffer
	if err := md.Convert([]byte(filteredContent), &buf); err != nil {
		return Post{}, err
	}

	return Post{
		Title:       title,
		Content:     template.HTML(buf.String()),
		RawContent:  rawContent,
		PublishDate: publishDate,
	}, nil
}
