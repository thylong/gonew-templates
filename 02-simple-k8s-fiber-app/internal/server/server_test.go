package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thylong/go-templates/02-simple-k8s-fiber-app/internal/core"
)

func TestCreateApp(t *testing.T) {
	tests := []struct {
		name         string
		httpTimeout  int64
		loggingLevel string
		production   bool
	}{
		{"Production mode with debug logging", 5000, "debug", true},
		{"Development mode with debug logging", 5000, "debug", false},
		{"Production mode without debug logging", 5000, "info", true},
		{"Development mode without debug logging", 5000, "info", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := CreateApp(tt.httpTimeout, tt.loggingLevel, tt.production)
			assert.NotNil(t, app)
			assert.IsType(t, &core.App{}, app)
			assert.IsType(t, &fiber.App{}, app.App)

			req := httptest.NewRequest("GET", "/healthz", nil)
			resp, err := app.App.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			req = httptest.NewRequest("GET", "/timeout", nil)
			time.Sleep(time.Duration(tt.httpTimeout+100) * time.Millisecond)
			resp, err = app.App.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		})
	}
}
