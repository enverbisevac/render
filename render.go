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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const (
	ContentTypeHeader = "Content-Type"
	AcceptHeader      = "Accept"
)

// Respond is a package-level variable set to our default Responder. We do this
// because it allows you to set render.Respond to another function with the
// same function signature, while also utilizing the render.Responder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to test if v is an error and respond
// differently, or log something before you respond.
var Respond = DefaultResponder

// DefaultResponder handles streaming JSON and XML responses, automatically setting the
// Content-Type based on request headers. It will default to a JSON response.
func DefaultResponder(w http.ResponseWriter, r *http.Request, v interface{}, params ...interface{}) {
	if v != nil && reflect.TypeOf(v).Kind() == reflect.Chan {
		v = channelIntoSlice(w, r, v)
	}

	// Format response based on request Accept header.
	switch GetAcceptedContentType(r) {
	case ContentTypePlainText, ContentTypeUnknown:
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%s", v)
		}
		PlainText(w, str, params...)
	case ContentTypeJSON:
		JSON(w, v, params...)
	case ContentTypeXML:
		XML(w, v, params...)
	case ContentTypeEventStream:
		Stream(w, r, v)
	case ContentTypeForm:
		// TBD
		fallthrough
	case ContentTypeHTML:
		// TBD
		fallthrough
	default:
		JSON(w, v, params...)
	}
}

// Bind decodes a request body and executes the Binder method of the
// payload structure.
func Bind(r *http.Request, v interface{}) error {
	return Decode(r, v)
}

// Render renders payload and respond to the client request.
func Render(w http.ResponseWriter, r *http.Request, v interface{}, params ...interface{}) {
	Respond(w, r, v, params...)
}

// Blob writes raw bytes to the response, the default Content-Type as
// application/octet-stream, params is optional which can be int or string type.
// Int will provide status code and string is for header pair values
//
// for example:
//
// Blob(w, v)
// Blob(w, v, http.StatusOK)
// Blob(w, v, http.StatusOK, "Content-Type", "application/json")
// Blob(w, v, "Content-Type", "application/json")
// Blob(w, v, "Content-Type", "application/json", http.StatusOK)
//
// the order of the parameters does not matter.
func Blob(w http.ResponseWriter, v []byte, params ...interface{}) {
	w.Header().Set(ContentTypeHeader, "application/octet-stream")
	status, key, value := 0, "", ""
	for _, param := range params {
		if rv := reflect.ValueOf(param); rv.Kind() == reflect.Ptr {
			param = rv.Elem().Interface()
		}
		switch arg := param.(type) {
		case int:
			if status == 0 && arg != 0 {
				// when status is set and there are more int values in params
				// ignore all values
				status = arg
			}
		case string:
			if key == "" {
				key = arg
			} else {
				value = arg
			}

			if key != "" && value != "" {
				w.Header().Set(key, value)
				key, value = "", ""
			}
		}
	}
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	w.Write(v) //nolint:errcheck
}

// PlainText writes a string to the response, setting the Content-Type as
// text/plain.
func PlainText(w http.ResponseWriter, v string, args ...interface{}) {
	Blob(w, []byte(v), append(args, ContentTypeHeader, "text/plain; charset=utf-8")...)
}

// HTML writes a string to the response, setting the Content-Type as text/html.
func HTML(w http.ResponseWriter, v string, args ...interface{}) {
	Blob(w, []byte(v), append(args, ContentTypeHeader, "text/html; charset=utf-8")...)
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, v interface{}, args ...interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Blob(w, buf.Bytes(), append(args, ContentTypeHeader, ApplicationJSONExt)...)
}

// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func XML(w http.ResponseWriter, v interface{}, args ...interface{}) {
	b, err := xml.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Try to find <?xml header in first 100 bytes (just in case there're some XML comments).
	findHeaderUntil := len(b)
	if findHeaderUntil > 100 {
		findHeaderUntil = 100
	}
	if !bytes.Contains(b[:findHeaderUntil], []byte("<?xml")) {
		// No header found. Print it out first.
		w.Write([]byte(xml.Header)) //nolint:errcheck
	}

	Blob(w, b, append(args, ContentTypeHeader, "application/xml; charset=utf-8")...)
}

// File sends a response with the content of the file.
func File(w http.ResponseWriter, r *http.Request, fullPath string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fullPath))
	w.Header().Set(ContentTypeHeader, "application/octet-stream")
	http.ServeFile(w, r, fullPath)
}

// Attachment sends a response as attachment, prompting client to save the
// file.
func Attachment(w http.ResponseWriter, r *http.Request, fullPath string) {
	w.Header().Set("Content-Disposition", "attachment")
	w.Header().Set(ContentTypeHeader, "application/octet-stream")
	http.ServeFile(w, r, fullPath)
}

// Inline sends a response as inline, opening the file in the browser.
func Inline(w http.ResponseWriter, r *http.Request, fullPath string) {
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set(ContentTypeHeader, "application/octet-stream")
	http.ServeFile(w, r, fullPath)
}

// NoContent returns a HTTP 204 "No Content" response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Stream sends a streaming response with status code and content type.
func Stream(w http.ResponseWriter, r *http.Request, v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Chan {
		panic(fmt.Sprintf("render: event stream expects a channel, not %v", reflect.TypeOf(v).Kind()))
	}

	w.Header().Set(ContentTypeHeader, "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	if r.ProtoMajor == 1 {
		// An endpoint MUST NOT generate an HTTP/2 message containing connection-specific header fields.
		// Source: RFC7540
		w.Header().Set("Connection", "keep-alive")
	}

	w.WriteHeader(http.StatusOK)

	ctx := r.Context()
	for {
		switch chosen, recv, ok := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(v)},
		}); chosen {
		case 0: // equivalent to: case <-ctx.Done()
			w.Write([]byte("event: error\ndata: {\"error\":\"Server Timeout\"}\n\n")) //nolint:errcheck
			return

		default: // equivalent to: case v, ok := <-stream
			if !ok {
				w.Write([]byte("event: EOF\n\n")) //nolint:errcheck
				return
			}
			v := recv.Interface()

			bytes, err := json.Marshal(v)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("event: error\ndata: {\"error\":\"%v\"}\n\n", err))) //nolint:errcheck
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				continue
			}
			w.Write([]byte(fmt.Sprintf("event: data\ndata: %s\n\n", bytes))) //nolint:errcheck
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

// channelIntoSlice buffers channel data into a slice.
func channelIntoSlice(w http.ResponseWriter, r *http.Request, from interface{}) interface{} {
	ctx := r.Context()

	var to []interface{}
	for {
		switch chosen, recv, ok := reflect.Select([]reflect.SelectCase{
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(from)},
		}); chosen {
		case 0: // equivalent to: case <-ctx.Done()
			http.Error(w, "Server Timeout", http.StatusGatewayTimeout)
			return nil
		default: // equivalent to: case v, ok := <-stream
			if !ok {
				return to
			}
			to = append(to, recv.Interface())
		}
	}
}
