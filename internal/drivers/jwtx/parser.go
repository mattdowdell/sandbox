package jwtx

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mattdowdell/sandbox/internal/domain/repositories"
	"github.com/mattdowdell/sandbox/internal/drivers/config/split"
)

// ParserConfig contains configuration for creating a Parser.
type ParserConfig struct {
	// The expected values of the `aud` claim to validate when parsing a JWT. Any of the provided
	// values will be accepted. Configuration values should be space delimited strings.
	Audience split.Space `koanf:"audience"`

	// The expected value of the `iss` claim to validate when parsing a JWT.
	Issuer string `koanf:"issuer"`

	// The signing algorithm methods accepted by the parser. Configuration values should be space
	// delimited strings.
	Methods split.Space `koanf:"methods"`
}

// Parser is used to validate and parse JWTs based on given configuration.
type Parser struct {
	parser *jwt.Parser
}

// NewParserFromConfig creates a new Parser using the given configuration.
func NewParserFromConfig(clock repositories.Clock, conf ParserConfig) (*Parser, error) {
	return NewParser(clock, conf.Audience, conf.Issuer, conf.Methods)
}

// NewParser creates a new Parser.
func NewParser(
	clock repositories.Clock,
	audience []string,
	issuer string,
	methods []string,
) (*Parser, error) {
	if len(audience) == 0 {
		return nil, errors.New("audience must not be empty for parser")
	}

	for i, value := range audience {
		if value == "" {
			return nil, fmt.Errorf("audience[%d] must not be empty for parser", i)
		}
	}

	if issuer == "" {
		return nil, errors.New("issuer must not be empty for parser")
	}

	if len(methods) == 0 {
		return nil, errors.New("methods must not be empty for parser")
	}

	for i, value := range methods {
		if value == "" {
			return nil, fmt.Errorf("methods[%d] must not be empty for parser", i)
		}
	}

	return &Parser{
		parser: jwt.NewParser(
			jwt.WithAudience(audience...),
			jwt.WithExpirationRequired(),
			jwt.WithIssuer(issuer),
			jwt.WithStrictDecoding(),
			jwt.WithValidMethods(methods),
			jwt.WithTimeFunc(clock.UTCNow),
		),
	}, nil
}

// Parse validates and parses the given token into claims.
func (p *Parser) Parse(input string) (jwt.Claims, error) {
	token, err := p.parser.Parse(input, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: invalid sub claim: %w", err)
	}

	if subject == "" {
		return nil, errors.New("failed to parse token: missing sub claim")
	}

	return token.Claims, err
}

func keyFunc(token *jwt.Token) (any, error) {
	switch token.Method {
	case jwt.SigningMethodHS256:
		return []byte("secret"), nil

	// case &jwt.SigningMethodEd25519{}:

	default:
		return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
	}
}
