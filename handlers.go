package glogger

import (
	"net/http"
	"strings"
)

// Handler returns an http.Handler that serves the blog.
// Mount with http.StripPrefix:
//
//	http.Handle("/blog/", http.StripPrefix("/blog", blog.Handler()))
func (b *Blog) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", b.handleListPosts)
	mux.HandleFunc("GET /_tags/{tag}", b.handleTaggedPosts)
	mux.HandleFunc("GET /_themes/{theme}", b.handleThemeCSS)
	mux.HandleFunc("GET /{slug}", b.handleSinglePost)
	return mux
}

func (b *Blog) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	for _, post := range b.posts {
		if post.Slug == slug {
			html, err := b.renderer.renderPost(post)
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
	html, err := b.renderer.renderPostList(b.posts, "")
	if err != nil {
		http.Error(w, "Error rendering post list: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (b *Blog) handleTaggedPosts(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")

	var filtered []Post
	for _, post := range b.posts {
		for _, t := range post.Tags {
			if t == tag {
				filtered = append(filtered, post)
				break
			}
		}
	}

	html, err := b.renderer.renderPostList(filtered, tag)
	if err != nil {
		http.Error(w, "Error rendering tag page: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (b *Blog) handleThemeCSS(w http.ResponseWriter, r *http.Request) {
	theme := r.PathValue("theme")
	theme = strings.TrimSuffix(theme, ".css")

	if !validateTheme(theme) {
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

// PostHandler returns a standalone handler for rendering a single markdown file.
// Useful for serving a specific post outside the blog structure.
func PostHandler(postPath string, theme string) http.HandlerFunc {
	if theme == "" {
		theme = "default"
	}

	md := newMarkdown()
	renderer, err := newTemplateRenderer(theme, "/blog", highlightJSStyleURL(defaultSyntaxTheme(theme)))
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Error creating renderer: "+err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		post, err := parsePost(postPath, md)
		if err != nil {
			http.Error(w, "Error parsing post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		html, err := renderer.renderPost(post)
		if err != nil {
			http.Error(w, "Error rendering post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}
