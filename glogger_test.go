package glogger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsePost(t *testing.T) {
	md := newMarkdown()

	t.Run("valid frontmatter", func(t *testing.T) {
		f := writeTempPost(t, `---
title: "Hello World"
date: 2025-01-15
description: "A test post"
tags: [go, test]
draft: false
---

Post body here.
`)
		post, err := parsePost(f, md)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if post.Title != "Hello World" {
			t.Errorf("title: got %q, want %q", post.Title, "Hello World")
		}
		if post.Description != "A test post" {
			t.Errorf("description: got %q, want %q", post.Description, "A test post")
		}
		if post.PublishDate.IsZero() {
			t.Error("expected non-zero publish date")
		}
		if len(post.Tags) != 2 {
			t.Errorf("tags: got %d, want 2", len(post.Tags))
		}
		if post.Draft {
			t.Error("expected draft=false")
		}
		if !strings.Contains(string(post.Content), "Post body here") {
			t.Error("expected body content in rendered HTML")
		}
	})

	t.Run("missing title default to Untitled Post", func(t *testing.T) {
		f := writeTempPost(t, "---\ndate: 2025-01-01\n---\n\nBody.\n")
		post, err := parsePost(f, md)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if post.Title != "Untitled Post" {
			t.Errorf("got %q, want %q", post.Title, "Untitled Post")
		}
	})

	t.Run("missing frontmatter delimiter returns error", func(t *testing.T) {
		f := writeTempPost(t, "# no frontmatter here\n\nblah blah.\n")
		_, err := parsePost(f, md)
		if err == nil {
			t.Error("expected error for missing frontmatter")
		}
	})

	t.Run("incomplete frontmatter returns error", func(t *testing.T) {
		f := writeTempPost(t, "---\ntitle: only one delimiter\n")
		_, err := parsePost(f, md)
		if err == nil {
			t.Error("expected error for incomplete frontmatter")
		}
	})

	t.Run("invalid date is ignored", func(t *testing.T) {
		f := writeTempPost(t, "---\ntitle: Test\ndate: not-a-date\n---\n\nBody.\n")
		post, err := parsePost(f, md)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !post.PublishDate.IsZero() {
			t.Error("expected no publish date for invalid date string")
		}
	})

	t.Run("draft field is parsed", func(t *testing.T) {
		f := writeTempPost(t, "---\ntitle: Draft\ndraft: true\n---\n\nBody.\n")
		post, err := parsePost(f, md)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !post.Draft {
			t.Error("expected draft=true")
		}
	})
}

// blog init

