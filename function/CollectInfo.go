package function

import (
	"fmt"
	"syscall"
	"unsafe"
	"windowsFinder/utils"
)

// 接受 NetUserEnumerate 返回的数据
var (
	bufptr        uintptr
	entriesread   uint32
	totalentries  uint32
	resume_handle uint32
)

// 参数 level 查询更多的信息
func collectSystemUserInfo(level uint32) {
	//加载 相关dll 以及 函数
	netapi32 := syscall.NewLazyDLL("netapi32.dll")
	netUserEnum := netapi32.NewProc("NetUserEnum")
	procNetApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	ret, _, _ := netUserEnum.Call(
		0,
		uintptr(level), // level 级别为 2  查询更多信息
		0,
		uintptr(unsafe.Pointer(&bufptr)),
		0xFFFFFFFF, // 自动选择缓冲区大小
		uintptr(unsafe.Pointer(&entriesread)),
		uintptr(unsafe.Pointer(&totalentries)),
		uintptr(unsafe.Pointer(&resume_handle)),
	)
	if ret != 0 {
		fmt.Printf("NetUserEnum failed with code: %d\n", ret)
		return
	}
	// 释放buffer
	defer procNetApiBufferFree.Call(bufptr)
	switch level {
	case uint32(0):
		utils.PrintUserInfoLevel0(bufptr, entriesread)
	case uint32(1):
		utils.PrintUserInfoLevel1(bufptr, entriesread)
	case uint32(2):
		utils.PrintUserInfoLevel2(bufptr, entriesread)
	default:
		return
	}
}
