package glogger

import (
	"net/http"

	"github.com/gorilla/mux"
)

// register blog routes
func (b *Blog) RegisterHandlers(router *mux.Router) {
	if len(b.posts) == 0 {
		err := b.initialize()
		if err != nil {
			panic("Failed to initialize blog: " + err.Error())
		}
	}

	blogRouter := router.PathPrefix(b.config.URLPrefix).Subrouter()

	blogRouter.HandleFunc("", b.handleListPosts).Methods("GET")
	blogRouter.HandleFunc("/{slug}", b.handleSinglePost).Methods("GET")
	blogRouter.HandleFunc("/_themes/{theme}.css", b.handleThemeCSS).Methods("GET")
}

func (b *Blog) handleSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

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
	vars := mux.Vars(r)
	theme := vars["theme"]

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

func PostHandler(postPath string, theme string) http.HandlerFunc {
	if theme == "" {
		theme = "default"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		post, err := parsePost(postPath)
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
