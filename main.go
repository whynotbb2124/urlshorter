package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/shorturl", shorturlHandler)

	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}

type URLRequest struct {
	URL string `json:"url"`
}

func shorturlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Received URL:", req.URL)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "URL received:", req.URL)
}
