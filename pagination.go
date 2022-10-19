package render

import (
	"fmt"
	"net/http"
	"strconv"
)

var (
	PageParam          = "page"
	PerPageParam       = "per_page"
	PerPageDefault     = 25
	Linkf              = `<%s>; rel="%s"`
	PaginationInHeader = true
	PaginationBody     = DefaultPaginationBody
)

type Pagination struct {
	page  int
	size  int
	prev  int
	next  int
	last  int
	Total int
}

func (p Pagination) Page() int {
	return p.page
}

func (p Pagination) Size() int {
	return p.size
}

func (p Pagination) Prev() int {
	return p.prev
}

func (p Pagination) Next() int {
	return p.next
}

func (p Pagination) Last() int {
	return p.last
}

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

	// if p.page > p.last {
	// 	p.page = p.last
	// 	redirect = true
	// }

	if redirect {
		uri := *r.URL

		params := uri.Query()
		params.Set(PageParam, strconv.Itoa(p.page))
		params.Set(PerPageParam, strconv.Itoa(p.size))
		uri.RawQuery = params.Encode()

		http.Redirect(w, r, uri.String(), http.StatusMovedPermanently)
		return
	}

	p.last = totalPages(p.size, p.Total)
	p.next = min(p.page+1, p.last)
	p.prev = max(p.page-1, 1)

	if PaginationInHeader {
		paginationHeader(w, r, p)
	} else {
		v = PaginationBody(r, p, v)
	}

	Render(w, r, v, params...)
}

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

func PaginationFromRequest(r *http.Request) Pagination {
	queryParams := r.URL.Query()
	strPage := queryParams.Get(PageParam)
	strPerPage := queryParams.Get(PerPageParam)

	if strPage == "" || strPerPage == "" {
		return Pagination{}
	}

	page, err := strconv.Atoi(strPage)
	if err != nil {
		return Pagination{}
	}
	size, err := strconv.Atoi(strPerPage)
	if err != nil {
		return Pagination{}
	}

	return Pagination{
		page: page,
		size: size,
	}
}

func paginationHeader(w http.ResponseWriter, r *http.Request, p Pagination) {
	uri := *r.URL

	params := uri.Query()
	params.Set(PageParam, strconv.Itoa(p.page))
	params.Set(PerPageParam, strconv.Itoa(p.size))

	w.Header().Set("x-page", strconv.Itoa(p.page))
	w.Header().Set("x-per-page", strconv.Itoa(p.size))

	if p.page != p.last {
		params.Set(PageParam, strconv.Itoa(p.next))
		uri.RawQuery = params.Encode()

		w.Header().Set("x-next-page", strconv.Itoa(p.next))
		w.Header().Add("Link", fmt.Sprintf(Linkf, uri.String(), "next"))
	}

	if p.page > 1 {
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
