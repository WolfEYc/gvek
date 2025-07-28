package gvek

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

var Add_f32 Apply_Args_Fn[float32]
var Sub_f32 Apply_Args_Fn[float32]
var Mul_f32 Apply_Args_Fn[float32]
var Div_f32 Apply_Args_Fn[float32]
var Pow_f32 Apply_Args_Fn[float32]

var AddNum_f32 Apply_Args_Num_Fn[float32]
var SubNum_f32 Apply_Args_Num_Fn[float32]
var MulNum_f32 Apply_Args_Num_Fn[float32]
var DivNum_f32 Apply_Args_Num_Fn[float32]
var PowNum_f32 Apply_Args_Num_Fn[float32]
var MinNum_f32 Apply_Args_Num_Fn[float32]
var MaxNum_f32 Apply_Args_Num_Fn[float32]

var NumAdd_f32 Num_Apply_Args_Fn[float32]
var NumSub_f32 Num_Apply_Args_Fn[float32]
var NumMul_f32 Num_Apply_Args_Fn[float32]
var NumDiv_f32 Num_Apply_Args_Fn[float32]
var NumPow_f32 Num_Apply_Args_Fn[float32]

var CumSum_f32 Apply_Cum_Fn[float32]
var CumProd_f32 Apply_Cum_Fn[float32]

var Ceil_f32 Apply_Args_Single_Fn[float32, float32]

var Neq_f32 Apply_Args_Bool_Fn[float32]
var Eq_f32 Apply_Args_Bool_Fn[float32]
var Lt_f32 Apply_Args_Bool_Fn[float32]
var Gt_f32 Apply_Args_Bool_Fn[float32]

var NeqNum_f32 Apply_Args_Num_Bool_Fn[float32]
var EqNum_f32 Apply_Args_Num_Bool_Fn[float32]
var LtNum_f32 Apply_Args_Num_Bool_Fn[float32]
var GtNum_f32 Apply_Args_Num_Bool_Fn[float32]

var Select_f32 Select_Apply_Args_Fn[float32]
var SelectNum_f32 Select_Apply_Args_Num_Fn[float32]
var NumSelect_f32 Num_Select_Apply_Args_Fn[float32]

var F32_to_f64 Apply_Args_Single_Fn[float32, float64]

var Set_f32 Set_Fn[float32]

func bind_f32_funcs() {
	Add_f32 = Register_apply_func[float32](f32, Add)
	Sub_f32 = Register_apply_func[float32](f32, Sub)
	Mul_f32 = Register_apply_func[float32](f32, Mul)
	Div_f32 = Register_apply_func[float32](f32, Div)
	Pow_f32 = Register_apply_func[float32](f32, Pow)

	AddNum_f32 = Register_apply_num_func[float32](f32, Add)
	SubNum_f32 = Register_apply_num_func[float32](f32, Sub)
	MulNum_f32 = Register_apply_num_func[float32](f32, Mul)
	DivNum_f32 = Register_apply_num_func[float32](f32, Div)
	PowNum_f32 = Register_apply_num_func[float32](f32, Pow)
	MinNum_f32 = Register_apply_num_func[float32](f32, Min)
	MaxNum_f32 = Register_apply_num_func[float32](f32, Max)

	NumAdd_f32 = Register_num_apply_func[float32](f32, Add)
	NumSub_f32 = Register_num_apply_func[float32](f32, Sub)
	NumMul_f32 = Register_num_apply_func[float32](f32, Mul)
	NumDiv_f32 = Register_num_apply_func[float32](f32, Div)
	NumPow_f32 = Register_num_apply_func[float32](f32, Pow)

	CumSum_f32 = Register_apply_cum_func[float32](f32, CumSum)
	CumProd_f32 = Register_apply_cum_func[float32](f32, CumProd)

	Ceil_f32 = Register_apply_single_func[float32, float32](f32, f32, Ceil)

	Neq_f32 = Register_apply_bool_func[float32](f32, Neq)
	Eq_f32 = Register_apply_bool_func[float32](f32, Eq)
	Lt_f32 = Register_apply_bool_func[float32](f32, Lt)
	Gt_f32 = Register_apply_bool_func[float32](f32, Gt)

	NeqNum_f32 = Register_apply_num_bool_func[float32](f32, Neq)
	EqNum_f32 = Register_apply_num_bool_func[float32](f32, Eq)
	LtNum_f32 = Register_apply_num_bool_func[float32](f32, Lt)
	GtNum_f32 = Register_apply_num_bool_func[float32](f32, Gt)

	Select_f32 = Register_select_apply_func[float32](f32)
	SelectNum_f32 = Register_select_apply_num_func[float32](f32)
	NumSelect_f32 = Register_num_select_apply_func[float32](f32)

	F32_to_f64 = Register_apply_single_func[float32, float64](f32, f64, Cast)

	Set_f32 = Register_set_func[float32](f32)
}

