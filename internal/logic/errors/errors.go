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
	ErrInternalServerError = &AppError{HTTPStatus: http.StatusInternalServerError, Slug: "error_internal"}
	ErrCVServer            = &AppError{HTTPStatus: http.StatusInternalServerError, Slug: "error_cv_server"}
	ErrServiceUnavailable  = &AppError{HTTPStatus: http.StatusServiceUnavailable, Slug: "error_service_unavailable"}
	ErrBadGateway          = &AppError{HTTPStatus: http.StatusBadGateway, Slug: "error_bad_gateway"}
	ErrMethodNotAllowed    = &AppError{HTTPStatus: http.StatusMethodNotAllowed, Slug: "error_method_not_allowed"}
	ErrInvalidInput        = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_invalid_input"}
	ErrUnsupportedLanguage = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_unsupported_language"}
	ErrInvalidPassword     = &AppError{HTTPStatus: http.StatusUnauthorized, Slug: "error_cv_auth"}
	ErrCVNotFound          = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_cv_not_found"}
	ErrCVExpired           = &AppError{HTTPStatus: http.StatusGone, Slug: "error_cv_expired"}
	ErrContentNotFound     = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_content_not_found"}
	ErrRouteNotFound       = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_route_not_found"}

	ErrPowSignature     = &AppError{HTTPStatus: http.StatusForbidden, Slug: "error_pow_signature"}
	ErrPowDoubleSpend   = &AppError{HTTPStatus: http.StatusConflict, Slug: "error_pow_double_spend"}
	ErrPowExpired       = &AppError{HTTPStatus: http.StatusGone, Slug: "error_pow_expired"}
	ErrPowWork          = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_pow_work"}
	ErrCaptchaNotFound  = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_captcha_not_found"}
	ErrCaptchaInvalid   = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_captcha_invalid"}
	ErrCaptchaNotSolved = &AppError{HTTPStatus: http.StatusForbidden, Slug: "error_captcha_not_solved"}
	ErrCaptchaExpired   = &AppError{HTTPStatus: http.StatusGone, Slug: "error_captcha_expired"}
	ErrCaptchaServer    = &AppError{HTTPStatus: http.StatusInternalServerError, Slug: "error_captcha_server"}
)

