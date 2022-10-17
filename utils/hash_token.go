package utils

import "crypto/sha512"

func HashToken(token string) []byte {
	checksum := sha512.Sum512([]byte(token))
	digest := make([]byte, sha512.Size)

	for i, n := 0, sha512.Size; i < n; i++ {
		digest[i] = checksum[i]
	}

	return digest
}
