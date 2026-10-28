package errors

import (
	"net/http"
	"testing"
)

func TestUpstreamErrorContract(t *testing.T) {
	for _, tt := range []struct {
		status int
		slug   string
		want   *AppError
	}{
		{http.StatusNotFound, "error_cv_not_found", ErrCVNotFound},
		{http.StatusNotFound, "error_route_not_found", ErrRouteNotFound},
		{http.StatusForbidden, "error_captcha_not_solved", ErrCaptchaNotSolved},
		{http.StatusBadGateway, "error_bad_gateway", ErrBadGateway},
		{http.StatusMethodNotAllowed, "error_message", ErrMethodNotAllowed},
		{http.StatusBadRequest, "error_message", ErrInvalidInput},
		{http.StatusGone, "error_cv_not_found", ErrBadGateway},
		{http.StatusNotFound, "", ErrBadGateway},
	} {
		if got := FromUpstream(tt.status, tt.slug); got != tt.want {
			t.Errorf("FromUpstream(%d, %q) = %v, want %v", tt.status, tt.slug, got, tt.want)
		}
	}
}
