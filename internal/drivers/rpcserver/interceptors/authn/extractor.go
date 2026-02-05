package authn

import (
	"context"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/drivers/jwtx"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Extractor uses claims stored in a context to create log attributes.
type Extractor struct{}

// NewExtractor creates a new Extractor.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts claims from the given context and creates log attributes.
func (e *Extractor) Extract(ctx context.Context) []slog.Attr {
	if subject, ok := jwtx.SubjectFromContext(ctx); ok {
		return []slog.Attr{
			slogx.Subject(subject),
		}
	}

	return nil
}
