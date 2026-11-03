package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceRequestVtbl struct {
	_IUnknownVtbl
	GetUri     ComProc
	PutUri     ComProc
	GetMethod  ComProc
	PutMethod  ComProc
	GetContent ComProc
	PutContent ComProc
	GetHeaders ComProc
}

type ICoreWebView2WebResourceRequest struct {
	vtbl *_ICoreWebView2WebResourceRequestVtbl
}

func (i *ICoreWebView2WebResourceRequest) AddRef() uintptr {
	// The `this` pointer is not optional: without it AddRef runs against
	// whatever happened to be in the register.
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

// Release drops the reference obtained from
// ICoreWebView2WebResourceRequestedEventArgs.GetRequest. Failing to call it
// leaks one COM object per intercepted request.
func (i *ICoreWebView2WebResourceRequest) Release() error {
	r1, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return hr(r1)
}

// PutUri rewrites the request URI before WebView2 sends it. Useful for
// inserting cache-busting query parameters from inside a
// WebResourceRequested handler — e.g. when a third-party backend caches
// responses by URL and would otherwise return a stale single-use token.
func (i *ICoreWebView2WebResourceRequest) PutUri(uri string) error {
	u, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}
	r1, _, _ := i.vtbl.PutUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u)),
	)
	return hr(r1)
}

func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16
	r1, _, _ := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if err := hr(r1); err != nil {
		return "", err
	} // Get result and cleanup
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

// GetHeaders returns the mutable request headers. The host can call
// SetHeader / RemoveHeader on the returned object inside a
// WebResourceRequested handler to rewrite the outgoing request — useful
// for spoofing UA Client Hints (sec-ch-ua*) that WebView2 otherwise
// always sends with "Microsoft Edge WebView2" branding.
func (i *ICoreWebView2WebResourceRequest) GetHeaders() (*ICoreWebView2HttpRequestHeaders, error) {
	var headers *ICoreWebView2HttpRequestHeaders
	_, _, err := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return headers, nil
}
