package crypto

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(id int) *string {
	secretKey := []byte(os.Getenv("JWT_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing the token.", err)
		return nil
	}

	return &tokenString
}

func ValidateJWT(tokenString string, id int, key string) bool {
	secretKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method not valid")
		}
		return secretKey, nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	// valid token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	switch key {
	case "id":
		return id == int(claims["id"].(float64))
	default:
		return false
	}
}

func GetIdFromJWT(tokenString string) *int {
	secretKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method not valid")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil
	}

	if !token.Valid {
		return nil
	}

	// valid token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	id := int(claims["id"].(float64))

	return &id
}

func GetJWTFromRequest(w http.ResponseWriter, r *http.Request) *string {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "'Authorization' Header missing.", http.StatusUnauthorized)
		return nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "'Authorization' Header Incorrect Format.", http.StatusUnauthorized)
		return nil
	}

	return &parts[1]
}
