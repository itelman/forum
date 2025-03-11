package validator

import "fmt"

var (
	ErrInputUnchanged = "You haven't made any changes"
)

func ErrInputRequired(field string) string {
	return fmt.Sprintf("Please enter a valid %s", field)
}

func ErrInputLength(min, max int) string {
	return fmt.Sprintf("Please provide input with length of %d-%d characters", min, max)
}
