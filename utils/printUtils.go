package utils

import (
	"fmt"
	"unsafe"
	"windowsFinder/structs"
)

func PrintUserInfoLevel2(bufptr uintptr, entriesread uint32) {
	infoSize := unsafe.Sizeof(structs.UserInfolevel2{})
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
	for i := uint32(0); i < entriesread; i++ {
		user := (*structs.UserInfolevel2)(unsafe.Pointer(bufptr + uintptr(i)*infoSize))
		fmt.Printf("─────────────────────────────────────%d──────────────────────────────────────────────\n", i+1)
		fmt.Printf("用户名        : %s\n", UTF16PtrToString(user.Usri2Name))
		fmt.Printf("用户全名      : %s\n", UTF16PtrToString(user.Usri2FullName))
		fmt.Printf("权限级别      : %d (%s)\n", user.Usri2Priv, structs.GetUserPriv(user.Usri2Priv))
		fmt.Printf("主目录        : %s\n", UTF16PtrToString(user.Usri2HomeDir))
		fmt.Printf("登录服务器    : %s\n", UTF16PtrToString(user.Usri2LogonServer))
		fmt.Printf("注释          : %s\n", UTF16PtrToString(user.Usri2Comment))
		fmt.Printf("账户状态      : %s\n", structs.GetUserFlags(user.Usri2Flags))
		fmt.Printf("登录次数      : %d\n", user.Usri2NumLogons)
		fmt.Printf("错误密码次数  : %d\n", user.Usri2BadPwCount)
		fmt.Printf("上次登录时间  : %s\n", structs.FormatUnixTime(user.Usri2LastLogon))
		fmt.Printf("账户过期时间  : %s\n", structs.FormatExpiry(user.Usri2AcctExpires))
		fmt.Println("────────────────────────────────────────────────────────────────────────────────────")
	}
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
}
