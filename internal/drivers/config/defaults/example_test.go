package defaults_test

import (
	"fmt"
	"net/netip"

	"github.com/mattdowdell/sandbox/internal/drivers/config/defaults"
)

type MyStruct struct {
	Int    int    `default:"10"`
	Bool   bool   `default:"true"`
	String string `default:"hello"`
	Child  ChildStruct
}

type ChildStruct struct {
	Float float64    `default:"1.2"`
	Uint  uint       `default:"5"`
	Addr  netip.Addr `default:"127.0.0.1"`
}

func Example() {
	s := &MyStruct{}

	if err := defaults.Set(s); err != nil {
		panic(err)
	}

	fmt.Println(s)
	// Output: &{10 true hello {1.2 5 127.0.0.1}}
}
