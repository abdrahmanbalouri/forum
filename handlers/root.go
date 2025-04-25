package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/config"
	"forum/database"
	"forum/models"
)

type Post struct {
	PostID               int
	Content              string
	Interest             string
	Title                string
	Username             string
	CreatedAt            string
	Likes                int
	Dislikes             int
	UserReaction         string
	MinutesSinceCreation int
	Photo                string
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		posts, err := GetAllPosts(0)
		if err != nil {
			http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		var postsWithComments []map[string]interface{}
		for _, post := range posts {
			comments, err := GetCommentsForPost(post.PostID, 0)
			if err != nil {
				continue
			}
			postData := map[string]interface{}{
				"Post":     post,
				"Comments": comments,
			}
			postsWithComments = append(postsWithComments, postData)
		}

		if len(postsWithComments) == 0 {
			http.Error(w, "No posts found for the selected interest.", http.StatusNotFound)
			return
		}

		templateData := map[string]interface{}{
			"Authenticated": false,
			"Posts":         postsWithComments,
		}
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", templateData)
		return
	}

	userID, exists := database.GetUserBySession(sessionCookie.Value)
	if !exists {
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]interface{}{"Authenticated": false})
		return
	}

	user := database.GetUserInfo(userID)

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		content := r.FormValue("content")
		interests := r.Form["interest"]
		title := r.FormValue("title")
		

		file, header, err := r.FormFile("photo")
		var photoURL string
		
		if err == nil {
			defer file.Close()
		
			photoDir := "uploads/" 
		
			if _, err := os.Stat(photoDir); os.IsNotExist(err) {
				err := os.MkdirAll(photoDir, 0755)
				if err != nil {
					http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
					return
				}
			}
		
			photoPath := photoDir + header.Filename
		
			dst, err := os.Create(photoPath)
			if err != nil {
				fmt.Println("Error saving file:", err)
				http.Error(w, "Error saving photo", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, file)
		
			photoURL = "/uploads/" + header.Filename
			fmt.Println("Photo URL:", photoURL)
		} else {
			fmt.Println("ddddd")
			photoURL =""
		}
		
		

		interest := strings.Join(interests, "#")
		if content != "" && title != "" {
			CreatePost(userID, content, interest, title, photoURL)
		}
	}

	posts, err := GetAllPosts(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		fmt.Println(err)
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

func CreatePost(userID int, content, interest, title, photo string) {
	stmt, err := database.DB.Prepare("INSERT INTO posts(user_id, content, interest, title, photo) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, content, interest, title, photo)
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
			p.title,
			p.interest, 
			p.created_at,
			p.photo,
			(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'like') as likes,
			(SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'dislike') as dislikes,
			(SELECT reaction_type FROM post_reactions WHERE post_id = p.id AND user_id = ?) as user_reaction
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		fmt.Println("Error querying posts:", err)
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
			&post.Title,
			&post.Interest,
			&post.CreatedAt,
			&post.Photo,
			&post.Likes,
			&post.Dislikes,
			&userReaction,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		post.UserReaction = userReaction.String
		createdAtTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing CreatedAt time:", err)
			continue
		}
		post.MinutesSinceCreation = int(time.Since(createdAtTime).Minutes())
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
		createdAtTime, err := time.Parse(time.RFC3339, comment.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing CreatedAt time:", err)
			continue
		}
		comment.MinutesSinceCreation = int(time.Since(createdAtTime).Minutes())
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetPostsByInterest(interest string, user int) ([]Post, error) {
	rows, err := database.DB.Query(`
        SELECT 
            p.id, 
            u.username, 
            p.content, 
            p.title,
            p.interest, 
            p.created_at,
			p.photo,
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'like') as likes,
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'dislike') as dislikes,
            (SELECT reaction_type FROM post_reactions WHERE post_id = p.id AND user_id = ?) as user_reaction
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.interest LIKE ?
        ORDER BY p.created_at DESC
    `, user, "%"+interest+"%")
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
			&post.Title,
			&post.Interest,
			&post.CreatedAt,
			&post.Photo,
			&post.Likes,
			&post.Dislikes,
			&userReaction,
		)
		if err != nil {
			continue
		}
		createdAtTime, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing CreatedAt time:", err)
			continue
		}
		// Calculate minutes since the post was created
		post.MinutesSinceCreation = int(time.Since(createdAtTime).Minutes())
		post.UserReaction = userReaction.String

		posts = append(posts, post)
	}
	return posts, nil
}
