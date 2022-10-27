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
	"net/url"
	"strconv"
)

var (
	// PageParam is query name param for current page
	PageParam = "page"
	// PerPageParam is number of items per page
	PerPageParam = "per_page"
	// PerPageDefault sets default number of items on response
	PerPageDefault = 25
	// Linkf is format for Link headers
	Linkf = `<%s>; rel="%s"`
	// PaginationInHeader write pagination in header
	PaginationInHeader = true
	// PaginationHeader generates pagination in header
	PaginationHeader = DefaultPaginationHeader
	// PaginationBody generates pagination in body
	PaginationBody = DefaultPaginationBody
)

// Pagination holds all page related data
type Pagination struct {
	url   *url.URL
	page  int
	size  int
	prev  int
	next  int
	last  int
	Total int
}

func PaginationFromRequest(r *http.Request) Pagination {
	return NewPagination(r.URL)
}

// NewPagination parses url and return new pagination object
func NewPagination(url *url.URL) Pagination {
	queryParams := url.Query()
	strPage := queryParams.Get(PageParam)
	strPerPage := queryParams.Get(PerPageParam)

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(strPerPage)
	if err != nil {
		size = PerPageDefault
	}

	return Pagination{
		url:  url,
		page: page,
		size: size,
	}
}

// URL returns non exported page value
func (p Pagination) URL() *url.URL {
	return p.url
}

// Page returns non exported page value
func (p Pagination) Page() int {
	return p.page
}

// Size returns size (per_page) value
func (p Pagination) Size() int {
	return p.size
}

// Prev page
func (p Pagination) Prev() int {
	return p.prev
}

// Next page
func (p Pagination) Next() int {
	return p.next
}

// Last page
func (p Pagination) Last() int {
	return p.last
}

// Render renders payload and respond to the client request.
func (p Pagination) Render(w http.ResponseWriter, r *http.Request, v interface{}, params ...interface{}) {
	redirect := false
	if p.page == 0 {
		p.page = 1
		redirect = true
	}
	if p.size == 0 {
		p.size = PerPageDefault
		redirect = true
	}

	p.last = totalPages(p.size, p.Total)

	if p.page > p.last {
		p.page = p.last
		redirect = true
	}

	if redirect {
		uri := *r.URL

		params := uri.Query()
		params.Set(PageParam, strconv.Itoa(p.page))
		params.Set(PerPageParam, strconv.Itoa(p.size))
		uri.RawQuery = params.Encode()

		http.Redirect(w, r, uri.String(), http.StatusMovedPermanently)
		return
	}

	p.next = min(p.page+1, p.last)
	p.prev = max(p.page-1, 1)

	if PaginationInHeader {
		PaginationHeader(w, p)
	} else {
		v = PaginationBody(r, p, v)
	}

	Render(w, r, v, params...)
}

func DefaultPaginationHeader(w http.ResponseWriter, p Pagination) {
	if p.url == nil {
		return
	}
	uri := *p.url

	page := p.Page()
	size := p.Size()

	params := uri.Query()
	params.Set(PageParam, strconv.Itoa(page))
	params.Set(PerPageParam, strconv.Itoa(size))

	w.Header().Set("x-page", strconv.Itoa(page))
	w.Header().Set("x-per-page", strconv.Itoa(size))

	if page != p.last {
		params.Set(PageParam, strconv.Itoa(p.next))
		uri.RawQuery = params.Encode()

		w.Header().Set("x-next-page", strconv.Itoa(p.next))
		w.Header().Add("Link", fmt.Sprintf(Linkf, uri.String(), "next"))
	}

	if page > 1 {
		params.Set(PageParam, strconv.Itoa(p.prev))
		uri.RawQuery = params.Encode()

		w.Header().Set("x-prev-page", strconv.Itoa(p.prev))
		w.Header().Add("Link", fmt.Sprintf(Linkf, uri.String(), "prev"))
	}

	params.Set(PageParam, strconv.Itoa(p.last))
	uri.RawQuery = params.Encode()

	w.Header().Set("x-total", strconv.Itoa(p.Total))
	w.Header().Set("x-total-pages", strconv.Itoa(p.last))
	w.Header().Add("Link", fmt.Sprintf(Linkf, uri.String(), "last"))
}

// DefaultPaginationBody returns custom pagination body
func DefaultPaginationBody(r *http.Request, p Pagination, v interface{}) interface{} {
	var (
		next string
		prev string
		last string
	)

	uri := *r.URL

	params := uri.Query()
	params.Set(PageParam, strconv.Itoa(p.page))
	params.Set(PerPageParam, strconv.Itoa(p.size))

	if p.page != p.last {
		params.Set(PageParam, strconv.Itoa(p.next))
		uri.RawQuery = params.Encode()

		next = uri.String()
	}

	if p.page > 1 {
		params.Set(PageParam, strconv.Itoa(p.prev))
		uri.RawQuery = params.Encode()

		prev = uri.String()
	}

	params.Set(PageParam, strconv.Itoa(p.last))
	uri.RawQuery = params.Encode()

	last = uri.String()
	return struct {
		Page    int         `json:"page" xml:"page"`
		PerPage int         `json:"per_page" xml:"per_page"`
		Total   int         `json:"total" xml:"total"`
		Next    string      `json:"next,omitempty" xml:"next,omitempty"`
		Prev    string      `json:"prev,omitempty" xml:"prev,omitempty"`
		Last    string      `json:"last,omitempty" xml:"last,omitempty"`
		Items   interface{} `json:"items" xml:"items"`
	}{
		Page:    p.page,
		PerPage: p.size,
		Total:   p.Total,
		Next:    next,
		Prev:    prev,
		Last:    last,
		Items:   v,
	}
}
