package handlers

import (
	"net/http"

	"forum/config"
	"forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: No session token", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value

	email := database.GetUserEmailBySession(sessionToken)

	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "post.html", email)
}