var Add_f64 Apply_Args_Fn[float64]
var Sub_f64 Apply_Args_Fn[float64]
var Mul_f64 Apply_Args_Fn[float64]
var Div_f64 Apply_Args_Fn[float64]
var Pow_f64 Apply_Args_Fn[float64]
var Min_f64 Apply_Args_Fn[float64]
var Max_f64 Apply_Args_Fn[float64]

var AddNum_f64 Apply_Args_Num_Fn[float64]
var SubNum_f64 Apply_Args_Num_Fn[float64]
var MulNum_f64 Apply_Args_Num_Fn[float64]
var DivNum_f64 Apply_Args_Num_Fn[float64]
var PowNum_f64 Apply_Args_Num_Fn[float64]
var MinNum_f64 Apply_Args_Num_Fn[float64]
var MaxNum_f64 Apply_Args_Num_Fn[float64]

func bind_f64_funcs() {
	Add_f64 = Register_apply_func[float64](f64, Add)
	Sub_f64 = Register_apply_func[float64](f64, Sub)
	Mul_f64 = Register_apply_func[float64](f64, Mul)
	Div_f64 = Register_apply_func[float64](f64, Div)
	Pow_f64 = Register_apply_func[float64](f64, Pow)
	Min_f64 = Register_apply_func[float64](f64, Min)
	Max_f64 = Register_apply_func[float64](f64, Max)

	AddNum_f64 = Register_apply_num_func[float64](f64, Add)
	SubNum_f64 = Register_apply_num_func[float64](f64, Sub)
	MulNum_f64 = Register_apply_num_func[float64](f64, Mul)
	DivNum_f64 = Register_apply_num_func[float64](f64, Div)
	PowNum_f64 = Register_apply_num_func[float64](f64, Pow)
	MinNum_f64 = Register_apply_num_func[float64](f64, Min)
	MaxNum_f64 = Register_apply_num_func[float64](f64, Max)
}

func As_bytes[T Number](nums []T) (byte_slice []byte) {
	ptr_t := unsafe.SliceData(nums)
	ptr_b := (*byte)(unsafe.Pointer(ptr_t))
	var fake_element T
	sizeof_t := unsafe.Sizeof(fake_element)
	byte_slice = unsafe.Slice(ptr_b, len(nums)*int(sizeof_t))
	return
}
func Bytes_as[T Number](byte_slice []byte) (nums []T) {
	ptr_b := unsafe.SliceData(byte_slice)
	ptr_t := (*T)(unsafe.Pointer(ptr_b))
	var fake_element T
	sizeof_t := unsafe.Sizeof(fake_element)
	nums = unsafe.Slice(ptr_t, len(byte_slice)/int(sizeof_t))
	return
}

var Xor_bytes Apply_Args_Fn[byte]
var And_bytes Apply_Args_Fn[byte]
var Or_bytes Apply_Args_Fn[byte]

func bind_byte_funcs() {
	Xor_bytes = Register_apply_func[byte](u8, Xor)
	And_bytes = Register_apply_func[byte](u8, And)
	Or_bytes = Register_apply_func[byte](u8, Or)
}

type Apply_Args_Fn[T Number] func(c []T, a []T, b []T)
type Num_Apply_Args_Fn[T Number] func(c []T, a T, b []T)
type Apply_Args_Num_Fn[T Number] func(c []T, a []T, b T)
type Apply_Args_Bool_Fn[T Number] func(c []byte, a []T, b []T)
type Apply_Args_Num_Bool_Fn[T Number] func(c []byte, a []T, b T)
type Num_Apply_Args_Bool_Fn[T Number] func(c []byte, a T, b []T)
type Select_Apply_Args_Fn[T Number] func(c []T, a []T, b []T, pred []byte)
type Select_Apply_Args_Num_Fn[T Number] func(c []T, a []T, b T, pred []byte)
type Num_Select_Apply_Args_Fn[T Number] func(c []T, a T, b []T, pred []byte)
type Apply_Args_Single_Fn[T Number, O Number] func(y []O, x []T)
type Apply_Cum_Fn[T Number] func(x []T) T
type Set_Fn[T Number] func(dst []T, src T)

type Set_Args_C[T Number] struct {
	x   *T
	len uint
	y   T
}
type Apply_Cum_Args_C[T Number] struct {
	x   *T
	len uint
}

type Apply_Args_C[T Number] struct {
	a, b, c *T
	len     uint
}

type Num_Apply_Args_C[T Number] struct {
	a    T
	b, c *T
	len  uint
}

