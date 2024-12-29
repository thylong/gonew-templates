package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutHandler(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "timeout")
}

// TimeoutMiddleware returns a timeoutMiddleware with given timeout
func TimeoutMiddleware(customTimeout int64) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Duration(customTimeout)*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutHandler),
	)
}
