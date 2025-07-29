package jwtx

import (
	"github.com/golang-jwt/jwt/v5"
)

// ...
//
// TODO: rename to parser config?
type Config struct {
	// The expected values of the `aud` claim to validate when parsing a JWT. Any of the provided
	// values will be accepted.
	Audience []string `koanf:"audience"`

	// The expected value of the `iss` claim to validate when parsing a JWT.
	Issuer string `koanf:"issuer"`

	// ...
	Methods []string `koanf:"methods"`
}

// TODO: move to shared adapter
type Claims interface {
	GetSubject() (string, error)
}

// ...
type Parser struct {
	parser *jwt.Parser
}

// ...
func NewParserFromConfig(conf Config) *Parser {
	return NewParser(conf.Audience, conf.Issuer, conf.Methods)
}

// ...
func NewParser(audience []string, issuer string, methods []string) *Parser {
	parser := jwt.NewParser(
		jwt.WithAudience(audience...),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(issuer),
		jwt.WithStrictDecoding(),
		jwt.WithValidMethods(methods),
	)

	return &Parser{
		parser: parser,
	}
}

// ...
func (p *Parser) Parse(input string) (Claims, error) {
	token, err := p.parser.Parse(input, keyFunc)
	if err != nil {
		return nil, err
	}

	return token.Claims, err
}

func keyFunc(token *jwt.Token) (any, error) {
	return nil, nil
}
