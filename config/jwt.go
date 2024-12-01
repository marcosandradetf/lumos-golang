package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTKey []byte

func LoadJWTKey() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Erro ao carregar .env, usando variável de ambiente direta")
	}

	JWTKey = []byte(os.Getenv("JWT_SECRET"))
	if len(JWTKey) == 0 {
		log.Fatal("JWT_SECRET não definido ou vazio")
	}
}
