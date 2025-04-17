package database

import (
	"forum/models"
	"log"
)

func AddNewUser(username, email, hashedPass string) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, username, email, hashedPass)
	return err
}

func AlreadyExists(username, email string) bool {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	err := DB.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}
	return count > 0
}

func GetUserHash(username string) (int, string) {
	var hash string
	var id int
	query := `SELECT id,password_hash FROM users WHERE username = ? OR email = ?`
	err := DB.QueryRow(query, username, username).Scan(&id, &hash)
	if err != nil {
		log.Println(err)
	}
	return id, hash
}
func GetUserEmailBySession(sessionToken string) (string) {
	var email string

	err := DB.QueryRow(`
		SELECT users.email
		FROM sessions
		JOIN users ON sessions.user_id = users.id
		WHERE sessions.session_token = ? AND sessions.expires_at > datetime('now')
	`, sessionToken).Scan(&email)

	if err != nil {
		log.Println(err)
	}

	return email
}



func GetUserHashById(id int) string {
	var hash string
	query := `SELECT password_hash FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&hash)
	if err != nil {
		log.Println(err)
	}
	return hash
}

func SetSessionToken(id int, token string) {
	update_query := `UPDATE sessions SET expires_at = DATETIME('now', '+1 hour') , session_token = ? WHERE user_id = ?`
	insert_query := `INSERT INTO sessions (user_id,session_token,expires_at) VALUES (? , ? , DATETIME('now','+1 hour'))`
	res, err := DB.Exec(update_query, token, id)
	if err != nil {
		log.Println(err)
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return
	}
	_, err = DB.Exec(insert_query, id, token)
	if err != nil {
		log.Println(err)
	}
}

func GetUserBySession(session_id string) (int, bool) {
	var user_id int

	query := `SELECT user_id FROM sessions WHERE session_token = ?`
	err := DB.QueryRow(query, session_id).Scan(&user_id)
	if err != nil {
		log.Print(err)
		return user_id, false
	}
	return user_id, true
}

func DeleteUserBySession(session_id string) bool {
	query := `UPDATE sessions SET session_token = NULL WHERE session_token = ?`
	_, err := DB.Exec(query, session_id)
	if err != nil {
		log.Printf("Error updating session: %v", err)
		return false
	}
	return true
}

func GetUserInfo(user_id int) models.User {
	var user models.User

	user_query := `SELECT * FROM users WHERE id = ?`
	err := DB.QueryRow(user_query, user_id).Scan(&user.Id, &user.Username, &user.Email, &user.Password_hash, &user.Created_at, &user.Updated_at, &user.Email_verified)
	if err != nil {
		return user
	}
	return user
}

func UpdateUsernmae(id int, username string) {
	query := `UPDATE users SET updated_at = DATETIME('now') , username = ? WHERE id = ?`
	res, err := DB.Exec(query, username, id)
	if err != nil {
		log.Println(err)
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return
	}
}

func UpdateEmail(id int, email string) {
	query := `UPDATE users SET updated_at = DATETIME('now') , email = ? WHERE id = ?`
	res, err := DB.Exec(query, email, id)
	if err != nil {
		log.Println(err)
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return
	}
}

func UpdatePassword(id int, password string) {
	query := `UPDATE users SET updated_at = DATETIME('now') , password_hash = ? WHERE id = ?`
	res, err := DB.Exec(query, password, id)
	if err != nil {
		log.Println(err)
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return
	}
}

func DeleteUser(user_id int) {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.Exec(query, user_id)
	if err != nil {
		log.Println("delete", err)
	}
}

// func GetUserId(username string) int {
// 	var id int

// 	query := `SELECT id FROM users WHERE username = ? OR email = ?`
// 	err := DB.QueryRow(query, username, username).Scan(&id)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return id
// }

// func UserHasSession(id int) bool {
// 	query := `SELECT session_token FROM sessions WHERE user_id = ?`
// 	err := DB.QueryRow(query, id).Scan()
// 	return err == nil
// }
