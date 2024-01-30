package token

import (
	"auth-application/internal/domain/models"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"
)

type TokenService struct {
	tokenTtl time.Duration
	secret   string
	length   uint8
}

func New(secret string, tokenTtl time.Duration) *TokenService {
	return &TokenService{tokenTtl: tokenTtl, secret: secret, length: 64}
}

func (ts *TokenService) NewAccess(user models.User) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ts.tokenTtl).Unix()
	tokenString, err := token.SignedString([]byte(ts.secret))
	if err != nil {
		return ""
	}

	return tokenString
}

func (ts *TokenService) CreateRefresh(userId int64) string {
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	token := make([]rune, ts.length)
	for i := range token {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		token[i] = chars[randomIndex.Int64()]
	}

	return string(token)
}
