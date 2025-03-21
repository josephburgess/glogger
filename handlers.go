package glogger

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
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

func (b *Blog) GetCurrentTheme() string {
	return b.renderer.theme
}

func (b *Blog) handleDebug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeOverride := vars["theme"]

	if themeOverride != "" && !ValidateTheme(themeOverride) {
		http.Error(w, "Invalid theme name", http.StatusBadRequest)
		return
	}

	themeName := b.renderer.theme
	if themeOverride != "" {
		themeName = themeOverride
	}

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Glogger Theme Debug</title>
			<style>
				@import url("%s");
				body {
					font-family: "JetBrains Mono", monospace;
					max-width: 800px;
					margin: 0 auto;
					padding: 1rem;
					background-color: var(--background);
					color: var(--text);
				}
				a { color: var(--link); }
				a:hover { color: var(--link-hover); }
				pre { background-color: var(--code-bg); padding: 1rem; border-radius: 3px; }
				.muted { color: var(--muted); }
			</style>
		</head>
		<body>
			<h1>Glogger Theme Debugger</h1>
			<p>Current theme: <strong>%s</strong></p>
			<p class="muted">This text should be using the muted color.</p>
			<p>This is a <a href="#">link</a> using the theme colors.</p>
			<pre>This is a code block using the theme's code background.</pre>
			<h2>Available Themes</h2>
			<ul>
	`,
		GetThemePath(b.config.URLPrefix, themeName),
		themeName)

	for _, t := range availableThemes {
		html += fmt.Sprintf(`<li><a href="%s/_debug/%s">%s</a></li>`, b.config.URLPrefix, t, t)
	}

	html += `
			</ul>
			<div style="margin-top: 2rem;">
				<a href="javascript:history.back()">&larr; Back</a>
			</div>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (b *Blog) handleRawThemeCSS(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func (b *Blog) RegisterHandlers(router *mux.Router) {
	if len(b.posts) == 0 {
		err := b.initialize()
		if err != nil {
			panic("Failed to initialize blog: " + err.Error())
		}
	}

	blogRouter := router.PathPrefix(b.config.URLPrefix).Subrouter()

	// list route/index
	blogRouter.HandleFunc("", b.handleListPosts).Methods("GET")

	// individual posts
	blogRouter.HandleFunc("/{slug}", b.handleSinglePost).Methods("GET")

	// theme files
	blogRouter.HandleFunc("/_themes/{theme}.css", b.handleThemeCSS).Methods("GET")

	// debug
	blogRouter.HandleFunc("/_debug", b.handleDebug).Methods("GET")
	blogRouter.HandleFunc("/_debug/{theme}", b.handleDebug).Methods("GET")
	blogRouter.HandleFunc("/_raw_themes/{theme}.css", b.handleRawThemeCSS).Methods("GET")
}
