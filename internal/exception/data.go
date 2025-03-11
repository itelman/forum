package exception

import "net/http"

type errorData struct {
	Code int
	Text string
}

func newErrorData(code int) errorData {
	return errorData{Code: code, Text: http.StatusText(code)}
}

var (
	errBadRequestData      = newErrorData(http.StatusBadRequest)
	errUnauthorizedData    = newErrorData(http.StatusUnauthorized)
	errForbiddenData       = newErrorData(http.StatusForbidden)
	errNotFoundData        = newErrorData(http.StatusNotFound)
	errNotAllowedData      = newErrorData(http.StatusMethodNotAllowed)
	errTooManyRequestsData = newErrorData(http.StatusTooManyRequests)
	errInternalServerData  = newErrorData(http.StatusInternalServerError)
)
