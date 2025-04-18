package utils

import (
	"crypto/rand"
	"math/big"
)

// PrintableASCII returns a set of printable ASCII characters for use with RandomString.
func PrintableASCII() string {
	var output []rune

	for c := ' '; c <= '~'; c++ {
		output = append(output, c)
	}

	return string(output)
}

// PrintableNonASCII returns a set of printable non-ASCII characters for use with RandomString.
//
// Specifically, it returns the characters code from Latin Extended-A and Latin Extended-B.
func PrintableNonASCII() string {
	var output []rune

	for c := 'Ā'; c <= 'ɏ'; c++ {
		output = append(output, c)
	}

	return string(output)
}

// RandomString generates a string of the given length using the set of given characters.
func RandomString(chars string, length int) (string, error) {
	input := []rune(chars)
	output := make([]rune, 0, length)

	l := big.NewInt(int64(len(input)))

	for range length {
		i, err := rand.Int(rand.Reader, l)
		if err != nil {
			return "", err
		}

		output = append(output, input[i.Int64()])
	}

	return string(output), nil
}
