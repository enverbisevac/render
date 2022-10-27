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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/enverbisevac/render/utest"
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
				url:     makeURL("http://localhost/users?page=1&per_page=20"),
				page:    1,
				perPage: 20,
			},
		},
		{
			name: "no page in query should return page 1",
			args: args{
				url: makeURL("http://localhost/users?per_page=20"),
			},
			want: Pagination{
				url:     makeURL("http://localhost/users?per_page=20"),
				page:    1,
				perPage: 20,
			},
		},
		{
			name: "no per_page in query should return page default value",
			args: args{
				url: makeURL("http://localhost/users?page=1"),
			},
			want: Pagination{
				url:     makeURL("http://localhost/users?page=1"),
				page:    1,
				perPage: PerPageDefault,
			},
		},
		{
			name: "no query params should return default values",
			args: args{
				url: makeURL("http://localhost/users"),
			},
			want: Pagination{
				url:     makeURL("http://localhost/users"),
				page:    1,
				perPage: PerPageDefault,
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
		perPage: PerPageDefault,
		Total:   100,
	}

	utest.Equals(t, 4, p.Last())
}

func TestPagination_Next(t *testing.T) {
	p := Pagination{
		page:    1,
		perPage: PerPageDefault,
		Total:   100,
	}

	utest.Equals(t, 2, p.Next())
}

func TestPagination_Page(t *testing.T) {
	p := Pagination{
		page: 1,
	}

	utest.Equals(t, 1, p.Page())
}

func TestPagination_Prev(t *testing.T) {
	p := Pagination{
		page:    2,
		perPage: PerPageDefault,
		Total:   100,
	}

	utest.Equals(t, 1, p.Prev())
}

func TestDefaultPaginationHeader(t *testing.T) {
	f := func(uri string, size, total int) Pagination {
		_url, err := url.Parse(uri)
		if err != nil {
			return Pagination{}
		}
		pagination := NewPagination(_url)
		pagination.perPage = size
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
				p: f(fmt.Sprintf("http://localhost/users?%s=1&%s=20", PageParam, PerPageParam), 20, 100),
			},
			exp: map[string]string{
				PageHeader:       "1",
				PerPageHeader:    "20",
				"x-next-page":    "2",
				PrevPageHeader:   "",
				TotalItemsHeader: "100",
				TotalPagesHeader: "5",
				LinkHeader: fmt.Sprintf("<http://localhost/users?%s=2&%s=20>; rel=\"next\"",
					PageParam, PerPageParam),
			},
		},
		{
			name: "test last page",
			args: args{
				w: httptest.NewRecorder(),
				p: f(fmt.Sprintf("http://localhost/users?%s=5&%s=20", PageParam, PerPageParam), 20, 100),
			},
			exp: map[string]string{
				PageHeader:       "5",
				PerPageHeader:    "20",
				PrevPageHeader:   "4",
				TotalItemsHeader: "100",
				TotalPagesHeader: "5",
				LinkHeader: fmt.Sprintf("<http://localhost/users?%s=4&%s=20>; rel=\"prev\"",
					PageParam, PerPageParam),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultPaginationHeader(tt.args.w, tt.args.p)
			utest.Equals(t, tt.exp[PageHeader], tt.args.w.Header().Get(PageHeader))
			utest.Equals(t, tt.exp[PerPageHeader], tt.args.w.Header().Get(PerPageHeader))
			utest.Equals(t, tt.exp[NextPageHeader], tt.args.w.Header().Get(NextPageHeader))
			utest.Equals(t, tt.exp[TotalItemsHeader], tt.args.w.Header().Get(TotalItemsHeader))
			utest.Equals(t, tt.exp[TotalPagesHeader], tt.args.w.Header().Get(TotalPagesHeader))
			utest.Equals(t, tt.exp[LinkHeader], tt.args.w.Header().Get(LinkHeader))
		})
	}
}

func TestPagination_Size(t *testing.T) {
	p := Pagination{
		perPage: 10,
	}

	utest.Equals(t, 10, p.PerPage())
}

func TestPagination_Render(t *testing.T) {
	type fields struct {
		page  int
		size  int
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
				page:    tt.fields.page,
				perPage: tt.fields.size,
				Total:   tt.fields.Total,
			}
			p.Render(tt.args.w, tt.args.r, tt.args.v, tt.args.params...)
		})
	}
}
