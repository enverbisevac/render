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

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enverbisevac/render"
)

type User struct {
	Name string `json:"name"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	name := paths[len(paths)-1]
	user := User{
		Name: name,
	}
	render.Render(w, r, user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	if err := render.Decode(r, &user); err != nil {
		render.Error(w, r, err)
		return
	}
	render.Render(w, r, user)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	render.Error(w, r, render.ErrNotFound)
}

func errorHandler1(w http.ResponseWriter, r *http.Request) {
	render.Error(w, r, fmt.Errorf("file demo.txt %w", render.ErrNotFound))
}

func errorHandler2(w http.ResponseWriter, r *http.Request) {
	render.Error(w, r, fmt.Errorf("file demo.txt %w", render.ErrNotFound), http.StatusInternalServerError)
}

func errorHandler3(w http.ResponseWriter, r *http.Request) {
	render.Error(w, r, &render.HTTPError{
		Err:    render.ErrNotFound,
		Status: http.StatusNotFound,
	})
}

func errorHandler4(w http.ResponseWriter, r *http.Request) {
	render.Error(w, r, &render.HTTPError{
		Err:    render.ErrNotFound,
		Status: http.StatusNotFound,
	}, http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/hello/", helloHandler)
	http.HandleFunc("/create", createUser)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/error1", errorHandler1)
	http.HandleFunc("/error2", errorHandler2)
	http.HandleFunc("/error3", errorHandler3)
	http.HandleFunc("/error4", errorHandler4)
	http.ListenAndServe(":8088", nil)
}
