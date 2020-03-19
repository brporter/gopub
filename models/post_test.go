package models

import (
	"fmt"
	"testing"
)

func TestNewPost(t *testing.T) {
	p, _ := NewPost("", "")

	if p == nil {
		t.Error("Expected type of post.")
	}
}

func TestToHtml(t *testing.T) {
	p, err := NewPost("A Title", "A Body")

	if err != nil {
		t.Errorf("Error creating post: %v", err)
	}

	html := p.ToHtml()

	if *html == "" {
		t.Errorf("Expected an HTML body, got an empty string (%s) instead.", *html)
	}
}

func TestToJson(t *testing.T) {
	p, err := NewPost("A Title", "A Body")

	if err != nil {
		t.Errorf("Error creating post: %v", err)
	}

	encoded := p.ToJson()

	fmt.Println(*encoded)

	result, err := FromJson(encoded)

	if err != nil {
		t.Errorf("Error decoding JSON: %v", err)
	}

	if result == nil {
		t.Error("Expected post, got nil")
	}

	if p.Title != result.Title || p.Body != result.Body {
		t.Errorf("Expected title: %s and body: %s, got title: %s and body: %s", p.Title, p.Body, result.Title, result.Body)
	}
}
