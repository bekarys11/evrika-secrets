package auth_test

import (
	"github.com/bekarys11/evrika-secrets/internal/users"
	"github.com/bekarys11/evrika-secrets/pkg/auth"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	user := users.User{
		ID:       11,
		Name:     "testname",
		Email:    "testemail@mail.com",
		Password: "testpassword",
	}
	token, err := auth.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if token == "" {
		t.Error("token should not be empty string")
	}
}
