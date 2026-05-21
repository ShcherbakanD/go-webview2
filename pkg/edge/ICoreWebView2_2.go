package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2_2Vtbl struct {
	iCoreWebView2Vtbl
	AddWebResourceResponseReceived    ComProc
	RemoveWebResourceResponseReceived ComProc
	NavigateWithWebResourceRequest    ComProc
	AddDomContentLoaded               ComProc
	RemoveDomContentLoaded            ComProc
	GetCookieManager                  ComProc
	GetEnvironment                    ComProc
}

type ICoreWebView2_2 struct {
	vtbl *iCoreWebView2_2Vtbl
}

func (i *ICoreWebView2_2) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call()
	return r
}

// GetICoreWebView2_2 returns the ICoreWebView2_2 view of this WebView via
// QueryInterface, or nil if the underlying WebView2 runtime is too old to
// support it. Required for events introduced in WebView2 1.0.705.50:
// WebResourceResponseReceived, DOMContentLoaded, CookieManager, etc.
func (i *ICoreWebView2) GetICoreWebView2_2() *ICoreWebView2_2 {
	var result *ICoreWebView2_2
	iid := NewGUID("{9E8F0CF8-E670-4B5E-B2BC-73E061E3184C}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(&result)),
	)
	return result
}

// AddWebResourceResponseReceived subscribes a host handler to receive a
// post-fact notification for every response the WebView gets for a web
// resource. Useful for status-code/header diagnostics — the handler is
// observational and cannot alter the response.
func (i *ICoreWebView2_2) AddWebResourceResponseReceived(handler *iCoreWebView2WebResourceResponseReceivedEventHandler, token *_EventRegistrationToken) error {
	_, _, err := i.vtbl.AddWebResourceResponseReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(token)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
