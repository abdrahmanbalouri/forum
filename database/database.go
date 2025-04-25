package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			session_token TEXT,
			expires_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	content TEXT,
	title TEXT,
	interest TEXT,
	photo TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id)
);


		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS post_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction_type TEXT NOT NULL CHECK(reaction_type IN ('like', 'dislike')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(post_id, user_id)
		);

		CREATE TABLE IF NOT EXISTS comment_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			comment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction_type TEXT NOT NULL CHECK(reaction_type IN ('like', 'dislike')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(comment_id, user_id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully")
}

func CloseDB() {
	err := DB.Close()
	fmt.Println("closing the databse")
	if err != nil {
		log.Fatal("Error closing database:", err)
	}
}
