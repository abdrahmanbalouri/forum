package handlers

import (
	"net/http"
	"os"
)

func StaticHnadler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Oops!?, Method Not Allowed try again", http.StatusMethodNotAllowed)
	}
	url := r.URL.Path[1:]
	file, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Oops!?, Error Not Found", http.StatusNotFound)
		}
		http.Error(w, "Oops!?, Internal Server Error", http.StatusInternalServerError)
	}
	if file.IsDir() {
		http.Error(w, "Oops!?, Error Not Found", http.StatusNotFound)
	}
	http.ServeFile(w, r, url)
}

