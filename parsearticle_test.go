package main

import (
	"os"
	"testing"
	"time"
)

func TestFrontMatter_String(t *testing.T) {
	f := frontMatter{
		Title:           "hello world",
		Subtitle:        "",
		Date:            time.Now(),
		Author:          "John Doe",
		Attachments:     []string{"./hello.png", "./world.png"},
		ShowReadingTime: true,
		ShowLanguages:   true,
		ShowAuthor:      true,
		ShowDate:        true,
	}
	t.Log("\n", f)
}

func TestParseArticle(t *testing.T) {
	f, err := os.Open("./test/test.md")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	a, err := parseArticle(f)
	t.Log(a)
	t.Log(err)
}
