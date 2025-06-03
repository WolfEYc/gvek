package gvek

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

func get_dylib_path(target_microarch string) string {
	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		// only arm64 macos supported (you can build your own intel macos zvek but why?)
		case "arm":
			switch target_microarch {
			default:
				return "libzvek-arm64-m1.dylib"
			}
		default:
			panic("arm64 is only supported macos cpu arch")
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm":
			switch target_microarch {
			default: // maybe should include neoverse-n1 or perhaps even cortex or something for compat
				return "libzvek-arm64-neoverse-v2.so"
			}
		case "amd64":
			switch target_microarch {
			case "zen4": // speeeeeeeed!
				return "libzvek-amd64-zen4.so"
			default: // still AVX2 so should be decent!
				return "libzvek-amd64-haswell.so"
			}
		default:
			panic("amd64 and arm are only supported linux cpu arch")
		}
	default:
		if runtime.GOOS == "windows" {
			panic("Windows user DEEEEEETECTED, everyone quick! laugh at them!")
		}
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

func Init(target_microarch string) {
	lib, err := purego.Dlopen(get_dylib_path(target_microarch), purego.RTLD_LAZY)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&Add_f32, lib, "Add_f32")
}
