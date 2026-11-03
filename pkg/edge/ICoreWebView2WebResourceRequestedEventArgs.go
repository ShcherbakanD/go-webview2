package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetRequest         ComProc
	GetResponse        ComProc
	PutResponse        ComProc
	GetDeferral        ComProc
	GetResourceContext ComProc
}

type ICoreWebView2WebResourceRequestedEventArgs struct {
	vtbl *_ICoreWebView2WebResourceRequestedEventArgsVtbl
}

// AddRef extends the lifetime of the event args past the return of the
// WebResourceRequested handler. WebView2 passes them as a borrowed reference and
// releases its own once Invoke returns, so a handler that hands the args to a
// goroutine must AddRef here and Release when it is done with them.
func (i *ICoreWebView2WebResourceRequestedEventArgs) AddRef() uintptr {
	// The `this` pointer is not optional: without it AddRef runs against
	// whatever happened to be in the register.
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

// Release drops one reference taken with AddRef.
func (i *ICoreWebView2WebResourceRequestedEventArgs) Release() error {
	r1, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return hr(r1)
}

func (i *ICoreWebView2WebResourceRequestedEventArgs) PutResponse(response *ICoreWebView2WebResourceResponse) error {
	var err error

	_, _, err = i.vtbl.PutResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(response)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// GetRequest returns the request being made. The caller owns the returned
// reference and must Release it.
func (i *ICoreWebView2WebResourceRequestedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {
	var request *ICoreWebView2WebResourceRequest
	r1, _, _ := i.vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if err := hr(r1); err != nil {
		return nil, err
	}
	return request, nil
}

// GetDeferral returns a deferral object so the WebResourceRequested handler can
// finish asynchronously. The handler should hold the deferral while it builds a
// response (e.g. fetching from a remote server in a goroutine), call PutResponse
// on the EventArgs once the response is ready, and then call Complete() on the
// deferral to let Chromium resume the navigation.
func (i *ICoreWebView2WebResourceRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {
	var err error
	var deferral *ICoreWebView2Deferral
	_, _, err = i.vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return deferral, nil
}
