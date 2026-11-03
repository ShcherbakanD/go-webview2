//go:build windows
// +build windows

package edge

import (
	"errors"
	"fmt"
)

// Cookie is a Go-side snapshot of an ICoreWebView2Cookie.
type Cookie struct {
	Name   string
	Value  string
	Domain string
	Path   string
	// Expires is seconds since the Unix epoch, or -1 for a session cookie.
	Expires    float64
	IsHttpOnly bool
	IsSecure   bool
	IsSession  bool
	SameSite   COREWEBVIEW2_COOKIE_SAME_SITE_KIND
}

// ErrCookieManagerUnavailable is returned when the WebView2 runtime predates
// ICoreWebView2_2 (1.0.705.50) or the controller is not initialized yet.
var ErrCookieManagerUnavailable = errors.New("webview2: cookie manager unavailable")

// getCookiesCompletion is the per-call impl behind GetCookies. It unregisters
// itself from Chromium.pendingGetCookies before handing the result to cb, so
// the handler stays reachable exactly as long as WebView2 may still Invoke it.
type getCookiesCompletion struct {
	e    *Chromium
	self *iCoreWebView2GetCookiesCompletedHandler
	cb   func(cookies []Cookie, err error)
}

func (c *getCookiesCompletion) QueryInterface(_, _ uintptr) uintptr { return 0 }
func (c *getCookiesCompletion) AddRef() uintptr                     { return 1 }
func (c *getCookiesCompletion) Release() uintptr                    { return 1 }

func (c *getCookiesCompletion) GetCookiesCompleted(errorCode uintptr, list *ICoreWebView2CookieList) uintptr {
	c.e.pendingMu.Lock()
	delete(c.e.pendingGetCookies, c.self)
	c.e.pendingMu.Unlock()

	if err := hr(errorCode); err != nil {
		c.cb(nil, err)
		return 0
	}
	if list == nil {
		c.cb(nil, nil)
		return 0
	}
	cookies, err := readCookieList(list)
	c.cb(cookies, err)
	return 0
}

// readCookieList copies every entry of list into Go values, releasing each
// COM cookie as it goes. The list itself belongs to WebView2 for the duration
// of the callback and is not released here.
func readCookieList(list *ICoreWebView2CookieList) ([]Cookie, error) {
	n, err := list.GetCount()
	if err != nil {
		return nil, fmt.Errorf("cookie list count: %w", err)
	}
	out := make([]Cookie, 0, n)
	for idx := uint32(0); idx < n; idx++ {
		raw, err := list.GetValueAtIndex(idx)
		if err != nil {
			return out, fmt.Errorf("cookie list index %d: %w", idx, err)
		}
		if raw == nil {
			continue
		}
		c, err := readCookie(raw)
		_ = raw.Release()
		if err != nil {
			return out, fmt.Errorf("cookie list index %d: %w", idx, err)
		}
		out = append(out, c)
	}
	return out, nil
}

func readCookie(raw *ICoreWebView2Cookie) (Cookie, error) {
	var c Cookie
	var err error
	if c.Name, err = raw.GetName(); err != nil {
		return c, err
	}
	if c.Value, err = raw.GetValue(); err != nil {
		return c, err
	}
	if c.Domain, err = raw.GetDomain(); err != nil {
		return c, err
	}
	if c.Path, err = raw.GetPath(); err != nil {
		return c, err
	}
	if c.Expires, err = raw.GetExpires(); err != nil {
		return c, err
	}
	if c.IsHttpOnly, err = raw.GetIsHttpOnly(); err != nil {
		return c, err
	}
	if c.IsSecure, err = raw.GetIsSecure(); err != nil {
		return c, err
	}
	if c.IsSession, err = raw.GetIsSession(); err != nil {
		return c, err
	}
	if c.SameSite, err = raw.GetSameSite(); err != nil {
		return c, err
	}
	return c, nil
}

// GetCookies asynchronously fetches the cookies that would be sent with a
// request to uri (empty uri = whole profile). cb runs on the UI thread once
// WebView2 has the list; it is never called if this function returns an error.
func (e *Chromium) GetCookies(uri string, cb func(cookies []Cookie, err error)) error {
	if e.cookieManager == nil {
		return ErrCookieManagerUnavailable
	}
	if cb == nil {
		return errors.New("webview2: GetCookies: nil callback")
	}
	impl := &getCookiesCompletion{e: e, cb: cb}
	h := newICoreWebView2GetCookiesCompletedHandler(impl)
	impl.self = h

	e.pendingMu.Lock()
	if e.pendingGetCookies == nil {
		e.pendingGetCookies = make(map[*iCoreWebView2GetCookiesCompletedHandler]struct{})
	}
	e.pendingGetCookies[h] = struct{}{}
	e.pendingMu.Unlock()

	if err := e.cookieManager.GetCookies(uri, h); err != nil {
		e.pendingMu.Lock()
		delete(e.pendingGetCookies, h)
		e.pendingMu.Unlock()
		return err
	}
	return nil
}

// AddOrUpdateCookie inserts c into the cookie store of the profile as a
// SESSION cookie (Expires is ignored; see ICoreWebView2Cookie for why). Must be
// called on the UI thread, typically before the first Navigate so the page
// sees the cookie from its very first request.
func (e *Chromium) AddOrUpdateCookie(c Cookie) error {
	if e.cookieManager == nil {
		return ErrCookieManagerUnavailable
	}
	path := c.Path
	if path == "" {
		path = "/"
	}
	raw, err := e.cookieManager.CreateCookie(c.Name, c.Value, c.Domain, path)
	if err != nil {
		return fmt.Errorf("CreateCookie(%s): %w", c.Name, err)
	}
	defer func() { _ = raw.Release() }()

	if err := raw.PutIsSecure(c.IsSecure); err != nil {
		return fmt.Errorf("PutIsSecure(%s): %w", c.Name, err)
	}
	if err := raw.PutIsHttpOnly(c.IsHttpOnly); err != nil {
		return fmt.Errorf("PutIsHttpOnly(%s): %w", c.Name, err)
	}
	if err := raw.PutSameSite(c.SameSite); err != nil {
		return fmt.Errorf("PutSameSite(%s): %w", c.Name, err)
	}
	if err := e.cookieManager.AddOrUpdateCookie(raw); err != nil {
		return fmt.Errorf("AddOrUpdateCookie(%s): %w", c.Name, err)
	}
	return nil
}
