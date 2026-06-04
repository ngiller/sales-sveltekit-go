package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	Name     string `validate:"required"`
	Age      int    `validate:"gte=0,lte=150"`
}

func TestValidateStructValid(t *testing.T) {
	data := testStruct{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Age:      25,
	}

	errors := ValidateStruct(data)
	assert.Empty(t, errors)
}

func TestValidateStructInvalid(t *testing.T) {
	data := testStruct{
		Email:    "not-an-email",
		Password: "123",
		Name:     "",
		Age:      -1,
	}

	errors := ValidateStruct(data)
	assert.NotEmpty(t, errors)
	assert.Len(t, errors, 4)

	// Verify each field has a validation error
	fieldErrors := make(map[string]string)
	for _, e := range errors {
		fieldErrors[e.FailedField] = e.Tag
	}

	assert.Equal(t, "email", fieldErrors["testStruct.Email"])
	assert.Equal(t, "min", fieldErrors["testStruct.Password"])
	assert.Equal(t, "required", fieldErrors["testStruct.Name"])
	assert.Equal(t, "gte", fieldErrors["testStruct.Age"])
}

func TestValidateStructEmptyStruct(t *testing.T) {
	errors := ValidateStruct(testStruct{})
	assert.NotEmpty(t, errors)
	assert.Len(t, errors, 3) // email, password, name are required; age has default 0 which is valid
}

func TestValidateStructNil(t *testing.T) {
	errors := ValidateStruct(nil)
	assert.Empty(t, errors) // nil struct returns InvalidValidationError, which we ignore
}

func TestValidateStructWithSpecialFields(t *testing.T) {
	type customStruct struct {
		Code  string `validate:"required,min=3,max=10"`
		Value int    `validate:"required"`
	}

	t.Run("valid", func(t *testing.T) {
		errors := ValidateStruct(customStruct{Code: "ABC-123", Value: 42})
		assert.Empty(t, errors)
	})

	t.Run("invalid code too short", func(t *testing.T) {
		errors := ValidateStruct(customStruct{Code: "AB", Value: 42})
		require.Len(t, errors, 1)
		assert.Equal(t, "min", errors[0].Tag)
	})

	t.Run("invalid code too long", func(t *testing.T) {
		errors := ValidateStruct(customStruct{Code: "ABCDEFGHIJKLMN", Value: 42})
		require.Len(t, errors, 1)
		assert.Equal(t, "max", errors[0].Tag)
	})
}

// The ErrorResponseValidation struct just has fields, no methods to test
func TestErrorResponseValidationStruct(t *testing.T) {
	err := ErrorResponseValidation{
		FailedField: "Email",
		Tag:         "required",
		Value:       "",
	}

	assert.Equal(t, "Email", err.FailedField)
	assert.Equal(t, "required", err.Tag)
	assert.Equal(t, "", err.Value)
}

func TestLoginRequestValidation(t *testing.T) {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	t.Run("valid login request", func(t *testing.T) {
		req := LoginRequest{Email: "admin@test.com", Password: "admin123"}
		errors := ValidateStruct(req)
		assert.Empty(t, errors)
	})

	t.Run("missing email", func(t *testing.T) {
		req := LoginRequest{Email: "", Password: "password123"}
		errors := ValidateStruct(req)
		require.Len(t, errors, 1)
		assert.Equal(t, "required", errors[0].Tag)
		assert.Contains(t, errors[0].FailedField, "Email")
	})

	t.Run("invalid email format", func(t *testing.T) {
		req := LoginRequest{Email: "not-email", Password: "password123"}
		errors := ValidateStruct(req)
		require.Len(t, errors, 1)
		assert.Equal(t, "email", errors[0].Tag)
	})

	t.Run("password too short", func(t *testing.T) {
		req := LoginRequest{Email: "test@test.com", Password: "12345"}
		errors := ValidateStruct(req)
		require.Len(t, errors, 1)
		assert.Equal(t, "min", errors[0].Tag)
	})

	t.Run("all fields invalid", func(t *testing.T) {
		req := LoginRequest{Email: "", Password: "12"}
		errors := ValidateStruct(req)
		assert.Len(t, errors, 2) // both email and password fail
	})
}
