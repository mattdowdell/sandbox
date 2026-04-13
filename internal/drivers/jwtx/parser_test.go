package jwtx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/domain/repositories/mockrepositories"
	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
)

const (
	testMethod = "HS256"

	testTokenMissingAudience = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIi" +
		"OiJzdWJqZWN0IiwiZXhwIjoxNzU2NDcyNDAwLCJpYXQiOjE3NTY0Njg4MDAsImp0aSI6IjdhZGZhNGZkLWViNWUt" +
		"NGEzMS1iY2U5LTY3Y2UyN2I3NGIwZiJ9.Ss1Fn6ood6FR5cyp9kH8aq_V6L5g7ROIAEIWR39TGy0"
	testTokenMissingIssuer = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzdWJqZWN0IiwiYXVkIj" +
		"pbImF1ZGllbmNlIl0sImV4cCI6MTc1NjQ3MjQwMCwiaWF0IjoxNzU2NDY4ODAwLCJqdGkiOiI3YWRmYTRmZC1lYj" +
		"VlLTRhMzEtYmNlOS02N2NlMjdiNzRiMGYifQ.ymkbpzZIWenvSOUtgK8SoeC95imVv5cVwCg54j70Y3U"
	testTokenMissingSubject = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJhdWQiO" +
		"lsiYXVkaWVuY2UiXSwiZXhwIjoxNzU2NDcyNDAwLCJpYXQiOjE3NTY0Njg4MDAsImp0aSI6IjdhZGZhNGZkLWViN" +
		"WUtNGEzMS1iY2U5LTY3Y2UyN2I3NGIwZiJ9.AC1w2ia17qXPazGbyXRmeISh_CiEM0wyFxZJ8RCQj0Q"
	testTokenInvalidSubject = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIiO" +
		"jEyMzQ1LCJhdWQiOlsiYXVkaWVuY2UiXSwiZXhwIjoxNzU2NDcyNDAwLCJpYXQiOjE3NTY0Njg4MDAsImp0aSI6I" +
		"jdhZGZhNGZkLWViNWUtNGEzMS1iY2U5LTY3Y2UyN2I3NGIwZiJ9.hXYzEXSmXIROYKQB9VevXwnUHsEbVS-qXVu9" +
		"7npGhAQ"
	testTokenMissingExpiresAt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWI" +
		"iOiJzdWJqZWN0IiwiYXVkIjpbImF1ZGllbmNlIl0sImlhdCI6MTc1NjQ2ODgwMCwianRpIjoiN2FkZmE0ZmQtZWI" +
		"1ZS00YTMxLWJjZTktNjdjZTI3Yjc0YjBmIn0.CNfCmx1ew-i_ALuVNfbduyEIsPL8xhRQ-shXkeZ5XCI"
	testTokenExpired = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIiOiJzdWJq" +
		"ZWN0IiwiYXVkIjpbImF1ZGllbmNlIl0sImV4cCI6MTcyNDkzNjQwMCwiaWF0IjoxNzI0OTMyODAwLCJqdGkiOiI3" +
		"YWRmYTRmZC1lYjVlLTRhMzEtYmNlOS02N2NlMjdiNzRiMGYifQ.ISDEZI73SU4r9IgqjMCdiUx3QiS65Dpk9RlUl" +
		"blD9Xg"
	testTokenUnexpectedMethod = "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJz" +
		"dWIiOiJzdWJqZWN0IiwiYXVkIjpbImF1ZGllbmNlIl0sImV4cCI6MTc1NjQ3MjQwMCwiaWF0IjoxNzU2NDY4ODAw" +
		"LCJqdGkiOiI3YWRmYTRmZC1lYjVlLTRhMzEtYmNlOS02N2NlMjdiNzRiMGYifQ.I1CePjpsyfxHw49wZxdSTdmXu" +
		"FVeaj1mu05rAOlC94uKJN_UUJESoVUk1sPg26Wa"
)

