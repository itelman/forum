package validator

type Errors map[string]string

func (e Errors) Add(field, message string) {
	e[field] = message
}

func (e Errors) Get(field string) string {
	return e[field]
}
