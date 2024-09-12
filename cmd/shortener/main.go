package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ivanLFN/go-musthave-shortener-tpl.git/internal/randstring"
)

const length int = 10

var URLAlias = make(map[string]string)

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(handleRequest))
	if err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getURL(w, r)
	case http.MethodPost:
		postURL(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postURL(w http.ResponseWriter, r *http.Request) {
	var shortURL string
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	currentURLButes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	currentURL := string(currentURLButes)

	if URLAlias[currentURL] != "" {
		shortURL = URLAlias[currentURL]
	} else {
		shortURL = fmt.Sprintf(`/%s`, randstring.RandString(length))
		URLAlias[string(currentURL)] = shortURL
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(r.Host + shortURL))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	shortURL := r.URL.Query().Get("Location")

	if shortURL == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	URL, found := findKeyByValue(shortURL)
	if found {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(URL))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func findKeyByValue(shortURL string) (string, bool) {
	for k, v := range URLAlias {
		if v == shortURL {
			return k, true
		}
	}
	return "", false
}
