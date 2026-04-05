package glogger

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed assets/templates/*.html
var templatesFS embed.FS

func newTemplateRenderer(theme, urlPrefix, highlightCSSURL string) (*templateRenderer, error) {
	postTmpl, err := template.ParseFS(templatesFS, "assets/templates/post.html")
	if err != nil {
		return nil, err
	}

	listTmpl, err := template.ParseFS(templatesFS, "assets/templates/list.html")
	if err != nil {
		return nil, err
	}

	return &templateRenderer{
		postTemplate:    postTmpl,
		listTemplate:    listTmpl,
		theme:           theme,
		urlPrefix:       urlPrefix,
		highlightCSSURL: highlightCSSURL,
	}, nil
}

func (tr *templateRenderer) renderPost(post Post) (string, error) {
	data := PostTemplateData{
		Post:         post,
		BlogPrefix:   tr.urlPrefix,
		ThemeCSS:     getThemePath(tr.urlPrefix, tr.theme),
		HighlightCSS: tr.highlightCSSURL,
	}

	var buf bytes.Buffer
	if err := tr.postTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tr *templateRenderer) renderPostList(posts []Post, tag string) (string, error) {
	title := "Blog Posts"
	if tag != "" {
		title = "Posts tagged: " + tag
	}
	data := ListTemplateData{
		Posts:      posts,
		BlogPrefix: tr.urlPrefix,
		ThemeCSS:   getThemePath(tr.urlPrefix, tr.theme),
		Title:      title,
	}

	var buf bytes.Buffer
	if err := tr.listTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
