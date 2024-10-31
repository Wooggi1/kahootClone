package service

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"quiz.com/quiz/internal/entity"
)

var secret_key []byte

func init() {
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading environment variables")
	}

	secret_key = []byte(os.Getenv("JWT_SECRET"))
}

func createToken(user entity.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro carregando variaveis de ambiente")
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodES256,
	jwt.MapClaims{
		"id": user.Id,
		"email": user.Email,
		"name": user.Name,
		"exp": time.Now().Add((time.Hour * 24) * 30).Unix(),
	})

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handleLogin()
