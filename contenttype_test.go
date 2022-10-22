package render

import (
	"net/http"
	"testing"
)

func TestGetContentType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want ContentType
	}{
		{
			name: "tex/plain content type",
			args: args{
				s: TextPlain,
			},
			want: ContentTypePlainText,
		},
		{
			name: "text/html content type",
			args: args{
				s: TextHTML,
			},
			want: ContentTypeHTML,
		},
		{
			name: "application/xhtml+xml content type",
			args: args{
				s: ApplicationXHTML,
			},
			want: ContentTypeHTML,
		},
		{
			name: "application/json content type",
			args: args{
				s: ApplicationJSON,
			},
			want: ContentTypeJSON,
		},
		{
			name: "text/javascript content type",
			args: args{
				s: TextJavascript,
			},
			want: ContentTypeJSON,
		},
		{
			name: "text/xml content type",
			args: args{
				s: TextXML,
			},
			want: ContentTypeXML,
		},
		{
			name: "application/xml content type",
			args: args{
				s: ApplicationXML,
			},
			want: ContentTypeXML,
		},
		{
			name: "application/x-www-form-urlencoded content type",
			args: args{
				s: ApplicationFormURL,
			},
			want: ContentTypeForm,
		},
		{
			name: "text/event-stream content type",
			args: args{
				s: TextEventStream,
			},
			want: ContentTypeEventStream,
		},
		{
			name: "unknown content type",
			args: args{
				s: "unknown",
			},
			want: ContentTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContentType(tt.args.s); got != tt.want {
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
		want ContentType
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
			want: ContentTypeJSON,
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
			want: ContentTypeJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAcceptedContentType(tt.args.r); got != tt.want {
				t.Errorf("GetAcceptedContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
