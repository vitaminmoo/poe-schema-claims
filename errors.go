package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vitaminmoo/poe-schema-claims/ctxutil"
)

// strictRequestErrorHandler is called for handling errors hit while processing the request.
//
// It sets the response status code to 400. I don't actually know when this one would be called.
func strictRequestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	errPtr := ctxutil.GetErrorPtr(r.Context())
	*errPtr = err

	code := http.StatusBadRequest
	w.Header().Set("Content-Type", "application/json")
	ret := map[string]string{
		"code":  strconv.Itoa(code),
		"error": err.Error(),
	}
	b, jerr := json.MarshalIndent(ret, "", "  ")
	if jerr != nil {
		// Marshaling failed, write the finalError response with 500 status
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmt.Appendf(nil, `{
			"code": "500",
			"error": "error marshaling error to JSON: %s"
		}`, jerr.Error()))
		return
	}
	// Marshaling succeeded, write the original code and JSON body
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

// strictResponseErrorHandler is called for handling errors hit while processing
// the response.
//
// It sets the response status code to 500 and writes the error message returned
// by the handler method to the response message in JSON format.
//
// This is going to be the one hit for when the handler methods return an error.
func strictResponseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	errPtr := ctxutil.GetErrorPtr(r.Context())
	*errPtr = err

	code := http.StatusInternalServerError
	w.Header().Set("Content-Type", "application/json")
	// Delay WriteHeader until after successful marshaling
	ret := map[string]string{
		"code":  strconv.Itoa(code),
		"error": err.Error(),
	}
	b, jerr := json.MarshalIndent(ret, "", "  ")
	if jerr != nil {
		// Marshaling failed, write the finalError response with 500 status
		// Note: code is already 500 here, but we ensure WriteHeader is called.
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmt.Appendf(nil, `{
			"code": "500",
			"error": "error marshaling error to JSON: %s"
		}`, jerr.Error()))
		return
	}
	// Marshaling succeeded, write the original code and JSON body
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

// middlewareErrorHandler is called if automatic openapi input validation fails
//
// We don't have a real error or a request object here so we can't get the error
// to the logger. Honestly these are all client errors so we probably don't care
// anyway.
func middlewareErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	// Delay WriteHeader until after successful marshaling
	ret := map[string]string{
		"code":  strconv.Itoa(statusCode),
		"error": message,
	}
	b, jerr := json.MarshalIndent(ret, "", "  ")
	if jerr != nil {
		// Marshaling failed, write the finalError response with 500 status
		w.WriteHeader(http.StatusInternalServerError)
		// Note: Status code is already set to 500 by WriteHeader above
		_, _ = w.Write(fmt.Appendf(nil, `{
			"code": "500",
			"error": "error marshaling error to JSON: %s"
		}`, jerr.Error()))
		return
	}
	// Marshaling succeeded, write the original code and JSON body
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

// errorHandler is the non-strict error handler
//
// This is called if there's a failure somewhere in the middleware
func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	errPtr := ctxutil.GetErrorPtr(r.Context())
	*errPtr = err

	w.Header().Set("Content-Type", "application/json")
	var code int
	ret := map[string]string{}
	code = http.StatusBadRequest
	ret["error"] = err.Error()
	ret["code"] = strconv.Itoa(code)

	b, jerr := json.MarshalIndent(ret, "", "  ")
	if jerr != nil {
		// Marshaling failed, write the finalError response with 500 status
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(fmt.Appendf(nil, `{
			"code": "500",
			"error": "error marshaling error to JSON: %s"
		}`, jerr.Error()))
		return
	}
	// Marshaling succeeded, write the original code and JSON body
	w.WriteHeader(code)
	_, _ = w.Write(b)
}
