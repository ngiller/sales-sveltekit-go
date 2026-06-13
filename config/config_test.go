package config

import (
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	origVars := map[string]string{}
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_STOCK_NAME", "JWT_SECRET"} {
		origVars[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range origVars {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_STOCK_NAME", "test_stock")
	os.Setenv("JWT_SECRET", "test_secret_key")

	InitConfig()

	require.NotNil(t, AppConfig)
	assert.Equal(t, "testhost", AppConfig.DBHost)
	assert.Equal(t, "3307", AppConfig.DBPort)
	assert.Equal(t, "testuser", AppConfig.DBUser)
	assert.Equal(t, "testpass", AppConfig.DBPassword)
	assert.Equal(t, "testdb", AppConfig.DBName)
	assert.Equal(t, "test_stock", AppConfig.DBStockName)
	assert.Equal(t, "test_secret_key", AppConfig.JWTSecret)
}

func TestInitConfigDefaults(t *testing.T) {
	origVars := map[string]string{}
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_STOCK_NAME", "JWT_SECRET"} {
		origVars[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range origVars {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_STOCK_NAME", "JWT_SECRET"} {
		os.Unsetenv(k)
	}

	InitConfig()

	require.NotNil(t, AppConfig)
	assert.Equal(t, "localhost", AppConfig.DBHost)
	assert.Equal(t, "3306", AppConfig.DBPort)
	assert.Equal(t, "root", AppConfig.DBUser)
	assert.Equal(t, "Pass@w0rd", AppConfig.DBPassword)
	assert.Equal(t, "magnum_sales_svelte_go", AppConfig.DBName)
	assert.Equal(t, "magnum_stock_db", AppConfig.DBStockName)
	assert.Equal(t, "magnum_secret_key_2024", AppConfig.JWTSecret)
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "securePassword123!"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	assert.True(t, CheckPassword(hash, password))
	assert.False(t, CheckPassword(hash, "wrongPassword"))
	assert.False(t, CheckPassword(hash, ""))
}

func TestJWTLifecycle(t *testing.T) {
	InitConfig()
	AppConfig.JWTSecret = "test_jwt_secret_for_testing"

	var groupID uint = 1
	token, err := GenerateJWT(42, "test@example.com", "TT", &groupID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateJWT(token)
	require.NoError(t, err)
	assert.NotNil(t, claims)

	assert.InDelta(t, float64(42), claims["user_id"], 0)
	assert.Equal(t, "test@example.com", claims["email"])
	assert.Equal(t, "TT", claims["inisial"])

	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.Greater(t, exp, float64(0))

	iat, ok := claims["iat"].(float64)
	assert.True(t, ok)
	assert.Greater(t, iat, float64(0))
}

func TestValidateInvalidJWT(t *testing.T) {
	InitConfig()
	AppConfig.JWTSecret = "secret"

	_, err := ValidateJWT("malformed-token")
	assert.Error(t, err)

	_, err = ValidateJWT("")
	assert.Error(t, err)

	_, err = ValidateJWT("eyJhbGciOiJSUzI1NiJ9.dGhpcyBpcyBmYWtl.sig")
	assert.Error(t, err)
}

func TestJWTWrongSecret(t *testing.T) {
	InitConfig()
	AppConfig.JWTSecret = "original_secret"

	token, err := GenerateJWT(1, "a@b.com", "A", nil)
	require.NoError(t, err)

	AppConfig.JWTSecret = "different_secret"
	_, err = ValidateJWT(token)
	assert.Error(t, err)

	AppConfig.JWTSecret = "original_secret"
}

func TestGetEnv(t *testing.T) {
	assert.Equal(t, "default", getEnv("NONEXISTENT_VAR_12345_TEST", "default"))
	assert.Equal(t, "", getEnv("NONEXISTENT_VAR_12345_TEST", ""))
}

func TestErrorHandlerFiberError(t *testing.T) {
	InitConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Resource not found")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Resource not found")
}

func TestErrorHandlerGenericError(t *testing.T) {
	InitConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fmt.Errorf("something unexpected happened")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "something unexpected happened")
}