var FallbackTranslations = map[string]map[string]string{
	"pl": {
		"error_title":                "Wystąpił błąd",
		"error_route_title":          "Nie znaleziono strony",
		"error_route_not_found":      "Podany adres nie prowadzi do żadnej strony.",
		"error_home":                 "Strona główna",
		"error_header":               "BŁĄD",
		"error_status_key":           "status",
		"error_message_key":          "wiadomość",
		"error_language_label":       "Język strony",
		"error_cv_server":            "Wystąpił nieoczekiwany błąd serwera.",
		"error_internal":             "Wystąpił nieoczekiwany błąd serwera.",
		"error_message":              "Serwis jest chwilowo niedostępny. Spróbuj ponownie później.",
		"error_service_unavailable":  "Serwis jest chwilowo niedostępny. Spróbuj ponownie później.",
		"error_bad_gateway":          "Nie udało się uzyskać poprawnej odpowiedzi usługi. Spróbuj ponownie później.",
		"error_method_not_allowed":   "Ta metoda żądania nie jest obsługiwana.",
		"error_invalid_input":        "Podane dane są nieprawidłowe.",
		"error_unsupported_language": "Wybrany język nie jest obsługiwany.",
		"error_content_not_found":    "Nie znaleziono żądanej treści.",
		"error_captcha_server":       "Wystąpił błąd serwera CAPTCHA.",
		"error_captcha_not_solved":   "CAPTCHA nie została rozwiązana.",
		"error_cv_expired":           "Link do CV wygasł lub jest nieprawidłowy.",
		"error_cv_auth":              "Odmowa dostępu. Nieprawidłowe hasło.",
		"error_cv_not_found":         "Plik CV nie został znaleziony.",
		"error_pow_failed":           "Błąd autoryzacji wstępnej. Spróbuj ponownie.",
		"error_pow_signature":        "Nieprawidłowa sygnatura żądania.",
		"error_pow_work":             "Błędne obliczenia PoW.",
		"error_pow_double_spend":     "To żądanie zostało już wykorzystane.",
		"error_pow_expired":          "Żądanie wygasło. Spróbuj ponownie.",
		"error_captcha_invalid":      "Nieprawidłowy kod CAPTCHA.",
		"error_captcha_expired":      "Sesja CAPTCHA wygasła. Pobieranie nowej...",
		"error_captcha_not_found":    "Nie znaleziono sesji CAPTCHA.",
	},
	"en": {
		"error_title":                "An error occurred",
		"error_route_title":          "Page not found",
		"error_route_not_found":      "There is no page at this address.",
		"error_home":                 "Home page",
		"error_header":               "ERROR",
		"error_status_key":           "status",
		"error_message_key":          "message",
		"error_language_label":       "Page language",
		"error_cv_server":            "An unexpected server error occurred.",
		"error_internal":             "An unexpected server error occurred.",
		"error_message":              "Service is temporarily unavailable. Please try again later.",
		"error_service_unavailable":  "Service is temporarily unavailable. Please try again later.",
		"error_bad_gateway":          "The service returned an invalid response. Please try again later.",
		"error_method_not_allowed":   "This request method is not supported.",
		"error_invalid_input":        "The supplied data is invalid.",
		"error_unsupported_language": "The selected language is not supported.",
		"error_content_not_found":    "The requested content was not found.",
		"error_captcha_server":       "A CAPTCHA server error occurred.",
		"error_captcha_not_solved":   "The CAPTCHA has not been solved.",
		"error_cv_expired":           "CV link has expired or is invalid.",
		"error_cv_auth":              "Access denied. Invalid password.",
		"error_cv_not_found":         "CV file was not found.",
		"error_pow_failed":           "Pre-authorization error. Please try again.",
		"error_pow_signature":        "Invalid request signature.",
		"error_pow_work":             "Incorrect PoW calculation.",
		"error_pow_double_spend":     "This request has already been used.",
		"error_pow_expired":          "The request expired. Please try again.",
		"error_captcha_invalid":      "Invalid CAPTCHA code.",
		"error_captcha_expired":      "CAPTCHA session expired. Fetching new one...",
		"error_captcha_not_found":    "CAPTCHA session not found.",
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
	case "error_route_not_found":
		return ErrRouteNotFound
	case "error_cv_server":
		return ErrCVServer
	case "error_internal":
		return ErrInternalServerError
	case "error_captcha_server":
		return ErrCaptchaServer
	case "error_message":
		return ErrServiceUnavailable
	case "error_service_unavailable":
		return ErrServiceUnavailable
	case "error_bad_gateway":
		return ErrBadGateway
	case "error_method_not_allowed":
		return ErrMethodNotAllowed
	case "error_invalid_input":
		return ErrInvalidInput
	case "error_unsupported_language":
		return ErrUnsupportedLanguage
	case "error_content_not_found":
		return ErrContentNotFound
	case "error_pow_signature":
		return ErrPowSignature
	case "error_pow_double_spend":
		return ErrPowDoubleSpend
	case "error_pow_expired":
		return ErrPowExpired
	case "error_pow_work":
		return ErrPowWork
	case "error_captcha_not_found":
		return ErrCaptchaNotFound
	case "error_captcha_invalid":
		return ErrCaptchaInvalid
	case "error_captcha_not_solved":
		return ErrCaptchaNotSolved
	case "error_captcha_expired":
		return ErrCaptchaExpired
	default:
		return ErrInternalServerError
	}
}

func FromUpstream(status int, slug string) *AppError {
	if slug == "error_message" {
		switch status {
		case http.StatusBadRequest:
			return ErrInvalidInput
		case http.StatusMethodNotAllowed:
			return ErrMethodNotAllowed
		case http.StatusNotFound:
			return ErrContentNotFound
		case http.StatusServiceUnavailable:
			return ErrServiceUnavailable
		}
	}
	appErr := FromSlug(slug)
	if slug == appErr.Slug && status == appErr.HTTPStatus {
		return appErr
	}
	return ErrBadGateway
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
