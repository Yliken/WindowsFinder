package utils

import (
	"syscall"
	"unsafe"
)

func UTF16PtrToString(ptr *uint16) string {
	if ptr == nil {
		return ""
	}
	return syscall.UTF16ToString((*[4096]uint16)(unsafe.Pointer(ptr))[:])
}
func StringToUTF16Ptr(str string) *uint16 {
	fromString, _ := syscall.UTF16PtrFromString(str)
	return fromString
}
