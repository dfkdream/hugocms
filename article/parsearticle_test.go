package article

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	dummyFront = `{
    "title": "hello world",
    "subtitle": "",
    "date": "2019-12-03T00:00:00Z",
    "author": "John Doe",
    "attachments": [
        "hello_world.jpg"
    ],
    "showReadingTime": true,
    "showLanguages": true,
    "showAuthor": true,
    "showDate": true
}`
	dummyArticle = "hello world"
	dummyFile    = dummyFront + "\n" + dummyArticle
	f            = FrontMatter{
		Title:           "hello world",
		Subtitle:        "",
		Date:            MustParseTime(time.Parse(time.RFC3339, "2019-12-03T00:00:00Z")),
		Author:          "John Doe",
		Attachments:     []string{"hello_world.jpg"},
		ShowReadingTime: true,
		ShowLanguages:   true,
		ShowAuthor:      true,
		ShowDate:        true,
	}
)

func MustParseTime(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func TestFrontMatter_String(t *testing.T) {
	if f.String() != dummyFront {
		t.Error("string does not matches")
	}
}

func TestParseArticle(t *testing.T) {
	a, err := Parse(strings.NewReader(dummyFile))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(f, a.FrontMatter) {
		t.Error("frontMatter does not matches")
	}

	if a.Body != dummyArticle {
		t.Error("article does not matches")
	}
}
