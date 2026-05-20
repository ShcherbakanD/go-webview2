package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2NavigationCompletedEventArgsVtbl struct {
	_IUnknownVtbl
	GetIsSuccess      ComProc
	GetWebErrorStatus ComProc
	GetNavigationId   ComProc
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	vtbl *_ICoreWebView2NavigationCompletedEventArgsVtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

// IsSuccess reports whether the navigation completed successfully (HTTP 2xx/3xx
// without network errors). false on 4xx/5xx HTTP responses too, in which case
// WebErrorStatus carries the reason (HttpInvalidServerResponse for 5xx).
func (i *ICoreWebView2NavigationCompletedEventArgs) IsSuccess() (bool, error) {
	var v int32
	_, _, err := i.vtbl.GetIsSuccess.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return v != 0, nil
}

// WebErrorStatus returns the error category when IsSuccess=false. See the
// COREWEBVIEW2_WEB_ERROR_STATUS enum in WebView2 IDL for known values.
func (i *ICoreWebView2NavigationCompletedEventArgs) WebErrorStatus() (uint32, error) {
	var v uint32
	_, _, err := i.vtbl.GetWebErrorStatus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return v, nil
}
