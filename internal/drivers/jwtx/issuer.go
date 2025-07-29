package jwtx

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

const (
	claimExpiresAt = "exp"
	claimIssuedAt  = "iat"
	claimIssuer    = "iss"
	claimSubject   = "sub"
	claimAudience  = "aud"
)

// ...
type Issuer struct {
	clock     repositories.Clock
	generator repositories.UUIDGenerator
	expiresIn time.Duration
	issuer    string
	audience  string
}

// ...
func NewIssuer() *Issuer {
	return &Issuer{}
}

// ...
func (i *Issuer) Issue(subject string) (string, error) {
	now := i.clock.UTCNow()

	id, err := i.generator.NewV4()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    i.issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{i.audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(i.expiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id.String(),
	})

	return token.SignedString("secret")
}