func TestInitialize_FiltersDrafts(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "published.md", "---\ntitle: Published\ndate: 2025-01-01\ndraft: false\n---\n\nContent.\n")
	writePost(t, dir, "draft.md", "---\ntitle: Draft\ndate: 2025-01-02\ndraft: true\n---\n\nContent.\n")

	blog, err := New(Config{ContentDir: dir, URLPrefix: "/blog", Theme: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	posts := blog.GetPosts()
	if len(posts) != 1 {
		t.Fatalf("got %d posts, want 1", len(posts))
	}
	if posts[0].Title != "Published" {
		t.Errorf("got %q, want %q", posts[0].Title, "Published")
	}
}

func TestInitialize_SortDateDescending(t *testing.T) {
	dir := t.TempDir()

	writePost(t, dir, "older.md", "---\ntitle: Older\ndate: 2024-01-01\n---\n\nContent.\n")
	writePost(t, dir, "newer.md", "---\ntitle: Newer\ndate: 2025-06-01\n---\n\nContent.\n")

	blog, err := New(Config{ContentDir: dir, URLPrefix: "/blog", Theme: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	posts := blog.GetPosts()
	if len(posts) != 2 {
		t.Fatalf("got %d posts, want 2", len(posts))
	}
	if posts[0].Title != "Newer" {
		t.Errorf("first post: got %q, want %q", posts[0].Title, "Newer")
	}
}

// themes

func TestValidateTheme(t *testing.T) {
	valid := []string{"default", "dark", "light", "rosepine"}
	for _, theme := range valid {
		if !validateTheme(theme) {
			t.Errorf("expected %q to be valid", theme)
		}
	}
	if validateTheme("invalid") {
		t.Error("expected \"invalid\" to be invalid")
	}
}

func TestDefaultSyntaxTheme(t *testing.T) {
	cases := []struct {
		theme  string
		syntax string
	}{
		{"rosepine", "rose-pine"},
		{"dark", "github-dark"},
		{"light", "github"},
		{"default", "github"},
		{"unknown", "github"},
	}
	for _, c := range cases {
		got := defaultSyntaxTheme(c.theme)
		if got != c.syntax {
			t.Errorf("defaultSyntaxTheme(%q): got %q, want %q", c.theme, got, c.syntax)
		}
	}
}

//handlers

func TestHandler_ListPosts(t *testing.T) {
	blog := blogWithPosts(t, "---\ntitle: Test Post\ndate: 2025-01-01\n---\n\nContent.\n")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "Test Post") {
		t.Error("expected post title in response body")
	}
}

func TestHandler_TaggedPosts(t *testing.T) {
	dir := t.TempDir()
	writePost(t, dir, "tagged.md", "---\ntitle: Tagged Post\ndate: 2025-01-01\ntags: [go, test]\n---\n\nContent.\n")
	writePost(t, dir, "other.md", "---\ntitle: Other Post\ndate: 2025-01-02\ntags: [rust]\n---\n\nContent.\n")

	blog, err := New(Config{ContentDir: dir, URLPrefix: "/blog", Theme: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest("GET", "/_tags/go", nil)
	req.SetPathValue("tag", "go")
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Tagged Post") {
		t.Error("expected tagged post in response")
	}
	if strings.Contains(body, "Other Post") {
		t.Error("expected other post to be excluded")
	}
	if !strings.Contains(body, "Posts tagged: go") {
		t.Error("expected tag title in response")
	}
}

func TestHandler_SinglePost(t *testing.T) {
	blog := blogWithPosts(t, "---\ntitle: Hello\ndate: 2025-01-01\n---\n\nBody content.\n")

	req := httptest.NewRequest("GET", "/hello", nil)
	req.SetPathValue("slug", "hello")
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "Body content") {
		t.Error("expected post content in response body")
	}
}

func TestHandler_SinglePost_NotFound(t *testing.T) {
	blog := blogWithPosts(t, "---\ntitle: Hello\ndate: 2025-01-01\n---\n\nContent.\n")

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	req.SetPathValue("slug", "nonexistent")
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_ThemeCSS(t *testing.T) {
	blog := blogWithPosts(t, "---\ntitle: Post\ndate: 2025-01-01\n---\n\nContent.\n")

	req := httptest.NewRequest("GET", "/_themes/default", nil)
	req.SetPathValue("theme", "default")
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/css" {
		t.Errorf("Content-Type: got %q, want %q", ct, "text/css")
	}
}

func TestHandler_ThemeCSS_InvalidTheme(t *testing.T) {
	blog := blogWithPosts(t, "---\ntitle: Post\ndate: 2025-01-01\n---\n\nContent.\n")

	req := httptest.NewRequest("GET", "/_themes/invalid", nil)
	req.SetPathValue("theme", "invalid")
	w := httptest.NewRecorder()
	blog.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

// helpers

func writeTempPost(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.md")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func writePost(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writing post %s: %v", name, err)
	}
}

func blogWithPosts(t *testing.T, postContent string) *Blog {
	t.Helper()
	dir := t.TempDir()
	writePost(t, dir, "hello.md", postContent)
	blog, err := New(Config{ContentDir: dir, URLPrefix: "/blog", Theme: "default"})
	if err != nil {
		t.Fatalf("creating blog: %v", err)
	}
	return blog
}
