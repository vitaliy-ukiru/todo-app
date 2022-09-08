package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	Claims           = jwt.Claims
	RegisteredClaims = jwt.RegisteredClaims
)

type Token = jwt.Token

type SigningMethod = jwt.SigningMethod

var (
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS512 = jwt.SigningMethodHS512
)

var (
	SigningMethodRS256 = jwt.SigningMethodRS256
	SigningMethodRS512 = jwt.SigningMethodRS512
)

func AddDurationToNow(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}
