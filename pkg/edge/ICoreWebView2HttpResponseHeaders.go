package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2HttpResponseHeadersVtbl struct {
	_IUnknownVtbl
	AppendHeader ComProc
	Contains     ComProc
	GetHeader    ComProc
	GetHeaders   ComProc
	GetIterator  ComProc
}

// ICoreWebView2HttpResponseHeaders is the read-side wrapper around an HTTP
// response's header collection. Returned by
// ICoreWebView2WebResourceResponseView.GetHeaders.
type ICoreWebView2HttpResponseHeaders struct {
	vtbl *_ICoreWebView2HttpResponseHeadersVtbl
}

// GetHeader returns the value of the named header, or empty string if absent.
// Useful for fishing out diagnostic headers (e.g. X-Request-Id,
// X-Internal-Error-Code) when investigating server-side rejections.
func (i *ICoreWebView2HttpResponseHeaders) GetHeader(name string) (string, error) {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return "", err
	}
	var raw *uint16
	_, _, err = i.vtbl.GetHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(&raw)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	if raw == nil {
		return "", nil
	}
	s := windows.UTF16PtrToString(raw)
	windows.CoTaskMemFree(unsafe.Pointer(raw))
	return s, nil
}
