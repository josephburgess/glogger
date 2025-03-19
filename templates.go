package glogger

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed assets/templates/*.html
var templatesFS embed.FS

type templateRenderer struct {
	postTemplate *template.Template
	listTemplate *template.Template
	theme        string
	urlPrefix    string
}

func newTemplateRenderer(theme string, urlPrefix string) (*templateRenderer, error) {
	// parse post template html
	postTmpl, err := template.ParseFS(templatesFS, "assets/templates/post.html")
	if err != nil {
		return nil, err
	}

	// parse list template html
	listTmpl, err := template.ParseFS(templatesFS, "assets/templates/list.html")
	if err != nil {
		return nil, err
	}

	return &templateRenderer{
		postTemplate: postTmpl,
		listTemplate: listTmpl,
		theme:        theme,
		urlPrefix:    urlPrefix,
	}, nil
}

func (tr *templateRenderer) renderPost(post Post, blogPrefix string) (string, error) {
    data := struct {
        Post
        BlogPrefix string
        ThemeCSS   string
    }{
        Post:       post,
        BlogPrefix: blogPrefix,
        ThemeCSS:   GetThemePath(blogPrefix, tr.theme),
    }

    var buf bytes.Buffer
    if err := tr.postTemplate.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}
func (tr *templateRenderer) renderPostList(posts []Post, blogPrefix string) (string, error) {
    data := struct {
        Posts      []Post
        BlogPrefix string
        ThemeCSS   string
        Title      string
    }{
        Posts:      posts,
        BlogPrefix: blogPrefix,
        ThemeCSS:   GetThemePath(blogPrefix, tr.theme),
        Title:      "Blog Posts", // default title
    }

    var buf bytes.Buffer
    if err := tr.listTemplate.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}
