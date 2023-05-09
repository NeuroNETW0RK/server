package userutils

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
)

func EncodePassword(passWord string) string {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(passWord, options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}
