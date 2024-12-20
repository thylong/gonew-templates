package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/pkg/db"
)

func TestSetupRoutes(t *testing.T) {
	app := fiber.New()
	queries := &db.Queries{}
	SetupRoutes(app, queries)

	tests := []struct {
		description  string
		method       string
		route        string
		expectedCode int
	}{
		{
			description:  "register user",
			method:       "POST",
			route:        "/api/v1/register",
			expectedCode: fiber.StatusOK,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedCode, resp.StatusCode)
	}
}
