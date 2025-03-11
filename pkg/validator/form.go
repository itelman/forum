package validator

import "net/url"

type Form struct {
	url.Values
	Errors Errors
}

func NewForm(data url.Values, errors Errors) *Form {
	return &Form{data, errors}
}
