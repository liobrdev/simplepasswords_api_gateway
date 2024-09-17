package utils

import (
	"bytes"

	"golang.org/x/crypto/argon2"
)

func CompareHashAndPassword(hashedPassword, password, salt []byte) bool {
	hash := argon2.IDKey(password, salt, 1, 64 * 1024, 4, 64)

	return bytes.Equal(hashedPassword, hash)
}
