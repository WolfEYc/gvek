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

// size in bytes for this ctx

var Make_Ctx func(len uint) Ctx

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
	var make_ctx_raw func(len uint) Ctx
	purego.RegisterLibFunc(&make_ctx_raw, lib, "make_ctx")
	Make_Ctx = func(len uint) Ctx {
		var tracker Ctx_Tracker
		tracker.Size = len
		tracker.Ctx = make_ctx_raw(len)
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
var New_f32 New_Stream_Func[float32]
var Set_f32 Set_Stream_Func[float32]
var To_f32 To_Stream_Func[float32]

var Add_f32 func(*Apply_Args[float32])
var Sub_f32 func(*Apply_Args[float32])
var Mul_f32 func(*Apply_Args[float32])
var Div_f32 func(*Apply_Args[float32])
var Pow_f32 func(*Apply_Args[float32])

func bind_f32_funcs() {
	New_f32 = register_new_stream_func[float32](f32)
	Set_f32 = register_set_stream_func[float32](f32)
	To_f32 = register_to_stream_func[float32](f32)

	Add_f32 = Register_apply_func[float32, Apply_Args[float32]](f32, Add, Vector, Vector)
	Sub_f32 = Register_apply_func[float32, Apply_Args[float32]](f32, Sub, Vector, Vector)
	Mul_f32 = Register_apply_func[float32, Apply_Args[float32]](f32, Mul, Vector, Vector)
	Div_f32 = Register_apply_func[float32, Apply_Args[float32]](f32, Div, Vector, Vector)
	Pow_f32 = Register_apply_func[float32, Apply_Args[float32]](f32, Pow, Vector, Vector)
}

var New_f64 New_Stream_Func[float64]
var Set_f64 Set_Stream_Func[float64]
var To_f64 To_Stream_Func[float64]

var Add_f64 func(*Apply_Args[float64])
var Sub_f64 func(*Apply_Args[float64])
var Mul_f64 func(*Apply_Args[float64])
var Div_f64 func(*Apply_Args[float64])
var Pow_f64 func(*Apply_Args[float64])

func bind_f64_funcs() {
	New_f64 = register_new_stream_func[float64](f64)
	Set_f64 = register_set_stream_func[float64](f64)
	To_f64 = register_to_stream_func[float64](f64)

	Add_f64 = Register_apply_func[float64, Apply_Args[float64]](f64, Add, Vector, Vector)
	Sub_f64 = Register_apply_func[float64, Apply_Args[float64]](f64, Sub, Vector, Vector)
	Mul_f64 = Register_apply_func[float64, Apply_Args[float64]](f64, Mul, Vector, Vector)
	Div_f64 = Register_apply_func[float64, Apply_Args[float64]](f64, Div, Vector, Vector)
	Pow_f64 = Register_apply_func[float64, Apply_Args[float64]](f64, Pow, Vector, Vector)
}

type Stream[T Number] struct {
	veks        *T
	len_veks    uint
	len_scalars uint
}

func As_slice[T Number](stream Stream[T]) []T {
	return unsafe.Slice(stream.veks, stream.len_scalars)
}

type New_Stream_Func[T Number] func(ctx Ctx, len uint) (stream Stream[T])
type Set_Stream_Func[T Number] func(ctx Ctx, slice []T)
type To_Stream_Func[T Number] func(ctx Ctx, slice []T) (stream Stream[T])

func register_new_stream_func[T Number](n NumType) (new_stream_func New_Stream_Func[T]) {
	name := "New_" + string(n) + "_Stream"
	purego.RegisterLibFunc(&new_stream_func, lib, name)
	return
}
func register_set_stream_func[T Number](n NumType) (set_stream_func Set_Stream_Func[T]) {
	name := "Set_" + string(n) + "_Stream"
	var set_stream_func_raw func(ctx Ctx, ptr *T, len uint)
	purego.RegisterLibFunc(&set_stream_func_raw, lib, name)
	set_stream_func = func(ctx Ctx, slice []T) {
		set_stream_func_raw(ctx, unsafe.SliceData(slice), uint(len(slice)))
	}
	return
}

// just convenience func, it calls new and set stream in the lib internally
func register_to_stream_func[T Number](n NumType) (to_stream_func To_Stream_Func[T]) {
	name := "To_" + string(n) + "_Stream"
	var to_stream_func_raw func(ctx Ctx, ptr *T, len uint) Stream[T]
	purego.RegisterLibFunc(&to_stream_func_raw, lib, name)
	to_stream_func = func(ctx Ctx, slice []T) (stream Stream[T]) {
		stream = to_stream_func_raw(ctx, unsafe.SliceData(slice), uint(len(slice)))
		return
	}
	return
}

type Apply_Args[T Number] struct {
	A, B, C Stream[T]
}
type Num_Apply_Args[T Number] struct {
	A    T
	B, C Stream[T]
}
type Apply_Args_Num[T Number] struct {
	A Stream[T]
	B T
	C Stream[T]
}
type Apply_Args_Bool[T Number] struct {
	A, B Stream[T]
	C    Stream[uint8]
}
type Apply_Args_Num_Bool[T Number] struct {
	A Stream[T]
	B T
	C Stream[uint8]
}
type Num_Apply_Args_Bool[T Number] struct {
	A T
	B Stream[T]
	C Stream[uint8]
}
type Select_Apply_Args[T Number] struct {
	A, B Stream[T]
	Pred Stream[uint8]
	C    Stream[T]
}
type Num_Select_Apply_Args[T Number] struct {
	A    T
	B    Stream[T]
	Pred Stream[uint8]
	C    Stream[T]
}
type Select_Apply_Args_Num[T Number] struct {
	A    Stream[T]
	B    T
	Pred Stream[uint8]
	C    Stream[T]
}

// supports casts!
type Apply_Args_Single[T Number, O Number] struct {
	X Stream[T]
	Y Stream[O]
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

func Register_apply_func[T Number, AT Apply_Args_Type[T]](t NumType, op Op, a_t Operand_Variant, b_t Operand_Variant) (apply_func func(*AT)) {
	if a_t == Scalar && b_t == Scalar {
		log.Panicf("in Register_apply_func, num_type=%s, op=%s, at least one input to apply must be a vector", t, op)
	}
	name := string(op) + "_" + string(t)
	if a_t == Scalar {
		name = "Num_" + name
	} else if b_t == Scalar {
		name = name + "_Num"
	}
	purego.RegisterLibFunc(&apply_func, lib, name)
	return
}
func Register_apply_single_func[T Number, O Number](t NumType, o NumType, op Op1) (apply_func func(Apply_Args_Single[T, O])) {
	name := string(op) + "_" + string(t)
	if t != o {
		name = name + "_" + string(o)
	}
	purego.RegisterLibFunc(&apply_func, lib, name)
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
	if debug {
		bind_debug_ctx_funcs()
	} else {
		bind_ctx_funcs()
	}
	bind_f32_funcs()
	bind_f64_funcs()
}
