package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	url := r.URL.Path[1:]
	file, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Error Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if file.IsDir() {
		http.Error(w, "Error Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, url)
}

func StaticHandler2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	url := r.URL.Path[1:]
	file, err := os.Stat(url)
	fmt.Println(url)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Error Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if file.IsDir() {
		http.Error(w, "Error Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, url)
}
