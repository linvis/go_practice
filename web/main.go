package main

import (
	"web/goweb"
)

func home(c *goweb.Context) {
	c.String("hello world")
}

func main() {
	engine := goweb.New()

	engine.GET("/", home)

	engine.Run(":4000")
}
