package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTest() (*fiber.App, string) {
	config.InitConfig()
	config.AppConfig.JWTSecret = "test_jwt_secret_for_middleware_test"

	app := fiber.New()

	token, err := config.GenerateJWT(1, "test@example.com", "TT")
	if err != nil {
		panic(err)
	}

	return app, token
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	app, token := setupAuthTest()

	app.Get("/protected", AuthMiddleware(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		email := c.Locals("email").(string)
		inisial := c.Locals("inisial").(string)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"email":   email,
			"inisial": inisial,
		})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"user_id":1`)
	assert.Contains(t, string(body), `"email":"test@example.com"`)
	assert.Contains(t, string(body), `"inisial":"TT"`)
}

func TestAuthMiddlewareNoToken(t *testing.T) {
	app, _ := setupAuthTest()

	app.Get("/protected", AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Missing or invalid JWT in request")
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	app, _ := setupAuthTest()

	app.Get("/protected", AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-that-is-clearly-wrong")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 401, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid or expired token")
}

func TestAuthMiddlewareBearerTokenWithoutPrefix(t *testing.T) {
	app, token := setupAuthTest()

	app.Get("/protected", AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user_id": c.Locals("user_id")})
	})

	// Send token directly without Bearer prefix
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthMiddlewareWithCookie(t *testing.T) {
	app, token := setupAuthTest()

	app.Get("/protected", AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user_id": c.Locals("user_id")})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Cookie", "token="+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"user_id":1`)
}
