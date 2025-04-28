package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Safar1997/urlshortener/internal/app"
)

func handler(w http.ResponseWriter, r *http.Request, urlMap map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "по этому пути только POST-запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	val, exists := urlMap[string(body)]

	if exists {
		fmt.Println("Найдено:", val)
	} else {
		fmt.Println("Ключа нет")
	}

	id := app.GenerateID()

	// Метод
	fmt.Println("Метод запроса:", r.Method)

	// Путь и параметры
	fmt.Println("Путь:", r.URL.Path)
	fmt.Println("Параметры:", r.URL.Query())

	// Заголовки
	fmt.Println("User-Agent:", r.Header.Get("User-Agent"))

	// Тело запроса
}

func main() {
	urlMap := make(map[string]string) // Пустая карта

	mux := http.NewServeMux()
	// mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, urlMap)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

// как запустить из iter1 директории
// go run cmd/shortener/main.go
