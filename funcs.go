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
	"html/template"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	day  = 24 * time.Hour
	year = 365 * day
)

var printer = message.NewPrinter(language.English)

// TemplateFuncs is map of basic functions to use in templates
var TemplateFuncs = template.FuncMap{
	// Time functions
	"now":            time.Now,
	"timeSince":      time.Since,
	"timeUntil":      time.Until,
	"formatTime":     formatTime,
	"approxDuration": approxDuration,

	// String functions
	"uppercase": strings.ToUpper,
	"lowercase": strings.ToLower,
	"pluralize": pluralize,
	"slugify":   slugify,
	"safeHTML":  safeHTML,

	// Slice functions
	"join": strings.Join,
	// "containsString": slices.Contains[string],

	// Number functions
	"incr":        incr,
	"decr":        decr,
	"formatInt":   formatInt,
	"formatFloat": formatFloat,

	// Boolean functions
	"yesno": yesno,

	// URL functions
	"urlSetParam": urlSetParam,
	"urlDelParam": urlDelParam,
}

func formatTime(format string, t time.Time) string {
	return t.Format(format)
}

func approxDuration(d time.Duration) string {
	if d < time.Second {
		return "less than 1 second"
	}

	ds := int(math.Round(d.Seconds()))
	if ds == 1 {
		return "1 second"
	} else if ds < 60 {
		return fmt.Sprintf("%d seconds", ds)
	}

	dm := int(math.Round(d.Minutes()))
	if dm == 1 {
		return "1 minute"
	} else if dm < 60 {
		return fmt.Sprintf("%d minutes", dm)
	}

	dh := int(math.Round(d.Hours()))
	if dh == 1 {
		return "1 hour"
	} else if dh < 24 {
		return fmt.Sprintf("%d hours", dh)
	}

	dd := int(math.Round(float64(d / day)))
	if dd == 1 {
		return "1 day"
	} else if dd < 365 {
		return fmt.Sprintf("%d days", dd)
	}

	dy := int(math.Round(float64(d / year)))
	if dy == 1 {
		return "1 year"
	}

	return fmt.Sprintf("%d years", dy)
}

func pluralize(count interface{}, singular string, plural string) (string, error) {
	n, err := toInt64(count)
	if err != nil {
		return "", err
	}

	if n == 1 {
		return singular, nil
	}

	return plural, nil
}

func slugify(s string) string {
	var buf bytes.Buffer

	for _, r := range s {
		switch {
		case r > unicode.MaxASCII:
			continue
		case unicode.IsLetter(r):
			buf.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r), r == '_', r == '-':
			buf.WriteRune(r)
		case unicode.IsSpace(r):
			buf.WriteRune('-')
		}
	}

	return buf.String()
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func incr(i interface{}) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n++
	return n, nil
}

func decr(i interface{}) (int64, error) {
	n, err := toInt64(i)
	if err != nil {
		return 0, err
	}

	n--
	return n, nil
}

func formatInt(i interface{}) (string, error) {
	n, err := toInt64(i)
	if err != nil {
		return "", err
	}

	return printer.Sprintf("%d", n), nil
}

func formatFloat(f float64, dp int) string {
	format := "%." + strconv.Itoa(dp) + "f"
	return printer.Sprintf(format, f)
}

func yesno(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

func urlSetParam(u *url.URL, key string, value interface{}) *url.URL {
	nu := *u
	values := nu.Query()

	values.Set(key, fmt.Sprintf("%v", value))

	nu.RawQuery = values.Encode()
	return &nu
}

func urlDelParam(u *url.URL, key string) *url.URL {
	nu := *u
	values := nu.Query()

	values.Del(key)

	nu.RawQuery = values.Encode()
	return &nu
}

func toInt64(i interface{}) (int64, error) {
	switch v := i.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	// Note: uint64 not supported due to risk of truncation.
	case string:
		return strconv.ParseInt(v, 10, 64)
	}

	return 0, fmt.Errorf("unable to convert type %T to int", i)
}
