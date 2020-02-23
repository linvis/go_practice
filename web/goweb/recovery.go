package goweb

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Error("%s", message)
				c.Status(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
