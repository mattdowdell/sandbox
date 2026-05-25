package defaults_test

import (
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config/defaults"
)

type (
	MyInt     int
	MyInt8    int8
	MyInt16   int16
	MyInt32   int32
	MyInt64   int64
	MyUint    uint
	MyUint8   uint8
	MyUint16  uint16
	MyUint32  uint32
	MyUint64  uint64
	MyUintptr uintptr
	MyFloat32 float32
	MyFloat64 float64
	MyBool    bool
	MyString  string
)

type Primitives struct {
	Int       int     `default:"1"`
	Int8      int8    `default:"8"`
	Int16     int16   `default:"16"`
	Int32     int32   `default:"32"`
	Int64     int64   `default:"64"`
	Uint      uint    `default:"1"`
	Uint8     uint8   `default:"8"`
	Uint16    uint16  `default:"16"`
	Uint32    uint32  `default:"32"`
	Uint64    uint64  `default:"64"`
	Uintptr   uintptr `default:"1"`
	Float32   float32 `default:"1.32"`
	Float64   float64 `default:"1.64"`
	BoolTrue  bool    `default:"true"`
	BoolFalse bool    `default:"false"`
	String    string  `default:"hello"`

	IntOct    int    `default:"0o1"`
	Int8Oct   int8   `default:"0o10"`
	Int16Oct  int16  `default:"0o20"`
	Int32Oct  int32  `default:"0o40"`
	Int64Oct  int64  `default:"0o100"`
	UintOct   uint   `default:"0o1"`
	Uint8Oct  uint8  `default:"0o10"`
	Uint16Oct uint16 `default:"0o20"`
	Uint32Oct uint32 `default:"0o40"`
	Uint64Oct uint64 `default:"0o100"`

	IntHex    int    `default:"0x1"`
	Int8Hex   int8   `default:"0x8"`
	Int16Hex  int16  `default:"0x10"`
	Int32Hex  int32  `default:"0x20"`
	Int64Hex  int64  `default:"0x40"`
	UintHex   uint   `default:"0x1"`
	Uint8Hex  uint8  `default:"0x8"`
	Uint16Hex uint16 `default:"0x10"`
	Uint32Hex uint32 `default:"0x20"`
	Uint64Hex uint64 `default:"0x40"`

	IntBin    int    `default:"0b1"`
	Int8Bin   int8   `default:"0b1000"`
	Int16Bin  int16  `default:"0b10000"`
	Int32Bin  int32  `default:"0b100000"`
	Int64Bin  int64  `default:"0b1000000"`
	UintBin   uint   `default:"0b1"`
	Uint8Bin  uint8  `default:"0b1000"`
	Uint16Bin uint16 `default:"0b10000"`
	Uint32Bin uint32 `default:"0b100000"`
	Uint64Bin uint64 `default:"0b1000000"`
}

type PrimitivePointers struct {
	IntPtr     *int     `default:"1"`
	UintPtr    *uint    `default:"1"`
	Float32Ptr *float32 `default:"1"`
	BoolPtr    *bool    `default:"true"`
	StringPtr  *string  `default:"hello"`
}

type Aliases struct {
	MyInt       MyInt         `default:"1"`
	MyInt8      MyInt8        `default:"8"`
	MyInt16     MyInt16       `default:"16"`
	MyInt32     MyInt32       `default:"32"`
	MyInt64     MyInt64       `default:"64"`
	MyUint      MyUint        `default:"1"`
	MyUint8     MyUint8       `default:"8"`
	MyUint16    MyUint16      `default:"16"`
	MyUint32    MyUint32      `default:"32"`
	MyUint64    MyUint64      `default:"64"`
	MyUintptr   MyUintptr     `default:"1"`
	MyFloat32   MyFloat32     `default:"1.32"`
	MyFloat64   MyFloat64     `default:"1.64"`
	MyBoolTrue  MyBool        `default:"true"`
	MyBoolFalse MyBool        `default:"false"`
	MyString    MyString      `default:"hello"`
	Duration    time.Duration `default:"10s"`
}

type Unmarshalable struct {
	Struct    netip.Addr  `default:"10.0.0.1"`
	StructPtr *netip.Addr `default:"10.0.0.2"`
}

type Nested struct {
	Child
	Other Child
}

type Child struct {
	Int int `default:"2"`
}

type Untagged struct {
	Int int
}

//nolint:unused // field is not meant to be populated
type Private struct {
	private int
}

type UnparseableInt struct {
	Int int `default:"a"`
}

type UnparseableUint struct {
	Uint uint `default:"a"`
}

type UnparseableFloat struct {
	Float64 float64 `default:"a"`
}

type UnparseableBool struct {
	Bool bool `default:"a"`
}

type UnparseableDuration struct {
	Duration time.Duration `default:"a"`
}

type UnmarshalableError struct {
	Struct netip.Addr `default:"invalid"`
}

