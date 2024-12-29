package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thylong/go-templates/05-gin-templ-htmx/web/view"
)

func CounterState(c *gin.Context) {
	c.Request.ParseForm()
	if c.Request.Form.Has("global") {
		global.Count++
	}

	c.HTML(http.StatusOK, "", view.Page(global.Count, 0))
}
