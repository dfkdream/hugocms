package main

import (
	"os"
	"testing"
	"time"
)

func TestFrontMatter_String(t *testing.T) {
	f := frontMatter{
		Title:           "Hello world!",
		Subtitle:        "",
		Date:            time.Now(),
		Author:          "John Doe",
		ShowReadingTime: true,
		ShowLanguages:   true,
		ShowAuthor:      true,
		ShowDate:        true,
	}
	t.Log(f)
}

func TestParseArticle(t *testing.T) {
	f, err := os.Open("./test/test.md")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fm, body, err := parseArticle(f)
	t.Log(fm)
	t.Log(body)
	t.Log(err)
}
