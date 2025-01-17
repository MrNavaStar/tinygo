// CGo errors:
//     testdata/errors.go:4:2: warning: some warning
//     testdata/errors.go:11:9: error: unknown type name 'someType'
//     testdata/errors.go:26:5: warning: another warning
//     testdata/errors.go:13:23: unexpected token ), expected end of expression
//     testdata/errors.go:21:26: unexpected token ), expected end of expression
//     testdata/errors.go:16:33: unexpected token ), expected end of expression
//     testdata/errors.go:17:34: unexpected token ), expected end of expression
//     -: unexpected token INT, expected end of expression

// Type checking errors after CGo processing:
//     testdata/errors.go:102: cannot use 2 << 10 (untyped int constant 2048) as C.char value in variable declaration (overflows)
//     testdata/errors.go:105: unknown field z in struct literal
//     testdata/errors.go:108: undefined: C.SOME_CONST_1
//     testdata/errors.go:110: cannot use C.SOME_CONST_3 (untyped int constant 1234) as byte value in variable declaration (overflows)
//     testdata/errors.go:112: undefined: C.SOME_CONST_4
//     testdata/errors.go:114: undefined: C.SOME_CONST_b
//     testdata/errors.go:116: undefined: C.SOME_CONST_startspace
//     testdata/errors.go:119: undefined: C.SOME_PARAM_CONST_invalid

package main

import "unsafe"

var _ unsafe.Pointer

//go:linkname C.CString runtime.cgo_CString
func C.CString(string) *C.char

//go:linkname C.GoString runtime.cgo_GoString
func C.GoString(*C.char) string

//go:linkname C.__GoStringN runtime.cgo_GoStringN
func C.__GoStringN(*C.char, uintptr) string

func C.GoStringN(cstr *C.char, length C.int) string {
	return C.__GoStringN(cstr, uintptr(length))
}

//go:linkname C.__GoBytes runtime.cgo_GoBytes
func C.__GoBytes(unsafe.Pointer, uintptr) []byte

func C.GoBytes(ptr unsafe.Pointer, length C.int) []byte {
	return C.__GoBytes(ptr, uintptr(length))
}

//go:linkname C.__CBytes runtime.cgo_CBytes
func C.__CBytes([]byte) unsafe.Pointer

func C.CBytes(b []byte) unsafe.Pointer {
	return C.__CBytes(b)
}

type (
	C.char      uint8
	C.schar     int8
	C.uchar     uint8
	C.short     int16
	C.ushort    uint16
	C.int       int32
	C.uint      uint32
	C.long      int32
	C.ulong     uint32
	C.longlong  int64
	C.ulonglong uint64
)
type C.struct_point_t struct {
	x C.int
	y C.int
}
type C.point_t = C.struct_point_t

const C.SOME_CONST_3 = 1234
const C.SOME_PARAM_CONST_valid = 3 + 4
