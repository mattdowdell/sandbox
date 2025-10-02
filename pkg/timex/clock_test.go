package timex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/pkg/timex"
)

func assertWithinRange(t *testing.T, expected, actual, delta time.Duration) {
	t.Helper()

	assert.GreaterOrEqual(t, expected+delta, actual)
	assert.LessOrEqual(t, expected-delta, actual)
}

func Test_Clock_UTCNow(t *testing.T) {
	// arrange
	clock := timex.NewClock()

	// act
	now := clock.UTCNow()

	// assert
	assert.WithinDuration(t, time.Now(), now, time.Second)

	zone, offset := now.Zone()

	assert.Equal(t, "UTC", zone)
	assert.Equal(t, 0, offset)
}

func Test_Clock_Now(t *testing.T) {
	// arrange
	clock := timex.NewClock()

	// act
	now := clock.Now()

	// assert
	assert.WithinDuration(t, time.Now(), now, time.Second)
}

func Test_Clock_Since(t *testing.T) {
	// arrange
	clock := timex.NewClock()

	// act
	got := clock.Since(time.Now().Add(time.Hour * -1))

	// assert
	assertWithinRange(t, time.Hour, got, time.Second)
}

func Test_Clock_Until(t *testing.T) {
	// arrange
	clock := timex.NewClock()

	// act
	got := clock.Until(time.Now().Add(time.Hour))

	// assert
	assertWithinRange(t, time.Hour, got, time.Second)
}
