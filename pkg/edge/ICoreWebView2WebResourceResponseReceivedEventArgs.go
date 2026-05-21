package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceResponseReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetRequest  ComProc
	GetResponse ComProc
}

// ICoreWebView2WebResourceResponseReceivedEventArgs holds the request/response
// pair for a WebResourceResponseReceived event.
type ICoreWebView2WebResourceResponseReceivedEventArgs struct {
	vtbl *_ICoreWebView2WebResourceResponseReceivedEventArgsVtbl
}

// GetRequest returns the request object as committed. Modifications have no
// effect — the request has already been sent.
func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {
	var request *ICoreWebView2WebResourceRequest
	_, _, err := i.vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return request, nil
}

// GetResponse returns a view of the response WebView2 actually received,
// including the real HTTP status code.
func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetResponse() (*ICoreWebView2WebResourceResponseView, error) {
	var response *ICoreWebView2WebResourceResponseView
	_, _, err := i.vtbl.GetResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&response)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return response, nil
}
