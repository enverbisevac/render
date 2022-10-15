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
	"net/http"
	"strings"
)

// DefaultContentType is a package-level variable set to our default content type
var DefaultContentType = ContentTypeJSON

// ContentType is an enumeration of common HTTP content types.
type ContentType int

// ContentTypes handled by this package.
const (
	ContentTypeUnknown ContentType = iota
	ContentTypePlainText
	ContentTypeHTML
	ContentTypeJSON
	ContentTypeXML
	ContentTypeForm
	ContentTypeEventStream
)

// GetContentType returns ContentType value based on input s
func GetContentType(s string) ContentType {
	s = strings.TrimSpace(strings.Split(s, ";")[0])
	switch s {
	case "text/plain":
		return ContentTypePlainText
	case "text/html", "application/xhtml+xml":
		return ContentTypeHTML
	case "application/json", "text/javascript":
		return ContentTypeJSON
	case "text/xml", "application/xml":
		return ContentTypeXML
	case "application/x-www-form-urlencoded":
		return ContentTypeForm
	case "text/event-stream":
		return ContentTypeEventStream
	default:
		return ContentTypeUnknown
	}
}

// GetRequestContentType is a helper function that returns ContentType based on
// context or request headers.
func GetRequestContentType(r *http.Request) ContentType {
	return GetContentType(r.Header.Get("Content-Type"))
}

// GetAcceptedContentType reads Accept header from request and returns ContentType
func GetAcceptedContentType(r *http.Request) ContentType {
	var contentType ContentType

	// Parse request Accept header.
	fields := strings.Split(r.Header.Get("Accept"), ",")
	if len(fields) > 0 {
		contentType = GetContentType(strings.TrimSpace(fields[0]))
	}

	if contentType == ContentTypeUnknown {
		contentType = DefaultContentType
	}
	return contentType
}
