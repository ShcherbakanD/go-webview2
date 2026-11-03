package edge

import "fmt"

// hr converts a COM method's return value into an error.
//
// ComProc.Call is syscall.SyscallN, whose third result is GetLastError(). COM
// methods report their status as an HRESULT in the FIRST result and do not touch
// LastError, so checking the third result — as most wrappers in this package
// still do — reads a stale thread-local error code. That both misses real
// failures and invents failures out of nothing, which is how a successful call
// ends up reported as broken.
//
// An HRESULT is a failure exactly when its sign bit is set.
func hr(r1 uintptr) error {
	if h := int32(uint32(r1)); h < 0 {
		return fmt.Errorf("HRESULT 0x%08x", uint32(r1))
	}
	return nil
}
