package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword uses crypto package's GenerateFromPassword to convert plain text password into hash format.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword takes plain text password and corresponding hash as parameters
// and compare both by using crypto package's CompareHashAndPassword funciton.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
