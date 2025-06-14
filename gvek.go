package gvek

import (
	"fmt"
	"log"
	"runtime"
	"slices"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// not threadsafe
type Ctx unsafe.Pointer

var Make_Ctx func(bytes uint) Ctx

const TARGET_SIMD_BYTES = 64

func bool_to_int(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
func bool_to_uint(b bool) uint {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i uint
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

func ceildiv(x uint, y uint) uint {
	return x/y + bool_to_uint(x%y != 0)
}

func Make_Ctx_Smart(elem_size uint, slice_len uint, num_slices uint) Ctx {
	lanes := TARGET_SIMD_BYTES / elem_size
	slice_padded_len := ceildiv(slice_len, lanes) * lanes
	bytes := num_slices * slice_padded_len * elem_size
	return Make_Ctx(bytes)
}

// reset arena allocator
var Reset_Ctx func(Ctx)

// frees mem
var Free_Ctx func(Ctx)

type Ctx_Tracker struct {
	Ctx       Ctx
	File_name string
	Line_no   int
	Size      uint
}

var active_ctxs []Ctx_Tracker
var active_ctxs_rwlock sync.RWMutex

func bind_ctx_funcs() {
	purego.RegisterLibFunc(&Make_Ctx, lib, "make_ctx")
	purego.RegisterLibFunc(&Reset_Ctx, lib, "reset_ctx")
	purego.RegisterLibFunc(&Free_Ctx, lib, "free_ctx")
}

// only availablle in debug mode
func List_active_ctxs() (num_active int) {
	active_ctxs_rwlock.RLock()
	defer active_ctxs_rwlock.RUnlock()
	num_active = len(active_ctxs)
	if num_active == 0 {
		return
	}
	log.Printf("Active gvek contexts:")
	for _, tracker := range active_ctxs {
		log.Printf("%s:%d size: %d", tracker.File_name, tracker.Line_no, tracker.Size)
	}
	return
}

func bind_debug_ctx_funcs() {
	var make_ctx_raw func(bytes uint) Ctx
	purego.RegisterLibFunc(&make_ctx_raw, lib, "make_ctx")
	Make_Ctx = func(bytes uint) Ctx {
		var tracker Ctx_Tracker
		tracker.Size = bytes
		tracker.Ctx = make_ctx_raw(bytes)
		var ok bool
		_, tracker.File_name, tracker.Line_no, ok = runtime.Caller(1)
		if !ok {
			panic("could not obtain caller information for tracking ctx allocations")
		}
		active_ctxs_rwlock.Lock()
		active_ctxs = append(active_ctxs, tracker)
		active_ctxs_rwlock.Unlock()
		return tracker.Ctx
	}
	purego.RegisterLibFunc(&Reset_Ctx, lib, "reset_ctx")
	var free_ctx_raw func(ctx Ctx)
	purego.RegisterLibFunc(&free_ctx_raw, lib, "free_ctx")
	Free_Ctx = func(c Ctx) {
		free_ctx_raw(c)
		active_ctxs_rwlock.Lock()
		active_ctxs = slices.DeleteFunc(active_ctxs, func(e Ctx_Tracker) bool {
			return e.Ctx == c
		})
		active_ctxs_rwlock.Unlock()
	}
}

// starter kit, if you need more then these, register them using Register_apply_func or Register_apply_single_func
var New_f32 func(ctx Ctx, len uint) []float32

var Add_f32 func(Apply_Args[float32])
var Sub_f32 func(Apply_Args[float32])
var Mul_f32 func(Apply_Args[float32])
var Div_f32 func(Apply_Args[float32])
var Pow_f32 func(Apply_Args[float32])

func bind_f32_funcs() {
	New_f32 = register_new_stream_func[float32](f32)

	Add_f32 = Register_apply_func[float32](f32, Add)
	Sub_f32 = Register_apply_func[float32](f32, Sub)
	Mul_f32 = Register_apply_func[float32](f32, Mul)
	Div_f32 = Register_apply_func[float32](f32, Div)
	Pow_f32 = Register_apply_func[float32](f32, Pow)
}

var New_f64 func(ctx Ctx, len uint) []float64

var Add_f64 func(Apply_Args[float64])
var Sub_f64 func(Apply_Args[float64])
var Mul_f64 func(Apply_Args[float64])
var Div_f64 func(Apply_Args[float64])
var Pow_f64 func(Apply_Args[float64])

func bind_f64_funcs() {
	New_f64 = register_new_stream_func[float64](f64)

	Add_f64 = Register_apply_func[float64](f64, Add)
	Sub_f64 = Register_apply_func[float64](f64, Sub)
	Mul_f64 = Register_apply_func[float64](f64, Mul)
	Div_f64 = Register_apply_func[float64](f64, Div)
	Pow_f64 = Register_apply_func[float64](f64, Pow)
}

var New_bytes func(ctx Ctx, len uint) []byte

var And_bytes func(Apply_Args[byte])
var Or_bytes func(Apply_Args[byte])
var Xor_bytes func(Apply_Args[byte])

func bind_byte_funcs() {
	New_bytes = register_new_stream_func[byte](u8)

	And_bytes = Register_apply_func[byte](u8, And)
	Or_bytes = Register_apply_func[byte](u8, Or)
	Xor_bytes = Register_apply_func[byte](u8, Xor)
}

func As_bytes[T Number](nums []T) (byte_slice []byte) {
	ptr_t := unsafe.SliceData(nums)
	ptr_b := (*byte)(unsafe.Pointer(ptr_t))
	var fake_element T
	sizeof_t := unsafe.Sizeof(fake_element)
	byte_slice = unsafe.Slice(ptr_b, len(nums)*int(sizeof_t))
	return
}
func Bytes_as[T Number](bytes []byte) (nums []T) {
	ptr_b := unsafe.SliceData(bytes)
	ptr_t := (*T)(unsafe.Pointer(ptr_b))
	var fake_element T
	sizeof_t := unsafe.Sizeof(fake_element)
	nums = unsafe.Slice(ptr_t, len(bytes)/int(sizeof_t))
	return
}
func register_new_stream_func[T Number](n NumType) (new_stream_func func(ctx Ctx, len uint) []T) {
	name := "New_" + string(n) + "_Stream"
	var new_stream_func_raw func(ctx Ctx, len uint) (stream *T)
	purego.RegisterLibFunc(&new_stream_func_raw, lib, name)
	new_stream_func = func(ctx Ctx, len uint) []T {
		stream_raw := new_stream_func_raw(ctx, len)
		return unsafe.Slice(stream_raw, len)
	}
	return
}

type Apply_Args[T Number] struct {
	A, B, C []T
}
type Apply_Args_C[T Number] struct {
	a, b, c *T
	len     uint
}

type Num_Apply_Args[T Number] struct {
	A    T
	B, C []T
}
type Num_Apply_Args_C[T Number] struct {
	a    T
	b, c *T
	len  uint
}

type Apply_Args_Num[T Number] struct {
	A []T
	B T
	C []T
}
type Apply_Args_Num_C[T Number] struct {
	a   *T
	b   T
	c   *T
	len uint
}

type Apply_Args_Bool[T Number] struct {
	A, B []T
	C    []uint8
}
type Apply_Args_Bool_C[T Number] struct {
	a, b *T
	c    *uint8
	len  uint
}

type Apply_Args_Num_Bool[T Number] struct {
	A []T
	B T
	C []uint8
}
type Apply_Args_Num_Bool_C[T Number] struct {
	a *T
	b T
	c *uint8
}

type Num_Apply_Args_Bool[T Number] struct {
	A T
	B []T
	C []uint8
}
type Num_Apply_Args_Bool_C[T Number] struct {
	a   T
	b   *T
	c   *uint8
	len uint
}

type Select_Apply_Args[T Number] struct {
	A, B []T
	Pred []uint8
	C    []T
}
type Select_Apply_Args_C[T Number] struct {
	a, b *T
	pred *uint8
	c    *T
	len  uint
}

type Num_Select_Apply_Args[T Number] struct {
	A    T
	B    []T
	Pred []uint8
	C    []T
}
type Num_Select_Apply_Args_C[T Number] struct {
	a    T
	b    *T
	pred *uint8
	c    *T
	len  uint
}

type Select_Apply_Args_Num[T Number] struct {
	A    []T
	B    T
	Pred []uint8
	C    []T
}
type Select_Apply_Args_Num_C[T Number] struct {
	a    *T
	b    T
	pred *uint8
	c    *T
	len  uint
}

// supports casts!
type Apply_Args_Single[T Number, O Number] struct {
	X []T
	Y []O
}
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

type Op1 string // for apply single

const (
	Ceil  = "Ceil"
	Floor = "Floor"
	Round = "Round"
	Sqrt  = "Sqrt"
	Not   = "Not"
	Neg   = "Neg"
	Abs   = "Abs"
	Ln    = "Ln"
	Log2  = "Log2"
	Log10 = "Log10"
	Exp   = "Exp"
	Exp2  = "Exp2"
	Cast  = "Cast"
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

type Operand_Variant uint

const (
	Vector Operand_Variant = iota
	Scalar
)

type Apply_Args_Type[T Number] interface {
	Apply_Args[T] |
		Num_Apply_Args[T] |
		Apply_Args_Num[T] |
		Apply_Args_Bool[T] |
		Num_Apply_Args_Bool[T] |
		Apply_Args_Num_Bool[T] |
		Select_Apply_Args[T] |
		Select_Apply_Args_Num[T] |
		Num_Select_Apply_Args[T]
}

type Number interface {
	uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | float32 | float64
}

var lib uintptr

func Register_apply_func[T Number](t NumType, op Op) (apply_func func(Apply_Args[T])) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Apply_Args[T]) {
		if len(args.A) != len(args.B) || len(args.B) != len(args.C) {
			panic("apply func called with slices of differing length")
		}
		c_func(&Apply_Args_C[T]{
			a:   unsafe.SliceData(args.A),
			b:   unsafe.SliceData(args.B),
			c:   unsafe.SliceData(args.C),
			len: uint(len(args.C)),
		})
	}
	return
}
func Register_num_apply_func[T Number](t NumType, op Op) (apply_func func(Num_Apply_Args[T])) {
	name := "Num_" + string(op) + "_" + string(t)
	var c_func func(*Num_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Num_Apply_Args[T]) {
		if len(args.B) != len(args.C) {
			panic("Num_Apply_Args: slice lengths differ")
		}
		c_func(&Num_Apply_Args_C[T]{
			a:   args.A,
			b:   unsafe.SliceData(args.B),
			c:   unsafe.SliceData(args.C),
			len: uint(len(args.C)),
		})
	}
	return
}

func Register_apply_num_func[T Number](t NumType, op Op) (apply_func func(Apply_Args_Num[T])) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_Num_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Apply_Args_Num[T]) {
		if len(args.A) != len(args.C) {
			panic("Apply_Args_Num: slice lengths differ")
		}
		c_func(&Apply_Args_Num_C[T]{
			a:   unsafe.SliceData(args.A),
			b:   args.B,
			c:   unsafe.SliceData(args.C),
			len: uint(len(args.C)),
		})
	}
	return
}

