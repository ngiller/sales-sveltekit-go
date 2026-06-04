package utils

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.StatusOK, "success")
	})
	app.Post("/test-create", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.StatusCreated, fiber.Map{"id": 1, "name": "Test"})
	})
	app.Get("/test-nil", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.StatusOK, nil)
	})
	app.Get("/test-slice", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.StatusOK, []int{1, 2, 3})
	})

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantData   string
	}{
		{
			name:       "OK with string",
			method:     "GET",
			path:       "/test",
			wantStatus: 200,
			wantData:   `"data":"success"`,
		},
		{
			name:       "Created with map",
			method:     "POST",
			path:       "/test-create",
			wantStatus: 201,
			wantData:   `"data":{"id":1,"name":"Test"}`,
		},
		{
			name:       "OK with nil",
			method:     "GET",
			path:       "/test-nil",
			wantStatus: 200,
			wantData:   `"data":null`,
		},
		{
			name:       "OK with slice",
			method:     "GET",
			path:       "/test-slice",
			wantStatus: 200,
			wantData:   `"data":[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), tt.wantData)
		})
	}
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/bad-request", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	})
	app.Get("/unauthorized", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid JWT in request")
	})
	app.Get("/forbidden", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusForbidden, "Access Denied")
	})
	app.Get("/not-found", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusNotFound, "Resource not found")
	})
	app.Get("/internal-error", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Something went wrong")
	})
	app.Get("/validation-error", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, []interface{}{
			fiber.Map{"field": "email", "tag": "required"},
		})
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "Bad request",
			path:       "/bad-request",
			wantStatus: 400,
			wantMsg:    "Invalid input",
		},
		{
			name:       "Unauthorized",
			path:       "/unauthorized",
			wantStatus: 401,
			wantMsg:    "Missing or invalid JWT in request",
		},
		{
			name:       "Forbidden",
			path:       "/forbidden",
			wantStatus: 403,
			wantMsg:    "Access Denied",
		},
		{
			name:       "Not found",
			path:       "/not-found",
			wantStatus: 404,
			wantMsg:    "Resource not found",
		},
		{
			name:       "Internal error",
			path:       "/internal-error",
			wantStatus: 500,
			wantMsg:    "Something went wrong",
		},
		{
			name:       "Validation error with array",
			path:       "/validation-error",
			wantStatus: 400,
			wantMsg:    `[{"field":"email","tag":"required"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), tt.wantMsg)
			assert.NotContains(t, string(body), `"data"`)
		})
	}
}

func TestResponseContentType(t *testing.T) {
	app := fiber.New()

	app.Get("/success", func(c *fiber.Ctx) error {
		return SuccessResponse(c, fiber.StatusOK, fiber.Map{"key": "value"})
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, "error")
	})

	for _, path := range []string{"/success", "/error"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			ct := resp.Header.Get("Content-Type")
			assert.Contains(t, strings.ToLower(ct), "application/json")
		})
	}
}
