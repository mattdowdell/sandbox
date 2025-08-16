package jwtx_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/mocks/domain/mockrepositories"
)

const (
	testExpiresInSeconds = 3600
	testIssuer           = "issuer"
	testAudience         = "audience"
	testSubject          = "subject"

	testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIiOiJzdWJqZWN0Iiwi" +
		"YXVkIjpbImF1ZGllbmNlIl0sImV4cCI6MTc1NjQ3MjQwMCwiaWF0IjoxNzU2NDY4ODAwLCJqdGkiOiI3YWRmYTRmZ" +
		"C1lYjVlLTRhMzEtYmNlOS02N2NlMjdiNzRiMGYifQ.967KdDGa_Le6rKo7-bdyNlgvXng3-IAMgiLvKQkt45U"
)

var (
	testIssuedAt = time.Date(2025, time.August, 29, 12, 0, 0, 0, time.UTC)
	testID       = uuid.Must(uuid.FromString("7adfa4fd-eb5e-4a31-bce9-67ce27b74b0f"))
)

func Test_NewIssuerFromConfig(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)
	generator := mockrepositories.NewUUIDGenerator(t)

	conf := jwtx.IssuerConfig{
		Issuer:   testIssuer,
		Audience: testAudience,
	}

	// act
	issuer, err := jwtx.NewIssuerFromConfig(clock, generator, conf)

	// assert
	assert.NotNil(t, issuer)
	assert.NoError(t, err)
}

func Test_NewIssuer_Success(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)
	generator := mockrepositories.NewUUIDGenerator(t)

	// act
	issuer, err := jwtx.NewIssuer(clock, generator, testIssuer, testAudience)

	// assert
	assert.NotNil(t, issuer)
	assert.NoError(t, err)
}

func Test_NewIssuer_Error(t *testing.T) {
	testCases := []struct {
		name     string
		issuer   string
		audience string
		want     string
	}{
		{
			name:     "missing issuer",
			issuer:   "",
			audience: testAudience,
			want:     "issuer must not be empty for issuer",
		},
		{
			name:     "missing audience",
			issuer:   testIssuer,
			audience: "",
			want:     "audience must not be empty for issuer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			clock := mockrepositories.NewClock(t)
			generator := mockrepositories.NewUUIDGenerator(t)

			// act
			issuer, err := jwtx.NewIssuer(clock, generator, tc.issuer, tc.audience)

			// assert
			assert.Nil(t, issuer)
			assert.EqualError(t, err, tc.want)
		})
	}
}

func Test_Issuer_Issue_Success(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)
	clock.EXPECT().UTCNow().Return(testIssuedAt).Once()

	generator := mockrepositories.NewUUIDGenerator(t)
	generator.EXPECT().NewV4().Return(testID, nil).Once()

	issuer, err := jwtx.NewIssuer(clock, generator, testIssuer, testAudience)
	require.NoError(t, err)

	// act
	got, err := issuer.Issue(testSubject, testExpiresInSeconds)

	// assert
	assert.Equal(t, testToken, got)
	assert.NoError(t, err)
}

func Test_Issuer_Issue_Error(t *testing.T) {
	// arrange
	clock := mockrepositories.NewClock(t)
	clock.EXPECT().UTCNow().Return(testIssuedAt).Once()

	generator := mockrepositories.NewUUIDGenerator(t)
	generator.EXPECT().NewV4().Return(uuid.Nil, errors.New("example")).Once()

	issuer, err := jwtx.NewIssuer(clock, generator, testIssuer, testAudience)
	require.NoError(t, err)

	// act
	got, err := issuer.Issue(testSubject, testExpiresInSeconds)

	// assert
	assert.Empty(t, got)
	assert.EqualError(t, err, "failed to generate token id: example")
}
