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
	// PageHeader represents x-page key in header
	PageHeader = "x-page"
	// PerPageHeader represents x-per-page key in header
	PerPageHeader = "x-per-page"
	// NextPageHeader represents x-next key in header
	NextPageHeader = "x-next-page"
	// PrevPageHeader represents x-pprev key in header
	PrevPageHeader = "x-prev-page"
	// TotalItemsHeader represents x-total key in header
	TotalItemsHeader = "x-total"
	// TotalPagesHeader represents x-total-pages key in header
	TotalPagesHeader = "x-total-pages"
	// LinkHeader represents Link key in header
	LinkHeader = "Link"
	// Linkf is format for Link headers
	Linkf = `<%s>; rel="%s"`

	// PaginationInHeader write pagination in header
	PaginationInHeader = true
	// PaginationHeader generates pagination in header
	PaginationHeader = DefaultPaginationHeader
	// PaginationBody generates pagination in body
	PaginationBody = DefaultPaginationBody
)

// Pagination holds all page related data.
type Pagination struct {
	url     *url.URL
	page    int
	perPage int
	last    int
	total   int
}

// PaginationOption is prototype for functional options.
type PaginationOption func(p *Pagination)

// WithPerPage set perPage value.
func WithPerPage(val int) PaginationOption {
	return func(p *Pagination) {
		p.perPage = val
		p.last = totalPages(p.perPage, p.total)
	}
}

// PaginationFromRequest returns pagination object from parsed request url field
func PaginationFromRequest(r *http.Request, totalItems int, options ...PaginationOption) Pagination {
	return NewPagination(r.URL, totalItems, options...)
}

// NewPagination parses url and return new pagination object.
func NewPagination(url *url.URL, totalItems int, options ...PaginationOption) Pagination {
	queryParams := url.Query()
	strPage := queryParams.Get(PageParam)
	strPerPage := queryParams.Get(PerPageParam)

	page, err := strconv.Atoi(strPage)
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(strPerPage)
	if err != nil {
		perPage = PerPageDefault
	}

	last := totalPages(perPage, totalItems)

	pagination := Pagination{
		url:     url,
		page:    page,
		perPage: perPage,
		last:    last,
		total:   totalItems,
	}

	for _, option := range options {
		option(&pagination)
	}

	return pagination
}

// URL returns non exported page value
func (p Pagination) URL() *url.URL {
	return p.url
}

// Page returns non exported page value
func (p Pagination) Page() int {
	return p.page
}

// PerPage returns perPage (per_page) value
func (p Pagination) PerPage() int {
	return p.perPage
}

// Prev page
func (p Pagination) Prev() int {
	return max(p.page-1, 1)
}

// PrevURL page
func (p Pagination) PrevURL() string {
	if p.url == nil {
		return ""
	}
	params := p.url.Query()
	params.Set(PageParam, strconv.Itoa(p.page))
	params.Set(PerPageParam, strconv.Itoa(p.perPage))

	if p.page > 1 {
		params.Set(PageParam, strconv.Itoa(p.Prev()))
		p.url.RawQuery = params.Encode()

		return p.url.String()
	}
	return ""
}

// Next page
func (p Pagination) Next() int {
	return min(p.page+1, p.last)
}

// NextURL page
func (p Pagination) NextURL() string {
	if p.url == nil {
		return ""
	}
	params := p.url.Query()
	params.Set(PageParam, strconv.Itoa(p.page))
	params.Set(PerPageParam, strconv.Itoa(p.perPage))

	if p.page != p.last {
		params.Set(PageParam, strconv.Itoa(p.Next()))
		p.url.RawQuery = params.Encode()

		return p.url.String()
	}
	return ""
}

// Last page
func (p Pagination) Last() int {
	return p.last
}

// LastURL page
func (p Pagination) LastURL() string {
	if p.url == nil {
		return ""
	}
	params := p.url.Query()
	params.Set(PageParam, strconv.Itoa(p.page))
	params.Set(PerPageParam, strconv.Itoa(p.perPage))

	params.Set(PageParam, strconv.Itoa(p.last))
	p.url.RawQuery = params.Encode()

	return p.url.String()
}

// Total returns total number of elements
func (p Pagination) Total() int {
	return p.total
}

func (p Pagination) shouldRedirect() bool {
	last := p.last
	switch {
	case p.page == 0:
		return true
	case p.page > last:
		return true
	case p.perPage == 0:
		return true
	}
	return false
}

func (p Pagination) redirect(w http.ResponseWriter, r *http.Request) {
	uri := *r.URL

	last := p.last
	page := p.page
	perPage := p.perPage

	if page == 0 {
		page = 1
	}

	if page > last {
		page = last
	}

	if perPage == 0 {
		perPage = PerPageDefault
	}

	params := uri.Query()
	params.Set(PageParam, strconv.Itoa(page))
	params.Set(PerPageParam, strconv.Itoa(perPage))
	uri.RawQuery = params.Encode()

	http.Redirect(w, r, uri.String(), http.StatusMovedPermanently)
}

// Render renders payload and respond to the client request.
func (p Pagination) Render(w http.ResponseWriter, r *http.Request, v interface{}, params ...interface{}) {
	if p.shouldRedirect() {
		p.redirect(w, r)
		return
	}

	if PaginationInHeader {
		PaginationHeader(w, p)
	} else {
		v = PaginationBody(p, v)
	}

	Render(w, r, v, params...)
}

// DefaultPaginationHeader returns pagination metadata in header.
func DefaultPaginationHeader(w http.ResponseWriter, p Pagination) {
	w.Header().Set(PageHeader, strconv.Itoa(p.page))
	w.Header().Set(PerPageHeader, strconv.Itoa(p.perPage))

	last := p.last

	if p.page != last {
		w.Header().Set(NextPageHeader, strconv.Itoa(p.Next()))
		w.Header().Add(LinkHeader, fmt.Sprintf(Linkf, p.NextURL(), "next"))
	}

	if p.page > 1 {
		w.Header().Set(PrevPageHeader, strconv.Itoa(p.Prev()))
		w.Header().Add(LinkHeader, fmt.Sprintf(Linkf, p.PrevURL(), "prev"))
	}

	w.Header().Set(TotalItemsHeader, strconv.Itoa(p.total))
	w.Header().Set(TotalPagesHeader, strconv.Itoa(last))
	w.Header().Add(LinkHeader, fmt.Sprintf(Linkf, p.LastURL(), "last"))
}

type simpleBody struct {
	Page    int         `json:"page" xml:"page"`
	PerPage int         `json:"per_page" xml:"per_page"`
	Total   int         `json:"total" xml:"total"`
	Next    string      `json:"next,omitempty" xml:"next,omitempty"`
	Prev    string      `json:"prev,omitempty" xml:"prev,omitempty"`
	Last    string      `json:"last,omitempty" xml:"last,omitempty"`
	Items   interface{} `json:"items" xml:"items"`
}

// DefaultPaginationBody returns custom pagination body.
func DefaultPaginationBody(p Pagination, v interface{}) interface{} {
	return simpleBody{
		Page:    p.page,
		PerPage: p.perPage,
		Total:   p.total,
		Next:    p.NextURL(),
		Prev:    p.PrevURL(),
		Last:    p.LastURL(),
		Items:   v,
	}
}
