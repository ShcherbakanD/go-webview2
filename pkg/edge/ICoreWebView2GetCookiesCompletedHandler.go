package edge

// ICoreWebView2GetCookiesCompletedHandler (IID 5A4F5069-5C15-47C3-8646-F4DE1C116670)
// is the completion callback of ICoreWebView2CookieManager.GetCookies.
//
// Unlike the long-lived event handlers in this package (which share the
// Chromium singleton as impl), a completion handler is created per call and
// fires exactly once. The caller is responsible for keeping it reachable until
// Invoke runs; WebView2 holds only the raw pointer.
type _ICoreWebView2GetCookiesCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2GetCookiesCompletedHandler struct {
	vtbl *_ICoreWebView2GetCookiesCompletedHandlerVtbl
	impl _ICoreWebView2GetCookiesCompletedHandlerImpl
}

func _ICoreWebView2GetCookiesCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2GetCookiesCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2GetCookiesCompletedHandlerIUnknownAddRef(this *iCoreWebView2GetCookiesCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2GetCookiesCompletedHandlerIUnknownRelease(this *iCoreWebView2GetCookiesCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2GetCookiesCompletedHandlerInvoke(this *iCoreWebView2GetCookiesCompletedHandler, errorCode uintptr, list *ICoreWebView2CookieList) uintptr {
	return this.impl.GetCookiesCompleted(errorCode, list)
}

type _ICoreWebView2GetCookiesCompletedHandlerImpl interface {
	_IUnknownImpl
	GetCookiesCompleted(errorCode uintptr, list *ICoreWebView2CookieList) uintptr
}

var _ICoreWebView2GetCookiesCompletedHandlerFn = _ICoreWebView2GetCookiesCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2GetCookiesCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2GetCookiesCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2GetCookiesCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2GetCookiesCompletedHandlerInvoke),
}

func newICoreWebView2GetCookiesCompletedHandler(impl _ICoreWebView2GetCookiesCompletedHandlerImpl) *iCoreWebView2GetCookiesCompletedHandler {
	return &iCoreWebView2GetCookiesCompletedHandler{
		vtbl: &_ICoreWebView2GetCookiesCompletedHandlerFn,
		impl: impl,
	}
}
