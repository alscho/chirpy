package auth

import(
	"github.com/google/uuid"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"errors"
)

// expiration shouldn't be longer than a month
const expirationDurationOfConcern = time.Hour*24

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	const expiresIn = time.Hour
	
	if tokenSecret == "" {
		return "", errors.New("secret token needed")
	}

	if userID == uuid.Nil {
		return "", errors.New("valid uuid needed")
	}

	if expiresIn <= 0{
		return "", errors.New("positive expiration time needed")
	}

	if expiresIn > expirationDurationOfConcern {
		return "", errors.New("expiration time too long")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error at signing token with tokenSecret: %v", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Double-check that the signing method is what you expect
			return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(tokenSecret), nil
	}

	_, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		var msg string
		if errors.Is(err, jwt.ErrTokenExpired) {
			msg = fmt.Sprintln("Token is expired!")
		} else if errors.Is(err, jwt.ErrSignatureInvalid) {
			msg = fmt.Sprintln("Token signature is invalid!")
		} else {
			msg = fmt.Sprintf("An error occurred: %v\n", err)
		}
		return uuid.Nil, fmt.Errorf(msg+": %v", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Propably invalid user id, problem when converting to uuid type: %v", err)
	}

	if userID == uuid.Nil {
		return uuid.Nil, errors.New("no valid user id - id is nil")
	}

	return userID, nil
}