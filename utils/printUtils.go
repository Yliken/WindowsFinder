package utils

import (
	"fmt"
	"unsafe"
	"windowsFinder/structs"
)

func PrintUserInfoLevel0(buffer uintptr, entriesread uint32) {
	infosize := unsafe.Sizeof(structs.UserInfolevel0{})
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
	fmt.Printf("[+] 用户: ")
	for i := uint32(0); i < entriesread; i++ {
		user := (*structs.UserInfolevel0)(unsafe.Pointer(uintptr(buffer) + infosize*uintptr(i)))
		fmt.Printf("%d: %s ", i+1, UTF16PtrToString(user.UsriName))
	}
	fmt.Println("")
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
}

func PrintUserInfoLevel1(bufptr uintptr, entriesread uint32) {
	infosize := unsafe.Sizeof(structs.UserInfolevel1{})
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
	for i := uint32(0); i < entriesread; i++ {
		user := (*structs.UserInfolevel1)(unsafe.Pointer(bufptr + uintptr(i)*infosize))
		fmt.Printf("─────────────────────────────────────%d──────────────────────────────────────────────\n", i+1)
		fmt.Printf("用户名        : %s\n", UTF16PtrToString(user.UsriName))
		fmt.Printf("权限级别      : %d (%s)\n", user.UsriPriv, structs.GetUserPriv(user.UsriPriv))
		fmt.Printf("主目录        : %s\n", UTF16PtrToString(user.UsriHomeDir))
		fmt.Printf("密码年龄      : %d 秒\n", user.UsriPasswordAge)
		fmt.Printf("用户注释      : %s\n", UTF16PtrToString(user.UsriUsrComment))
		fmt.Printf("登录脚本路径  : %s\n", UTF16PtrToString(user.Usriscriptpath))
		fmt.Printf("账户状态      : %s\n", structs.GetUserFlags(user.UsriFlags))
		fmt.Println("────────────────────────────────────────────────────────────────────────────────────")
	}
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
}

func PrintUserInfoLevel2(bufptr uintptr, entriesread uint32) {
	infoSize := unsafe.Sizeof(structs.UserInfolevel2{})
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")
	for i := uint32(0); i < entriesread; i++ {
		user := (*structs.UserInfolevel2)(unsafe.Pointer(bufptr + uintptr(i)*infoSize))
		fmt.Printf("─────────────────────────────────────%d──────────────────────────────────────────────\n", i+1)
		fmt.Printf("用户名        : %s\n", UTF16PtrToString(user.UsriName))
		fmt.Printf("用户全名      : %s\n", UTF16PtrToString(user.UsriFullName))
		fmt.Printf("权限级别      : %d (%s)\n", user.UsriPriv, structs.GetUserPriv(user.UsriPriv))
		fmt.Printf("主目录        : %s\n", UTF16PtrToString(user.UsriHomeDir))
		fmt.Printf("登录服务器    : %s\n", UTF16PtrToString(user.UsriLogonServer))
		fmt.Printf("注释          : %s\n", UTF16PtrToString(user.UsriComment))
		fmt.Printf("账户状态      : %s\n", structs.GetUserFlags(user.UsriFlags))
		fmt.Printf("登录次数      : %d\n", user.UsriNumLogons)
		fmt.Printf("错误密码次数  : %d\n", user.UsriBadPwCount)
		fmt.Printf("上次登录时间  : %s\n", structs.FormatUnixTime(user.UsriLastLogon))
		fmt.Printf("账户过期时间  : %s\n", structs.FormatExpiry(user.UsriAcctExpires))
		fmt.Printf("密码年龄      : %d 秒\n", user.UsriPasswordAge)
		fmt.Printf("登录脚本路径  : %s\n", UTF16PtrToString(user.UsriScriptPath))
		fmt.Printf("用户注释      : %s\n", UTF16PtrToString(user.UsriUsrComment))
		fmt.Printf("允许工作站    : %s\n", UTF16PtrToString(user.UsriWorkstations))
		fmt.Printf("上次注销时间  : %s\n", structs.FormatUnixTime(user.UsriLastLogoff))
		fmt.Printf("账户最大存储  : %d 字节\n", user.UsriMaxStorage)
		fmt.Printf("国家代码      : %d\n", user.UsriCountryCode)
		fmt.Printf("代码页        : %d\n", user.UsriCodePage)
		fmt.Println("────────────────────────────────────────────────────────────────────────────────────")
	}
	fmt.Println("══════════════════════════════════════════════用户名列表══════════════════════════════════════════════")

}

