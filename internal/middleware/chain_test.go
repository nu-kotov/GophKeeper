package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_Order(t *testing.T) {
	calls := []string{}

	m1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m1 before")
			next(w, r)
			calls = append(calls, "m1 after")
		}
	}

	m2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "m2 before")
			next(w, r)
			calls = append(calls, "m2 after")
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}

	chained := Chain(m1, m2)(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chained(rec, req)

	expected := []string{
		"m1 before",
		"m2 before",
		"handler",
		"m2 after",
		"m1 after",
	}

	assert.Equal(t, expected, calls, "error middlewares order")
}

func TestChain_SingleMiddleware(t *testing.T) {
	called := false

	m := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			called = true
			next(w, r)
		}
	}

	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}

	chained := Chain(m)(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chained(rec, req)

	assert.True(t, called, "middleware must be called")
	assert.True(t, handlerCalled, "handler must be called")
}

func TestChain_NoMiddleware(t *testing.T) {
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}

	chained := Chain()(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chained(rec, req)

	assert.True(t, handlerCalled, "handler must be called, without middleware")
}

func TestChain_MiddlewareStopsChain(t *testing.T) {
	called := []string{}

	m1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			called = append(called, "m1 before")
		}
	}

	m2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			called = append(called, "m2 before")
			next(w, r)
			called = append(called, "m2 after")
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		called = append(called, "handler")
	}

	chained := Chain(m1, m2)(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	chained(rec, req)

	assert.Equal(t, []string{"m1 before"}, called, "m1 must stop chain")
}
