package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthz(t *testing.T) {
	// Create a new Fiber instance
	app := fiber.New()

	// Define a test route
	app.Get("/healthz", Healthz)

	// Perform a GET request to /healthz
	resp, err := app.Test(httptest.NewRequest("GET", "/healthz", nil))

	// Verify the status code
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("Expected status code 200, but got %d", resp.StatusCode)
	}
}
