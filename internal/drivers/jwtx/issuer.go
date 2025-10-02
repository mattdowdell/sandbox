package jwtx

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// IssuerConfig contains configuration for creating an Issuer.
type IssuerConfig struct {
	// The value to set the `iss` claim to for new JWTs.
	Issuer string `koanf:"issuer"`

	// The value to set the `aud` claim to for new JWTs.
	Audience string `koanf:"audience"`
}

// Issuer is used to create new JWTs based on given configuration.
type Issuer struct {
	timer     repositories.Timer
	generator repositories.UUIDGenerator
	issuer    string
	audience  string
}

// NewIssuerFromConfig creates a new Issuer using the given configuration.
func NewIssuerFromConfig(
	timer repositories.Timer,
	generator repositories.UUIDGenerator,
	conf IssuerConfig,
) (*Issuer, error) {
	return NewIssuer(
		timer,
		generator,
		conf.Issuer,
		conf.Audience,
	)
}

// NewIssuer creates a new Issuer.
func NewIssuer(
	timer repositories.Timer,
	generator repositories.UUIDGenerator,
	issuer string,
	audience string,
) (*Issuer, error) {
	if issuer == "" {
		return nil, errors.New("issuer must not be empty for issuer")
	}

	if audience == "" {
		return nil, errors.New("audience must not be empty for issuer")
	}

	return &Issuer{
		timer:     timer,
		generator: generator,
		issuer:    issuer,
		audience:  audience,
	}, nil
}

// Issue creates a new JWT with the given subject as the subject claim.
func (i *Issuer) Issue(subject string, expiresInSeconds uint32) (string, error) {
	now := i.timer.UTCNow()
	expiresIn := time.Second * time.Duration(expiresInSeconds)

	id, err := i.generator.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate token id: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    i.issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{i.audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id.String(),
	})

	// TODO: use private key
	// TODO: make this configurable
	return token.SignedString([]byte("secret"))
}
