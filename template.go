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
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type engine interface {
	executeTemplate(w io.Writer, name string, v interface{}) error
	execute(w io.Writer, v interface{}) error
	parse(tmpl string) (engine, error)
	funcs(funcs htmltemplate.FuncMap) engine
	set(v interface{})
}

type templateWrapper struct {
	text *template.Template
	html *htmltemplate.Template
}

func newTemplateWrapper(name string) *templateWrapper {
	if name == "text" {
		return &templateWrapper{
			text: template.New(name),
		}
	}
	return &templateWrapper{
		html: htmltemplate.New(name),
	}
}

func (t *templateWrapper) executeTemplate(w io.Writer, name string, v interface{}) error {
	if t.text != nil {
		return t.text.ExecuteTemplate(w, name, v)
	}
	return t.html.ExecuteTemplate(w, name, v)
}

func (t *templateWrapper) execute(w io.Writer, v interface{}) error {
	if t.text != nil {
		return t.text.Execute(w, v)
	}
	return t.html.Execute(w, v)
}

func (t *templateWrapper) parse(tmpl string) (engine, error) {
	if t.text != nil {
		parse, err := t.text.Parse(tmpl)
		if err != nil {
			return nil, err
		}
		t.text = parse
		return t, nil
	}

	parse, err := t.html.Parse(tmpl)
	if err != nil {
		return nil, err
	}
	t.html = parse
	return t, nil
}

func (t *templateWrapper) funcs(funcs htmltemplate.FuncMap) engine {
	if t.text != nil {
		t.text = t.text.Funcs(funcs)
		return t
	}

	t.html = t.html.Funcs(funcs)
	return t
}

func (t *templateWrapper) set(v interface{}) {
	switch value := v.(type) {
	case *template.Template:
		t.text = value
	case *htmltemplate.Template:
		t.html = value
	}
}

func templateFactory(w http.ResponseWriter, factory engine, v interface{}, ct string, params ...interface{}) {
	var (
		t         engine
		tmpl      string
		buf       bytes.Buffer
		err       error
		newParams = make([]interface{}, 0, len(params))
	)
	switch value := v.(type) {
	case string:
		_, _ = buf.WriteString(value)
	case *string:
		if value != nil {
			_, _ = buf.WriteString(*value)
		}
	default:
		// check params for template input
		for _, param := range params {
			switch value := param.(type) {
			case string:
				if tmpl == "" {
					tmpl = value
				}
			case *string:
				if tmpl == "" {
					tmpl = *value
				}
			case *template.Template, *htmltemplate.Template:
				if t == nil {
					t.set(v)
				}
			default:
				newParams = append(newParams, value)
			}
		}
	}

	switch {
	case t != nil && tmpl != "":
		if strings.HasPrefix(tmpl, "tmpl://") {
			err = t.executeTemplate(&buf, strings.Replace(tmpl, "tmpl://", "", 1), v)
		} else if t, err = t.parse(tmpl); err == nil {
			err = t.execute(&buf, v)
		}
	case tmpl != "":
		if t, err = factory.funcs(TemplateFuncs).parse(tmpl); err == nil {
			err = t.execute(&buf, v)
		}
	default:
		_, _ = buf.WriteString(fmt.Sprintf("%v", v))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Blob(w, buf.Bytes(), append(newParams, ContentTypeHeader, ct)...)
}
