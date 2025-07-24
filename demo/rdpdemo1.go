package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	wevtapi       = syscall.NewLazyDLL("wevtapi.dll")
	procEvtQuery  = wevtapi.NewProc("EvtQuery")
	procEvtNext   = wevtapi.NewProc("EvtNext")
	procEvtRender = wevtapi.NewProc("EvtRender")
	procEvtClose  = wevtapi.NewProc("EvtClose")
)

const (
	EvtQueryChannelPath = 0x00000001
	EvtRenderEventXml   = 1
	BufferSize          = 1 << 16 // 64KB
	batchSize           = 6       // 每批读取事件数
	ERROR_NO_MORE_ITEMS = 259     // Windows 错误码：没有更多项
)

// 将字符串转为 *uint16（UTF-16 指针）
func utf16Ptr(s string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(s)
	return ptr
}

func main() {
	channels := []string{
		"Microsoft-Windows-TerminalServices-RemoteConnectionManager/Operational",
		"Microsoft-Windows-TerminalServices-LocalSessionManager/Operational",
	}
	query := "*[System[(EventID=21 or EventID=25 or EventID=1149)]]"

	for _, channel := range channels {
		fmt.Printf("\n==== 查询通道: %s ====\n", channel)

		// 1. 打开查询句柄
		handle, _, err := procEvtQuery.Call(
			0,
			uintptr(unsafe.Pointer(utf16Ptr(channel))),
			uintptr(unsafe.Pointer(utf16Ptr(query))),
			uintptr(EvtQueryChannelPath),
		)
		if handle == 0 {
			fmt.Printf("[!] EvtQuery 失败: %v\n", err)
			continue
		}
		defer procEvtClose.Call(handle)

		total := 0

		// 2. 循环读取所有事件
		for {
			var events [batchSize]syscall.Handle
			var returned uint32

			ret, _, _ := procEvtNext.Call(
				handle,
				uintptr(batchSize),
				uintptr(unsafe.Pointer(&events[0])),
				0,
				0,
				uintptr(unsafe.Pointer(&returned)),
			)

			if ret == 0 {
				lastErr := syscall.GetLastError()
				if errno, ok := lastErr.(syscall.Errno); ok && errno == ERROR_NO_MORE_ITEMS {
					fmt.Println("[*] 所有事件已读取完毕")
					break
				}
				fmt.Printf("[!] EvtNext 失败: %v\n", lastErr)
				break
			}

			fmt.Printf("[+] 本次读取 %d 条事件\n", returned)
			total += int(returned)

			// 3. 遍历处理每条事件
			for i := uint32(0); i < returned; i++ {
				// 渲染事件 XML
				buf := make([]uint16, BufferSize)
				var used, propCount uint32

				ret, _, err := procEvtRender.Call(
					0,
					uintptr(events[i]),
					EvtRenderEventXml,
					uintptr(len(buf)*2),
					uintptr(unsafe.Pointer(&buf[0])),
					uintptr(unsafe.Pointer(&used)),
					uintptr(unsafe.Pointer(&propCount)),
				)
				if ret == 0 {
					fmt.Printf("[!] 第 %d 条事件渲染失败: %v\n", total-(int(returned)-int(i)), err)
					procEvtClose.Call(uintptr(events[i]))
					continue
				}

				xml := syscall.UTF16ToString(buf[:used/2])
				fmt.Printf("\n--- 第 %d 条事件 ---\n%s\n", total-(int(returned)-int(i))+1, xml)

				// 释放事件句柄
				procEvtClose.Call(uintptr(events[i]))
			}
		}
	}
}
