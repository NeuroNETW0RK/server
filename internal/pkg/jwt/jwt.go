package auth

import (
	"github.com/dgrijalva/jwt-go"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"time"
)

type JwtClaims struct {
	ID       int64   `json:"id"`
	SystemID int64   `json:"system_id"`
	Account  string  `json:"account"`
	RoleIDs  []int64 `json:"role_ids"`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		//todo: 先写死，之后会存入环境变量或者其他地方
		[]byte("neuronetwork666"),
	}
}

type JWT struct {
	SigningKey []byte
}

func (j *JWT) CreateToken(claim *JwtClaims, expireAt time.Duration) (string, error) {
	claim.ExpiresAt = time.Now().Add(expireAt).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.WithCode(code.ErrTokenInvalid, "token invalid")
	}

	if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	}
	return nil, errors.WithCode(code.ErrTokenInvalid, "token invalid")

}
