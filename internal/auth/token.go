package auth

import(
	"net/http"
	"errors"
	"strings"
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