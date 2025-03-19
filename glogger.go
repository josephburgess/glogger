package glogger

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	Slug        string
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
	if config.ContentDir == "" {
		config.ContentDir = "content/posts"
	}

	if config.URLPrefix == "" {
		config.URLPrefix = "/blog"
	}

	if config.PageSize == 0 {
		config.PageSize = 10
	}

	return &Blog{
		config: config,
		posts:  []Post{},
	}, nil
}

func (b *Blog) Initialize() error {
	return b.loadPosts()
}

// load all posts from the dir
func (b *Blog) loadPosts() error {
	b.posts = []Post{}

	err := filepath.Walk(b.config.ContentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		post, err := ParsePost(path)
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

func RenderPostList(posts []Post, urlPrefix string) (string, error) {
	const listTmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog Posts</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 1rem;
        }
        h1 { margin-bottom: 1.5rem; }
        .post-list { list-style: none; padding: 0; }
        .post-item { margin-bottom: 2rem; }
        .post-title { font-size: 1.4rem; margin-bottom: 0.2rem; }
        .post-date { color: #666; font-size: 0.9rem; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .home-link { margin-top: 2rem; display: inline-block; }
    </style>
</head>
<body>
    <h1>Blog Posts</h1>
    {{if .Posts}}
    <ul class="post-list">
        {{range .Posts}}
        <li class="post-item">
            <div class="post-title">
                <a href="{{$.URLPrefix}}/{{.Slug}}">{{.Title}}</a>
            </div>
            <div class="post-date">{{.PublishDate.Format "January 2, 2006"}}</div>
        </li>
        {{end}}
    </ul>
    {{else}}
    <p>No posts found.</p>
    {{end}}
    <a href="/" class="home-link">&larr; Back to home</a>
</body>
</html>
`
	data := struct {
		Posts     []Post
		URLPrefix string
	}{
		Posts:     posts,
		URLPrefix: urlPrefix,
	}

	tmpl, err := template.New("list").Parse(listTmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (b *Blog) HandleSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, post := range b.posts {
		if post.Slug == slug {
			html, err := RenderPost(post)
			if err != nil {
				http.Error(w, "Error rendering post: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
			return
		}
	}

	http.NotFound(w, r)
}

func (b *Blog) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	html, err := RenderPostList(b.posts, b.config.URLPrefix)
	if err != nil {
		http.Error(w, "Error rendering post list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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
	if len(b.posts) == 0 {
		err := b.Initialize()
		if err != nil {
			panic("Failed to initialize blog: " + err.Error())
		}
	}

	// create subrouter
	blogRouter := router.PathPrefix(b.config.URLPrefix).Subrouter()

	// list all
	blogRouter.HandleFunc("", b.HandleListPosts).Methods("GET")

	// single post
	blogRouter.HandleFunc("/{slug}", b.HandleSinglePost).Methods("GET")
}
