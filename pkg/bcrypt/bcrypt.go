package bcryptx

import (
	unsafe "github.com/go-playground/pkg/v5/unsafe"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(unsafe.StringToBytes(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// here return via builtin conversion for copy
	// i.e. this value travel via network.
	return string(hash), nil
}

func CompareHash(src, input string) error {
	return bcrypt.CompareHashAndPassword(unsafe.StringToBytes(src), unsafe.StringToBytes(input))
}
