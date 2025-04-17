package database

import "forum/models"

func GetComments() ([]models.Comment, error) {
	rows, err := DB.Query("SELECT id, post_id, user_id, content, created_at FROM comments ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
		if err != nil {
			continue 
		}
		comments = append(comments, c)
	}
	return comments, nil
}


func AddComment(userID, postID int, content string) error {
	_, err := DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", userID, postID, content)
	return err
}

