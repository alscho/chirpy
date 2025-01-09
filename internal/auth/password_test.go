package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "your_password" // example password
	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword() failed with error: %v", err)
	}

	if len(hash) == 0 {
		t.Errorf("expected non-empty hash, got empty string")
	}

	if password == hash {
		t.Errorf("hash should not match the original password")
	}

	// Further, bcrypt ensures that even the same password results in different hashes.
	hash2, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() failed on second hash with error: %v", err)
	}

	if hash == hash2 {
		t.Errorf("expected hashes to differ due to salt, but they are the same")
	}
}

func TestHashPassword_LongPassword(t *testing.T) {
    longPassword := "a_really_long_password_that_is_more_than_72_characters_long..................................."
    _, err := HashPassword(longPassword)

    if err == nil {
        t.Errorf("expected error for password longer than 72 characters, got none")
    }
}

func TestHashPassword_EmptyPassword(t *testing.T) {
    emptyPassword := ""
    _, err := HashPassword(emptyPassword)
    
    if err == nil {
        t.Errorf("expected error for empty password, got none")
    }
}

func TestHashPassword_InvalidUTF8(t *testing.T) {
    invalidPassword := string([]byte{0xff, 0xfe, 0xfd}) // invalid UTF-8 bytes
    _, err := HashPassword(invalidPassword)

    if err == nil {
        t.Errorf("expected error for invalid UTF-8 password, got none")
    } else if err.Error() != "password contains invalid UTF-8 characters" {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestHashPassword_ValidUTF8(t *testing.T) {
    validPassword := "安全!password" // mixed valid UTF-8 characters
    hash, err := HashPassword(validPassword)

    if err != nil {
        t.Errorf("expected valid hash for UTF-8 password, got error: %v", err)
    }
    if len(hash) == 0 {
        t.Errorf("expected non-empty hash for valid UTF-8 password")
    }
}

func TestHashPassword_WhitespaceOnly(t *testing.T) {
    whitespacePassword := "    "
    hash, err := HashPassword(whitespacePassword)

    if err != nil {
        t.Errorf("expected password with only whitespace to hash successfully, but got error: %v", err)
    }
    if len(hash) == 0 {
        t.Errorf("expected non-empty hash for password with only whitespace")
    }
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
    specialCharPassword := "!@#$%^&*()"
    hash, err := HashPassword(specialCharPassword)

    if err != nil {
        t.Errorf("expected special character password to hash successfully, but got error: %v", err)
	}
	if len(hash) == 0 {
        t.Errorf("expected non-empty hash for password with special characters")
    }
}

func TestCheckPasswordHash_Valid(t *testing.T) {
    password := "securepassword"
    hash, _ := HashPassword(password)
    
    err := CheckPasswordHash(password, hash)
    if err != nil {
        t.Errorf("expected passwords to match, but got error: %v", err)
    }
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
    password := "securepassword"
    incorrectPassword := "wrongpassword"
    hash, _ := HashPassword(password)
    
    err := CheckPasswordHash(incorrectPassword, hash)
    if err == nil {
        t.Errorf("expected error for incorrect password, got none")
    }
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
    password := "securepassword"
    invalidHash := "notAValidHash"
    
    err := CheckPasswordHash(password, invalidHash)
    if err == nil {
        t.Errorf("expected error for invalid hash, got none")
	}
}