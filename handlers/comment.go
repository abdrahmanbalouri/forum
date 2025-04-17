package handlers

import (
	// "database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/database"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	// Get user ID from session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, exists := database.GetUserBySession(sessionCookie.Value)
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Insert comment into database
	stmt, err := database.DB.Prepare(`
		INSERT INTO comments (post_id, user_id, content, created_at) 
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, userID, content, time.Now())
	if err != nil {
		fmt.Println("Error executing statement:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
