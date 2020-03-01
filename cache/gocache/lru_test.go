package gocache

import "testing"

func TestLRU(t *testing.T) {
	cache := NewLRUCache(100)

	cache.Add("1", 1)
	cache.Add("2", 2)
	cache.Add("3", 3)

	v, _ := cache.Back()
	if v != 3 {
		t.Error("error")
	}

	v, _ = cache.Get("1")
	if v != 1 {
		t.Error("error")
	}

	v, _ = cache.Back()
	if v != 1 {
		t.Error("error")
	}

}
