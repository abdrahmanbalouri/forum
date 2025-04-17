package handlers

import (
	"encoding/json"
	"net/http"

	"forum/database"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := database.AddLike(1, 1)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Post liked successfully"))
}

func LikesCountHandler(w http.ResponseWriter, r *http.Request) {
	count := database.GetLikesCount()
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func UnlikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := database.RemoveLike(1, 1)
	if err != nil {
		http.Error(w, "Failed to unlike post", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Post unliked successfully"))
}
