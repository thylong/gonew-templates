package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thylong/go-templates/05-gin-templ-htmx/web/view"
)

type GlobalState struct {
	Count int
}

var global GlobalState

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "", view.Page(global.Count, 0))
}
