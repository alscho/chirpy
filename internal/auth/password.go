package auth

import(
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"unicode/utf8"
	"errors"
)

const maxPasswordLength = 72

func HashPassword(password string) (string, error) {
	
	if len(password) == 0 {
        return "", errors.New("password cannot be empty")
    }
	if len(password) > maxPasswordLength {
        return "", errors.New("password exceeds maximum allowed length")
    }
	if !utf8.ValidString(password) {
        return "", errors.New("password contains invalid UTF-8 characters")
    }

	// cost = 0 < MinCost, so GenerateFromPassword sets it to the default value
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", fmt.Errorf("Hashing didn't work: %v", err)
	}
	if bcrypt.CompareHashAndPassword(hashed_password, []byte(password)) != nil {
		return "", fmt.Errorf("Hashed password didn't match hashed_password: %v", err)
	}

	return string(hashed_password), nil
}

func CheckPasswordHash(password, hash string) error {
    if len(password) == 0 {
        return errors.New("password cannot be empty")
    }
    if len(hash) == 0 {
        return errors.New("hash cannot be empty")
    }
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
