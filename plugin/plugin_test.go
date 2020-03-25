package plugin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func get(url string) io.Reader {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return res.Body
}

func TestPlugin_ServeHTTP(t *testing.T) {
	p := New(Info{
		Name:        "TestPlugin",
		Author:      "Test",
		Description: "Test Plugin",
		Version:     "0.0.1",
	}, "test")

	p.HandleAdminPage("/hello", "hello", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("Hello, world!"))
	}))

	s := httptest.NewServer(p)
	defer s.Close()

	var j Metadata
	err := json.NewDecoder(get(s.URL + "/metadata")).Decode(&j)
	if err != nil {
		t.Error(j)
		t.FailNow()
	}

	if !reflect.DeepEqual(j, Metadata{
		Identifier:     p.metadata.Identifier,
		Info:           p.metadata.Info,
		AdminMenuItems: []adminMenuItem{{"hello", "/hello"}},
	}) {
		t.Error("metadata not equals")
	}
}
