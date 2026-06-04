package repository

import (
	"fmt"
	"regexp"
	"testing"

	"backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupUserRepoTest(t *testing.T) (*UserRepository, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := NewUserRepository(db)
	return repo, mock
}

func TestUserRepositoryGetDB(t *testing.T) {
	repo, _ := setupUserRepoTest(t)
	assert.NotNil(t, repo.GetDB())
}

func TestUserRepositoryFindByEmailFound(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT 1").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "user_group_id", "enable"}).
			AddRow(1, "Test User", "test@example.com", 1, true))

	mock.ExpectQuery("SELECT \\* FROM `user_groups` WHERE `user_groups`.`id` = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Admin"))

	user, err := repo.FindByEmail("test@example.com")
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "Admin", user.RoleName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindByEmailNotFound(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? ORDER BY `users`.`id` LIMIT 1").
		WithArgs("notfound@example.com").
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail("notfound@example.com")
	require.NoError(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindByIDFound(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(uint(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "user_group_id", "enable"}).
			AddRow(5, "John Doe", "john@example.com", 2, true))

	mock.ExpectQuery("SELECT \\* FROM `user_groups` WHERE `user_groups`.`id` = \\?").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "User"))

	user, err := repo.FindByID(5)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, uint(5), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "User", user.RoleName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindByIDNotFound(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? ORDER BY `users`.`id` LIMIT 1").
		WithArgs(uint(999)).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByID(999)
	require.NoError(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryCreate(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &models.User{
		Name:  "New User",
		Email: "new@example.com",
	}

	err := repo.Create(user)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdate(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := &models.User{
		ID:    1,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	err := repo.Update(user)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryDelete(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = \\?").
		WithArgs(uint(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindAllWithSearch(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	joinSQL := "FROM `users` left join user_groups on user_groups.id = users.user_group_id left join master_departements on master_departements.id = users.departement_id"
	whereSQL := " WHERE users.name LIKE ? OR users.email LIKE ?"
	orderSQL := " ORDER BY users.name asc"
	limit := 50

	// Count query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) "+joinSQL+whereSQL)).
		WithArgs("%search%", "%search%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Data query — Scan() substitutes LIMIT directly
	selectSQL := "SELECT users.*, user_groups.name as role_name, master_departements.name as dept_name " + joinSQL + whereSQL + orderSQL + fmt.Sprintf(" LIMIT %d", limit)
	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WithArgs("%search%", "%search%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "role_name"}).
			AddRow(1, "Search Result", "search@example.com", "Admin"))

	users, total, err := repo.FindAll("search", 1, 50, "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, users, 1)
	assert.Equal(t, "Search Result", users[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindAllNoSearch(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	joinSQL := "FROM `users` left join user_groups on user_groups.id = users.user_group_id left join master_departements on master_departements.id = users.departement_id"
	orderSQL := " ORDER BY users.name asc"
	limit := 50

	// Count query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) " + joinSQL)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// Data query
	selectSQL := "SELECT users.*, user_groups.name as role_name, master_departements.name as dept_name " + joinSQL + orderSQL + fmt.Sprintf(" LIMIT %d", limit)
	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "Alice", "alice@test.com").
			AddRow(2, "Bob", "bob@test.com").
			AddRow(3, "Charlie", "charlie@test.com"))

	users, total, err := repo.FindAll("", 1, 50, "", "")
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, users, 3)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryFindAllWithSortAndOffset(t *testing.T) {
	repo, mock := setupUserRepoTest(t)

	joinSQL := "FROM `users` left join user_groups on user_groups.id = users.user_group_id left join master_departements on master_departements.id = users.departement_id"
	orderSQL := " ORDER BY users.email desc"
	limit := 10
	offset := 20

	// Count query
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) " + joinSQL)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Data query with custom sort and offset
	selectSQL := "SELECT users.*, user_groups.name as role_name, master_departements.name as dept_name " + joinSQL + orderSQL + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(3, "Zoe", "zoe@test.com"))

	users, total, err := repo.FindAll("", 3, 10, "email", "desc")
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	require.Len(t, users, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}
