package glogger

import (
	"net/http"
	"strings"
)

// Handler returns an http.Handler that serves the blog.
// mount with http.StripPrefix
// ex: http.Handle("/blog/", http.StripPrefix("/blog", blog.Handler()))
func (b *Blog) Handler() http.Handler {
	if len(b.posts) == 0 {
		if err := b.Initialize(); err != nil {
			panic("Failed to initialize blog: " + err.Error())
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", b.handleListPosts)
	mux.HandleFunc("GET /{slug}", b.handleSinglePost)
	mux.HandleFunc("GET /_themes/{theme}", b.handleThemeCSS)

	return mux
}

func (b *Blog) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	for _, post := range b.posts {
		if post.Slug == slug {
			html, err := b.renderer.renderPost(post, b.config.URLPrefix)
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

func (b *Blog) handleListPosts(w http.ResponseWriter, r *http.Request) {
	html, err := b.renderer.renderPostList(b.posts, b.config.URLPrefix)
	if err != nil {
		http.Error(w, "Error rendering post list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (b *Blog) handleThemeCSS(w http.ResponseWriter, r *http.Request) {
	theme := r.PathValue("theme")
	theme = strings.TrimSuffix(theme, ".css")

	if !ValidateTheme(theme) {
		http.NotFound(w, r)
		return
	}

	content, err := themeFS.ReadFile("assets/themes/" + theme + ".css")
	if err != nil {
		http.Error(w, "Theme not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.Write(content)
}

// PostHandler returns a standalone handler for rendering a single post file
// useful if you want to serve a specific markdown file somewhere
func PostHandler(postPath string, theme string) http.HandlerFunc {
	if theme == "" {
		theme = "default"
	}

	md := newMarkdown(highlightStyleForTheme(theme))

	return func(w http.ResponseWriter, r *http.Request) {
		post, err := parsePost(postPath, md)
		if err != nil {
			http.Error(w, "Error parsing post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		renderer, err := newTemplateRenderer(theme, "/blog")
		if err != nil {
			http.Error(w, "Error creating renderer: "+err.Error(), http.StatusInternalServerError)
			return
		}

		html, err := renderer.renderPost(post, "/blog")
		if err != nil {
			http.Error(w, "Error rendering post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}
