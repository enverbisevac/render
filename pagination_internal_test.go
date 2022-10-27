package render

import (
	"fmt"
	"github.com/enverbisevac/render/utest"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewPagination(t *testing.T) {
	makeURL := func(uri string) *url.URL {
		parse, err := url.Parse(uri)
		if err != nil {
			return nil
		}
		return parse
	}
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "happy path",
			args: args{
				url: makeURL("http://localhost/users?page=1&per_page=20"),
			},
			want: Pagination{
				url:  makeURL("http://localhost/users?page=1&per_page=20"),
				page: 1,
				size: 20,
			},
		},
		{
			name: "no page in query should return page 1",
			args: args{
				url: makeURL("http://localhost/users?per_page=20"),
			},
			want: Pagination{
				url:  makeURL("http://localhost/users?per_page=20"),
				page: 1,
				size: 20,
			},
		},
		{
			name: "no per_page in query should return page default value",
			args: args{
				url: makeURL("http://localhost/users?page=1"),
			},
			want: Pagination{
				url:  makeURL("http://localhost/users?page=1"),
				page: 1,
				size: PerPageDefault,
			},
		},
		{
			name: "no query params should return default values",
			args: args{
				url: makeURL("http://localhost/users"),
			},
			want: Pagination{
				url:  makeURL("http://localhost/users"),
				page: 1,
				size: PerPageDefault,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPagination(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_Last(t *testing.T) {

	p := Pagination{
		last: 10,
	}

	utest.Equals(t, 10, p.Last())
}

func TestPagination_Next(t *testing.T) {
	p := Pagination{
		next: 10,
	}

	utest.Equals(t, 10, p.Next())
}

func TestPagination_Page(t *testing.T) {
	p := Pagination{
		page: 1,
	}

	utest.Equals(t, 1, p.Page())
}

func TestPagination_Prev(t *testing.T) {
	p := Pagination{
		prev: 15,
	}

	utest.Equals(t, 15, p.Prev())
}

func TestDefaultPaginationHeader(t *testing.T) {
	f := func(uri string, next, prev, size, last, total int) Pagination {
		_url, err := url.Parse(uri)
		if err != nil {
			return Pagination{}
		}
		pagination := NewPagination(_url)
		pagination.next = next
		pagination.prev = prev
		pagination.size = size
		pagination.last = last
		pagination.Total = total
		return pagination
	}
	type args struct {
		w http.ResponseWriter
		p Pagination
	}
	tests := []struct {
		name string
		args args
		exp  map[string]string
	}{
		{
			name: "basic test",
			args: args{
				w: httptest.NewRecorder(),
				p: f(fmt.Sprintf("http://localhost/users?%s=1&%s=20", PageParam, PerPageParam),
					2, 0, 20, 5, 100),
			},
			exp: map[string]string{
				"x-page":        "1",
				"x-per-page":    "20",
				"x-next-page":   "2",
				"x-prev-page":   "",
				"x-total":       "100",
				"x-total-pages": "5",
				"Link": fmt.Sprintf("<http://localhost/users?%s=2&%s=20>; rel=\"next\"",
					PageParam, PerPageParam),
			},
		},
		{
			name: "test last page",
			args: args{
				w: httptest.NewRecorder(),
				p: f(fmt.Sprintf("http://localhost/users?%s=5&%s=20", PageParam, PerPageParam),
					0, 4, 20, 5, 100),
			},
			exp: map[string]string{
				"x-page":        "5",
				"x-per-page":    "20",
				"x-prev-page":   "4",
				"x-total":       "100",
				"x-total-pages": "5",
				"Link": fmt.Sprintf("<http://localhost/users?%s=4&%s=20>; rel=\"prev\"",
					PageParam, PerPageParam),
			},
		},
		{
			name: "uri is nil",
			args: args{
				w: httptest.NewRecorder(),
				p: Pagination{},
			},
			exp: map[string]string{
				"x-page":        "",
				"x-per-page":    "",
				"x-prev-page":   "",
				"x-total":       "",
				"x-total-pages": "",
				"Link":          "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultPaginationHeader(tt.args.w, tt.args.p)
			utest.Equals(t, tt.exp["x-page"], tt.args.w.Header().Get("x-page"))
			utest.Equals(t, tt.exp["x-per-page"], tt.args.w.Header().Get("x-per-page"))
			utest.Equals(t, tt.exp["x-next-page"], tt.args.w.Header().Get("x-next-page"))
			utest.Equals(t, tt.exp["x-total"], tt.args.w.Header().Get("x-total"))
			utest.Equals(t, tt.exp["x-total-pages"], tt.args.w.Header().Get("x-total-pages"))
			utest.Equals(t, tt.exp["Link"], tt.args.w.Header().Get("Link"))
		})
	}
}

func TestPagination_Size(t *testing.T) {
	p := Pagination{
		size: 10,
	}

	utest.Equals(t, 10, p.Size())
}

func TestPagination_Render(t *testing.T) {
	type fields struct {
		page  int
		size  int
		prev  int
		next  int
		last  int
		Total int
	}
	type args struct {
		w      http.ResponseWriter
		r      *http.Request
		v      interface{}
		params []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				page:  tt.fields.page,
				size:  tt.fields.size,
				prev:  tt.fields.prev,
				next:  tt.fields.next,
				last:  tt.fields.last,
				Total: tt.fields.Total,
			}
			p.Render(tt.args.w, tt.args.r, tt.args.v, tt.args.params...)
		})
	}
}
