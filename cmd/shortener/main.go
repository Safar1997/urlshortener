package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Safar1997/urlshortener/internal/app"
	"github.com/go-chi/chi/v5"
)

func postHandler(w http.ResponseWriter, r *http.Request, urlMap map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "неверный путь", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)
	val, exists := urlMap[originalURL]

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated) // 201

	if exists {
		w.Write([]byte("http://localhost:8080/" + val))
	} else {
		id := app.GenerateID()
		urlMap[originalURL] = id
		w.Write([]byte("http://localhost:8080/" + id))
	}
}

func getHandler(w http.ResponseWriter, r *http.Request, urlMap map[string]string) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "не указан ID", http.StatusBadRequest)
		return
	}
	// Проверяем, существует ли ID в urlMap
	for originalURL, shortID := range urlMap {
		if shortID == id {
			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect) // 307
			return
		}
	}

	http.Error(w, "ID не найден", http.StatusBadRequest)
}

func main() {
	urlMap := make(map[string]string)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/" {
			postHandler(w, r, urlMap)
		} else {
			getHandler(w, r, urlMap)
		}
	})

	// создаём новый роутер chi
	r := chi.NewRouter()

	// добавляем простой обработчик на GET /
	r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		postHandler(rw, r, urlMap)
	})

	// добавляем обработчик с параметром id
	r.Get("/{id}", func(rw http.ResponseWriter, r *http.Request) {
		getHandler(rw, r, urlMap)
	})

	// запускаем HTTP-сервер на порту 8080
	fmt.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", r)
}
