package helpers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HassPass(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password should not be empty")
	}

	bytePassword := []byte(pass)
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	
	if err != nil {
		fmt.Printf("[UserHandler.SetPassword] Error when generate password with error: %v\n", err)
		return "", nil
	}

	return string(hashPassword), nil
}

func ComparePassword(hassPass, inputPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hassPass, inputPassword)
	if err != nil {
		return false
	}
	return true
}

// func ComparePassword(hashPass, inputPasss []byte) bool {
// 	err := bcrypt.CompareHashAndPassword(hashPass, inputPasss)
// 	return err == nil
// }
