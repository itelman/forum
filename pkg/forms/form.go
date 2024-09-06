package forms

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX    = regexp.MustCompile(`^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
	NameRX     = regexp.MustCompile(`^[a-zA-Z]{5,}([._-]{0,1}[a-zA-Z0-9]{2,})*$`)
	PasswordRX = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

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
			f.Errors.Add(field, "This field must not be empty")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Max length exceeded (max %d)", d))
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
	f.Errors.Add(field, "Please provide input that matches the requested format")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("Please provide input with minimum %d characters", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Please provide input that matches the requested format")
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
		f.Errors.Add(fields[0], "This field must not be empty")
	}
}

func (f *Form) ProvidedAtLeastOne(fields ...string) bool {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) != "" {
			return true
		}
	}

	return false
}

func (f *Form) ImgMaxSize(handler *multipart.FileHeader, d int) {
	if handler.Size > int64(d) {
		f.Errors.Add("image", fmt.Sprintf("Max size exceeded (max %d MB)", d/1048576))
	}
}

func (f *Form) ImgExtension(handler *multipart.FileHeader, exts ...string) {
	for _, ext := range exts {
		if filepath.Ext(handler.Filename) == ext {
			return
		}
	}

	f.Errors.Add("image", fmt.Sprintf("File extension should be one of the following: %s", exts))
}
