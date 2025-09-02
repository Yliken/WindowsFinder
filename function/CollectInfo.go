package function

import (
	"encoding/xml"
	"fmt"
	"syscall"
	"time"
	"unsafe"
	"windowsFinder/structs"
	"windowsFinder/utils"
)

// 接受 NetUserEnumerate 返回的数据
// 查询 系统用户 信息
var (
	bufptr        uintptr
	entriesread   uint32
	totalentries  uint32
	resume_handle uint32
)

// 查询 系统用户 信息
// 参数 level 查询更多的信息
func CollectSystemUserInfo(level uint32) {
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
	case uint32(3):
		utils.PrintUserInfoLevel3(bufptr, entriesread)
	default:
		return
	}
}

// 搜集RDP 相关的日志信息 时间 1149 事件21 事件 25
func CollectRDPInfo() {
	const (
		EvtQueryChannelPath = 0x00000001
		EvtRenderEventXml   = 1
		BufferSize          = 1 << 16
		MaxEvents           = 6
		ERROR_NO_MORE_ITEMS = 259
	)

	type Event21Or25Info struct {
		ID        int
		User      string
		Address   string
		SessionID int
		Time      string
	}

	wevtapidll := syscall.NewLazyDLL("wevtapi.dll")
	procEvtQuery := wevtapidll.NewProc("EvtQuery")
	procEvtNext := wevtapidll.NewProc("EvtNext")
	procEvtRender := wevtapidll.NewProc("EvtRender")
	procEvtClose := wevtapidll.NewProc("EvtClose")

	channels := []string{
		"Microsoft-Windows-TerminalServices-RemoteConnectionManager/Operational",
		"Microsoft-Windows-TerminalServices-LocalSessionManager/Operational",
	}
	query := "*[System[(EventID=21 or EventID=25 or EventID=1149)]]"

	for _, channel := range channels {
		fmt.Printf("\n==== 查询通道: %s ====\n", channel)
		handle, _, err := procEvtQuery.Call(
			0,
			uintptr(unsafe.Pointer(utils.StringToUTF16Ptr(channel))),
			uintptr(unsafe.Pointer(utils.StringToUTF16Ptr(query))),
			uintptr(EvtQueryChannelPath),
		)
		if handle == 0 {
			fmt.Printf("[!] EvtQuery 失败: %v\n", err)
			continue
		}
		defer procEvtClose.Call(handle)

		var stored21 []Event21Or25Info
		var stored25 []Event21Or25Info
		total := 0

		for {
			var events [MaxEvents]syscall.Handle
			var returned uint32

			ret, _, err := procEvtNext.Call(
				handle,
				uintptr(MaxEvents),
				uintptr(unsafe.Pointer(&events)),
				0,
				0,
				uintptr(unsafe.Pointer(&returned)),
			)
			if ret == 0 || returned == 0 {
				lastErr := syscall.GetLastError()
				if errno, ok := lastErr.(syscall.Errno); ok && errno == ERROR_NO_MORE_ITEMS {
					fmt.Println("[*] 所有事件已读取完毕")
					break
				}
				fmt.Println("[!] 没有事件:", err)
				break
			}

			total += int(returned)
			for i := uint32(0); i < returned; i++ {
				buf := make([]uint16, BufferSize)
				var used, propCount uint32
				procEvtRender.Call(
					0,
					uintptr(events[i]),
					EvtRenderEventXml,
					uintptr(len(buf)*2),
					uintptr(unsafe.Pointer(&buf[0])),
					uintptr(unsafe.Pointer(&used)),
					uintptr(unsafe.Pointer(&propCount)),
				)

				xmlStr := syscall.UTF16ToString(buf[:used/2])
				var generic structs.Generic
				err := xml.Unmarshal([]byte(xmlStr), &generic)
				if err != nil {
					fmt.Println("[!] 无法识别事件 ID:", err)
					continue
				}

				switch generic.System.EventID {
				case 1149:
					var evt structs.Event1149
					if err := xml.Unmarshal([]byte(xmlStr), &evt); err != nil {
						fmt.Println("[!] 解析 Event1149 失败:", err)
						continue
					}
					fmt.Printf("🟢 1149事件：用户=%s IP=%s 时间=%s\n",
						evt.UserData.EventXML.Param1,
						evt.UserData.EventXML.Param3,
						evt.System.TimeCreated.SystemTime,
					)

				case 21:
					var evt structs.Event21
					if err := xml.Unmarshal([]byte(xmlStr), &evt); err != nil {
						fmt.Println("[!] 解析 Event21 失败:", err)
						continue
					}
					stored21 = append(stored21, Event21Or25Info{
						ID:        21,
						User:      evt.UserData.EventXML.User,
						Address:   evt.UserData.EventXML.Address,
						SessionID: evt.UserData.EventXML.SessionID,
						Time:      evt.System.TimeCreated.SystemTime,
					})

				case 25:
					var evt structs.Event21
					if err := xml.Unmarshal([]byte(xmlStr), &evt); err != nil {
						fmt.Println("[!] 解析 Event25 失败:", err)
						continue
					}
					stored25 = append(stored25, Event21Or25Info{
						ID:        25,
						User:      evt.UserData.EventXML.User,
						Address:   evt.UserData.EventXML.Address,
						SessionID: evt.UserData.EventXML.SessionID,
						Time:      evt.System.TimeCreated.SystemTime,
					})

				default:
					fmt.Printf("未知事件 ID: %d\n", generic.System.EventID)
				}
				procEvtClose.Call(uintptr(events[i]))
			}
		}

		// 分别打印事件21和事件25
		fmt.Printf("\n--- 汇总打印 事件21，共 %d 条 ---\n", len(stored21))
		for _, evt := range stored21 {
			fmt.Printf("🔵 事件21：用户=%s 地址=%s SessionID=%d 时间=%s\n",
				evt.User, evt.Address, evt.SessionID, evt.Time)
		}

		fmt.Printf("\n--- 汇总打印 事件25，共 %d 条 ---\n", len(stored25))
		for _, evt := range stored25 {
			fmt.Printf("🔵 事件25：用户=%s 地址=%s SessionID=%d 时间=%s\n",
				evt.User, evt.Address, evt.SessionID, evt.Time)
		}
	}
}

