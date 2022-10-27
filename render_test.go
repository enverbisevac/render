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
	"github.com/enverbisevac/render/utest"
	"net/http"
	"testing"

	"github.com/enverbisevac/render"
)

func TestBlob(t *testing.T) {
	var (
		buffer []byte
		status int
		header http.Header
	)

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
		v      []byte
		params []interface{}
	}
	tests := []struct {
		name    string
		args    args
		status  int
		content []byte
		header  http.Header
	}{
		{
			name: "happy path",
			args: args{
				w: writer,
				v: []byte("Some content"),
			},
			status:  http.StatusOK,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{"application/octet-stream"},
			},
		},
		{
			name: "happy path - status 201",
			args: args{
				w:      writer,
				v:      []byte("Some content"),
				params: []interface{}{http.StatusCreated},
			},
			status:  http.StatusCreated,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{"application/octet-stream"},
			},
		},
		{
			name: "happy path - content type application/json",
			args: args{
				w:      writer,
				v:      []byte("Some content"),
				params: []interface{}{render.ContentTypeHeader, render.ApplicationJSON},
			},
			status:  http.StatusOK,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{render.ApplicationJSON},
			},
		},
		{
			name: "happy path - content type status code and application/json",
			args: args{
				w:      writer,
				v:      []byte("Some content"),
				params: []interface{}{http.StatusCreated, render.ContentTypeHeader, render.ApplicationJSON},
			},
			status:  http.StatusCreated,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{render.ApplicationJSON},
			},
		},
		{
			name: "happy path - http.Header as argument",
			args: args{
				w: writer,
				v: []byte("Some content"),
				params: []interface{}{http.StatusCreated, http.Header{
					render.ContentTypeHeader: []string{render.ApplicationJSON},
				}},
			},
			status:  http.StatusCreated,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{render.ApplicationJSON},
			},
		},
		{
			name: "incomplete header pair should return default content type",
			args: args{
				w:      writer,
				v:      []byte("Some content"),
				params: []interface{}{http.StatusCreated, render.ContentTypeHeader},
			},
			status:  http.StatusCreated,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{"application/octet-stream"},
			},
		},
		{
			name: "incomplete header pair and no status specified should return default content type and OK status",
			args: args{
				w:      writer,
				v:      []byte("Some content"),
				params: []interface{}{render.ContentTypeHeader},
			},
			status:  http.StatusOK,
			content: []byte("Some content"),
			header: http.Header{
				render.ContentTypeHeader: []string{"application/octet-stream"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			render.Blob(tt.args.w, tt.args.v, tt.args.params...)
			utest.Equals(t, tt.status, status)
			utest.Equals(t, tt.content, buffer)
			utest.Equals(t, tt.header, header)
			status = 0
			buffer = []byte{}
			header = http.Header{}
		})
	}
}
