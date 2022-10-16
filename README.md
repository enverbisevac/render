[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/enverbisevac/render)
[![Coverage Status](https://coveralls.io/repos/github/enverbisevac/render/badge.svg?branch=main)](https://coveralls.io/github/enverbisevac/render?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/enverbisevac/render)](https://goreportcard.com/report/github.com/enverbisevac/render)
[![CodeQL](https://github.com/enverbisevac/render/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/enverbisevac/render/actions/workflows/codeql-analysis.yaml)

# Render

The `render` package helps manage HTTP request / response payloads. The motivation and
ideas for making this package come from [go-chi/render](https://github.com/go-chi/render).

Every well-designed, robust and maintainable Web Service / REST API also needs
well-_defined_ request and response payloads. Together with the endpoint handlers,
the request and response payloads make up the contract between your server and the
clients calling on it.

Typically in a REST API application, you will have your data models (objects/structs)
that hold lower-level runtime application state, and at times you need to assemble,
decorate, hide or transform the representation before responding to a client. That
server output (response payload) structure, is also likely the input structure to
another handler on the server.

This is where `render` comes in - offering a few simple helpers to provide a simple
pattern for managing payload encoding and decoding.

Render is also combined with some helpers for responding to content types and parsing
request bodies. Please have a look at the examples. [examples](./examples/getting_started/README.md).

All feedback is welcome, thank you!

## Features

- Very simple API
- Render based on Accept header
- Custom render functions JSON, XML, PlainText ...
- Map of defaults error and statusCode
- Customizable error handling
- Easy decoding request body based on content type
- Switch encoders/decoders with some popular open source lib

## Installation

Install **render** with go get:

```go
  go get github.com/enverbisevac/render
```

## Usage/Examples

```go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enverbisevac/render"
)

type Person struct {
	Name string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	name := paths[len(paths)-1]
	render.Render(w, r, Person{
		Name: name,
	})
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

func main() {
    http.HandleFunc("/hello/", helloHandler)
    http.HandleFunc("/create", createUser)
	http.HandleFunc("/error/", errorHandler)
	http.ListenAndServe(":8088", nil)
}
```

## API Reference

#### Bind request body to data type `v`

```go
  func Bind(r *http.Request, v interface{}) error
```

| Parameter | Type            | Description                          |
| :-------- | :-------------- | :----------------------------------- |
| `r`       | `*http.Request` | **Required**. Handler request param. |
| `v`       | `interface{}`   | **Required**. Pointer to variable.   |

error will be returned if binding fails

#### Render responses based on request `r` headers

```go
  func Render(w http.ResponseWriter, r *http.Request, v interface{}, params ...interface{})
```

| Parameter | Type                  | Description                             |
| :-------- | :-------------------- | :-------------------------------------- |
| `w`       | `http.WriterResponse` | **Required**. Writer.                   |
| `r`       | `*http.Request`       | **Required**. Handler request param.    |
| `v`       | `interface{}`         | **Required**. Pointer to variable.      |
| `params`  | `...interface{}`      | Variadic number of params. (int/string) |

#### Render error response and status code based on request `r` headers

```go
  func Error(w http.ResponseWriter, r *http.Request, err error, params ...interface{})
```

| Parameter | Type                  | Description                             |
| :-------- | :-------------------- | :-------------------------------------- |
| `w`       | `http.WriterResponse` | **Required**. Writer.                   |
| `r`       | `*http.Request`       | **Required**. Handler request param.    |
| `err`     | `error`.              | **Required**. Error value.              |
| `params`  | `...interface{}`      | Variadic number of params. (int/string) |

#### Params variadic function parameter

`params` can be found in almost any function. When param type is string then API
will set header (key, value) pair or in case of int it will set status code.

```go
render.Render(w, v, http.StatusOK, "Content-Type", "application/json")
// or you can use const ContentTypeHeader, ApplicationJSON
render.Render(w, v, http.StatusOK, render.ContentTypeHeader, render.ApplicationJSON)
render.Render(w, v, "Content-Type", "application/json")
render.Render(w, v, "Content-Type", "application/json", http.StatusOK)
```

#### Integrate 3rd party JSON/XML lib

in this example we will replace standard encoder with goccy/go-json.

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"

	"github.com/enverbisevac/render"
)

func init() {
	render.JSONEncoder = func(w io.Writer) render.Encoder {
		return json.NewEncoder(w)
	}
}
```

#### Other functions

```go
func Blob(w http.ResponseWriter, v []byte, params ...interface{})
func PlainText(w http.ResponseWriter, v string, args ...interface{})
func HTML(w http.ResponseWriter, v string, args ...interface{})
func JSON(w http.ResponseWriter, v interface{}, args ...interface{})
func XML(w http.ResponseWriter, v interface{}, args ...interface{})
func File(w http.ResponseWriter, r *http.Request, fullPath string)
func Attachment(w http.ResponseWriter, r *http.Request, fullPath string)
func Inline(w http.ResponseWriter, r *http.Request, fullPath string)
func NoContent(w http.ResponseWriter)
func Stream(w http.ResponseWriter, r *http.Request, v interface{})
```

more help on API can be found in Documentation.

## Documentation

[Documentation](https://linktodocumentation)

## FAQ

#### Pass status code or response header in Render or Error function

Last parameter is variadic parameter in `Render()` or `Error()` function. Then we can
use

```go
render.Render(w, r, object, http.StatusAccepted)
// or
render.Render(w, r, object, http.StatusAccepted, "Content-Type", "application/json")
// or
render.Render(w, r, object, "Content-Type", "application/json")
```

#### Register global error/status codes

```go
render.ErrorMap[ErrConflict] = http.StatusConflict
```

we can even wrap the error:

```go
render.Error(w, r, fmt.Errorf("file %s %w", "demo.txt", render.ErrConflict))
```

#### Inline error rendering with status code

```go
render.Error(w, r, errors.New("some error"), http.StatusBadRequest)
```

or

```go
render.Error(w, r, &render.HTTPError{
    Err: errors.New("some error"),
    Status: http.StatusBadRequest,
})
```

#### Customize error response

```go
type CustomError struct {
    Module  string `json:"module"`
    Message string `json:"message"`
    Version string `json:"version"`
}

func (e CustomError) Error() string {
    return e.Message
}

render.TreatError = func(r *http.Request, err error) interface{} {
	cerr := &CustomError{}
	if errors.As(err, &cerr) {
		return cerr
    }
    // always return default one
    return render.DefaultErrorRespond(err)
}
```

## Running Tests

To run tests, run the following command

```bash
  make test
```

## Acknowledgements

- [go-chi/render](https://github.com/go-chi/render)

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Feedback

If you have any feedback, please reach out enver[@]bisevac.com