func Register_apply_bool_func[T Number](t NumType, op Op) (apply_func func(Apply_Args_Bool[T])) {
	name := string(op) + "_" + string(t)
	var c_func func(*Apply_Args_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Apply_Args_Bool[T]) {
		if len(args.A) != len(args.B) || len(args.B) != len(args.C) {
			panic("Apply_Args_Bool: slice lengths differ")
		}
		c_func(&Apply_Args_Bool_C[T]{
			a:   unsafe.SliceData(args.A),
			b:   unsafe.SliceData(args.B),
			c:   unsafe.SliceData(args.C),
			len: uint(len(args.C)),
		})
	}
	return
}

func Register_apply_num_bool_func[T Number](t NumType, op Op) (apply_func func(Apply_Args_Num_Bool[T])) {
	name := string(op) + "_" + string(t) + "_Num"
	var c_func func(*Apply_Args_Num_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Apply_Args_Num_Bool[T]) {
		if len(args.A) != len(args.C) {
			panic("Apply_Args_Num_Bool: slice lengths differ")
		}
		c_func(&Apply_Args_Num_Bool_C[T]{
			a: unsafe.SliceData(args.A),
			b: args.B,
			c: unsafe.SliceData(args.C),
		})
	}
	return
}

func Register_num_apply_bool_func[T Number](t NumType, op Op) (apply_func func(Num_Apply_Args_Bool[T])) {
	name := "Num_" + string(op) + "_" + string(t)
	var c_func func(*Num_Apply_Args_Bool_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Num_Apply_Args_Bool[T]) {
		if len(args.B) != len(args.C) {
			panic("Num_Apply_Args_Bool: slice lengths differ")
		}
		c_func(&Num_Apply_Args_Bool_C[T]{
			a:   args.A,
			b:   unsafe.SliceData(args.B),
			c:   unsafe.SliceData(args.C),
			len: uint(len(args.C)),
		})
	}
	return
}

