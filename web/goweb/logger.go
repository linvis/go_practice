package goweb

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()

		log.Printf("[%d]  %s %v", c.StatusCode, c.Req.URL.Path, time.Since(t))

	}
}

func init() {
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}
