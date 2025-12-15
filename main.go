package main

import (
	"URL-Shortener-App/utils"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {

	})

	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {
		url := req.FormValue("url")
		shortCode := utils.GetShortCode()
		shortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortCode)
		fmt.Printf("Generated short URL: %s\n", shortURL)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
