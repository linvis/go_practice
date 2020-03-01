package gocache

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer()

	g := NewGroup("test", 100, func(key string) (interface{}, error) {
		return key, nil
	})
	g.Set("tom", "say hello")

	r, _ := http.NewRequest("GET", "/gocache/test/tom", nil)

	s.Run(":4000")
}
