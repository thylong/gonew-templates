package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/pkg/db"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/pkg/utils"
)

type AuthHandler struct {
	db *db.Queries
}

func NewAuthHandler(db *db.Queries) *AuthHandler {
	return &AuthHandler{db}
}

func (ac *AuthHandler) SignUpUser(ctx *fiber.Ctx) error {
	var credentials *db.User

	if err := ctx.BodyParser(&credentials); err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return nil
	}

	hashedPassword := utils.HashPassword(credentials.Password)

	args := &db.CreateUserParams{
		Name:      credentials.Name,
		Email:     credentials.Email,
		Password:  hashedPassword,
		Photo:     "default.jpeg",
		Verified:  true,
		Role:      "user",
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	user, err := ac.db.CreateUser(context.Background(), *args)
	if err != nil {
		ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		return nil
	}

	ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
	return nil
}
