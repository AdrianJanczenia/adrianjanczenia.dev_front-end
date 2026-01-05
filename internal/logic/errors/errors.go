package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AppError struct {
	HTTPStatus int
	Slug       string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Slug
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrInternalServerError = &AppError{HTTPStatus: http.StatusInternalServerError, Slug: "error_cv_server"}
	ErrServiceUnavailable  = &AppError{HTTPStatus: http.StatusServiceUnavailable, Slug: "error_message"}
	ErrMethodNotAllowed    = &AppError{HTTPStatus: http.StatusMethodNotAllowed, Slug: "error_message"}
	ErrInvalidInput        = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_message"}
	ErrUnsupportedLanguage = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_message"}
	ErrInvalidPassword     = &AppError{HTTPStatus: http.StatusUnauthorized, Slug: "error_cv_auth"}
	ErrCVNotFound          = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_cv_not_found"}
	ErrCVExpired           = &AppError{HTTPStatus: http.StatusGone, Slug: "error_cv_expired"}
	ErrContentNotFound     = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_message"}
)

var FallbackTranslations = map[string]map[string]string{
	"pl": {
		"error_cv_server":    "Wystąpił nieoczekiwany błąd serwera.",
		"error_message":      "Serwis jest chwilowo niedostępny. Spróbuj ponownie później.",
		"error_cv_expired":   "Link do CV wygasł lub jest nieprawidłowy.",
		"error_cv_auth":      "Odmowa dostępu. Nieprawidłowe hasło.",
		"error_cv_not_found": "Plik CV nie został znaleziony.",
	},
	"en": {
		"error_cv_server":    "An unexpected server error occurred.",
		"error_message":      "Service is temporarily unavailable. Please try again later.",
		"error_cv_expired":   "CV link has expired or is invalid.",
		"error_cv_auth":      "Access denied. Invalid password.",
		"error_cv_not_found": "CV file was not found.",
	},
}

func FromSlug(slug string) *AppError {
	switch slug {
	case "error_cv_auth":
		return ErrInvalidPassword
	case "error_cv_expired":
		return ErrCVExpired
	case "error_cv_not_found":
		return ErrCVNotFound
	case "error_cv_server":
		return ErrInternalServerError
	case "error_message":
		return ErrServiceUnavailable
	default:
		return ErrInternalServerError
	}
}

func FromHTTPStatus(status int) *AppError {
	switch status {
	case http.StatusBadRequest:
		return ErrInvalidInput
	case http.StatusUnauthorized:
		return ErrInvalidPassword
	case http.StatusNotFound:
		return ErrCVNotFound
	case http.StatusGone:
		return ErrCVExpired
	case http.StatusMethodNotAllowed:
		return ErrMethodNotAllowed
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return ErrInternalServerError
	}
}

func WriteJSON(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = ErrInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"error": appErr.Slug,
	})
}
