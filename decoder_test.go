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
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/enverbisevac/render"
)

func TestDefaultDecoder(t *testing.T) {
	type User struct {
		Name string `json:"name" form:"name"`
	}
	var user User
	type args struct {
		r *http.Request
		v interface{}
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "decode json data to user object",
			args: args{
				r: &http.Request{
					Header: http.Header{
						render.ContentTypeHeader: []string{render.ApplicationJSON},
					},
					Body: io.NopCloser(strings.NewReader("{\"name\":\"Enver\"}")),
				},
				v: &user,
			},
			err: nil,
		},
		{
			name: "decode xml data to user object",
			args: args{
				r: &http.Request{
					Header: http.Header{
						render.ContentTypeHeader: []string{"application/xml"},
					},
					Body: io.NopCloser(strings.NewReader("<name>Enver</name>")),
				},
				v: &user,
			},
			err: nil,
		},
		{
			name: "decode form data to user object",
			args: args{
				r: &http.Request{
					Header: http.Header{
						render.ContentTypeHeader: []string{"application/x-www-form-urlencoded"},
					},
					Body: io.NopCloser(strings.NewReader("name=Enver")),
				},
				v: &user,
			},
			err: nil,
		},
		{
			name: "decode error",
			args: args{
				r: &http.Request{
					Header: http.Header{
						render.ContentTypeHeader: []string{"some content type"},
					},
					Body: io.NopCloser(strings.NewReader("")),
				},
				v: &user,
			},
			err: render.ErrUnableToParseContentType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := render.DefaultDecoder(tt.args.r, tt.args.v); !errors.Is(err, tt.err) {
				t.Errorf("DefaultDecoder() error = %v, wantErr %v", err, tt.err)
			} else {
				user, ok := tt.args.v.(*User)
				if ok {
					equals(t, user.Name, "Enver")
				}
			}
		})
	}
}
