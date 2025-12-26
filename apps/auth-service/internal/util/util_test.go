package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	secret := "test-secret"
	userID := "user-123"
	email := "test@example.com"
	duration := time.Hour

	t.Run("Generate and Validate Token", func(t *testing.T) {
		token, err := GenerateToken(userID, email, secret, duration)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateToken(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := ValidateToken("invalid", secret)
		assert.Error(t, err)
	})

	t.Run("Expired Token", func(t *testing.T) {
		token, err := GenerateToken(userID, email, secret, -time.Hour)
		assert.NoError(t, err)

		_, err = ValidateToken(token, secret)
		assert.Error(t, err)
	})
}

func TestPassword(t *testing.T) {
	password := "password123"

	t.Run("Hash and Check Password", func(t *testing.T) {
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		assert.True(t, CheckPasswordHash(password, hash))
		assert.False(t, CheckPasswordHash("wrong", hash))
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("Valid Passwords", func(t *testing.T) {
		validPasswords := []string{
			"Password123",
			"MySecure1Pass",
			"Test123ABC",
			"Abcdefgh1",
		}

		for _, password := range validPasswords {
			assert.True(t, ValidatePassword(password), "Password should be valid: %s", password)
		}
	})

	t.Run("Invalid Passwords", func(t *testing.T) {
		invalidPasswords := []string{
			"short1A",        // Too short
			"password123",    // No uppercase
			"PASSWORD123",    // No lowercase
			"PasswordABC",    // No number
			"Pass1",          // Too short
			"",               // Empty
			"12345678",       // Only numbers
			"abcdefgh",       // Only lowercase
			"ABCDEFGH",       // Only uppercase
		}

		for _, password := range invalidPasswords {
			assert.False(t, ValidatePassword(password), "Password should be invalid: %s", password)
		}
	})
}
