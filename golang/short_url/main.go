package main

import "net/http"

func main() {

	http.HandleFunc("/shorten", ShortUrlHandler)
	http.ListenAndServe(":8080", nil)
}
