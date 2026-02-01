package splitter_test

import (
	"fmt"
	"os"

	"github.com/mattdowdell/sandbox/internal/drivers/config"
	"github.com/mattdowdell/sandbox/internal/drivers/config/splitter"
)

type ExampleConfig struct {
	Comma splitter.Comma `koanf:"foo"`
	Space splitter.Space `koanf:"bar"`
}

func Example() {
	// arrange
	os.Setenv("FOO", "foo,bar,baz")
	os.Setenv("BAR", "foo bar baz")

	loaded, _ := config.Load[ExampleConfig](&config.Options{})

	fmt.Printf("comma: %s -> %#v\n", loaded.Comma.String(), loaded.Comma.Unwrap())
	fmt.Printf("space: %s -> %#v\n", loaded.Space.String(), loaded.Space.Unwrap())

	// Output:
	// comma: foo,bar,baz -> []string{"foo", "bar", "baz"}
	// space: foo bar baz -> []string{"foo", "bar", "baz"}
}
