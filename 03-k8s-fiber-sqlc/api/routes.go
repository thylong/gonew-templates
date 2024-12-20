package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/pkg/db"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/pkg/handler"
)

func SetupRoutes(app *fiber.App, queries *db.Queries) {
	api := app.Group("/api/v1")
	// Post
	api.Post("/register", handler.NewAuthHandler(queries).SignUpUser)
}