type Unsupported struct {
	Slice []string `default:"a,b,c"`
}

func Test_Set_Success(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		have any
		want any
	}{
		"primitives": {
			have: &Primitives{},
			want: &Primitives{
				Int:       1,
				Int8:      8,
				Int16:     16,
				Int32:     32,
				Int64:     64,
				Uint:      1,
				Uint8:     8,
				Uint16:    16,
				Uint32:    32,
				Uint64:    64,
				Uintptr:   0x1,
				Float32:   1.32,
				Float64:   1.64,
				BoolTrue:  true,
				BoolFalse: false,
				String:    "hello",
				IntOct:    1,
				Int8Oct:   8,
				Int16Oct:  16,
				Int32Oct:  32,
				Int64Oct:  64,
				UintOct:   1,
				Uint8Oct:  8,
				Uint16Oct: 16,
				Uint32Oct: 32,
				Uint64Oct: 64,
				IntHex:    1,
				Int8Hex:   8,
				Int16Hex:  16,
				Int32Hex:  32,
				Int64Hex:  64,
				UintHex:   1,
				Uint8Hex:  8,
				Uint16Hex: 16,
				Uint32Hex: 32,
				Uint64Hex: 64,
				IntBin:    1,
				Int8Bin:   8,
				Int16Bin:  16,
				Int32Bin:  32,
				Int64Bin:  64,
				UintBin:   1,
				Uint8Bin:  8,
				Uint16Bin: 16,
				Uint32Bin: 32,
				Uint64Bin: 64,
			},
		},
		"primitive pointer": {
			have: &PrimitivePointers{},
			want: &PrimitivePointers{
				IntPtr:     new(1),
				UintPtr:    new(uint(1)),
				Float32Ptr: new(float32(1.0)),
				BoolPtr:    new(true),
				StringPtr:  new("hello"),
			},
		},
		"aliases": {
			have: &Aliases{},
			want: &Aliases{
				MyInt:       1,
				MyInt8:      8,
				MyInt16:     16,
				MyInt32:     32,
				MyInt64:     64,
				MyUint:      1,
				MyUint8:     8,
				MyUint16:    16,
				MyUint32:    32,
				MyUint64:    64,
				MyUintptr:   1,
				MyFloat32:   1.32,
				MyFloat64:   1.64,
				MyBoolTrue:  true,
				MyBoolFalse: false,
				MyString:    "hello",
				Duration:    time.Second * 10,
			},
		},
		"unmarshalable": {
			have: &Unmarshalable{},
			want: &Unmarshalable{
				Struct:    netip.MustParseAddr("10.0.0.1"),
				StructPtr: new(netip.MustParseAddr("10.0.0.2")),
			},
		},
		"nested": {
			have: &Nested{},
			want: &Nested{
				Child: Child{Int: 2},
				Other: Child{Int: 2},
			},
		},
		"untagged": {
			have: &Untagged{},
			want: &Untagged{},
		},
		"private": {
			have: &Private{},
			want: &Private{},
		},
		"non-zero": {
			have: &Child{Int: 1},
			want: &Child{Int: 1},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// arrange

			// act
			err := defaults.Set(tt.have)

			// assert
			assert.Equal(t, tt.want, tt.have)
			assert.NoError(t, err)
		})
	}
}

func Test_Set_Error(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		have any
		want string
	}{
		"non-pointer struct": {
			have: Sample{},
			want: "expected struct pointer, found: struct",
		},
		"non-struct pointer": {
			have: new(int),
			want: "expected struct pointer, found pointer to: int",
		},
		"non-pointer non-struct": {
			have: 5,
			want: "expected struct pointer, found: int",
		},
		"invalid int": {
			have: &UnparseableInt{},
			want: `set field failed: Int: strconv.ParseInt: parsing "a": invalid syntax`,
		},
		"invalid uint": {
			have: &UnparseableUint{},
			want: `set field failed: Uint: strconv.ParseUint: parsing "a": invalid syntax`,
		},
		"invalid float": {
			have: &UnparseableFloat{},
			want: `set field failed: Float64: strconv.ParseFloat: parsing "a": invalid syntax`,
		},
		"invalid bool": {
			have: &UnparseableBool{},
			want: `set field failed: Bool: strconv.ParseBool: parsing "a": invalid syntax`,
		},
		"invalid duration": {
			have: &UnparseableDuration{},
			want: `set field failed: Duration: time: invalid duration "a"`,
		},
		"invalid unmarshalable": {
			have: &UnmarshalableError{},
			want: `set field failed: Struct: ParseAddr("invalid"): unable to parse IP`,
		},
		"unsupported": {
			have: &Unsupported{},
			want: `set field failed: Slice: unsupported kind: slice`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// arrange

			// act
			err := defaults.Set(tt.have)

			// assert
			assert.EqualError(t, err, tt.want)
		})
	}
}
