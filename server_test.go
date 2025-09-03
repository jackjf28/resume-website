package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootEndpoint(t *testing.T) {
	t.Run("/ returns 404 not found", func(t *testing.T) {
		server := NewServer(context.Background())
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestBaseApi(t *testing.T) {
	server := NewServer(context.Background())
	t.Run("/api/v1 returns 200", func(t *testing.T) {
		request := newGetBaseApiRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
	t.Run("/api/v1 returns json when called", func(t *testing.T) {
		request := newGetBaseApiRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		assertContentType(t, response, jsonContentType)
	})
}

func newGetBaseApiRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/api/v1", nil)
	return request
}

func assertContentType(t testing.TB, got *httptest.ResponseRecorder, want string) {
	t.Helper()
	actual := got.Result().Header.Get("content-type")
	if actual != want {
		t.Errorf("wanted content-type %q, got content-type %q", want, actual)
	}
}

func assertStatus(t testing.TB, got int, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
