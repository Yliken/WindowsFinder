package structs

import "time"

type UserInfolevel0 struct {
	UsriName *uint16 // 用户密码（仅用于设置，不会返回）
}

type UserInfolevel1 struct {
	UsriName        *uint16
	UsriPassword    *uint16 // 用户密码（仅用于设置，不会返回）
	UsriPasswordAge uint32  // 密码年龄（单位：秒）
	UsriPriv        uint32  // 权限级别（0=访客，1=普通用户，2=管理员）
	UsriHomeDir     *uint16 // 主目录（用户登录后的目录）
	UsriUsrComment  *uint16 // 用户注释（用户可读的）
	UsriFlags       uint32  // 用户账户控制标志（如启用/禁用）
	Usriscriptpath  *uint16 // 登录脚本路径
}

type UserInfolevel2 struct {
	UsriName         *uint16 // 用户名
	UsriPassword     *uint16 // 用户密码（仅用于设置，不会返回）
	UsriPasswordAge  uint32  // 密码年龄（单位：秒）
	UsriPriv         uint32  // 权限级别（0=访客，1=普通用户，2=管理员）
	UsriHomeDir      *uint16 // 主目录（用户登录后的目录）
	UsriComment      *uint16 // 注释（通常为管理员备注）
	UsriFlags        uint32  // 用户账户控制标志（如启用/禁用）
	UsriScriptPath   *uint16 // 登录脚本路径
	UsriAuthFlags    uint32  // 认证标志（通常未使用）
	UsriFullName     *uint16 // 用户全名
	UsriUsrComment   *uint16 // 用户注释（用户可读的）
	UsriParms        *uint16 // 管理参数（自定义字段）
	UsriWorkstations *uint16 // 允许登录的工作站（以逗号分隔）
	UsriLastLogon    uint32  // 上次登录时间（UNIX时间戳）
	UsriLastLogoff   uint32  // 上次注销时间（很少使用）
	UsriAcctExpires  uint32  // 账户过期时间（0xffffffff 表示永不过期）
	UsriMaxStorage   uint32  // 最大存储空间（以字节为单位，0xffffffff 表示无限制）
	UsriUnitsPerWeek uint32  // 一周的时间单位数（用于 logon_hours）
	UsriLogonHours   *byte   // 登录时间限制（bitmask，按小时排列）
	UsriBadPwCount   uint32  // 错误密码尝试次数（连续错误次数）
	UsriNumLogons    uint32  // 登录次数（成功的）
	UsriLogonServer  *uint16 // 登录的服务器名称
	UsriCountryCode  uint32  // 国家/地区代码（如 86 表示中国）
	UsriCodePage     uint32  // 代码页（如 936 表示简体中文）
}
type UserInfoLevel3 struct {
	Usri3Name            *uint16
	Usri3Password        *uint16 // 通常只用于设置，不会返回
	Usri3PasswordAge     uint32
	Usri3Priv            uint32
	Usri3HomeDir         *uint16
	Usri3Comment         *uint16
	Usri3Flags           uint32
	Usri3ScriptPath      *uint16
	Usri3AuthFlags       uint32
	Usri3FullName        *uint16
	Usri3UsrComment      *uint16
	Usri3Parms           *uint16
	Usri3Workstations    *uint16
	Usri3LastLogon       uint32
	Usri3LastLogoff      uint32
	Usri3AcctExpires     uint32
	Usri3MaxStorage      uint32
	Usri3UnitsPerWeek    uint32
	Usri3LogonHours      *byte
	Usri3BadPwCount      uint32
	Usri3NumLogons       uint32
	Usri3LogonServer     *uint16
	Usri3CountryCode     uint32
	Usri3CodePage        uint32
	Usri3UserId          uint32
	Usri3PrimaryGroupId  uint32
	Usri3Profile         *uint16
	Usri3HomeDirDrive    *uint16
	Usri3PasswordExpired uint32
}

// userInfo 信息辅助函数
func GetUserPriv(priv uint32) string {
	switch priv {
	case 0:
		return "访客"
	case 1:
		return "普通用户"
	case 2:
		return "管理员"
	default:
		return "未知"
	}
}

func GetUserFlags(flags uint32) string {
	if flags&0x0001 != 0 {
		return "账户禁用"
	}
	return "正常"
}

func FormatUnixTime(ts uint32) string {
	if ts == 0 {
		return "从未登录"
	}
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}

func FormatExpiry(expiry uint32) string {
	if expiry == 0xFFFFFFFF {
		return "永不过期"
	}
	return time.Unix(int64(expiry), 0).Format("2006-01-02 15:04:05")
}