type Apply_Args_Num_C[T Number] struct {
	a   *T
	b   T
	c   *T
	len uint
}

type Apply_Args_Bool_C[T Number] struct {
	a, b *T
	c    *byte
	len  uint
}

type Apply_Args_Num_Bool_C[T Number] struct {
	a *T
	b T
	c *byte
}

type Num_Apply_Args_Bool_C[T Number] struct {
	a   T
	b   *T
	c   *byte
	len uint
}

type Select_Apply_Args_C[T Number] struct {
	a, b *T
	pred *byte
	c    *T
	len  uint
}

type Num_Select_Apply_Args_C[T Number] struct {
	a    T
	b    *T
	pred *byte
	c    *T
	len  uint
}

type Select_Apply_Args_Num_C[T Number] struct {
	a    *T
	b    T
	pred *byte
	c    *T
	len  uint
}

// supports casts!
type Apply_Args_Single_C[T Number, O Number] struct {
	x   *T
	y   *O
	len uint
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
	And    Op = "And"
	Or     Op = "Or"
	Xor    Op = "Xor"
	Select Op = "Select"
)

type Cum_Op string

const (
	CumSum  Cum_Op = "CumSum"
	CumProd Cum_Op = "CumProd"
)

type Op1 string // for apply single

const (
	Ceil  Op1 = "Ceil"
	Floor Op1 = "Floor"
	Round Op1 = "Round"
	Sqrt  Op1 = "Sqrt"
	Not   Op1 = "Not"
	Neg   Op1 = "Neg"
	Abs   Op1 = "Abs"
	Ln    Op1 = "Ln"
	Log2  Op1 = "Log2"
	Log10 Op1 = "Log10"
	Exp   Op1 = "Exp"
	Exp2  Op1 = "Exp2"
	Cast  Op1 = "Cast"
)

type Bool_Op string

const (
	Gt  Bool_Op = "Gt"
	Gte Bool_Op = "Gte"
	Lt  Bool_Op = "Lt"
	Lte Bool_Op = "Lte"
	Eq  Bool_Op = "Eq"
	Neq Bool_Op = "Neq"
)

type NumType string

const (
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

var numtype_to_size = map[NumType]uint{
	u8:  1,
	i8:  1,
	u16: 2,
	i16: 2,
	i32: 4,
	u32: 4,
	u64: 8,
	i64: 8,
	f32: 4,
	f64: 8,
}

type Operand_Variant uint

const (
	Vector Operand_Variant = iota
	Scalar
)

type Number interface {
	uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

var lib uintptr

func Register_apply_func[T Number](t NumType, op Op) (apply_func Apply_Args_Fn[T]) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c, a, b []T) {
		if len(a) != len(b) || len(b) != len(c) {
			panic("apply func called with slices of differing length")
		}
		c_func(&Apply_Args_C[T]{
			a:   unsafe.SliceData(a),
			b:   unsafe.SliceData(b),
			c:   unsafe.SliceData(c),
			len: uint(len(c)),
		})
	}
	return
}
func Register_apply_cum_func[T Number](t NumType, op Cum_Op) (apply_func Apply_Cum_Fn[T]) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Cum_Args_C[T]) T
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(x []T) T {
		return c_func(&Apply_Cum_Args_C[T]{
			x:   unsafe.SliceData(x),
			len: uint(len(x)),
		})
	}
	return
}
func Register_set_func[T Number](t NumType) (set_func Set_Fn[T]) {
	name := "Set_" + string(t)
	var c_func func(*Set_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	set_func = func(x []T, y T) {
		c_func(&Set_Args_C[T]{
			x:   unsafe.SliceData(x),
			len: uint(len(x)),
			y:   y,
		})
	}
	return
}
func Register_num_apply_func[T Number](t NumType, op Op) (apply_func Num_Apply_Args_Fn[T]) {
	name := "Num_" + string(op) + "_" + string(t)
	var c_func func(*Num_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c []T, a T, b []T) {
		if len(b) != len(c) {
			panic("Num_Apply_Args: slice lengths differ")
		}
		c_func(&Num_Apply_Args_C[T]{
			a:   a,
			b:   unsafe.SliceData(b),
			c:   unsafe.SliceData(c),
			len: uint(len(c)),
		})
	}
	return
}

func Register_apply_num_func[T Number](t NumType, op Op) (apply_func Apply_Args_Num_Fn[T]) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_Num_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c, a []T, b T) {
		if len(a) != len(c) {
			panic("Apply_Args_Num: slice lengths differ")
		}
		c_func(&Apply_Args_Num_C[T]{
			a:   unsafe.SliceData(a),
			b:   b,
			c:   unsafe.SliceData(c),
			len: uint(len(c)),
		})
	}
	return
}

