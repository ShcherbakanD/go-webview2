package edge

// FrameNavigationCompleted uses the same NavigationCompletedEventArgs as the
// top-level NavigationCompleted event but is fired for sub-frame (iframe)
// navigations. We need a separate handler type because the COM event uses
// its own IID; the impl interface is also separate so the host can wire
// distinct callbacks for top-level vs frame completion.

type _ICoreWebView2FrameNavigationCompletedEventHandlerImpl interface {
	_IUnknownImpl
	FrameNavigationCompleted(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

type iCoreWebView2FrameNavigationCompletedEventHandler struct {
	vtbl *_ICoreWebView2NavigationCompletedEventHandlerVtbl
	impl _ICoreWebView2FrameNavigationCompletedEventHandlerImpl
}

func _ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface(this *iCoreWebView2FrameNavigationCompletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef(this *iCoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease(this *iCoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2FrameNavigationCompletedEventHandlerInvoke(this *iCoreWebView2FrameNavigationCompletedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	return this.impl.FrameNavigationCompleted(sender, args)
}

var _ICoreWebView2FrameNavigationCompletedEventHandlerFn = _ICoreWebView2NavigationCompletedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2FrameNavigationCompletedEventHandlerInvoke),
}

func newICoreWebView2FrameNavigationCompletedEventHandler(impl _ICoreWebView2FrameNavigationCompletedEventHandlerImpl) *iCoreWebView2FrameNavigationCompletedEventHandler {
	return &iCoreWebView2FrameNavigationCompletedEventHandler{
		vtbl: &_ICoreWebView2FrameNavigationCompletedEventHandlerFn,
		impl: impl,
	}
}
