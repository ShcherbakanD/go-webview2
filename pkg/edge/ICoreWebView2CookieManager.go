package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// _ICoreWebView2CookieManagerVtbl mirrors the COM vtable of
// ICoreWebView2CookieManager (IID 177CD9E7-B6F5-451A-94A0-5D7A3A4C4141).
// Slot order must match the IDL exactly — reordering breaks ABI.
type _ICoreWebView2CookieManagerVtbl struct {
	_IUnknownVtbl
	CreateCookie                   ComProc
	CopyCookie                     ComProc
	GetCookies                     ComProc
	AddOrUpdateCookie              ComProc
	DeleteCookie                   ComProc
	DeleteCookies                  ComProc
	DeleteCookiesWithDomainAndPath ComProc
	DeleteAllCookies               ComProc
}

// ICoreWebView2CookieManager is the profile-wide cookie store of a WebView2.
// Obtained via ICoreWebView2_2.GetCookieManager (runtime 1.0.705.50+).
type ICoreWebView2CookieManager struct {
	vtbl *_ICoreWebView2CookieManagerVtbl
}

func (i *ICoreWebView2CookieManager) AddRef() uintptr {
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return r
}

func (i *ICoreWebView2CookieManager) Release() error {
	r1, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return hr(r1)
}

// CreateCookie builds a cookie object that is not yet part of the store. The
// result has no expiry (a session cookie) until put_Expires is called, which
// this package cannot do; see ICoreWebView2Cookie. Hand it to AddOrUpdateCookie
// and Release it afterwards.
func (i *ICoreWebView2CookieManager) CreateCookie(name, value, domain, path string) (*ICoreWebView2Cookie, error) {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	v, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return nil, err
	}
	d, err := windows.UTF16PtrFromString(domain)
	if err != nil {
		return nil, err
	}
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var c *ICoreWebView2Cookie
	r1, _, _ := i.vtbl.CreateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(d)),
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&c)),
	)
	if err := hr(r1); err != nil {
		return nil, err
	}
	return c, nil
}

// GetCookies asynchronously fetches the cookies that would be sent with a
// request to uri. An empty uri returns every cookie in the profile. The
// handler is invoked on the UI thread once the list is ready; the caller must
// keep the handler reachable until then.
func (i *ICoreWebView2CookieManager) GetCookies(uri string, handler *iCoreWebView2GetCookiesCompletedHandler) error {
	var u *uint16
	if uri != "" {
		var err error
		u, err = windows.UTF16PtrFromString(uri)
		if err != nil {
			return err
		}
	}
	r1, _, _ := i.vtbl.GetCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(u)),
		uintptr(unsafe.Pointer(handler)),
	)
	return hr(r1)
}

// AddOrUpdateCookie inserts the cookie into the store, replacing any cookie
// with the same name, domain and path.
func (i *ICoreWebView2CookieManager) AddOrUpdateCookie(c *ICoreWebView2Cookie) error {
	r1, _, _ := i.vtbl.AddOrUpdateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(c)),
	)
	return hr(r1)
}

// DeleteCookies removes every cookie with the given name that would be sent
// to uri.
func (i *ICoreWebView2CookieManager) DeleteCookies(name, uri string) error {
	n, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	u, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}
	r1, _, _ := i.vtbl.DeleteCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(n)),
		uintptr(unsafe.Pointer(u)),
	)
	return hr(r1)
}

// DeleteAllCookies wipes the cookie store of the profile.
func (i *ICoreWebView2CookieManager) DeleteAllCookies() error {
	r1, _, _ := i.vtbl.DeleteAllCookies.Call(uintptr(unsafe.Pointer(i)))
	return hr(r1)
}