func PrintUserInfoLevel3(bufptr uintptr, entriesread uint32) {
	infoSize := unsafe.Sizeof(structs.UserInfoLevel3{})
	fmt.Println("══════════════════════════════════════════════ 用户信息列表 ══════════════════════════════════════════════")
	for i := uint32(0); i < entriesread; i++ {
		user := (*structs.UserInfoLevel3)(unsafe.Pointer(bufptr + uintptr(i)*infoSize))
		fmt.Printf("───────────────────────────────────── %d ──────────────────────────────────────────────\n", i+1)
		fmt.Printf("用户名           : %s\n", UTF16PtrToString(user.Usri3Name))
		// 密码不打印
		fmt.Printf("密码年龄         : %d 秒\n", user.Usri3PasswordAge)
		fmt.Printf("权限级别         : %d (%s)\n", user.Usri3Priv, structs.GetUserPriv(user.Usri3Priv))
		fmt.Printf("主目录           : %s\n", UTF16PtrToString(user.Usri3HomeDir))
		fmt.Printf("注释             : %s\n", UTF16PtrToString(user.Usri3Comment))
		fmt.Printf("账户状态         : %d\n", user.Usri3Flags) // 可做标志位解析
		fmt.Printf("登录脚本路径     : %s\n", UTF16PtrToString(user.Usri3ScriptPath))
		fmt.Printf("认证标志         : %d\n", user.Usri3AuthFlags)
		fmt.Printf("用户全名         : %s\n", UTF16PtrToString(user.Usri3FullName))
		fmt.Printf("用户注释         : %s\n", UTF16PtrToString(user.Usri3UsrComment))
		fmt.Printf("管理参数         : %s\n", UTF16PtrToString(user.Usri3Parms))
		fmt.Printf("允许登录工作站   : %s\n", UTF16PtrToString(user.Usri3Workstations))
		fmt.Printf("上次登录时间     : %s\n", structs.FormatUnixTime(user.Usri3LastLogon))
		fmt.Printf("上次注销时间     : %s\n", structs.FormatUnixTime(user.Usri3LastLogoff))
		fmt.Printf("账户过期时间     : %s\n", structs.FormatExpiry(user.Usri3AcctExpires))
		fmt.Printf("最大存储         : %d 字节\n", user.Usri3MaxStorage)
		fmt.Printf("一周时间单位数   : %d\n", user.Usri3UnitsPerWeek)
		fmt.Printf("错误密码次数     : %d\n", user.Usri3BadPwCount)
		fmt.Printf("登录次数         : %d\n", user.Usri3NumLogons)
		fmt.Printf("登录服务器       : %s\n", UTF16PtrToString(user.Usri3LogonServer))
		fmt.Printf("国家代码         : %d\n", user.Usri3CountryCode)
		fmt.Printf("代码页           : %d\n", user.Usri3CodePage)
		fmt.Printf("用户ID           : %d\n", user.Usri3UserId)
		fmt.Printf("主组ID           : %d\n", user.Usri3PrimaryGroupId)
		fmt.Printf("用户配置文件路径 : %s\n", UTF16PtrToString(user.Usri3Profile))
		fmt.Printf("用户主目录驱动器 : %s\n", UTF16PtrToString(user.Usri3HomeDirDrive))
		fmt.Printf("密码是否过期     : %d\n", user.Usri3PasswordExpired)
		fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────────")
	}
	fmt.Println("══════════════════════════════════════════════ 用户信息列表 ══════════════════════════════════════════════")
}
