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

type Stream[T Number] struct {
	veks        *T
	len_veks    uint
	len_scalars uint
}

func as_slice[T Number](stream Stream[T]) []T {
	return unsafe.Slice(stream.veks, stream.len_scalars)
}

type Apply_Args[T Number] struct {
	a, b, c Stream[T]
}
type Num_Apply_Args[T Number] struct {
	a    T
	b, c Stream[T]
}
type Apply_Args_Num[T Number] struct {
	a Stream[T]
	b T
	c Stream[T]
}
type Apply_Args_Bool[T Number] struct {
	a, b Stream[T]
	c    Stream[bool]
}
type Apply_Args_Single[T Number] struct {
	x, y Stream[T]
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

type NumType uint

type Number interface {
	bool | uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

// Unsigned integers
var apply_vtable_u8 []func(Apply_Args[uint8])
var apply_vtable_u16 []func(Apply_Args[uint16])
var apply_vtable_u32 []func(Apply_Args[uint32])
var apply_vtable_u64 []func(Apply_Args[uint64])

// Signed integers
var apply_vtable_i8 []func(Apply_Args[int8])
var apply_vtable_i16 []func(Apply_Args[int16])
var apply_vtable_i32 []func(Apply_Args[int32])
var apply_vtable_i64 []func(Apply_Args[int64])

// Floats
var apply_vtable_f32 []func(Apply_Args[float32])
var apply_vtable_f64 []func(Apply_Args[float64])

// Bool
var apply_vtable_bool []func(Apply_Args[bool])

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