func Register_apply_bool_func[T Number](t NumType, op Bool_Op) (apply_func Apply_Args_Bool_Fn[T]) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c []byte, a, b []T) {
		if len(a) != len(b) || len(b) <= len(c)*8 {
			panic("Apply_Args_Bool: slice lengths differ")
		}
		c_func(&Apply_Args_Bool_C[T]{
			a:   unsafe.SliceData(a),
			b:   unsafe.SliceData(b),
			c:   unsafe.SliceData(c),
			len: uint(len(c)),
		})
	}
	return
}

func Register_apply_num_bool_func[T Number](t NumType, op Bool_Op) (apply_func Apply_Args_Num_Bool_Fn[T]) {
	name := string(op) + "_" + string(t) + "_Num"
	var c_func func(*Apply_Args_Num_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c []byte, a []T, b T) {
		if len(a) <= len(c)*8 {
			panic("Apply_Args_Num_Bool: slice lengths differ")
		}
		c_func(&Apply_Args_Num_Bool_C[T]{
			a: unsafe.SliceData(a),
			b: b,
			c: unsafe.SliceData(c),
		})
	}
	return
}

func Register_num_apply_bool_func[T Number](t NumType, op Bool_Op) (apply_func Num_Apply_Args_Bool_Fn[T]) {
	name := "Num_" + string(op) + "_" + string(t)
	var c_func func(*Num_Apply_Args_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c []byte, a T, b []T) {
		if len(b) <= len(c)*8 {
			panic("Num_Apply_Args_Bool: slice lengths differ")
		}
		c_func(&Num_Apply_Args_Bool_C[T]{
			a:   a,
			b:   unsafe.SliceData(b),
			c:   unsafe.SliceData(c),
			len: uint(len(c)),
		})
	}
	return
}

func Register_select_apply_func[T Number](t NumType) (apply_func Select_Apply_Args_Fn[T]) {
	name := "Select_" + string(t)
	var c_func func(*Select_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c, a, b []T, pred []byte) {
		if len(a) != len(b) || len(b) != len(c) || len(c) <= len(pred)*8 {
			panic("Select_Apply_Args: slice lengths differ")
		}
		c_func(&Select_Apply_Args_C[T]{
			a:    unsafe.SliceData(a),
			b:    unsafe.SliceData(b),
			pred: unsafe.SliceData(pred),
			c:    unsafe.SliceData(c),
			len:  uint(len(c)),
		})
	}
	return
}

func Register_num_select_apply_func[T Number](t NumType) (apply_func Num_Select_Apply_Args_Fn[T]) {
	name := "Num_Select_" + string(t)
	var c_func func(*Num_Select_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c []T, a T, b []T, pred []byte) {
		if len(b) != len(pred) || len(c) <= len(pred)*8 {
			panic("Num_Select_Apply_Args: slice lengths differ")
		}
		c_func(&Num_Select_Apply_Args_C[T]{
			a:    a,
			b:    unsafe.SliceData(b),
			pred: unsafe.SliceData(pred),
			c:    unsafe.SliceData(c),
			len:  uint(len(c)),
		})
	}
	return
}

func Register_select_apply_num_func[T Number](t NumType) (apply_func Select_Apply_Args_Num_Fn[T]) {
	name := "Select_" + string(t) + "_Num"
	var c_func func(*Select_Apply_Args_Num_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(c, a []T, b T, pred []byte) {
		if len(a) != len(pred) || len(c) <= len(pred)*8 {
			panic("Select_Apply_Args_Num: slice lengths differ")
		}
		c_func(&Select_Apply_Args_Num_C[T]{
			a:    unsafe.SliceData(a),
			b:    b,
			pred: unsafe.SliceData(pred),
			c:    unsafe.SliceData(c),
			len:  uint(len(c)),
		})
	}
	return
}
func Register_apply_single_func[T Number, O Number](t NumType, o NumType, op Op1) (apply_func Apply_Args_Single_Fn[T, O]) {
	name := string(op) + "_" + string(t)
	if t != o {
		name = name + "_" + string(o)
	}
	var c_func func(*Apply_Args_Single_C[T, O])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(y []O, x []T) {
		if len(x) != len(y) {
			panic("Apply_Args_Single: slice lengths differ")
		}
		c_func(&Apply_Args_Single_C[T, O]{
			x:   unsafe.SliceData(x),
			y:   unsafe.SliceData(y),
			len: uint(len(x)),
		})
	}
	return
}

var dylib_ext_map = map[string]string{
	"darwin": "dylib",
	"linux":  "so",
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
	bind_f32_funcs()
	bind_f64_funcs()
	bind_byte_funcs()
}
