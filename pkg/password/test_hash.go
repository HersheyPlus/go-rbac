package password

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		shouldSucceed bool
	}{
		{
			name:          "Valid password",
			password:      "TestPass123!",
			shouldSucceed: true,
		},
		{
			name:          "Password too short",
			password:      "Test1!",
			shouldSucceed: false,
		},
		{
			name:          "Password too long",
			password:      strings.Repeat("a", 73) + "A1!",
			shouldSucceed: false,
		},
		{
			name:          "Password missing uppercase",
			password:      "testpass123!",
			shouldSucceed: false,
		},
		{
			name:          "Password missing lowercase",
			password:      "TESTPASS123!",
			shouldSucceed: false,
		},
		{
			name:          "Password missing number",
			password:      "TestPass!!!",
			shouldSucceed: false,
		},
		{
			name:          "Password missing special char",
			password:      "TestPass123",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("HashPassword() error = %v, want nil", err)
				}
				if !IsPasswordHash(hash) {
					t.Error("HashPassword() did not return a valid bcrypt hash")
				}
			} else {
				if err == nil {
					t.Error("HashPassword() error = nil, want error")
				}
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "TestPass123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name          string
		password      string
		hash          string
		shouldSucceed bool
	}{
		{
			name:          "Correct password",
			password:      password,
			hash:          hash,
			shouldSucceed: true,
		},
		{
			name:          "Incorrect password",
			password:      "WrongPass123!",
			hash:          hash,
			shouldSucceed: false,
		},
		{
			name:          "Invalid hash",
			password:      password,
			hash:          "invalid_hash",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hash, tt.password)
			if tt.shouldSucceed && err != nil {
				t.Errorf("ComparePassword() error = %v, want nil", err)
			}
			if !tt.shouldSucceed && err == nil {
				t.Error("ComparePassword() error = nil, want error")
			}
		})
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "Minimum length",
			length: minPasswordLength,
		},
		{
			name:   "Medium length",
			length: 16,
		},
		{
			name:   "Maximum length",
			length: maxPasswordLength,
		},
		{
			name:   "Too short (should use minimum)",
			length: minPasswordLength - 1,
		},
		{
			name:   "Too long (should use maximum)",
			length: maxPasswordLength + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GenerateRandomPassword(tt.length)
			if err != nil {
				t.Errorf("GenerateRandomPassword() error = %v, want nil", err)
			}

			// Validate the generated password
			if err := ValidatePassword(password); err != nil {
				t.Errorf("Generated password failed validation: %v", err)
			}
		})
	}
}