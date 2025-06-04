package gvek

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

var dylib_ext_map = map[string]string{
	"darwin": "dylib",
	"linux":  "so",
}

type Stream struct {
	veks        unsafe.Pointer
	len_veks    uint
	len_scalars uint
}

type Apply_Args struct {
	a, b, c Stream
}

type Op uint

const (
	Add Op = iota
	Sub
	Div
	Mul
	Mod
	Min
	Max
	Log
	Pow
	// LShift
	// RShift
	// SLShift
	And
	Or
	Xor
)

var Add_f32 func(args *Apply_Args)

type Number interface {
	uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}
type Floats interface {
	float32 | float64
}
type Signed_Numbers interface {
	int8 | int16 | int32 | int64 | float32 | float64
}

func Apply[T Number, OP Op](args *Apply_Args) {

}

func Init() {
	dylib_ext, ok := dylib_ext_map[runtime.GOOS]
	if !ok {
		panic(fmt.Errorf("os: %s is not supported", runtime.GOOS))
	}
	dylib_name := "libzvek." + dylib_ext
	lib, err := purego.Dlopen(dylib_name, purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&Add_f32, lib, "Add_f32")
}
