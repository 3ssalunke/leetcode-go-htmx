package util

import (
	"net/http"
)

func ExtractCookieFromHeader(r *http.Request) (string, error) {
	cookie, err := r.Cookie("leetcode_auth")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
