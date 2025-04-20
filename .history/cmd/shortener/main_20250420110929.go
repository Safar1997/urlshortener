package main

import (
	"fmt"
	"io"
	"net/http"
)

// func mainPage(res http.ResponseWriter, req *http.Request) {
// 	res.Write([]byte("Привет!"))
// }

// func apiPage(res http.ResponseWriter, req *http.Request) {
// 	res.Write([]byte("Это страница /api."))
// }

func handler(w http.ResponseWriter, r *http.Request) {
	// Метод
	fmt.Println("Метод запроса:", r.Method)

	// Путь и параметры
	fmt.Println("Путь:", r.URL.Path)
	fmt.Println("Параметры:", r.URL.Query())

	// Заголовки
	fmt.Println("User-Agent:", r.Header.Get("User-Agent"))

	// Тело запроса
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println("Тело запроса:", string(body))
}

func main() {
	mux := http.NewServeMux()
	// mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc(`/`, handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
