package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceResponseViewVtbl struct {
	_IUnknownVtbl
	GetHeaders      ComProc
	GetStatusCode   ComProc
	GetReasonPhrase ComProc
	GetContent      ComProc
}

// ICoreWebView2WebResourceResponseView is the response side of the
// WebResourceResponseReceived event. Observational only — the response has
// already been delivered to the renderer by the time this is raised.
type ICoreWebView2WebResourceResponseView struct {
	vtbl *_ICoreWebView2WebResourceResponseViewVtbl
}

// GetStatusCode returns the HTTP status code WebView2 actually received for
// the response.
func (i *ICoreWebView2WebResourceResponseView) GetStatusCode() (int, error) {
	var statusCode int32
	_, _, err := i.vtbl.GetStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&statusCode)),
	)
	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return int(statusCode), nil
}

// GetReasonPhrase returns the HTTP reason phrase (e.g. "OK", "Not Found").
func (i *ICoreWebView2WebResourceResponseView) GetReasonPhrase() (string, error) {
	var raw *uint16
	_, _, err := i.vtbl.GetReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&raw)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	s := windows.UTF16PtrToString(raw)
	windows.CoTaskMemFree(unsafe.Pointer(raw))
	return s, nil
}

// GetHeaders returns the response's header collection. Callers can use
// GetHeader to fetch individual entries (e.g. X-Request-Id) for diagnostics.
func (i *ICoreWebView2WebResourceResponseView) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {
	var headers *ICoreWebView2HttpResponseHeaders
	_, _, err := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return headers, nil
}
