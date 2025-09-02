package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"windowsFinder/config"
	"windowsFinder/function"
)

func init() {
	config.Banner()
}
func list() {
	fmt.Println()
	fmt.Println(" ****Windows finder menu**** ")
	fmt.Println("1. 进行一些基础信息收集")
	fmt.Println("2. 进行用户信息信息收集")
	fmt.Println("3. 收集RDP相关的信息(事件 21, 事件25, 事件 1149)")
	fmt.Println("4. 收集系统进程信息")
}

func Menu() {
	for {
		list()
		fmt.Print("请选择操作编号：")
		reader := bufio.NewReader(os.Stdin)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("输入无效，请输入数字")
			continue
		}
		switch choice {
		case 1:
			function.BasicCollect()
		case 2:
			UserInfoMean()
		case 3:
			function.CollectRDPInfo()
		case 4:
			function.ProcessEnum()
		default:
			return
		}
	}
}

func UserInfoMean() {
	fmt.Println("tips:")
	fmt.Println("level => 0 只包含用户名")
	fmt.Println("level => 1 在 level 0 的基础上，增加了密码、权限等级、账号状态（启用/禁用）、注释、主目录和登录脚本路径")
	fmt.Println("level => 2 在 level 1 的基础上，新增了用户全名、登录服务器、登录次数、上次登录时间、账户过期时间、密码年龄、允许登录工作站、用户注释、账户最大存储、国家代码和代码页等详细信息")
	fmt.Println("level => 3 在 level 2 的基础上，进一步增加了认证标志、管理参数、主组ID、用户ID、用户配置文件路径、用户主目录驱动器、密码是否过期、登录失败次数、一周时间单位数等高级账户属性信息")

	// 丢弃输入缓冲区中残留的换行符（仅第一次调用时需要）
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // 吞掉残留的 Enter，防止后面立即 return

	fmt.Print("请输入 level 值: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	level, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("输入无效，请输入 0 到 3 的整数")
		return
	}
	fmt.Println("你输入的是 level:", level)
	function.CollectSystemUserInfo(uint32(level))
}
