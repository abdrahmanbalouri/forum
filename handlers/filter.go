package handlers

import (
	"net/http"

	"forum/config"
	"forum/database"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: No session token", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value
	userID, _ := database.GetUserBySession(sessionToken)

	interest := r.URL.Query().Get("interest")
	if interest == "" {
		http.Error(w, "Interest is required", http.StatusBadRequest)
		return
	}

	posts, err := GetPostsByInterest(interest, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	var postsWithComments []map[string]interface{}
	for _, post := range posts {
		comments, err := GetCommentsForPost(post.PostID, userID) // Pass 0 for unauthenticated users
		if err != nil {
			continue
		}
		postData := map[string]interface{}{
			"Post":     post,
			"Comments": comments,
		}
		postsWithComments = append(postsWithComments, postData)
	}


	templateData := map[string]interface{}{
		"Posts": postsWithComments,
	}

	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "filter.html", templateData)
}
