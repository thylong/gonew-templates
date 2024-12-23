package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/db"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/utils"
)

type AuthHandler struct {
	db *db.Queries
}

func NewAuthHandler(db *db.Queries) *AuthHandler {
	return &AuthHandler{db}
}

// SignUpUser godoc
// @Summary      Signup flow for a new user
// @Description  Create a user entry based on given data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email       body    string     true  "user email"
// @Param        name        body    string     true  "username"
// @Param        password    body    string     true  "user password"
// @Success      200  {object}  db.User
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /auth/register [post]
func (ac *AuthHandler) SignUpUser(ctx *gin.Context) {
	var credentials *db.User

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
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

	user, err := ac.db.CreateUser(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user}})
}
