package lru

import "testing"

type String string

func (s String) Len() int {
	return len(s)
}

func TestLRU(t *testing.T) {
	cache := New(100)

	cache.Add("1", String("1"))
	cache.Add("2", String("2"))
	cache.Add("3", String("3"))

	v, _ := cache.Back()
	if v.(String) != "3" {
		t.Error("error")
	}

	v, _ = cache.Get("1")
	if v.(String) != "1" {
		t.Error("error")
	}

	v, _ = cache.Back()
	if v.(String) != "1" {
		t.Error("error")
	}

}
