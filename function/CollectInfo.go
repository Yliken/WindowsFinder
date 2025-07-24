package function

import (
	"encoding/xml"
	"fmt"
	"syscall"
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
//搜集RDP 相关的日志信息 时间 1149 事件21 事件 25
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