func Test_NewParserFromConfig(t *testing.T) {
	// arrange
	clock := mockrepositories.NewMockClock(t)
	conf := jwtx.ParserConfig{
		Audience: []string{testAudience},
		Issuer:   testIssuer,
		Methods:  []string{testMethod},
	}

	// act
	parser, err := jwtx.NewParserFromConfig(clock, conf)

	// assert
	assert.NotNil(t, parser)
	assert.NoError(t, err)
}

func Test_NewParser_Success(t *testing.T) {
	// arrange
	clock := mockrepositories.NewMockClock(t)

	// act
	parser, err := jwtx.NewParser(clock, []string{testAudience}, testIssuer, []string{testMethod})

	// assert
	assert.NotNil(t, parser)
	assert.NoError(t, err)
}

func Test_NewParser_Error(t *testing.T) {
	tests := map[string]struct {
		audience []string
		issuer   string
		methods  []string
		want     string
	}{
		"empty audience": {
			audience: []string{},
			issuer:   testIssuer,
			methods:  []string{testMethod},
			want:     "audience must not be empty for parser",
		},
		"empty audience value": {
			audience: []string{""},
			issuer:   testIssuer,
			methods:  []string{testMethod},
			want:     "audience[0] must not be empty for parser",
		},
		"empty issuer": {
			audience: []string{testAudience},
			issuer:   "",
			methods:  []string{testMethod},
			want:     "issuer must not be empty for parser",
		},
		"empty methods": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			methods:  []string{},
			want:     "methods must not be empty for parser",
		},
		"empty methods value": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			methods:  []string{""},
			want:     "methods[0] must not be empty for parser",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewMockClock(t)

			// act
			parser, err := jwtx.NewParser(clock, tt.audience, tt.issuer, tt.methods)

			// assert
			assert.Nil(t, parser)
			assert.EqualError(t, err, tt.want)
		})
	}
}

func Test_Parser_Parse_Success(t *testing.T) {
	// arrange
	clock := mockrepositories.NewMockClock(t)
	clock.EXPECT().UTCNow().Return(testIssuedAt).Once()

	parser, err := jwtx.NewParser(clock, []string{testAudience}, testIssuer, []string{testMethod})
	require.NoError(t, err)

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
	tests := map[string]struct {
		name     string
		audience []string
		issuer   string
		input    string
		want     string
	}{
		"invalid token": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    "invalid",
			want:     "failed to parse token: token is malformed: token contains an invalid number of segments",
		},
		"wrong audience": {
			audience: []string{"other"},
			issuer:   testIssuer,
			input:    testToken,
			want:     "failed to parse token: token has invalid claims: token has invalid audience",
		},
		"missing audience": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenMissingAudience,
			want: "failed to parse token: token has invalid claims: " +
				"token is missing required claim: aud claim is required",
		},
		"wrong issuer": {
			audience: []string{testAudience},
			issuer:   "other",
			input:    testToken,
			want: "failed to parse token: token has invalid claims: " +
				"token has invalid issuer",
		},
		"missing issuer": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenMissingIssuer,
			want: "failed to parse token: token has invalid claims: " +
				"token is missing required claim: iss claim is required",
		},
		"missing subject": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenMissingSubject,
			want:     "failed to parse token: missing sub claim",
		},
		"invalid subject": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenInvalidSubject,
			want: "failed to parse token: invalid sub claim: invalid type for claim: " +
				"sub is invalid",
		},
		"missing expires at": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenMissingExpiresAt,
			want: "failed to parse token: token has invalid claims: " +
				"token is missing required claim: exp claim is required",
		},
		"expired": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenExpired,
			want:     "failed to parse token: token has invalid claims: token is expired",
		},
		"unexpected signing method": {
			audience: []string{testAudience},
			issuer:   testIssuer,
			input:    testTokenUnexpectedMethod,
			want: "failed to parse token: token signature is invalid: " +
				"signing method HS384 is invalid",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewMockClock(t)
			clock.EXPECT().UTCNow().Return(testIssuedAt).Maybe()

			parser, err := jwtx.NewParser(clock, tt.audience, tt.issuer, []string{testMethod})
			require.NoError(t, err)

			// act
			claims, err := parser.Parse(tt.input)

			// assert
			assert.Nil(t, claims)
			assert.EqualError(t, err, tt.want)
		})
	}
}
