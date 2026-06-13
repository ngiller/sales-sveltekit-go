package middleware

import (
	"io"
	"net/http/httptest"
	"regexp"
	"testing"

	"backend/config"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type policyTest struct {
	app   *fiber.App
	mock  sqlmock.Sqlmock
	db    *gorm.DB
	token string
}

func setupPolicyTest(t *testing.T, groupID interface{}) policyTest {
	config.InitConfig()
	config.AppConfig.JWTSecret = "test_secret"

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	app := fiber.New()

	var jwtGroupID *uint
	if gID, ok := groupID.(int); ok {
		v := uint(gID)
		jwtGroupID = &v
	}
	token, err := config.GenerateJWT(1, "admin@test.com", "AD", jwtGroupID)

	mock.ExpectQuery("SELECT `user_group_id` FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"user_group_id"}).AddRow(groupID))

	return policyTest{app, mock, db, token}
}

func TestRequirePolicyAdminBypass(t *testing.T) {
	pt := setupPolicyTest(t, 1)

	pt.app.Get("/api/users", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "OK", string(body))

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}

func TestRequirePolicyNoUserGroup(t *testing.T) {
	pt := setupPolicyTest(t, nil)

	pt.app.Get("/api/users", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 403, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Access Denied")

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}

func TestRequirePolicyEndpointNotFound(t *testing.T) {
	pt := setupPolicyTest(t, 2)

	pt.mock.ExpectQuery("SELECT \\* FROM `menu_access` WHERE endpoint = \\? ORDER BY `menu_access`.`id` LIMIT 1").
		WithArgs("users").
		WillReturnError(gorm.ErrRecordNotFound)

	pt.app.Get("/api/users", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 403, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Endpoint not registered")

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}

func TestRequirePolicyAccessDenied(t *testing.T) {
	pt := setupPolicyTest(t, 2)

	pt.mock.ExpectQuery("SELECT \\* FROM `menu_access` WHERE endpoint = \\? ORDER BY `menu_access`.`id` LIMIT 1").
		WithArgs("users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Users"))

	pt.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `group_policies` WHERE group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)")).
		WithArgs(uint(2), "read", "Users", uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	pt.app.Get("/api/users", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 403, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Access Denied")

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}

func TestRequirePolicyAccessGranted(t *testing.T) {
	pt := setupPolicyTest(t, 2)

	pt.mock.ExpectQuery("SELECT \\* FROM `menu_access` WHERE endpoint = \\? ORDER BY `menu_access`.`id` LIMIT 1").
		WithArgs("users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Users"))

	pt.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `group_policies` WHERE group_id = ? AND action = ? AND (table_name = ? OR table_id = ?)")).
		WithArgs(uint(2), "read", "Users", uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	pt.app.Get("/api/users", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("Access Granted")
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Access Granted", string(body))

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}

func TestRequirePolicyInvalidPathFormat(t *testing.T) {
	pt := setupPolicyTest(t, 2)

	pt.app.Get("/short", AuthMiddleware(), RequirePolicy(pt.db, "read"), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest("GET", "/short", nil)
	req.Header.Set("Authorization", "Bearer "+pt.token)

	resp, err := pt.app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 403, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid API path format")

	assert.NoError(t, pt.mock.ExpectationsWereMet())
}
