package handlers

import (
	"net/http"
	"strconv"

	"forum/database"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	reactionType := r.FormValue("type")
	if reactionType != "like" && reactionType != "dislike" {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	targetType := r.FormValue("target_type")
	targetID, err := strconv.Atoi(r.FormValue("target_id"))
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}

	var query string
	if targetType == "post" {
		query = `
			INSERT INTO post_reactions (post_id, user_id, reaction_type)
			VALUES (?, ?, ?)
			ON CONFLICT(post_id, user_id) DO UPDATE SET reaction_type = ?
		`
	} else if targetType == "comment" {
		query = `
			INSERT INTO comment_reactions (comment_id, user_id, reaction_type)
			VALUES (?, ?, ?)
			ON CONFLICT(comment_id, user_id) DO UPDATE SET reaction_type = ?
		`
	} else {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(query, targetID, userID, reactionType, reactionType)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
