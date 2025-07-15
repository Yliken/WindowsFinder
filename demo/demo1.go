package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"windowsFinder/utils"
)

var (
	modNetapi32          = syscall.NewLazyDLL("Netapi32.dll")
	procNetUserEnum      = modNetapi32.NewProc("NetUserEnum")
	procNetApiBufferFree = modNetapi32.NewProc("NetApiBufferFree")
)

// USER_INFO_0 结构体表示用户条目，字段仅包含用户名
type USER_INFO_0 struct {
	Usri0_name *uint16
}

func main() {
	var (
		bufPtr       uintptr
		entriesRead  uint32
		totalEntries uint32
		resumeHandle uint32
	)

	// 调用 NetUserEnum
	// 参数说明：
	// servername = 0（本地机器）
	// level = 0（只返回用户名）
	// filter = 2（普通用户）
	// bufPtr = 输出缓冲区地址
	// prefmaxlen = 最大缓冲长度，-1 表示自动分配
	// entriesRead = 返回实际读取数量
	// totalEntries = 返回总条目数
	// resumeHandle = 继续句柄（用于分页查询）

	ret, _, _ := procNetUserEnum.Call(
		0, // local server
		0, // level 0
		0, // FILTER_NORMAL_ACCOUNT
		uintptr(unsafe.Pointer(&bufPtr)),
		0xFFFFFFFF, // pref max len (unlimited)
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	const NERR_Success = 0
	if ret != NERR_Success {
		fmt.Printf("NetUserEnum failed with code: %d\n", ret)
		return
	}

	defer procNetApiBufferFree.Call(bufPtr)

	fmt.Println("Windows 用户名列表：")
	entrySize := unsafe.Sizeof(USER_INFO_0{})
	for i := uint32(0); i < entriesRead; i++ {
		user := (*USER_INFO_0)(unsafe.Pointer(bufPtr + uintptr(i)*entrySize))
		username := utils.UTF16PtrToString(user.Usri0_name)
		fmt.Printf(" - %s\n", username)
	}

}
