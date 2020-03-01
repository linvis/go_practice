package gocache

import "testing"

func TestGoCache(t *testing.T) {
	g := NewGroup("test", 100, func(key string) (interface{}, error) {
		return 1, nil
	})

	v, _ := g.Get("1")
	if v != 1 {
		t.Error("error")
	}

	v, _ = g.Get("1")
	if v != 1 {
		t.Error("error")
	}

	g.Set("2", 2)
	v, _ = g.Get("2")
	if v != 2 {
		t.Error("error")
	}
}
