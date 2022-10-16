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

package render_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/enverbisevac/render"
)

func TestError(t *testing.T) {
	var (
		buffer []byte
		status int
		header http.Header
	)

	jsonEncoder := func(msg string) []byte {
		resErr := render.ErrorResponse{
			Message: msg,
		}

		data, err := json.Marshal(resErr)
		if err != nil {
			t.Fatal(err)
		}
		return append(data, '\n')
	}

	writer := &mockWriter{
		WriteFunc: func(b []byte) (int, error) {
			buffer = make([]byte, len(b))
			copy(buffer, b)
			return len(buffer), nil
		},
		WriteHeaderFunc: func(statusCode int) {
			status = statusCode
		},
		HeaderFunc: func() http.Header {
			header = http.Header{}
			return header
		},
	}

	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		err    error
		params []interface{}
	}
	tests := []struct {
		name   string
		args   args
		body   []byte
		status int
	}{
		{
			name: "default error 500 - Internal Server Error",
			args: args{
				w: writer,
				r: &http.Request{
					Header: http.Header{
						render.AcceptHeader: []string{render.ApplicationJSON},
					},
				},
				err: errors.New("some error"),
			},
			body:   jsonEncoder("some error"),
			status: http.StatusInternalServerError,
		},
		{
			name: "default mapped error 404 - file not found",
			args: args{
				w: writer,
				r: &http.Request{
					Header: http.Header{
						render.AcceptHeader: []string{render.ApplicationJSON},
					},
				},
				err: fmt.Errorf("file %s %w", "demo.txt", render.ErrNotFound),
			},
			body:   jsonEncoder(fmt.Sprintf("file %s %v", "demo.txt", render.ErrNotFound)),
			status: http.StatusNotFound,
		},
		{
			name: "set optional param to BadRequest status",
			args: args{
				w: writer,
				r: &http.Request{
					Header: http.Header{
						render.AcceptHeader: []string{render.ApplicationJSON},
					},
				},
				err:    errors.New("bad input data"),
				params: []interface{}{http.StatusBadRequest},
			},
			body:   jsonEncoder("bad input data"),
			status: http.StatusBadRequest,
		},
		{
			name: "provide http error",
			args: args{
				w: writer,
				r: &http.Request{
					Header: http.Header{
						render.AcceptHeader: []string{render.ApplicationJSON},
					},
				},
				err: &render.HTTPError{
					Err:    errors.New("conflict data"),
					Status: http.StatusConflict,
				},
			},
			body:   jsonEncoder("conflict data"),
			status: http.StatusConflict,
		},
		{
			name: "provide map error, http error and param status should return param status",
			args: args{
				w: writer,
				r: &http.Request{
					Header: http.Header{
						render.AcceptHeader: []string{render.ApplicationJSON},
					},
				},
				err: &render.HTTPError{
					Err:    render.ErrForbidden,
					Status: http.StatusBadGateway,
				},
				params: []interface{}{http.StatusGatewayTimeout},
			},
			body:   jsonEncoder(render.ErrForbidden.Error()),
			status: http.StatusGatewayTimeout,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			render.Error(tt.args.w, tt.args.r, tt.args.err, tt.args.params...)
			equals(t, tt.status, status)
			equals(t, string(tt.body), string(buffer))
			status = 0
			buffer = []byte{}
			header = http.Header{}
		})
	}
}
