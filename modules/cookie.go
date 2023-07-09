package modules

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func WriteCookie(w http.ResponseWriter, cookieName, cookieValue string) error {

	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   cookieValue,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}

	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, cookie)

	return nil
}

func ReadCookie(r *http.Request, cookieName string) (string, error) {

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}
