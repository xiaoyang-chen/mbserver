package mbserver

import "unsafe"

func CopyBytes(src []byte) (dst []byte) {

	dst = make([]byte, len(src))
	copy(dst, src)
	return
}

func CopyUint16(src []uint16) (dst []uint16) {

	dst = make([]uint16, len(src))
	copy(dst, src)
	return
}

func ConcatStrings(ss ...string) string {

	var length = len(ss)
	if length == 0 {
		return ""
	}
	var i, n = 0, 0
	for i = 0; i < length; i++ {
		n += len(ss[i])
	}
	var b = make([]byte, 0, n)
	for i = 0; i < length; i++ {
		b = append(b, ss[i]...)
	}
	return Bytes2Str(b)
}

func Str2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
	// return *(*[]byte)(unsafe.Pointer(&struct {
	// 	string
	// 	Cap int
	// }{s, len(s)}))
}

func Bytes2Str(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
	// return *(*string)(unsafe.Pointer(&b))
}
