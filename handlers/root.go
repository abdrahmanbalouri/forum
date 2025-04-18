package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/config"
	"forum/database"
	"forum/models"
)

type Post struct {
	PostID       int
	Content      string
	Interest     string
	Username     string
	CreatedAt    string
	Likes        int
	Dislikes     int
	UserReaction string
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]interface{}{"Authenticated": false})
		return
	}

	userID, exists := database.GetUserBySession(sessionCookie.Value)
	if !exists {
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]interface{}{"Authenticated": false})
		return
	}

	user := database.GetUserInfo(userID)

	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		interest := r.FormValue("interest")

		if content != "" {
			CreatePost(userID, content, interest)
		}
	}

	// Récupérer tous les posts
	posts, err := GetAllPosts(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	var postsWithComments []map[string]interface{}
	for _, post := range posts {
		comments, err := GetCommentsForPost(post.PostID, userID)
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
		"Authenticated": true,
		"Username":      user.Username,
		"Posts":         postsWithComments,
	}

	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", templateData)
}

func CreatePost(userID int, content, interest string) {
	stmt, err := database.DB.Prepare("INSERT INTO posts(user_id, content, interest) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, content, interest)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	fmt.Println("Post created:", content, interest)
}

func GetAllPosts(userID int) ([]Post, error) {
	rows, err := database.DB.Query(`
		SELECT 
			p.id, 
			u.username, 
			p.content, 
			p.interest, 
			p.created_at,
			(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'like') as likes,
			(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'dislike') as dislikes,
			(SELECT reaction_type FROM post_reactions WHERE post_id = p.id AND user_id = ?) as user_reaction
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var userReaction sql.NullString
		err := rows.Scan(
			&post.PostID,
			&post.Username,
			&post.Content,
			&post.Interest,
			&post.CreatedAt,
			&post.Likes,
			&post.Dislikes,
			&userReaction,
		)
		if err != nil {
			continue
		}
		post.UserReaction = userReaction.String
		posts = append(posts, post)
	}
	return posts, nil
}

func GetCommentsForPost(postID int, currentUserID int) ([]models.Comment, error) {
	rows, err := database.DB.Query(`
		SELECT 
			c.id,
			u.username, 
			c.content, 
			c.created_at,
			(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.id AND reaction_type = 'like') as like_count,
			(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.id AND reaction_type = 'dislike') as dislike_count,
			(SELECT reaction_type FROM comment_reactions WHERE comment_id = c.id AND user_id = ?) as user_reaction
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, currentUserID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var userReaction sql.NullString
		err := rows.Scan(
			&comment.ID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
			&comment.LikeCount,
			&comment.DislikeCount,
			&userReaction,
		)
		if err != nil {
			continue
		}
		comment.UserReaction = userReaction.String
		comments = append(comments, comment)
	}
	return comments, nil
}
