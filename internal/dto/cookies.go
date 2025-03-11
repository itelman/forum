package dto

import (
	"net/http"
	"time"
)

var (
	FlashSignupSuccessful = "Signup successful! Please log in."
	FlashSessionExpired   = "Your session has expired. Please sign in again."
	FlashPostRemoved      = "Post successfully removed!"
	FlashCommentEnter     = "Please enter a valid comment."
	FlashFilterSelect     = "Please select at least one filter."
)

func NewCookie(name, val string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

func DeleteCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}
