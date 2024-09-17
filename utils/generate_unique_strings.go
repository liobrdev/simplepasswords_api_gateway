package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/argon2"
)

var lettersSize = big.NewInt(int64(len(UPPERCASE_LETTERS)))
var digitSize = big.NewInt(int64(len(DIGITS)))
var specialCharsSize = big.NewInt(int64(len(SPECIAL_CHARS)))
var slugAlphabetSize = big.NewInt(int64(len(SLUG_ALPHABET)))
var passwordAlphabetSize = big.NewInt(int64(len(PASSWORD_ALPHABET)))
var otpAlphabetSize = big.NewInt(int64(len(OTP_ALPHABET)))

func init() {
	buffer := make([]byte, 1)

	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		panic(fmt.Sprint("crypto/rand is unavailable:\n", err.Error()))
	}
}

func GeneratePassword(n int) (string, error) {
	byte_password := make([]byte, n)
	firstChars := make([]byte, 4)

	var num *big.Int
	var err error

	if num, err = rand.Int(rand.Reader, lettersSize); err != nil {
		return "", err
	} else {
		firstChars[0] = UPPERCASE_LETTERS[num.Int64()]
	}

	if num, err = rand.Int(rand.Reader, lettersSize); err != nil {
		return "", err
	} else {
		firstChars[1] = LOWERCASE_LETTERS[num.Int64()]
	}

	if num, err = rand.Int(rand.Reader, digitSize); err != nil {
		return "", err
	} else {
		firstChars[2] = DIGITS[num.Int64()]
	}

	if num, err = rand.Int(rand.Reader, specialCharsSize); err != nil {
		return "", err
	} else {
		firstChars[3] = SPECIAL_CHARS[num.Int64()]
	}

	for i := 0; i < 4; {
		if num, err = rand.Int(rand.Reader, big.NewInt(int64(n))); err != nil {
			return "", err
		} else if index := num.Int64(); byte_password[index] == 0 {
			byte_password[index] = firstChars[i]
			i++
		}
	}

	for i := 0; i < n; i++ {
		if byte_password[i] != 0 {
			continue
		} else if num, err = rand.Int(rand.Reader, passwordAlphabetSize); err != nil {
			return "", err
		} else {
			byte_password[i] = PASSWORD_ALPHABET[num.Int64()]
		}
	}

	return string(byte_password), nil
}

func GenerateSlug(n int) (string, error) {
	byte_slug := make([]byte, n)

	for i := 0; i < n; i++ {
		if num, err := rand.Int(rand.Reader, slugAlphabetSize); err != nil {
			return "", err
		} else {
			byte_slug[i] = SLUG_ALPHABET[num.Int64()]
		}
	}

	return string(byte_slug), nil
}

func GenerateOTP() ([]string, error) {
	blocks := make([]string, 5)

	for n := 0; n < 5; n++ {
		byte_block := make([]byte, 4)

		for i := 0; i < 4; i++ {
			if num, err := rand.Int(rand.Reader, otpAlphabetSize); err != nil {
				return nil, err
			} else {
				byte_block[i] = OTP_ALPHABET[num.Int64()]
			}
		}

		blocks[n] = string(byte_block)
	}

	return blocks, nil
}

func GenerateSalt(n int) (salt []byte, err error) {
	salt = make([]byte, n)

	for i := 0; i < n; i++ {
		if num, err := rand.Int(rand.Reader, passwordAlphabetSize); err != nil {
			return nil, err
		} else {
			salt[i] = PASSWORD_ALPHABET[num.Int64()]
		}
	}

	return salt, nil
}

func GenerateUserCredentials(password []byte) (hash, salt []byte, err error) {
	if salt, err = GenerateSalt(16); err != nil {
		return nil, nil, err
	}

	return argon2.IDKey(password, salt, 1, 64 * 1024, 4, 64), salt, nil
}
