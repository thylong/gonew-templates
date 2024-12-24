package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/db"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/handlers"
)

func SetupRoutes(r *gin.Engine, queries *db.Queries) {
	// Healthcheck endpoint for k8s probes
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	router := api.Group("/auth")
	router.POST("/register", handlers.NewAuthHandler(queries).SignUpUser)
}
