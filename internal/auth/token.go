package auth

import(
	"net/http"
	"errors"
	"strings"
	"crypto/rand"
	"fmt"
	"encoding/hex"
)

func GetBearerToken(headers http.Header) (string, error) {
	const authorizationPrefix = "Bearer "

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no Authorization header included")
	}

    if !strings.HasPrefix(authHeader, authorizationPrefix) {
        return "", errors.New("invalid authorization header format")
    }

	token := strings.TrimPrefix(authHeader, authorizationPrefix)
	if token == "" {
		return "", errors.New("No token in header")
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	length := 32
	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("Couldn't create random token: %v", err)
	}
	return hex.EncodeToString(tokenBytes), nil
}