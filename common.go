package webview2

import (
	"unsafe"

	"github.com/jchv/go-webview2/pkg/edge"
)

// This is copied from webview/webview.
// The documentation is included for convenience.

// Hint is used to configure window sizing and resizing behavior.
type Hint int

const (
	// HintNone specifies that width and height are default size
	HintNone Hint = iota

	// HintFixed specifies that window size can not be changed by a user
	HintFixed

	// HintMin specifies that width and height are minimum bounds
	HintMin

	// HintMax specifies that width and height are maximum bounds
	HintMax
)

// WebView is the interface for the webview.
type WebView interface {

	// Run runs the main loop until it's terminated. After this function exits -
	// you must destroy the webview.
	Run()

	// Terminate stops the main loop. It is safe to call this function from
	// a background thread.
	Terminate()

	// Dispatch posts a function to be executed on the main thread. You normally
	// do not need to call this function, unless you want to tweak the native
	// window.
	Dispatch(f func())

	// Destroy destroys a webview and closes the native window.
	Destroy()

	// Window returns a native window handle pointer. When using GTK backend the
	// pointer is GtkWindow pointer, when using Cocoa backend the pointer is
	// NSWindow pointer, when using Win32 backend the pointer is HWND pointer.
	Window() unsafe.Pointer

	// SetTitle updates the title of the native window. Must be called from the UI
	// thread.
	SetTitle(title string)

	// SetSize updates native window size. See Hint constants.
	SetSize(w int, h int, hint Hint)

	// Navigate navigates webview to the given URL. URL may be a data URI, i.e.
	// "data:text/text,<html>...</html>". It is often ok not to url-encode it
	// properly, webview will re-encode it for you.
	Navigate(url string)

	// SetHtml sets the webview HTML directly.
	// The origin of the page is `about:blank`.
	SetHtml(html string)

	// Init injects JavaScript code at the initialization of the new page. Every
	// time the webview will open a the new page - this initialization code will
	// be executed. It is guaranteed that code is executed before window.onload.
	Init(js string)

	// Eval evaluates arbitrary JavaScript code. Evaluation happens asynchronously,
	// also the result of the expression is ignored. Use RPC bindings if you want
	// to receive notifications about the results of the evaluation.
	Eval(js string)

	// Bind binds a callback function so that it will appear under the given name
	// as a global JavaScript function. Internally it uses webview_init().
	// Callback receives a request string and a user-provided argument pointer.
	// Request string is a JSON array of all the arguments passed to the
	// JavaScript function.
	//
	// f must be a function
	// f must return either value and error or just error
	Bind(name string, f interface{}) error

	// OnNewWindowRequested installs a handler that is called when the page
	// requests a new window (window.open, target="_blank", ctrl+click, etc).
	// The default WebView2 popup is automatically suppressed before the
	// handler runs; the handler decides what to do with the URI (open in a
	// new webview2 window, in the system browser, navigate in place, etc).
	// Pass nil to remove a previously installed handler.
	//
	// Use OnNewWindowRequestedRaw when you need to hand a real popup window
	// back to WebView2 (so JavaScript's window.open returns a live Window).
	OnNewWindowRequested(handler func(uri string))

	// OnNewWindowRequestedRaw is the low-level variant of
	// OnNewWindowRequested. The handler receives the raw event args and is
	// responsible for one of: args.PutHandled(true) (suppress popup),
	// args.PutNewWindow(child) (host-managed popup so window.open returns a
	// live Window ref), or doing nothing (default WebView2 behavior).
	// Pass nil to remove a previously installed handler.
	OnNewWindowRequestedRaw(handler func(args *edge.ICoreWebView2NewWindowRequestedEventArgs))

	// CoreWebView2 returns the raw ICoreWebView2 pointer backing this
	// webview. Mainly useful for handing the webview to another instance
	// via NewWindowRequestedEventArgs.PutNewWindow.
	CoreWebView2() *edge.ICoreWebView2

	// OnWebResourceRequested installs a handler invoked before any matching
	// request leaves WebView2. The handler can mutate request headers
	// (req.GetHeaders().SetHeader / RemoveHeader) to spoof or strip
	// fingerprint headers like sec-ch-ua* that WebView2 hard-codes with
	// "Microsoft Edge WebView2" branding. Must be paired with at least one
	// AddWebResourceRequestedFilter — without a filter the event never fires.
	// Pass nil to remove a previously installed handler.
	OnWebResourceRequested(handler func(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs))

	// AddWebResourceRequestedFilter narrows which requests fire the
	// WebResourceRequested handler. Pattern supports * wildcards; context
	// filters by request type (use COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL
	// to match everything).
	AddWebResourceRequestedFilter(uri string, ctx edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT)
}
