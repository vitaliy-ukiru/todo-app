package jwt

import (
	"errors"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"
	CookieAccessToken   = "access_token"
)

var (
	ErrInvalidHeader = errors.New("invalid header")
)

//GetTokenFromHeader returns access token header like "Bearer <Data>".
// Param header is authorization header value.
func GetTokenFromHeader(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", ErrInvalidHeader
	}
	if parts[0] != "Bearer" {
		return "", ErrInvalidHeader
	}
	return parts[1], nil
}
