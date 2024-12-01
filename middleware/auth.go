package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var publicKey []byte

func init() {
	var err error
	publicKey, err = ioutil.ReadFile("keys/app.pub")
	if err != nil {
		log.Fatal("Erro ao carregar a chave pública:", err)
	}

}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return key, nil
	})

	if err != nil {
		return nil, errors.New("token inválido")
	}

	return token, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token ausente ou inválido"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Valida o token
		_, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Continua se o token for válido
		c.Next()
	}
}
