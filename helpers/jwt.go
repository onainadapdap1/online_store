package helpers

import "github.com/golang-jwt/jwt/v5"

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