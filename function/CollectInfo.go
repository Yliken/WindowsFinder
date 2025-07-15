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

func collectSystemUserInfo() {
	//加载 相关dll 以及 函数
	netapi32 := syscall.NewLazyDLL("netapi32.dll")
	netUserEnum := netapi32.NewProc("NetUserEnum")
	procNetApiBufferFree := netapi32.NewProc("NetApiBufferFree")

	ret, _, _ := netUserEnum.Call(
		0,
		2, // level 级别为 2  查询更多信息
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

	infoSize := unsafe.Sizeof(userInfo{})
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
	for i := uint32(0); i < entriesread; i++ {
		user := (*userInfo)(unsafe.Pointer(bufptr + uintptr(i)*infoSize))
		fmt.Printf("─────────────────────────────────────%d──────────────────────────────────────────────\n", i+1)
		fmt.Printf("用户名        : %s\n", utils.UTF16PtrToString(user.Usri2Name))
		fmt.Printf("用户全名      : %s\n", utils.UTF16PtrToString(user.Usri2FullName))
		fmt.Printf("权限级别      : %d (%s)\n", user.Usri2Priv, getUserPriv(user.Usri2Priv))
		fmt.Printf("主目录        : %s\n", utils.UTF16PtrToString(user.Usri2HomeDir))
		fmt.Printf("登录服务器    : %s\n", utils.UTF16PtrToString(user.Usri2LogonServer))
		fmt.Printf("注释          : %s\n", utils.UTF16PtrToString(user.Usri2Comment))
		fmt.Printf("账户状态      : %s\n", getUserFlags(user.Usri2Flags))
		fmt.Printf("登录次数      : %d\n", user.Usri2NumLogons)
		fmt.Printf("错误密码次数  : %d\n", user.Usri2BadPwCount)
		fmt.Printf("上次登录时间  : %s\n", formatUnixTime(user.Usri2LastLogon))
		fmt.Printf("账户过期时间  : %s\n", formatExpiry(user.Usri2AcctExpires))
		fmt.Println("────────────────────────────────────────────────────────────────────────────────────")
	}
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
}
