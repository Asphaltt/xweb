package xweb

import (
	"net/http"
)

/* ------------------ headers ------------------ */

// Header is alias of map[string]string
type Header map[string]string

// header prepared for response
var headerSetting = make(Header)
var headerAdding = make(Header)

// SetHeaders prepare headers to set when response
func SetHeaders(headers Header) {
	for k, v := range headers {
		headerSetting[k] = v
	}
}

// AddHeaders prepare headers to add when response
func AddHeaders(headers Header) {
	for k, v := range headers {
		headerAdding[k] = v
	}
}

/* ------------------ methods ------------------ */

// Callback is alias of http function
type Callback func(w http.ResponseWriter, r *http.Request)

// methods of all endpoints. uri -> method -> callback
var methods = make(map[string]map[string]Callback)

// Register registers the http method of uri into xweb
func Register(uri, meth string, cb Callback) {
	register(uri, meth, cb, false)
}

// RegisterWithHeader registers the http method of uri into xweb
// when response, set and add prepared headers
func RegisterWithHeader(uri, meth string, cb Callback) {
	register(uri, meth, cb, true)
}

func register(uri, meth string, cb Callback, withHeader bool) {
	if _, ok := methods[uri]; !ok {
		methods[uri] = make(map[string]Callback)
		methods[uri][http.MethodOptions] = methOptions
		http.HandleFunc(uri, response(uri, withHeader))
		registerUnexpectedMethods(uri)
	}
	methods[uri][meth] = cb
}

/* ------------------ normally response ------------------ */

func response(uri string, withHeader bool) Callback {
	return func(w http.ResponseWriter, req *http.Request) {
		if withHeader {
			for k, v := range headerSetting {
				w.Header().Set(k, v)
			}
			for k, v := range headerAdding {
				w.Header().Add(k, v)
			}
		}
		// Oh yes, there is no if and switch; they are ugly
		methods[uri][req.Method](w, req)
	}
}

/* ------------------ unexpected methods ------------------ */

var allHTTPMethods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

var methOptions = func(w http.ResponseWriter, r *http.Request) {}
var unGet = func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
var unCb = func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func registerUnexpectedMethods(uri string) {
	if _, existing := methods[uri][http.MethodGet]; !existing { // GET 404
		methods[uri][http.MethodGet] = unGet
	}
	for _, m := range allHTTPMethods {
		if _, existing := methods[uri][m]; !existing {
			methods[uri][m] = unCb // others 405
		}
	}
}
