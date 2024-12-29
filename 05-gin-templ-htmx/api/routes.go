package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	view "github.com/thylong/go-templates/05-gin-templ-htmx/web/view"
)

func SetupRoutes(r *gin.Engine) {
	// Healthcheck endpoint for k8s probes
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", view.Counter("testong"))
	})
}
