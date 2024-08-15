package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field should not be empty")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Max length exceeded (max %d", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "Input does not match the format")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("Please provide input with minimum length (min %d)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Input does not match the format")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) RequiredAtLeastOne(fields ...string) {
	var count int

	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			count++
		}
	}

	if count == len(fields) {
		f.Errors.Add(fields[0], "This field should not be empty")
	}
}

func (f *Form) ProvidedAtLeastOne(fields ...string) bool {
	var count int

	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			count++
		}
	}

	if count < len(fields) {
		return true
	}

	return false
}
