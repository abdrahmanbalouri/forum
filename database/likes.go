package database

import "log"

func GetLikesCount() int {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM likes").Scan(&count)
	if err != nil {
		log.Printf("Error fetching likes count: %v", err)
		return 0
	}
	return count
}

func RemoveLike(userID, postID int) error {
	_, err := DB.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
	return err
}

func AddLike(userID, postID int) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	return err
}
