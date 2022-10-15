// Copyright (c) 2022 Enver Bisevac
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package render

import (
	"errors"
	"net/http"
)

var (
	// ErrInvalidToken is returned when the api request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")

	// ErrUnauthorized is returned when the user is not authorized.
	ErrUnauthorized = errors.New("Unauthorized")

	// ErrForbidden is returned when user access is forbidden.
	ErrForbidden = errors.New("Forbidden")

	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("not found")
)

// ErrorMap contains predefined errors with assigned status code.
var ErrorMap = map[error]int{
	ErrInvalidToken: http.StatusBadRequest,
	ErrUnauthorized: http.StatusUnauthorized,
	ErrForbidden:    http.StatusForbidden,
	ErrNotFound:     http.StatusNotFound,
}

// TreatError is a package-level variable set to default function with basic
// error message response. Any error provided will have just a simple struct
// with field message describing the error. Developer can create custom function
// for treating error responses, for example:
//
//	render.TreatError = func(r *http.Request, err error) interface{} {
//		cerr := &CustomError{}
//		if errors.As(err, &cerr) {
//			return &customResponse{
//				Message: cerr.Err,
//				Version: "1.0",
//				...
//		}
//
//	    return &HTTPError{
//	      Message: "some error message",
//	      Status: http.StatusCreated,
//	   }
//	}
//
// and render.Error(w, r, err) will create response based of your treat function.
var TreatError = DefaultErrorRespond

// ErrorResponse represents a json-encoded API error.
type ErrorResponse struct {
	Message string `json:"message" xml:"message"`
}

// HTTPError helper structure used as error with status code.
type HTTPError struct {
	Err    error
	Status int
}

// Error method returns error from HTTPError
func (h HTTPError) Error() string {
	return h.Err.Error()
}

// DefaultErrorRespond returns ErrorResponse object for later processing
func DefaultErrorRespond(r *http.Request, err error) interface{} {
	return ErrorResponse{
		Message: err.Error(),
	}
}

// Error renders response body with content type based on Accept header of request.
// Status codes must be >= 400.
func Error(w http.ResponseWriter, r *http.Request, err error, params ...interface{}) {
	status := http.StatusInternalServerError
	// find in map of default errors and return status
	for key, value := range ErrorMap {
		if errors.Is(err, key) {
			status = value
		}
	}
	// http error checking
	httpError := &HTTPError{}
	if errors.As(err, &httpError) {
		status = httpError.Status
		err = httpError.Err
	}
	v := TreatError(r, err)
	Respond(w, r, v, append(params, status)...)
}