// 收集进程相关的信息
const (
	TH32CS_SNAPPROCESS                = 0x00000002
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	TOKEN_QUERY                       = 0x0008
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	advapi32                     = syscall.NewLazyDLL("advapi32.dll")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procGetProcessTimes          = kernel32.NewProc("GetProcessTimes")
	procOpenProcessToken         = advapi32.NewProc("OpenProcessToken")
	procGetTokenInformation      = advapi32.NewProc("GetTokenInformation")
	procLookupAccountSidW        = advapi32.NewProc("LookupAccountSidW")
)

func filetimeToTime(ft syscall.Filetime) time.Time {
	return time.Unix(0, ft.Nanoseconds())
}

// 获取进程创建时间
func getProcessTimes(pid uint32) (time.Time, error) {
	h, _, _ := procOpenProcess.Call(PROCESS_QUERY_LIMITED_INFORMATION, 0, uintptr(pid))
	if h == 0 {
		return time.Time{}, fmt.Errorf("OpenProcess failed")
	}
	defer syscall.CloseHandle(syscall.Handle(h))

	var creation, exit, kernel, user syscall.Filetime
	ret, _, _ := procGetProcessTimes.Call(h,
		uintptr(unsafe.Pointer(&creation)),
		uintptr(unsafe.Pointer(&exit)),
		uintptr(unsafe.Pointer(&kernel)),
		uintptr(unsafe.Pointer(&user)),
	)
	if ret == 0 {
		return time.Time{}, fmt.Errorf("GetProcessTimes failed")
	}
	return filetimeToTime(creation), nil
}

