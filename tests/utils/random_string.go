package utils

import (
	"crypto/rand"
	"math/big"
)

// ...
func PrintableASCII() string {
	var output []rune

	for c := ' '; c <= '~'; c++ {
		output = append(output, c)
	}

	return string(output)
}

// ...
func NonASCII() string {
	return "TODO"
}

// ...
func RandomString(chars string, length int) (string, error) {
	output := make([]byte, 0, length)
	l := big.NewInt(int64(len(chars)))

	for _ = range length {
		i, err := rand.Int(rand.Reader, l)
		if err != nil {
			return "", err
		}

		output = append(output, chars[i.Int64()])
	}

	return string(output), nil
}
