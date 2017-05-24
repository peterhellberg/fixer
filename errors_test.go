package fixer

import (
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	for _, msg := range []string{"", "foo", "bar"} {
		if got := NewError(msg).Error(); got != msg {
			t.Fatalf("NewError(%q).Error() = %q, want %q", msg, got, msg)
		}
	}
}

func TestResponseError(t *testing.T) {
	for _, tt := range []struct {
		resp *http.Response
		want error
	}{
		{nil, ErrNilResponse},
		{&http.Response{}, ErrUnexpectedStatus},
		{&http.Response{StatusCode: 200}, nil},
		{&http.Response{StatusCode: 404}, ErrNotFound},
		{&http.Response{StatusCode: 422}, ErrUnprocessableEntity},
	} {
		if got := responseError(tt.resp); got != tt.want {
			t.Fatalf("responseError(tt.resp) = %v, want %v", got, tt.want)
		}
	}
}
