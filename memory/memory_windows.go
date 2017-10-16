// +build windows

package memory

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
)

// Get memory statistics
func Get() (*MemoryStats, error) {
	var memoryStatus memoryStatusEx
	memoryStatus.Length = uint32(unsafe.Sizeof(memoryStatus))

	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memoryStatus)))
	if ret == 0 {
		return nil, fmt.Errorf("failed in GlobalMemoryStatusEx: %s", err)
	}

	var memory MemoryStats
	memory.Free = memoryStatus.AvailPhys
	memory.Total = memoryStatus.TotalPhys
	memory.Used = memory.Total - memory.Free
	memory.PageFileTotal = memoryStatus.TotalPageFile
	memory.PageFileFree = memoryStatus.AvailPageFile
	memory.VirtualTotal = memoryStatus.TotalVirtual
	memory.VirtualFree = memoryStatus.AvailVirtual

	return &memory, nil
}

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// MemoryStats represents memory statistics for Windows
type MemoryStats struct {
	Total, Used, Free, PageFileTotal, PageFileFree, VirtualTotal, VirtualFree uint64
}
