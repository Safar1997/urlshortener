package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// вспомогательная функция для создания роутера с текущей map[string]string
func setupRouter(urlMap map[string]string) http.Handler {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		postHandler(w, r, urlMap)
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		getHandler(w, r, urlMap)
	})

	return r
}

func TestHandlers(t *testing.T) {
	urlMap := make(map[string]string)
	router := setupRouter(urlMap)

	t.Run("POST / — успешное создание новой ссылки", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusCreated, res.StatusCode)
		require.Equal(t, "text/plain", res.Header.Get("Content-Type"))

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(string(body), "http://localhost:8080/"))
	})

	t.Run("POST / — повторный запрос с тем же URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusCreated, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(string(body), "http://localhost:8080/"))
	})

	t.Run("GET /{id} — успешный редирект", func(t *testing.T) {
		id := urlMap["https://example.com"]
		req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, "https://example.com", res.Header.Get("Location"))
	})

	t.Run("GET /wrongid — несуществующий ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wrongid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "ID не найден\n", string(body))
	})

	t.Run("GET / — пустой ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		// chi по умолчанию отдаст 404
		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	})
}
