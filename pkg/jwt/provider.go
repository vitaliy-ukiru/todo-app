package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var ErrInvalidToken = errors.New("invalid token")

type Provider struct {
	method     SigningMethod
	privateKey []byte
}

func New(method SigningMethod, privateKey []byte) *Provider {
	return &Provider{method: method, privateKey: privateKey}
}

// VerifyTokenWithClaims parse claims to dest (must be pointer) and check on valid.
func (s Provider) VerifyTokenWithClaims(tokenString string, dest Claims) error {
	token, err := jwt.ParseWithClaims(
		tokenString,
		dest,
		func(_ *jwt.Token) (interface{}, error) {
			return s.privateKey, nil
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if !token.Valid || token.Claims.Valid() != nil {
		return ErrInvalidToken
	}

	return nil
}

func (s Provider) VerifyToken(tokenString string) error {
	// in jwt.Parse() call jwt.ParseWithClaims with MapClaims{} as claims.
	return errors.WithStack(
		s.VerifyTokenWithClaims(
			tokenString,
			jwt.MapClaims{},
		),
	)
}

func (s Provider) CreateToken(claims Claims) (string, error) {
	token, err := jwt.
		NewWithClaims(s.method, claims).
		SignedString(s.privateKey)
	return token, errors.WithStack(err)
}

// UnverifiedClaims only unmarshall json claims to dest.
func UnverifiedClaims(tokenString string, dest Claims) error {
	_, _, err := jwt.NewParser().ParseUnverified(tokenString, dest)
	return errors.WithStack(err)
}
