package main

import (
	"cache/gocache"
)

func main() {
	s := gocache.NewServer()

	g := gocache.NewGroup("test", 100, func(key string) (interface{}, error) {
		return key, nil
	})
	g.Set("tom", "say hello")

	s.Run(":4000")
}
