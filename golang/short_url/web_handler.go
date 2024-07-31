package main

import (
	"encoding/json"
	"net/http"
)

var maxRetries = 3
var urlMapCache = make(map[string]string)
var urlSet = make(map[string]struct{})
var commonPrefix = "http://short.url/"

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Url string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = UrlValidator(input.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shorten := commonPrefix + GenerateByMurmurHash(input.Url)

	response := struct {
		ShortUrl string `json:"short_url"`
	}{
		ShortUrl: shorten,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
