package auth

import(
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
    t.Run("valid token creation and validation", func(t *testing.T) {
        originalID := uuid.New()
        secret := "test-secret"
        duration := time.Hour

        token, err := MakeJWT(originalID, secret, duration)
        if err != nil {
            t.Fatalf("Failed to create token: %v", err)
        }

        returnedID, err := ValidateJWT(token, secret)
        if err != nil {
            t.Fatalf("Failed to validate token: %v", err)
        }

        if returnedID != originalID {
            t.Errorf("User IDs don't match. Got %v, want %v", returnedID, originalID)
        }
    })

    t.Run("expired token", func(t *testing.T) {
        originalID := uuid.New()
        secret := "test-secret"
        duration := -time.Hour // negative duration creates an expired token

        token, _ := MakeJWT(originalID, secret, duration)
        _, err := ValidateJWT(token, secret)
        if err == nil {
            t.Error("Expected error for expired token, got nil")
        }
    })

    t.Run("invalid signature", func(t *testing.T) {
        originalID := uuid.New()
        token, _ := MakeJWT(originalID, "correct-secret", time.Hour)
        
        _, err := ValidateJWT(token, "wrong-secret")
        if err == nil {
            t.Error("Expected error for invalid signature, got nil")
        }
    })

	t.Run("missing required claims", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256) // no claims attached
		tokenString, _ := token.SignedString([]byte("secret"))
		
		_, err := ValidateJWT(tokenString, "secret")
		if err == nil {
			t.Error("Expected error for missing required claims, got nil")
		}
	})

	t.Run("valid claims but expired", func(t *testing.T) {
		originalID := uuid.New()
		secret := "test-secret"
	
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // issued 2 hours ago
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),    // expired 1 hour ago
			Subject:   originalID.String(),
		})
		tokenString, _ := token.SignedString([]byte(secret))
		
		_, err := ValidateJWT(tokenString, secret)
		if err == nil {
			t.Error("Expected error for valid claims but expired token, got nil")
		}
	})

	t.Run("manipulated token", func(t *testing.T) {
		token, _ := MakeJWT(uuid.New(), "secret", time.Hour)
		manipulatedToken := token[:len(token)-1] + "a" // tamper with token
		
		_, err := ValidateJWT(manipulatedToken, "secret")
		if err == nil {
			t.Error("Expected error for manipulated token, got nil")
		}
	})

	t.Run("empty secret", func(t *testing.T) {
        originalID := uuid.New()
        token, _ := MakeJWT(originalID, "", time.Hour)
        
        _, err := ValidateJWT(token, "")
        if err == nil {
            t.Error("Expected error for empty secret, got nil")
        }
    })

	t.Run("invalid uuid in token", func(t *testing.T) {
		// Create a token with invalid UUID string as subject
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject: "not-a-uuid",
			// ... other claims
		})
		tokenString, _ := token.SignedString([]byte("secret"))
		
		_, err := ValidateJWT(tokenString, "secret")
		if err == nil {
			t.Error("Expected error for invalid UUID, got nil")
		}
	})

	t.Run("malformed token string", func(t *testing.T) {
		_, err := ValidateJWT("not.a.jwt", "secret")
		if err == nil {
			t.Error("Expected error for malformed token, got nil")
		}
	})
	
	t.Run("zero duration", func(t *testing.T) {
		originalID := uuid.New()
		_, err := MakeJWT(originalID, "secret", 0)
		if err == nil {
			t.Error("Expected error for not positive duration, got nil")
		}
	})
	
	t.Run("nil uuid", func(t *testing.T) {
		_, err := MakeJWT(uuid.Nil, "secret", time.Hour)
		if err == nil {
			t.Error("Expected error for nil id, got nil")
		}
	})
	
	t.Run("very long duration", func(t *testing.T) {
		originalID := uuid.New()
		_, err := MakeJWT(originalID, "secret", time.Hour * 24 * 365 * 100) // 100 years
		if err == nil {
			t.Error("Expected error for expiration time too large, got nil")
		}
	})
}
