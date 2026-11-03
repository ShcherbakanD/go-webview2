package edge

import "unsafe"

// _ICoreWebView2CookieListVtbl mirrors the COM vtable of ICoreWebView2CookieList
// (IID F7F6F714-5D2A-43C6-9503-346ECE02D186). Slot order must match the IDL.
type _ICoreWebView2CookieListVtbl struct {
	_IUnknownVtbl
	GetCount        ComProc
	GetValueAtIndex ComProc
}

// ICoreWebView2CookieList is the result of CookieManager.GetCookies. It is an
// in-parameter of the completion handler and is owned by WebView2 for the
// duration of the callback; do not Release it.
type ICoreWebView2CookieList struct {
	vtbl *_ICoreWebView2CookieListVtbl
}

func (i *ICoreWebView2CookieList) GetCount() (uint32, error) {
	var n uint32
	r1, _, _ := i.vtbl.GetCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&n)),
	)
	if err := hr(r1); err != nil {
		return 0, err
	}
	return n, nil
}

// GetValueAtIndex returns an owned reference; the caller must Release it.
func (i *ICoreWebView2CookieList) GetValueAtIndex(index uint32) (*ICoreWebView2Cookie, error) {
	var c *ICoreWebView2Cookie
	r1, _, _ := i.vtbl.GetValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(index),
		uintptr(unsafe.Pointer(&c)),
	)
	if err := hr(r1); err != nil {
		return nil, err
	}
	return c, nil
}
