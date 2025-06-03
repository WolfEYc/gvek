package gvek

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

func get_dylib_path() string {
	switch runtime.GOOS {
	case "darwin":
		return "./libzvek.dylib"
	case "linux":
		return "./libzvek.so"
	default:
		// haha screw windows users, not that anything is going to run fast on windows anyways.
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

type Stream struct {
	veks        unsafe.Pointer
	len_veks    uint
	len_scalars uint
}

type Apply_Args struct {
	a, b, c Stream
}

var Add_f32 func(args *Apply_Args)

func Init() {
	lib, err := purego.Dlopen(get_dylib_path(), purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&Add_f32, lib, "Add_f32")
}
