package main

import (
	"net/http"
	"web/goweb"
)

func home(c *goweb.Context) {
	// c.String(http.StatusOK, "hello world")
	c.JSON(http.StatusOK, map[string]int{
		"hello": 1,
	})
}

func hello(c *goweb.Context) {
	name := c.Query("name")
	c.String(http.StatusOK, name)
}

func login(c *goweb.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, name)
}

func getFile(c *goweb.Context) {
	name := c.Param("file")
	c.String(http.StatusOK, name)
}

func html(c *goweb.Context) {
	c.HTML(http.StatusOK, "<h1>Hello world</h1>")
}

func main() {
	engine := goweb.New()

	engine.GET("/", home)
	engine.GET("/hello", hello)
	engine.GET("/html", html)
	engine.GET("/getfile/*file", getFile)

	engine.Run(":4000")
}
