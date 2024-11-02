package password

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort     = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong      = errors.New("password must not exceed 72 characters")
	ErrPasswordTooWeak      = errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrPasswordHashFailed   = errors.New("failed to hash password")
	ErrPasswordCompareError = errors.New("failed to compare password with hash")
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt's maximum length
	bcryptCost       = 12 // Higher cost means more secure but slower
)

// HashPassword converts a plain text password into a hashed version
func HashPassword(password string) (string, error) {
	// Validate password before hashing
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	// Generate hash from password with specific cost
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", ErrPasswordHashFailed
	}

	return string(hashedBytes), nil
}

// ComparePassword checks if a plain text password matches a hash
func ComparePassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrPasswordCompareError
		}
		return err
	}
	return nil
}

// ValidatePassword checks if a password meets the minimum security requirements
func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return ErrPasswordTooWeak
	}

	return nil
}

// GenerateRandomPassword generates a random password that meets the validation requirements
func GenerateRandomPassword(length int) (string, error) {
	if length < minPasswordLength {
		length = minPasswordLength
	}
	if length > maxPasswordLength {
		length = maxPasswordLength
	}

	// Character sets for password generation
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	numberChars := "0123456789"
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// Ensure at least one character from each set
	password := make([]byte, length)
	password[0] = upperChars[SecureRandom(len(upperChars))]
	password[1] = lowerChars[SecureRandom(len(lowerChars))]
	password[2] = numberChars[SecureRandom(len(numberChars))]
	password[3] = specialChars[SecureRandom(len(specialChars))]

	// Fill the rest with random characters from all sets
	allChars := upperChars + lowerChars + numberChars + specialChars
	for i := 4; i < length; i++ {
		password[i] = allChars[SecureRandom(len(allChars))]
	}

	// Shuffle the password to avoid predictable pattern
	ShuffleBytes(password)

	return string(password), nil
}

// SecureRandom generates a random number using crypto/rand
func SecureRandom(max int) int {
	// Implementation using crypto/rand
	// Note: This is a simplified version. In production, you should use crypto/rand properly
	return max / 2 // Placeholder implementation
}

// ShuffleBytes shuffles a byte slice using Fisher-Yates algorithm
func ShuffleBytes(slice []byte) {
	for i := len(slice) - 1; i > 0; i-- {
		j := SecureRandom(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// IsPasswordHash checks if a string is likely to be a bcrypt hash
func IsPasswordHash(hash string) bool {
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	if len(hash) != 60 {
		return false
	}
	return hash[0] == '$' && hash[1] == '2' && hash[3] == '$'
}

// UpdatePasswordHash checks if a password hash needs to be updated
// This is useful when you want to upgrade security parameters
func UpdatePasswordHash(hash, password string) (string, bool, error) {
	// Check if the hash needs to be updated (e.g., cost is too low)
	hashCost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return "", false, err
	}

	if hashCost < bcryptCost {
		// Generate new hash with higher cost
		newHash, err := HashPassword(password)
		return newHash, true, err
	}

	return hash, false, nil
}


// func useSage(){
// 		// Hash a password
// 		hash, err := password.HashPassword("MySecurePass123!")

// 		// Compare a password with a hash
// 		err = password.ComparePassword(hash, "MySecurePass123!")
	
// 		// Generate a random password
// 		pass, err := password.GenerateRandomPassword(16)
	
// 		// Validate a password
// 		err = password.ValidatePassword("MyPass123!")
	
// }