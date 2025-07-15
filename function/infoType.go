package function

import "time"

/*
	 *leval 级别为 2的时候返回的结构体类型
		typedef struct _USER_INFO_2 {
		  LPWSTR usri2_name;
		  LPWSTR usri2_password;
		  DWORD  usri2_password_age;
		  DWORD  usri2_priv;
		  LPWSTR usri2_home_dir;
		  LPWSTR usri2_comment;
		  DWORD  usri2_flags;
		  LPWSTR usri2_script_path;
		  DWORD  usri2_auth_flags;
		  LPWSTR usri2_full_name;
		  LPWSTR usri2_usr_comment;
		  LPWSTR usri2_parms;
		  LPWSTR usri2_workstations;
		  DWORD  usri2_last_logon;
		  DWORD  usri2_last_logoff;
		  DWORD  usri2_acct_expires;
		  DWORD  usri2_max_storage;
		  DWORD  usri2_units_per_week;
		  PBYTE  usri2_logon_hours;
		  DWORD  usri2_bad_pw_count;
		  DWORD  usri2_num_logons;
		  LPWSTR usri2_logon_server;
		  DWORD  usri2_country_code;
		  DWORD  usri2_code_page;
		} USER_INFO_2, *PUSER_INFO_2, *LPUSER_INFO_2;
*/
type userInfo struct {
	Usri2Name         *uint16 // 用户名
	Usri2Password     *uint16 // 用户密码（仅用于设置，不会返回）
	Usri2PasswordAge  uint32  // 密码年龄（单位：秒）
	Usri2Priv         uint32  // 权限级别（0=访客，1=普通用户，2=管理员）
	Usri2HomeDir      *uint16 // 主目录（用户登录后的目录）
	Usri2Comment      *uint16 // 注释（通常为管理员备注）
	Usri2Flags        uint32  // 用户账户控制标志（如启用/禁用）
	Usri2ScriptPath   *uint16 // 登录脚本路径
	Usri2AuthFlags    uint32  // 认证标志（通常未使用）
	Usri2FullName     *uint16 // 用户全名
	Usri2UsrComment   *uint16 // 用户注释（用户可读的）
	Usri2Parms        *uint16 // 管理参数（自定义字段）
	Usri2Workstations *uint16 // 允许登录的工作站（以逗号分隔）
	Usri2LastLogon    uint32  // 上次登录时间（UNIX时间戳）
	Usri2LastLogoff   uint32  // 上次注销时间（很少使用）
	Usri2AcctExpires  uint32  // 账户过期时间（0xffffffff 表示永不过期）
	Usri2MaxStorage   uint32  // 最大存储空间（以字节为单位，0xffffffff 表示无限制）
	Usri2UnitsPerWeek uint32  // 一周的时间单位数（用于 logon_hours）
	Usri2LogonHours   *byte   // 登录时间限制（bitmask，按小时排列）
	Usri2BadPwCount   uint32  // 错误密码尝试次数（连续错误次数）
	Usri2NumLogons    uint32  // 登录次数（成功的）
	Usri2LogonServer  *uint16 // 登录的服务器名称
	Usri2CountryCode  uint32  // 国家/地区代码（如 86 表示中国）
	Usri2CodePage     uint32  // 代码页（如 936 表示简体中文）
}

// userInfo 信息辅助函数
func getUserPriv(priv uint32) string {
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

func getUserFlags(flags uint32) string {
	if flags&0x0001 != 0 {
		return "账户禁用"
	}
	return "正常"
}

func formatUnixTime(ts uint32) string {
	if ts == 0 {
		return "从未登录"
	}
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}

func formatExpiry(expiry uint32) string {
	if expiry == 0xFFFFFFFF {
		return "永不过期"
	}
	return time.Unix(int64(expiry), 0).Format("2006-01-02 15:04:05")
}
