package plugin

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/proto"
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

	body, err := ioutil.ReadAll(get(s.URL + "/metadata"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	j := new(Metadata)

	err = proto.Unmarshal(body, j)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !(j.Identifier == "test" || j.Info.Name == "TestPlugin" ||
		j.Info.Author == "Test" || j.Info.Author == "Test Plugin" ||
		j.Info.Version == "1.0.0" || j.AdminMenuItems[0].MenuName == "hello" ||
		j.AdminMenuItems[0].Endpoint == "/hello") {
		t.Error("metadata not equals")
	}
}
