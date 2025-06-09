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

func As_slice[T Number](stream Stream[T]) []T {
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
	c    Stream[uint8]
}
type Apply_Args_Num_Bool[T Number] struct {
	a Stream[T]
	b T
	c Stream[uint8]
}
type Num_Apply_Args_Bool[T Number] struct {
	a T
	b Stream[T]
	c Stream[uint8]
}

// supports casts!
type Apply_Args_Single[T Number, O Number] struct {
	x Stream[T]
	y Stream[O]
}

type Op string

const (
	Add Op = "Add"
	Sub Op = "Sub"
	Div Op = "Div"
	Mul Op = "Mul"
	Mod Op = "Mod"
	Min Op = "Min"
	Max Op = "Max"
	Log Op = "Log"
	Pow Op = "Pow"
	//  LShift
	//  RShift
	//  SLShift
	And Op = "And"
	Or  Op = "Or"
	Xor Op = "Xor"
)

type NumType string

const (
	b8  NumType = "bool"
	u8  NumType = "u8"
	i8  NumType = "i8"
	u16 NumType = "u16"
	i16 NumType = "i16"
	u32 NumType = "u32"
	i32 NumType = "i32"
	u64 NumType = "u64"
	i64 NumType = "i64"
	f32 NumType = "f32"
	f64 NumType = "f64"
)

type Number interface {
	uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

var lib uintptr

func Register_apply_func[T Number](t NumType, op Op) (apply_func func(Apply_Args[T])) {
	name := string(op) + "_" + string(t)
	purego.RegisterLibFunc(&apply_func, lib, name)
	return
}
func Init() {
	dylib_ext, ok := dylib_ext_map[runtime.GOOS]
	if !ok {
		panic(fmt.Errorf("os: %s is not supported", runtime.GOOS))
	}
	dylib_name := "libzvek." + dylib_ext
	var err error
	lib, err = purego.Dlopen(dylib_name, purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
}
