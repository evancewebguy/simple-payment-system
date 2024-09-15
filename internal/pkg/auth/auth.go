package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// GenerateVerificationCode generates a random 5-digit verification code
func GenerateVerificationCode() int {
	seed := rand.NewSource(time.Now().UnixNano()) // Properly seed random generator
	r := rand.New(seed)
	return r.Intn(90000) + 10000 // Generate code between 10000 and 99999
}

// HashPassword hashes the password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plaintext password with its hashed version.
func CheckPasswordHash(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
