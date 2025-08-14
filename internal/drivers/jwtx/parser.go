package jwtx

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mattdowdell/sandbox/internal/drivers/config/splitter"
)

// ParserConfig contains configuration for creating a Parser.
type ParserConfig struct {
	// The expected values of the `aud` claim to validate when parsing a JWT. Any of the provided
	// values will be accepted. Configuration values should be space delimited strings.
	Audience splitter.Space `koanf:"audience"`

	// The expected value of the `iss` claim to validate when parsing a JWT.
	Issuer string `koanf:"issuer"`

	// The signing algorithm methods accepted by the parser. Configuration values should be space
	// delimited strings.
	Methods splitter.Space `koanf:"methods"`
}

// Parser is used to validate and parse JWTs based on given configuration.
type Parser struct {
	parser *jwt.Parser
}

// NewParserFromConfig creates a new Parser using the given configuration.
func NewParserFromConfig(conf ParserConfig) *Parser {
	return NewParser(conf.Audience, conf.Issuer, conf.Methods)
}

// NewParser creates a new Parser.
func NewParser(audience []string, issuer string, methods []string) *Parser {
	opts := []jwt.ParserOption{
		jwt.WithAudience(audience...),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(issuer),
		jwt.WithStrictDecoding(),
		jwt.WithValidMethods(methods),
	}

	// TODO: append jwt.WithTimeFunc to make unit tests deterministic

	return &Parser{
		parser: jwt.NewParser(opts...),
	}
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
