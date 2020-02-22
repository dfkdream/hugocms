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
	ShowReadingTime bool      `json:"showReadingTime"`
	ShowLanguages   bool      `json:"showLanguages"`
	ShowAuthor      bool      `json:"showAuthor"`
	ShowDate        bool      `json:"showDate"`
}

func (f frontMatter) String() string {
	res, _ := json.MarshalIndent(f, "", "    ")
	return string(res)
}

func parseArticle(reader io.Reader) (frontMatter, string, error) {
	var result frontMatter
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return frontMatter{}, "", err
	}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&result)
	if err != nil {
		return frontMatter{}, "", err
	}
	article := strings.TrimSpace(strings.TrimPrefix(string(body), result.String()))
	return result, article, nil
}
