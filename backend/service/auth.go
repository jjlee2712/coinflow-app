package service

import (
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassowrd(password, hash string) error{
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}

func CreateToken(userID int64, secret string) (string, error) {
	// Step 1 — define the claims (payload inside the token)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	// Step 2 — create an unsigned token with those claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Step 3 — sign it with your secret and return the final string
	return token.SignedString([]byte(secret));
}

func ValidateToken(tokenStr, secret string) (int64, error){
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token)(interface {}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid{
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	return userID, err
}