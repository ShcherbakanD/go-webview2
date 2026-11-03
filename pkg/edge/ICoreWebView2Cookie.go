package edge

import (
	"unsafe"

	"github.com/jchv/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

// _ICoreWebView2CookieVtbl mirrors the COM vtable of ICoreWebView2Cookie
// (IID AD26D6BE-1486-43E6-BF87-A2034006CA21) as declared in WebView2.h.
// The slot order must match the IDL exactly — reordering breaks ABI.
type _ICoreWebView2CookieVtbl struct {
	_IUnknownVtbl
	GetName       ComProc
	GetValue      ComProc
	PutValue      ComProc
	GetDomain     ComProc
	GetPath       ComProc
	GetExpires    ComProc
	PutExpires    ComProc // slot only — see the note on PutExpires below
	GetIsHttpOnly ComProc
	PutIsHttpOnly ComProc
	GetSameSite   ComProc
	PutSameSite   ComProc
	GetIsSecure   ComProc
	PutIsSecure   ComProc
	GetIsSession  ComProc
}

// ICoreWebView2Cookie is one cookie as seen by the WebView2 cookie manager.
//
// put_Expires(double) is deliberately not wrapped: it takes a double BY VALUE,
// which on x64 travels in an XMM register. ComProc.Call is syscall.SyscallN and
// only fills integer registers, so any value passed through it would arrive as
// garbage. Cookies created through CreateCookie have no expiry (session
// cookies), which is the only case this package needs.
type ICoreWebView2Cookie struct {
	vtbl *_ICoreWebView2CookieVtbl
}

func (i *ICoreWebView2Cookie) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

// Release drops the reference obtained from CookieManager.CreateCookie or
// CookieList.GetValueAtIndex.
func (i *ICoreWebView2Cookie) Release() error {
	r1, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return hr(r1)
}

func (i *ICoreWebView2Cookie) getString(p ComProc) (string, error) {
	var raw *uint16
	r1, _, _ := p.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&raw)),
	)
	if err := hr(r1); err != nil {
		return "", err
	}
	if raw == nil {
		return "", nil
	}
	out := w32.Utf16PtrToString(raw)
	windows.CoTaskMemFree(unsafe.Pointer(raw))
	return out, nil
}

func (i *ICoreWebView2Cookie) getBool(p ComProc) (bool, error) {
	var v int32
	r1, _, _ := p.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err := hr(r1); err != nil {
		return false, err
	}
	return v != 0, nil
}

func (i *ICoreWebView2Cookie) putBool(p ComProc, v bool) error {
	r1, _, _ := p.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(v)),
	)
	return hr(r1)
}

func (i *ICoreWebView2Cookie) GetName() (string, error)   { return i.getString(i.vtbl.GetName) }
func (i *ICoreWebView2Cookie) GetValue() (string, error)  { return i.getString(i.vtbl.GetValue) }
func (i *ICoreWebView2Cookie) GetDomain() (string, error) { return i.getString(i.vtbl.GetDomain) }
func (i *ICoreWebView2Cookie) GetPath() (string, error)   { return i.getString(i.vtbl.GetPath) }

func (i *ICoreWebView2Cookie) PutValue(value string) error {
	v, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	r1, _, _ := i.vtbl.PutValue.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(v)),
	)
	return hr(r1)
}

// GetExpires returns the expiry as seconds since the Unix epoch, or -1 for a
// session cookie.
func (i *ICoreWebView2Cookie) GetExpires() (float64, error) {
	var v float64
	r1, _, _ := i.vtbl.GetExpires.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err := hr(r1); err != nil {
		return 0, err
	}
	return v, nil
}

func (i *ICoreWebView2Cookie) GetIsHttpOnly() (bool, error) { return i.getBool(i.vtbl.GetIsHttpOnly) }
func (i *ICoreWebView2Cookie) PutIsHttpOnly(v bool) error   { return i.putBool(i.vtbl.PutIsHttpOnly, v) }
func (i *ICoreWebView2Cookie) GetIsSecure() (bool, error)   { return i.getBool(i.vtbl.GetIsSecure) }
func (i *ICoreWebView2Cookie) PutIsSecure(v bool) error     { return i.putBool(i.vtbl.PutIsSecure, v) }
func (i *ICoreWebView2Cookie) GetIsSession() (bool, error)  { return i.getBool(i.vtbl.GetIsSession) }

func (i *ICoreWebView2Cookie) GetSameSite() (COREWEBVIEW2_COOKIE_SAME_SITE_KIND, error) {
	var v uint32
	r1, _, _ := i.vtbl.GetSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&v)),
	)
	if err := hr(r1); err != nil {
		return 0, err
	}
	return COREWEBVIEW2_COOKIE_SAME_SITE_KIND(v), nil
}

func (i *ICoreWebView2Cookie) PutSameSite(v COREWEBVIEW2_COOKIE_SAME_SITE_KIND) error {
	r1, _, _ := i.vtbl.PutSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(v),
	)
	return hr(r1)
}
