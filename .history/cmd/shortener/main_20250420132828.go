package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Safar1997/urlshortener/internal/app"
	"github.com/gorilla/mux"
)

func postHandler(w http.ResponseWriter, r *http.Request, urlMap map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "по этому пути только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	val, exists := urlMap[string(body)]

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated) // 201

	if exists {
		fmt.Println("Найдено:", val)
		w.Write([]byte(val))
	} else {
		id := app.GenerateID()
		urlMap[string(body)] = id
		w.Write([]byte(id))
		fmt.Println("Создано:", id)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request, urlMap map[string]string) {
	vars := mux.Vars(r)
	id := vars["id"]

	for originalURL, shortID := range urlMap {
		if shortID == id {
			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect) // 307
			fmt.Println("✅ Найдено, Location:", originalURL)
			return
		}
	}

	// Не найдено — ошибка 400
	http.Error(w, "неверный или несуществующий ID", http.StatusBadRequest)
}

func main() {
	urlMap := make(map[string]string)

	r := mux.NewRouter()

	// POST /
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		postHandler(w, r, urlMap)
	}).Methods("POST")

	// GET /{id}
	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		getHandler(w, r, urlMap)
	}).Methods("GET")

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