func Register_select_apply_func[T Number](t NumType) (apply_func func(Select_Apply_Args[T])) {
	name := "Select_" + string(t)
	var c_func func(*Select_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Select_Apply_Args[T]) {
		if len(args.A) != len(args.B) || len(args.B) != len(args.Pred) || len(args.Pred) != len(args.C) {
			panic("Select_Apply_Args: slice lengths differ")
		}
		c_func(&Select_Apply_Args_C[T]{
			a:    unsafe.SliceData(args.A),
			b:    unsafe.SliceData(args.B),
			pred: unsafe.SliceData(args.Pred),
			c:    unsafe.SliceData(args.C),
			len:  uint(len(args.C)),
		})
	}
	return
}

func Register_num_select_apply_func[T Number](t NumType) (apply_func func(Num_Select_Apply_Args[T])) {
	name := "Num_Select_" + string(t)
	var c_func func(*Num_Select_Apply_Args_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Num_Select_Apply_Args[T]) {
		if len(args.B) != len(args.Pred) || len(args.Pred) != len(args.C) {
			panic("Num_Select_Apply_Args: slice lengths differ")
		}
		c_func(&Num_Select_Apply_Args_C[T]{
			a:    args.A,
			b:    unsafe.SliceData(args.B),
			pred: unsafe.SliceData(args.Pred),
			c:    unsafe.SliceData(args.C),
			len:  uint(len(args.C)),
		})
	}
	return
}

