package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type frontMatter struct {
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

type article struct {
	FrontMatter frontMatter `json:"frontMatter"`
	Body        string      `json:"body"`
}

func (f frontMatter) String() string {
	res, _ := json.MarshalIndent(f, "", "    ")
	return string(res)
}

func parseArticle(reader io.Reader) (*article, error) {
	var result frontMatter
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&result)
	if err != nil {
		return nil, err
	}
	bString := strings.TrimSpace(strings.TrimPrefix(string(body), result.String()))
	return &article{FrontMatter: result, Body: bString}, nil
}
