package glogger

import (
	"bytes"
	"html/template"
)

const (
	defaultPostTemplate = `
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
        pre {
            background-color: #f6f8fa;
            padding: 1rem;
            overflow: auto;
            border-radius: 3px;
        }
        code { font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace; }
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
    <a href="{{.BlogPrefix}}" class="back">&larr; Back to all posts</a>
</body>
</html>
`

	defaultListTemplate = `
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
                <a href="{{$.BlogPrefix}}/{{.Slug}}">{{.Title}}</a>
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
)

type templateRenderer struct {
	postTemplate *template.Template
	listTemplate *template.Template
}

func newTemplateRenderer() (*templateRenderer, error) {
	postTmpl, err := template.New("post").Parse(defaultPostTemplate)
	if err != nil {
		return nil, err
	}

	listTmpl, err := template.New("list").Parse(defaultListTemplate)
	if err != nil {
		return nil, err
	}

	return &templateRenderer{
		postTemplate: postTmpl,
		listTemplate: listTmpl,
	}, nil
}

func (tr *templateRenderer) renderPost(post Post, blogPrefix string) (string, error) {
	data := struct {
		Post
		BlogPrefix string
	}{
		Post:       post,
		BlogPrefix: blogPrefix,
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
	}{
		Posts:      posts,
		BlogPrefix: blogPrefix,
	}

	var buf bytes.Buffer
	if err := tr.listTemplate.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