func Register_select_apply_num_func[T Number](t NumType) (apply_func func(Select_Apply_Args_Num[T])) {
	name := "Select_" + string(t) + "_Num"
	var c_func func(*Select_Apply_Args_Num_C[T])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Select_Apply_Args_Num[T]) {
		if len(args.A) != len(args.Pred) || len(args.Pred) != len(args.C) {
			panic("Select_Apply_Args_Num: slice lengths differ")
		}
		c_func(&Select_Apply_Args_Num_C[T]{
			a:    unsafe.SliceData(args.A),
			b:    args.B,
			pred: unsafe.SliceData(args.Pred),
			c:    unsafe.SliceData(args.C),
			len:  uint(len(args.C)),
		})
	}
	return
}
func Register_apply_single_func[T Number, O Number](t NumType, o NumType, op Op1) (apply_func func(Apply_Args_Single[T, O])) {
	name := string(op) + "_" + string(t)
	if t != o {
		name = name + "_" + string(o)
	}
	var c_func func(*Apply_Args_Single_C[T, O])
	purego.RegisterLibFunc(&c_func, lib, name)
	apply_func = func(args Apply_Args_Single[T, O]) {
		if len(args.X) != len(args.Y) {
			panic("Apply_Args_Single: slice lengths differ")
		}
		c_func(&Apply_Args_Single_C[T, O]{
			x:   unsafe.SliceData(args.X),
			y:   unsafe.SliceData(args.Y),
			len: uint(len(args.X)),
		})
	}
	return
}

var dylib_ext_map = map[string]string{
	"darwin": "dylib",
	"linux":  "so",
}

func Init(debug bool) {
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
	var init_zvek func()
	purego.RegisterLibFunc(&init_zvek, lib, "init")
	init_zvek()
	if debug {
		bind_debug_ctx_funcs()
	} else {
		bind_ctx_funcs()
	}
	bind_f32_funcs()
	bind_f64_funcs()
	bind_byte_funcs()
}
