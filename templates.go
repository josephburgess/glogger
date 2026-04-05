package glogger

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed assets/templates/*.html
var templatesFS embed.FS

func newTemplateRenderer(config Config) (*templateRenderer, error) {
	postTmpl, err := template.ParseFS(templatesFS, "assets/templates/post.html")
	if err != nil {
		return nil, err
	}

	listTmpl, err := template.ParseFS(templatesFS, "assets/templates/list.html")
	if err != nil {
		return nil, err
	}

	return &templateRenderer{
		postTemplate: postTmpl,
		listTemplate: listTmpl,
		config:       config,
	}, nil
}

func (tr *templateRenderer) renderPost(post Post) (string, error) {
	data := PostTemplateData{
		Post:         post,
		BlogPrefix:   tr.config.URLPrefix,
		ThemeCSS:     getThemePath(tr.config.URLPrefix, tr.config.Theme),
		HighlightCSS: highlightJSStyleURL(tr.config.SyntaxTheme),
	}

	var buf bytes.Buffer
	if err := tr.postTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tr *templateRenderer) renderPostList(posts []Post, tag string) (string, error) {
	data := ListTemplateData{
		Posts:      posts,
		BlogPrefix: tr.config.URLPrefix,
		ThemeCSS:   getThemePath(tr.config.URLPrefix, tr.config.Theme),
		BlogTitle:  tr.config.Title,
		Tag:        tag,
	}

	var buf bytes.Buffer
	if err := tr.listTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
