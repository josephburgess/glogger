package glog

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Post struct {
	Title       string
	Content     template.HTML
	RawContent  string
	PublishDate time.Time
}

type Config struct {
	ContentDir    string // markdown files stored here
	URLPrefix     string // url prefix for the blog
	DefaultAuthor string // default author for posts
	PageSize      int    // post per page
}

type Blog struct {
	config Config
	posts  []Post
}

// create a new blog
func New(config Config) (*Blog, error) {
	return &Blog{
		config: config,
		posts:  []Post{},
	}, nil
}

func (b *Blog) Initialize() error {
	// tbc
	return nil
}

// parse md file with goldmark
func ParsePost(filename string) (Post, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Post{}, err
	}

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

// template for rendering a post
func RenderPost(post Post) (string, error) {
	const postTmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 1rem;
        }
        h1 { margin-bottom: 0.5rem; }
        .date { color: #666; font-size: 0.9rem; margin-bottom: 2rem; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .back { margin-top: 2rem; display: inline-block; }
    </style>
</head>
<body>
    <article>
        <h1>{{.Title}}</h1>
        <div class="date">{{.PublishDate.Format "January 2, 2006"}}</div>
        <div class="content">
            {{.Content}}
        </div>
    </article>
    <a href="/blog" class="back">&larr; Back to all posts</a>
</body>
</html>
`
	tmpl, err := template.New("post").Parse(postTmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// createa handler function for a single post
func PostHandler(postPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := ParsePost(postPath)
		if err != nil {
			http.Error(w, "Error parsing post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		html, err := RenderPost(post)
		if err != nil {
			http.Error(w, "Error rendering post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}

// register http handlers
func (b *Blog) RegisterHandlers(router *mux.Router) {
}
