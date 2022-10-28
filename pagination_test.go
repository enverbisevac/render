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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/enverbisevac/render"
	"github.com/enverbisevac/render/utest"
)

func defaultURL(page, perPage int) *url.URL {
	return &url.URL{
		Scheme:   "http",
		Host:     "localhost",
		Path:     "users",
		RawQuery: fmt.Sprintf("%s=%d&%s=%d", render.PageParam, page, render.PerPageParam, perPage),
	}
}

func request(page, perPage int) *http.Request {
	return &http.Request{
		URL: defaultURL(page, perPage),
	}
}

func TestWithPerPage(t *testing.T) {
	p := render.Pagination{}

	f := render.WithPerPage(20)
	f(&p)

	utest.Equals(t, 20, p.PerPage())
}

func TestPaginationFromRequest(t *testing.T) {
	r := request(1, 20)

	got := render.PaginationFromRequest(r, 100)

	utest.Equals(t, 1, got.Page())
	utest.Equals(t, 20, got.PerPage())
	utest.Equals(t, 100, got.Total())
}

func TestNewPagination(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, 1, got.Page())
	utest.Equals(t, 20, got.PerPage())
	utest.Equals(t, 100, got.Total())
}

func TestPagination_URL(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, uri, got.URL())
}

func TestPagination_Page(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, 1, got.Page())
}

func TestPagination_PerPage(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, uri, got.URL())
}

func TestPagination_Prev(t *testing.T) {
	uri := defaultURL(2, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, 1, got.Prev())
}

func TestPagination_PrevURL(t *testing.T) {
	uri := defaultURL(2, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, defaultURL(1, 20).String(), got.PrevURL())
}

func TestPagination_Next(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, 2, got.Next())
}

func TestPagination_NextURL(t *testing.T) {
	uri := defaultURL(1, 20)

	got := render.NewPagination(uri, 100)

	utest.Equals(t, defaultURL(2, 20).String(), got.NextURL())
}

func TestPagination_Last(t *testing.T) {
	perPage := 20
	total := 100
	uri := defaultURL(1, perPage)

	got := render.NewPagination(uri, total)

	utest.Equals(t, total/perPage, got.Last())
}

func TestPagination_LastURL(t *testing.T) {
	perPage := 20
	total := 100
	uri := defaultURL(1, perPage)

	got := render.NewPagination(uri, total)

	utest.Equals(t, defaultURL(total/perPage, perPage).String(), got.LastURL())
}

func TestPagination_Render(t *testing.T) {
}

func TestDefaultPaginationHeader(t *testing.T) {
	type user struct {
		Name string `json:"name"`
	}
	w := httptest.NewRecorder()
	r := request(1, 20)
	v := []user{
		{
			Name: "Enver",
		},
		{
			Name: "Joe",
		},
		{
			Name: "Dave",
		},
	}

	refPaginationInHeader := render.PaginationInHeader
	render.PaginationInHeader = true
	refHeaderFunc := render.DefaultPaginationHeader
	render.PaginationHeader = func(w http.ResponseWriter, p render.Pagination) {
		w.Header().Set("x-cur-page", strconv.Itoa(p.Page()))
		w.Header().Set("x-size", strconv.Itoa(p.PerPage()))

		last := p.Last()

		if p.Page() != last {
			w.Header().Set("next-page", strconv.Itoa(p.Next()))
			w.Header().Add(render.LinkHeader, fmt.Sprintf(render.Linkf, p.NextURL(), "next"))
		}

		if p.Page() > 1 {
			w.Header().Set("prev-page", strconv.Itoa(p.Prev()))
			w.Header().Add(render.LinkHeader, fmt.Sprintf(render.Linkf, p.PrevURL(), "prev"))
		}

		w.Header().Set("total-items", strconv.Itoa(p.Total()))
		w.Header().Set("total-pages", strconv.Itoa(last))
		w.Header().Add(render.LinkHeader, fmt.Sprintf(render.Linkf, p.LastURL(), "last"))
	}

	pagination := render.PaginationFromRequest(r, 100)
	pagination.Render(w, r, v)

	utest.Equals(t, "1", w.Header().Get("x-cur-page"))
	utest.Equals(t, "20", w.Header().Get("x-size"))
	utest.Equals(t, "100", w.Header().Get("total-items"))

	render.PaginationInHeader = refPaginationInHeader
	render.PaginationHeader = refHeaderFunc
}

func TestDefaultPaginationBody(t *testing.T) {
	type user struct {
		Name string `json:"name"`
	}
	w := httptest.NewRecorder()
	r := request(1, 20)
	v := []user{
		{
			Name: "Enver",
		},
		{
			Name: "Joe",
		},
		{
			Name: "Dave",
		},
	}
	type custom struct {
		Page  int `json:"page"`
		Size  int `json:"size"`
		Total int `json:"total"`
		Items interface{}
	}

	refPaginationInHeader := render.PaginationInHeader
	render.PaginationInHeader = false
	refBodyFunc := render.DefaultPaginationBody
	render.PaginationBody = func(p render.Pagination, v interface{}) interface{} {
		return custom{
			Page:  p.Page(),
			Size:  p.PerPage(),
			Total: p.Total(),
			Items: v,
		}
	}

	pagination := render.PaginationFromRequest(r, 100)
	pagination.Render(w, r, v)

	data, err := io.ReadAll(w.Body)
	utest.OK(t, err)

	cstRes := custom{}
	err = json.Unmarshal(data, &cstRes)
	utest.OK(t, err)

	utest.Equals(t, 1, cstRes.Page)
	utest.Equals(t, 20, cstRes.Size)
	utest.Equals(t, 100, cstRes.Total)

	render.PaginationInHeader = refPaginationInHeader
	render.PaginationBody = refBodyFunc
}
