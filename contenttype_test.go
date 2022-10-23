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
	"net/http"
	"testing"

	"github.com/enverbisevac/render"
)

func TestGetContentType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want render.ContentType
	}{
		{
			name: "tex/plain content type",
			args: args{
				s: render.TextPlain,
			},
			want: render.ContentTypePlainText,
		},
		{
			name: "text/html content type",
			args: args{
				s: render.TextHTML,
			},
			want: render.ContentTypeHTML,
		},
		{
			name: "application/xhtml+xml content type",
			args: args{
				s: render.ApplicationXHTML,
			},
			want: render.ContentTypeHTML,
		},
		{
			name: "application/json content type",
			args: args{
				s: render.ApplicationJSON,
			},
			want: render.ContentTypeJSON,
		},
		{
			name: "text/javascript content type",
			args: args{
				s: render.TextJavascript,
			},
			want: render.ContentTypeJSON,
		},
		{
			name: "text/xml content type",
			args: args{
				s: render.TextXML,
			},
			want: render.ContentTypeXML,
		},
		{
			name: "application/xml content type",
			args: args{
				s: render.ApplicationXML,
			},
			want: render.ContentTypeXML,
		},
		{
			name: "application/x-www-form-urlencoded content type",
			args: args{
				s: render.ApplicationFormURL,
			},
			want: render.ContentTypeForm,
		},
		{
			name: "text/event-stream content type",
			args: args{
				s: render.TextEventStream,
			},
			want: render.ContentTypeEventStream,
		},
		{
			name: "unknown content type",
			args: args{
				s: "unknown",
			},
			want: render.ContentTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render.GetContentType(tt.args.s); got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAcceptedContentType(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want render.ContentType
	}{
		{
			name: "happy path",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Accept": []string{"application/json; charset=utf-8"},
					},
				},
			},
			want: render.ContentTypeJSON,
		},
		{
			name: "unknown content type",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Accept": []string{"unknown path"},
					},
				},
			},
			want: render.ContentTypeJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render.GetAcceptedContentType(tt.args.r); got != tt.want {
				t.Errorf("GetAcceptedContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
