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

//nolint:unparam
func formatURL(page, perPage int) string {
	return fmt.Sprintf("http://localhost/users?%s=%d&%s=%d",
		PageParam, page, PerPageParam, perPage)
}

func TestNewPagination(t *testing.T) {
	makeURL := func(uri string) *url.URL {
		parse, err := url.Parse(uri)
		if err != nil {
			return nil
		}
		return parse
	}
	type args struct {
		url        *url.URL
		totalItems int
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "happy path",
			args: args{
				url:        makeURL(formatURL(1, 20)),
				totalItems: 100,
			},
			want: Pagination{
				url:     makeURL(formatURL(1, 20)),
				page:    1,
				perPage: 20,
				last:    5,
				total:   100,
			},
		},
		{
			name: "no page in query should return page 1",
			args: args{
				url:        makeURL("http://localhost/users?per_page=20"),
				totalItems: 100,
			},
			want: Pagination{
				url:     makeURL("http://localhost/users?per_page=20"),
				page:    1,
				perPage: 20,
				last:    5,
				total:   100,
			},
		},
		{
			name: "no per_page in query should return page default value",
			args: args{
				url:        makeURL("http://localhost/users?page=1"),
				totalItems: 100,
			},
			want: Pagination{
				url:     makeURL("http://localhost/users?page=1"),
				page:    1,
				perPage: PerPageDefault,
				last:    100 / PerPageDefault,
				total:   100,
			},
		},
		{
			name: "no query params should return default values",
			args: args{
				url:        makeURL("http://localhost/users"),
				totalItems: 100,
			},
			want: Pagination{
				url:     makeURL("http://localhost/users"),
				page:    1,
				perPage: PerPageDefault,
				last:    100 / PerPageDefault,
				total:   100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPagination(tt.args.url, tt.args.totalItems); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationFromRequest(t *testing.T) {
	uri := &url.URL{
		Scheme:   "http",
		Host:     "localhost",
		Path:     "users",
		RawQuery: "page=1&per_page=20",
	}
	type args struct {
		r          *http.Request
		totalItems int
		options    []PaginationOption
	}
	test := struct {
		args args
		want Pagination
	}{
		args: args{
			r: &http.Request{
				URL: uri,
			},
			totalItems: 60,
		},
		want: Pagination{
			url:     uri,
			page:    1,
			perPage: 20,
			total:   60,
			last:    3,
		},
	}

	if got := PaginationFromRequest(test.args.r, test.args.totalItems, test.args.options...); !reflect.DeepEqual(got, test.want) {
		t.Errorf("PaginationFromRequest() = %v, want %v", got, test.want)
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
		page: 1,
	}

	utest.Equals(t, 2, p.Next())
}

func TestPagination_Page(t *testing.T) {
	p := Pagination{
		page: 1,
	}

	utest.Equals(t, 1, p.Page())
}

func TestPagination_PerPage(t *testing.T) {
	p := Pagination{
		perPage: 20,
	}

	utest.Equals(t, 20, p.PerPage())
}

func TestPagination_Prev(t *testing.T) {
	p := Pagination{
		page: 2,
	}

	utest.Equals(t, 1, p.Prev())
}

func TestPagination_URL(t *testing.T) {
	uri := &url.URL{
		Scheme: "http",
		Host:   "localhost",
	}
	p := Pagination{
		url: uri,
	}

	utest.Equals(t, uri, p.URL())
}

//nolint:dupl
func TestPagination_PrevURL(t *testing.T) {
	type fields struct {
		url     *url.URL
		page    int
		perPage int
		last    int
		total   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "url is nil",
			fields: fields{
				page:    2,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: "",
		},
		{
			name: "happy path",
			fields: fields{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    2,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: formatURL(1, 20),
		},
		{
			name: "first page has no prev page",
			fields: fields{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				url:     tt.fields.url,
				page:    tt.fields.page,
				perPage: tt.fields.perPage,
				last:    tt.fields.last,
				total:   tt.fields.total,
			}
			if got := p.PrevURL(); got != tt.want {
				t.Errorf("PrevURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:dupl
func TestPagination_NextURL(t *testing.T) {
	type fields struct {
		url     *url.URL
		page    int
		perPage int
		last    int
		total   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "url is nil",
			fields: fields{
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: "",
		},
		{
			name: "happy path",
			fields: fields{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: formatURL(2, 20),
		},
		{
			name: "last page has no next page",
			fields: fields{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    3,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				url:     tt.fields.url,
				page:    tt.fields.page,
				perPage: tt.fields.perPage,
				last:    tt.fields.last,
				total:   tt.fields.total,
			}
			if got := p.NextURL(); got != tt.want {
				t.Errorf("NextURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_LastURL(t *testing.T) {
	type fields struct {
		url     *url.URL
		page    int
		perPage int
		last    int
		total   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "url is nil",
			fields: fields{
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: "",
		},
		{
			name: "happy path",
			fields: fields{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    3,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: formatURL(3, 20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				url:     tt.fields.url,
				page:    tt.fields.page,
				perPage: tt.fields.perPage,
				last:    tt.fields.last,
				total:   tt.fields.total,
			}
			if got := p.LastURL(); got != tt.want {
				t.Errorf("LastURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_shouldRedirect(t *testing.T) {
	type fields struct {
		url     *url.URL
		page    int
		perPage int
		last    int
		total   int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "happy path",
			fields: fields{
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: false,
		},
		{
			name: "page is 0 return true",
			fields: fields{
				page:    0,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: true,
		},
		{
			name: "page is greater then total pages return true",
			fields: fields{
				page:    4,
				perPage: 20,
				last:    3,
				total:   60,
			},
			want: true,
		},
		{
			name: "per page is 0 return true",
			fields: fields{
				page:    1,
				perPage: 0,
				last:    3,
				total:   60,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				url:     tt.fields.url,
				page:    tt.fields.page,
				perPage: tt.fields.perPage,
				last:    tt.fields.last,
				total:   tt.fields.total,
			}
			if got := p.shouldRedirect(); got != tt.want {
				t.Errorf("Pagination.shouldRedirect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_redirect(t *testing.T) {
	type fields struct {
		url     *url.URL
		page    int
		perPage int
		last    int
		total   int
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	type test struct {
		name   string
		fields fields
		args   args
		want   struct {
			url    string
			status int
		}
	}

	createTest := func(name string, page, perPage, total int, uri string, status int) test {
		last := total / PerPageDefault
		if perPage > 0 {
			last = total / perPage
		}
		return test{
			name: name,
			fields: fields{
				page:    page,
				perPage: perPage,
				last:    last,
				total:   total,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost",
						Path:   "users",
					},
				},
			},
			want: struct {
				url    string
				status int
			}{
				url:    uri,
				status: status,
			},
		}
	}
	tests := []test{
		createTest("happy path - redirect to the same address", 1, 20, 60,
			"http://localhost/users?page=1&per_page=20", http.StatusMovedPermanently),
		createTest("if page is 0 redirect to page 1", 0, 20, 60,
			"http://localhost/users?page=1&per_page=20", http.StatusMovedPermanently),
		createTest("if page is greater total pages the redirect to last one", 4, 20, 60,
			"http://localhost/users?page=3&per_page=20", http.StatusMovedPermanently),
		createTest("if per page is 0 redirect to same page with default per page", 5, 0, 100,
			fmt.Sprintf("http://localhost/users?page=4&per_page=%d", PerPageDefault),
			http.StatusMovedPermanently),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				url:     tt.fields.url,
				page:    tt.fields.page,
				perPage: tt.fields.perPage,
				last:    tt.fields.last,
				total:   tt.fields.total,
			}
			p.redirect(tt.args.w, tt.args.r)
			utest.Equals(t, tt.want.url, tt.args.w.Header().Get("Location"))
			utest.Equals(t, tt.want.status, tt.args.w.Code)
		})
	}
}

func TestDefaultPaginationHeader(t *testing.T) {
	f := func(uri string, size, total int) Pagination {
		_url, err := url.Parse(uri)
		if err != nil {
			return Pagination{}
		}
		pagination := NewPagination(_url, total, WithPerPage(size))
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
				p: f(formatURL(1, 20), 20, 100),
			},
			exp: map[string]string{
				PageHeader:       "1",
				PerPageHeader:    "20",
				NextPageHeader:   "2",
				PrevPageHeader:   "",
				TotalItemsHeader: "100",
				TotalPagesHeader: "5",
				LinkHeader:       fmt.Sprintf("<%s>; rel=\"next\"", formatURL(2, 20)),
			},
		},
		{
			name: "test last page",
			args: args{
				w: httptest.NewRecorder(),
				p: f(formatURL(5, 20), 20, 100),
			},
			exp: map[string]string{
				PageHeader:       "5",
				PerPageHeader:    "20",
				PrevPageHeader:   "4",
				TotalItemsHeader: "100",
				TotalPagesHeader: "5",
				LinkHeader:       fmt.Sprintf("<%s>; rel=\"prev\"", formatURL(4, 20)),
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

func TestDefaultPaginationBody(t *testing.T) {
	type args struct {
		p Pagination
		v interface{}
	}
	test := struct {
		args args
		want interface{}
	}{
		args: args{
			p: Pagination{
				url: &url.URL{
					Scheme: "http",
					Host:   "localhost",
					Path:   "users",
				},
				page:    1,
				perPage: 20,
				last:    3,
				total:   60,
			},
		},
		want: simpleBody{
			Page:    1,
			PerPage: 20,
			Total:   60,
			Prev:    "",
			Next:    formatURL(2, 20),
			Last:    formatURL(3, 20),
		},
	}

	if got := DefaultPaginationBody(test.args.p, test.args.v); !reflect.DeepEqual(got, test.want) {
		t.Errorf("DefaultPaginationBody() = %v, want %v", got, test.want)
	}
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
				total:   tt.fields.Total,
			}
			p.Render(tt.args.w, tt.args.r, tt.args.v, tt.args.params...)
		})
	}
}
