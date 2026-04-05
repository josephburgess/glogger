package glogger

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

// Handler returns an http.Handler that serves the blog.
// Mount with http.StripPrefix:
//
//	http.Handle("/blog/", http.StripPrefix("/blog", blog.Handler()))
func (b *Blog) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", b.handleListPosts)
	mux.HandleFunc("GET /feed.xml", b.handleFeed)
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

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type rssChannel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel
}

func (b *Blog) handleFeed(w http.ResponseWriter, r *http.Request) {
	baseURL := strings.TrimRight(b.config.BaseURL, "/")

	items := make([]rssItem, 0, len(b.posts))
	for _, post := range b.posts {
		link := baseURL + b.config.URLPrefix + "/" + post.Slug
		pubDate := ""
		if !post.PublishDate.IsZero() {
			pubDate = post.PublishDate.UTC().Format(time.RFC1123Z)
		}
		items = append(items, rssItem{
			Title:       post.Title,
			Link:        link,
			Description: post.Description,
			PubDate:     pubDate,
			GUID:        link,
		})
	}

	lastBuild := ""
	if len(b.posts) > 0 && !b.posts[0].PublishDate.IsZero() {
		lastBuild = b.posts[0].PublishDate.UTC().Format(time.RFC1123Z)
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         b.config.Title,
			Link:          baseURL + b.config.URLPrefix,
			Description:   b.config.Description,
			LastBuildDate: lastBuild,
			Items:         items,
		},
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		http.Error(w, "Error generating feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	w.Write(out)
}

// PostHandler returns a standalone handler for rendering a single markdown file.
// Useful for serving a specific post outside the blog structure.
func PostHandler(postPath string, theme string) http.HandlerFunc {
	cfg := Config{Theme: theme}
	cfg.setDefaults()

	renderer, err := newTemplateRenderer(cfg)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Error creating renderer: "+err.Error(), http.StatusInternalServerError)
		}
	}

	post, err := parsePost(postPath, newMarkdown())
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Error parsing post: "+err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		html, err := renderer.renderPost(post)
		if err != nil {
			http.Error(w, "Error rendering post: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}
