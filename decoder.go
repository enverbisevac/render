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
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/ajg/form"
)

// ErrUnableToParseContentType is an error for unknown content type
var ErrUnableToParseContentType = errors.New("render: unable to automatically decode the request content type")

var (
	// JSONDecoder is a package-level variable set to our default JSON decoder
	// function.
	JSONDecoder = DefaultJSONDecoder
	// XMLDecoder is a package-level variable set to our default XML decoder
	// function.
	XMLDecoder = DefaultXMLDecoder
	// FormDecoder is a package-level variable set to our default Form decoder
	// function.
	FormDecoder = DefaultFormDecoder
)

// Decoder decodes data from reader
type Decoder interface {
	Decode(v interface{}) error
}

// DefaultJSONDecoder returns new JSON decoder for decoding
// JSON data.
func DefaultJSONDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// DefaultXMLDecoder returns new XML decoder for decoding
// XML data.
func DefaultXMLDecoder(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

// DefaultFormDecoder returns new Form decoder for decoding
// form data.
func DefaultFormDecoder(r io.Reader) Decoder {
	return form.NewDecoder(r)
}

// Decode is a package-level variable set to our DefaultDecoder. We do this
// because it allows you to set render.Decode to another function with the
// same function signature, while also utilizing the render.DefaultDecoder()
// function itself. Effectively, allowing you to easily add your own logic
// to the package defaults. For example, maybe you want to impose a limit
// on the number of bytes allowed to be read from the request body.
var Decode = DefaultDecoder

// DefaultDecoder detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func DefaultDecoder(r *http.Request, v interface{}) (err error) {
	switch GetRequestContentType(r) {
	case ContentTypeJSON:
		err = DecodeJSON(r.Body, v)
	case ContentTypeXML:
		err = DecodeXML(r.Body, v)
	case ContentTypeForm:
		err = DecodeForm(r.Body, v)
	case ContentTypePlainText:
		// to consider (string for example)
	case ContentTypeEventStream, ContentTypeHTML:
		// event stream not used
	case ContentTypeUnknown: // this should be always on top of default
		fallthrough
	default:
		err = ErrUnableToParseContentType
	}
	return
}

// DecodeJSON decodes a given reader into an interface using the json decoder.
func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return JSONDecoder(r).Decode(v)
}

// DecodeXML decodes a given reader into an interface using the xml decoder.
func DecodeXML(r io.Reader, v interface{}) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return XMLDecoder(r).Decode(v)
}

// DecodeForm decodes a given reader into an interface using the form decoder.
func DecodeForm(r io.Reader, v interface{}) error {
	return FormDecoder(r).Decode(v)
}
