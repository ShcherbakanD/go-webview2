package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// _ICoreWebView2HttpRequestHeadersVtbl mirrors the COM vtable for
// ICoreWebView2HttpRequestHeaders. Slot order must match the IDL exactly.
type _ICoreWebView2HttpRequestHeadersVtbl struct {
	_IUnknownVtbl
	GetHeader     ComProc
	GetHeaders    ComProc
	Contains      ComProc
	SetHeader     ComProc
	RemoveHeader  ComProc
	GetIterator   ComProc
}

type ICoreWebView2HttpRequestHeaders struct {
	vtbl *_ICoreWebView2HttpRequestHeadersVtbl
}

// SetHeader replaces (or adds) a single header. Used by hosts to rewrite
// outgoing requests inside a WebResourceRequested handler.
func (h *ICoreWebView2HttpRequestHeaders) SetHeader(name, value string) error {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	v, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	_, _, callErr := h.vtbl.SetHeader.Call(
		uintptr(unsafe.Pointer(h)),
		uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(v)),
	)
	if callErr != windows.ERROR_SUCCESS {
		return callErr
	}
	return nil
}

// RemoveHeader drops a single header from the request.
func (h *ICoreWebView2HttpRequestHeaders) RemoveHeader(name string) error {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	_, _, callErr := h.vtbl.RemoveHeader.Call(
		uintptr(unsafe.Pointer(h)),
		uintptr(unsafe.Pointer(n)),
	)
	if callErr != windows.ERROR_SUCCESS {
		return callErr
	}
	return nil
}

// Contains reports whether a header with the given name is present.
func (h *ICoreWebView2HttpRequestHeaders) Contains(name string) (bool, error) {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return false, err
	}
	var contains int32
	_, _, callErr := h.vtbl.Contains.Call(
		uintptr(unsafe.Pointer(h)),
		uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(&contains)),
	)
	if callErr != windows.ERROR_SUCCESS {
		return false, callErr
	}
	return contains != 0, nil
}

func (h *ICoreWebView2HttpRequestHeaders) Release() error {
	_, _, err := h.vtbl.Release.Call(uintptr(unsafe.Pointer(h)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
