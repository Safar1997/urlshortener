package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		response    string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		urlMap map[string]string
		want   want
	}{
		{
			name:   "успешный POST-запрос — новая ссылка",
			method: http.MethodPost,
			path:   "/",
			body:   "https://example.com",
			urlMap: make(map[string]string),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				// ответ будет динамический (id генерируется внутри функции), проверим префикс
			},
		},
		{
			name:   "неправильный метод",
			method: http.MethodGet,
			path:   "/",
			body:   "",
			urlMap: make(map[string]string),
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				response:    "только POST-запросы разрешены\n",
			},
		},
		{
			name:   "неправильный путь",
			method: http.MethodPost,
			path:   "/wrongpath",
			body:   "",
			urlMap: make(map[string]string),
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    "неверный путь\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			postHandler(w, request, tt.urlMap)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			body := string(bodyBytes)

			// Проверка тела ответа
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, body)
			} else {
				// Если создаётся новая ссылка, проверяем что начинается на нужный префикс
				assert.True(t, strings.HasPrefix(body, "http://localhost:8080/"))
			}
		})
	}
}

func Test_getHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
		body       string
	}
	tests := []struct {
		name   string
		method string
		path   string
		urlMap map[string]string
		want   want
	}{
		{
			name:   "успешный GET-запрос — редирект",
			method: http.MethodGet,
			path:   "/abc123",
			urlMap: map[string]string{
				"https://example.com": "abc123",
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://example.com",
			},
		},
		{
			name:   "GET-запрос без ID",
			method: http.MethodGet,
			path:   "/",
			urlMap: map[string]string{
				"https://example.com": "abc123",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "не указан ID\n",
			},
		},
		{
			name:   "GET-запрос с несуществующим ID",
			method: http.MethodGet,
			path:   "/wrongid",
			urlMap: map[string]string{
				"https://example.com": "abc123",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "ID не найден\n",
			},
		},
		{
			name:   "неправильный метод",
			method: http.MethodPost,
			path:   "/abc123",
			urlMap: map[string]string{
				"https://example.com": "abc123",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       "только GET-запросы разрешены\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			getHandler(w, request, tt.urlMap)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			// Проверка заголовка Location, если это редирект
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, res.Header.Get("Location"))
			}

			// Проверка тела, если оно ожидается
			if tt.want.body != "" {
				bodyBytes, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				body := string(bodyBytes)
				assert.Equal(t, tt.want.body, body)
			}
		})
	}
}
