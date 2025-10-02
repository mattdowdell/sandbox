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

func Test_Timer_UTCNow(t *testing.T) {
	// arrange
	timer := timex.NewTimer()

	// act
	now := timer.UTCNow()

	// assert
	assert.WithinDuration(t, time.Now(), now, time.Second)

	zone, offset := now.Zone()

	assert.Equal(t, "UTC", zone)
	assert.Equal(t, 0, offset)
}

func Test_Timer_Now(t *testing.T) {
	// arrange
	timer := timex.NewTimer()

	// act
	now := timer.Now()

	// assert
	assert.WithinDuration(t, time.Now(), now, time.Second)
}

func Test_Timer_Since(t *testing.T) {
	// arrange
	timer := timex.NewTimer()

	// act
	got := timer.Since(time.Now().Add(time.Hour * -1))

	// assert
	assertWithinRange(t, time.Hour, got, time.Second)
}

func Test_Timer_Until(t *testing.T) {
	// arrange
	timer := timex.NewTimer()

	// act
	got := timer.Until(time.Now().Add(time.Hour))

	// assert
	assertWithinRange(t, time.Hour, got, time.Second)
}
