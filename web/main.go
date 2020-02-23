package main

import (
	"fmt"
	"net/http"
	"web/goweb"
)

func home(c *goweb.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
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
	c.HTML(http.StatusOK, "", "<h1>Hello world</h1>")
}

func panic(c *goweb.Context) {
	a := []int{1, 2, 3}
	fmt.Println(a[4])

	c.String(http.StatusOK, "test panic ok")
}

func main() {
	engine := goweb.New()

	engine.Use(goweb.Logger(), goweb.Recovery())
	engine.LoadHTMLGlob("static/*.html")
	engine.Static("/static", "./static")

	engine.GET("/", home)
	engine.GET("/hello", hello)
	engine.GET("/html", html)
	engine.GET("/getfile/*file", getFile)
	engine.GET("/panic", panic)

	v1 := engine.Group("/v1", func(c *goweb.Context) {
		c.String(http.StatusOK, "v1 middleware")
	})
	v1.GET("/v2", func(c *goweb.Context) {
		c.String(http.StatusOK, "v2 callback")
	})

	engine.Run(":4000")
}
