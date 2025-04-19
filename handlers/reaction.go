package handlers

import (
	"database/sql"
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

	var table string
	if targetType == "post" {
		table = "post_reactions"
	} else if targetType == "comment" {
		table = "comment_reactions"
	} else {
		http.Error(w, "Invalid target type", http.StatusBadRequest)
		return
	}

	var existingReaction string
	err = database.DB.QueryRow(
		"SELECT reaction_type FROM "+table+" WHERE post_id = ? AND user_id = ?",
		targetID, userID,
	).Scan(&existingReaction)

	if err != nil {
		if err == sql.ErrNoRows {
			// No existing reaction, insert the new reaction
			_, err = database.DB.Exec(
				"INSERT INTO "+table+" (post_id, user_id, reaction_type) VALUES (?, ?, ?)",
				targetID, userID, reactionType,
			)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		if existingReaction == reactionType {
			// If the same reaction exists, remove it
			_, err = database.DB.Exec(
				"DELETE FROM "+table+" WHERE post_id = ? AND user_id = ?",
				targetID, userID,
			)
		} else {
			// Update to the new reaction type
			_, err = database.DB.Exec(
				"UPDATE "+table+" SET reaction_type = ? WHERE post_id = ? AND user_id = ?",
				reactionType, targetID, userID,
			)
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