// 获取进程运行用户名
func getProcessUser(pid uint32) string {
	h, _, _ := procOpenProcess.Call(PROCESS_QUERY_INFORMATION, 0, uintptr(pid))
	if h == 0 {
		return "N/A"
	}
	defer syscall.CloseHandle(syscall.Handle(h))

	var token syscall.Handle
	ret, _, _ := procOpenProcessToken.Call(h, TOKEN_QUERY, uintptr(unsafe.Pointer(&token)))
	if ret == 0 {
		return "N/A"
	}
	defer syscall.CloseHandle(token)

	var tokenUser structs.TOKEN_USER
	var retLen uint32
	ret, _, _ = procGetTokenInformation.Call(uintptr(token),
		uintptr(1), // TokenUser = 1
		uintptr(unsafe.Pointer(&tokenUser)),
		uintptr(unsafe.Sizeof(tokenUser)),
		uintptr(unsafe.Pointer(&retLen)))
	if ret == 0 {
		// 内存不足可以用 retLen 再分配
		buf := make([]byte, retLen)
		ret, _, _ = procGetTokenInformation.Call(uintptr(token),
			uintptr(1),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(retLen),
			uintptr(unsafe.Pointer(&retLen)))
		if ret == 0 {
			return "N/A"
		}
		tokenUser = *(*structs.TOKEN_USER)(unsafe.Pointer(&buf[0]))
	}

	var name [256]uint16
	var cchName uint32 = 256
	var domain [256]uint16
	var cchDomain uint32 = 256
	var sidType uint32
	ret, _, _ = procLookupAccountSidW.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(tokenUser.User.Sid)),
		uintptr(unsafe.Pointer(&name[0])),
		uintptr(unsafe.Pointer(&cchName)),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&cchDomain)),
		uintptr(unsafe.Pointer(&sidType)),
	)
	if ret == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%s\\%s", syscall.UTF16ToString(domain[:]), syscall.UTF16ToString(name[:]))
}

// 获取父进程名称
func getParentName(parentPID uint32, allProcs map[uint32]string) string {
	if name, ok := allProcs[parentPID]; ok {
		return name
	}
	return "N/A"
}
func ProcessEnum() {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(TH32CS_SNAPPROCESS), uintptr(0))
	if handle < 0 {
		fmt.Println("CreateToolhelp32Snapshot failed")
		return
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	var entry structs.PROCESSENTRY32
	entry.DwSize = uint32(unsafe.Sizeof(entry))

	ret, _, _ := procProcess32FirstW.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		fmt.Println("Process32FirstW failed")
		return
	}

	// 先把所有进程 PID -> Name 存起来，用于获取父进程名称
	allProcs := make(map[uint32]string)
	for {
		allProcs[entry.Th32ProcessID] = syscall.UTF16ToString(entry.SzExeFile[:])
		ret, _, _ = procProcess32NextW.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	// 打印表头
	fmt.Printf("%-8s %-8s %-30s %-30s %-25s %-20s\n", "PID", "PPID", "ParentName", "Name", "User", "Created")

	// 重置快照
	ret, _, _ = procProcess32FirstW.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		fmt.Println("Process32FirstW failed")
		return
	}

	for {
		creationTime, _ := getProcessTimes(entry.Th32ProcessID)
		user := getProcessUser(entry.Th32ProcessID)
		parentName := getParentName(entry.Th32ParentProcessID, allProcs)

		fmt.Printf("%-8d %-8d %-30s %-30s %-25s %-20s\n",
			entry.Th32ProcessID,
			entry.Th32ParentProcessID,
			parentName,
			syscall.UTF16ToString(entry.SzExeFile[:]),
			user,
			creationTime.Format("2006-01-02 15:04:05"),
		)

		ret, _, _ = procProcess32NextW.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
}
