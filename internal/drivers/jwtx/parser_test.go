package jwtx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
)

const (
	testMethod = "HS256"
)

func Test_NewParserFromConfig(t *testing.T) {
	// arrange
	conf := jwtx.ParserConfig{
		Audience: []string{testAudience},
		Issuer:   testIssuer,
		Methods:  []string{testMethod},
	}

	// act
	parser := jwtx.NewParserFromConfig(conf)

	// assert
	assert.NotNil(t, parser)
}

func Test_NewParser(t *testing.T) {
	// arrange

	// act
	parser := jwtx.NewParser([]string{testAudience}, testIssuer, []string{testMethod})

	// assert
	assert.NotNil(t, parser)
}

func Test_Parser_Parse_Success(t *testing.T) {
	// arrange
	parser := jwtx.NewParser([]string{testAudience}, testIssuer, []string{testMethod})

	// act
	claims, err := parser.Parse(testToken)

	// assert
	if assert.NotNil(t, claims) {
		subject, err := claims.GetSubject()

		assert.Equal(t, testSubject, subject)
		assert.NoError(t, err)
	}

	assert.NoError(t, err)
}

func Test_Parser_Parse_Error(t *testing.T) {
	testCases := []struct {
		name     string
		audience []string
		issuer   string
		input    string
		want     string
	}{
		{
			name:     "invalid token",
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    "invalid",
			want:     "failed to parse token: token is malformed: token contains an invalid number of segments",
		},
		{
			name:     "wrong audience",
			audience: []string{"other"},
			issuer:   testIssuer,
			input:    testToken,
			want:     "failed to parse token: token has invalid claims: token has invalid audience",
		},
		// TODO: missing audience
		{
			name:     "wrong issuer",
			audience: []string{testAudience},
			issuer:   "other",
			input:    testToken,
			want:     "failed to parse token: token has invalid claims: token has invalid issuer",
		},
		// TODO: missing issuer
		// TODO: missing subject
		// TODO: missing expires at
		// TODO: expired
		// TODO: unexpected signing method
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			parser := jwtx.NewParser(tc.audience, tc.issuer, []string{testMethod})

			// act
			claims, err := parser.Parse(tc.input)

			// assert
			assert.Nil(t, claims)
			assert.EqualError(t, err, tc.want)
		})
	}
}
