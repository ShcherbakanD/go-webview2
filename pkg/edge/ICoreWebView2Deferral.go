package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2DeferralVtbl struct {
	_IUnknownVtbl
	Complete ComProc
}

type ICoreWebView2Deferral struct {
	vtbl *_ICoreWebView2DeferralVtbl
}

func (i *ICoreWebView2Deferral) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

func (i *ICoreWebView2Deferral) Release() error {
	_, _, err := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2Deferral) Complete() error {
	_, _, err := i.vtbl.Complete.Call(uintptr(unsafe.Pointer(i)))
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}
