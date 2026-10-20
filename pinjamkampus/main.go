// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}