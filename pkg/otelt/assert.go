package otelt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// AssertSpanEqual tests whether the given span has the expected name, attributes and events.
func AssertSpanEqual(
	tb testing.TB,
	span trace.ReadOnlySpan,
	name string,
	attributes []attribute.KeyValue,
	events []trace.Event,
) bool {
	tb.Helper()

	var failed bool

	if !assert.Equal(tb, name, span.Name()) {
		failed = true
	}

	if !assert.ElementsMatch(tb, attributes, span.Attributes(), "span attributes do not match") {
		failed = true
	}

	if !AssertEvents(tb, events, span.Events()) {
		failed = true
	}

	return !failed
}

// AssertEvents tests whether a span event has the expected name and attributes.
func AssertEvents(tb testing.TB, expected, actual []trace.Event) bool {
	tb.Helper()

	if !assert.Len(tb, actual, len(expected)) {
		return false
	}

	var failed bool

	for i := range actual {
		if !assert.Equal(
			tb,
			expected[i].Name,
			actual[i].Name,
			"event name does not match at index %d",
			i,
		) {
			failed = true
		}

		if !assert.ElementsMatchf(
			tb,
			expected[i].Attributes,
			actual[i].Attributes,
			"event attributes do not match at index %d: %s",
			i,
			expected[i].Name,
		) {
			failed = true
		}
	}

	return !failed
}
