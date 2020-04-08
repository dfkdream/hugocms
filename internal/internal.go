package internal

import (
	"crypto/rand"
	"fmt"
	"strings"
)

func SingleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}

func GenerateRandomKey(bytes int) string {
	buff := make([]byte, bytes)
	_, _ = rand.Read(buff)
	return fmt.Sprintf("%x", buff)
}
