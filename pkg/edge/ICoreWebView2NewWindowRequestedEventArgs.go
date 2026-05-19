package edge

import (
	"unsafe"

	"github.com/jchv/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

// _ICoreWebView2NewWindowRequestedEventArgsVtbl mirrors the COM vtable layout
// of ICoreWebView2NewWindowRequestedEventArgs as declared in the WebView2 IDL.
// The slot order below must match the IDL exactly — reordering breaks ABI.
type _ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	PutNewWindow       ComProc
	GetNewWindow       ComProc
	PutHandled         ComProc
	GetHandled         ComProc
	GetIsUserInitiated ComProc
	GetDeferral        ComProc
	GetWindowFeatures  ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	vtbl *_ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) Release() error {
	_, _, err := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// GetUri returns the URI of the requested popup.
func (i *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
	var uri *uint16
	_, _, err := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&uri)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	if uri == nil {
		return "", nil
	}
	out := w32.Utf16PtrToString(uri)
	windows.CoTaskMemFree(unsafe.Pointer(uri))
	return out, nil
}

// PutHandled marks the event as handled, suppressing WebView2's default popup behavior.
func (i *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {
	v := uintptr(0)
	if handled {
		v = 1
	}
	_, _, err := i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		v,
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetHandled() (bool, error) {
	var handled int32
	_, _, err := i.vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&handled)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return handled != 0, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	var v int32
	_, _, err := i.vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	return v != 0, nil
}

// PutNewWindow hands a host-created ICoreWebView2 back to WebView2 to use as
// the new popup window. After this call (and Deferral.Complete if used),
// WebView2 will treat newWindow as the popup target — JavaScript's window.open
// returns a live Window reference backed by newWindow, and subsequent
// popup.location / popup.postMessage / popup.close calls drive newWindow
// normally instead of failing with TypeError.
func (i *ICoreWebView2NewWindowRequestedEventArgs) PutNewWindow(newWindow *ICoreWebView2) error {
	_, _, err := i.vtbl.PutNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(newWindow)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

// GetDeferral lets the host complete the event asynchronously. Caller must
// call Deferral.Complete() and Release() when done.
func (i *ICoreWebView2NewWindowRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {
	var deferral *ICoreWebView2Deferral
	_, _, err := i.vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return deferral, nil
}
