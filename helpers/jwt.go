package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("ChallengeTestSynapsisBackEnd01")

func GenerateToken(id uint, email string) (string, error) {
	// menyimpan data user
	claims := jwt.MapClaims {
		"user_id": id,
		"email": email,
	}
	// enkripsi data user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// parsing menjadi string panjang
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}


func VerifyToken(c *gin.Context) (*jwt.Token, error) {
	headerToken := c.GetHeader("Authorization")
	fmt.Println(headerToken)
	if headerToken == "" {
		return nil, errors.New("no token provided")
	}

	bearer := strings.HasPrefix(headerToken, "Bearer ")

	if !bearer {
		return nil, errors.New("invalid token format")
	}

	stringToken := strings.TrimSpace(strings.TrimPrefix(headerToken, "Bearer "))

	// parsing token menjadi struct pointer dari jwt.Token
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token verification error: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return token, nil
}
