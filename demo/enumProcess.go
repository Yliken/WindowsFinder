package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	TH32CS_SNAPPROCESS                = 0x00000002
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	TOKEN_QUERY                       = 0x0008
)

type PROCESSENTRY32 struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [260]uint16
}

type SID_AND_ATTRIBUTES struct {
	Sid        *syscall.SID
	Attributes uint32
}

type TOKEN_USER struct {
	User SID_AND_ATTRIBUTES
}

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

	var tokenUser TOKEN_USER
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
		tokenUser = *(*TOKEN_USER)(unsafe.Pointer(&buf[0]))
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

func main() {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(TH32CS_SNAPPROCESS), uintptr(0))
	if handle < 0 {
		fmt.Println("CreateToolhelp32Snapshot failed")
		return
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	var entry PROCESSENTRY32
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
