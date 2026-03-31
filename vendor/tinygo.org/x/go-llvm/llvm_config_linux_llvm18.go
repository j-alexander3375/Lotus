//go:build !byollvm && linux && !llvm14 && !llvm15 && !llvm16 && !llvm17

package llvm

// #cgo CPPFLAGS: -I/usr/include/llvm -I/usr/include/llvm-c -D_GNU_SOURCE -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
// #cgo CXXFLAGS: -std=c++17
// #cgo LDFLAGS: -L/usr/lib -lLLVM-21
import "C"

type run_build_sh int
