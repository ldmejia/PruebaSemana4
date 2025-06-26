package main

import (
	"net/http"
	"prueba4/handler"
)

func main() {
	http.HandleFunc("/users", handler.HandleUsers)
	http.ListenAndServe(":8080", nil)
}