package article

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type FrontMatter struct {
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle"`
	Date            time.Time `json:"date"`
	Author          string    `json:"author"`
	Attachments     []string  `json:"attachments"`
	ShowReadingTime bool      `json:"showReadingTime"`
	ShowLanguages   bool      `json:"showLanguages"`
	ShowAuthor      bool      `json:"showAuthor"`
	ShowDate        bool      `json:"showDate"`
}

type Article struct {
	FrontMatter FrontMatter `json:"frontMatter"`
	Body        string      `json:"body"`
}

func (f FrontMatter) String() string {
	res, _ := json.MarshalIndent(f, "", "    ")
	return string(res)
}

func Parse(reader io.Reader) (*Article, error) {
	var result FrontMatter
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&result)
	if err != nil {
		return nil, err
	}
	bString := strings.TrimSpace(strings.TrimPrefix(string(body), result.String()))
	return &Article{FrontMatter: result, Body: bString}, nil
}
