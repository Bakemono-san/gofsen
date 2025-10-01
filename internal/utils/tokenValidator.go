package utils

import (
	"gofsen/internal/types"
	"strings"
)

func ValidateToken(token string) bool {
	// Simple validation: check if the token is "Bearer valid-token"
	return strings.TrimSpace(token) == "Bearer valid-token"
}

type tokenValidator struct{}

func (t *tokenValidator) ValidateToken(token string) bool {
	return ValidateToken(token)
}

func NewTokenValidator() types.TokenValidator {
	return &tokenValidator{}
}
