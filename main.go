package main

import (
	"fmt"
	"net/http"

	"forum/config"
	"forum/database"
	"forum/handlers"

	_ "github.com/mattn/go-sqlite3"
)

const (
	PORT      = ":8080"
	SERVERURL = "http://localhost:8080"
)

func main() {
	database.InitDB("./database/forum.db")
	config.InitTemplate()
	config.InitRegex()

	forumux := http.NewServeMux()
	forumux.HandleFunc("/login", handlers.SwitchLogin)
	forumux.HandleFunc("/register", handlers.SwitchRegister)
	forumux.HandleFunc("/logout", handlers.LogoutHandler)

	forumux.HandleFunc("/profile", handlers.AuthMidleware(handlers.ProfilHandler))
	forumux.HandleFunc("/profile/update/{value}", handlers.AuthMidleware(handlers.UpddateProfile))
	forumux.HandleFunc("/profile/update/{value}/save", handlers.AuthMidleware(handlers.SaveChanges))
	forumux.HandleFunc("/profile/delete", handlers.AuthMidleware(handlers.ServeDelete))
	forumux.HandleFunc("/profile/delete/confirm", handlers.AuthMidleware(handlers.DeleteConfirmation))

	forumux.HandleFunc("/post", handlers.AuthMidleware(handlers.PostHandler))
	forumux.HandleFunc("/comment", handlers.AuthMidleware(handlers.CommentHandler))
	forumux.HandleFunc("/reaction", handlers.AuthMidleware(handlers.ReactionHandler))
	forumux.HandleFunc("/filter", handlers.AuthMidleware(handlers.FilterHandler))

	forumux.HandleFunc("/",handlers.RootHandler)
	forumux.HandleFunc("/static/", handlers.StaticHandler)
	forumux.HandleFunc("uploads/", handlers.StaticHandler2)


	fmt.Println("Server running on ", SERVERURL)
	err := http.ListenAndServe(PORT, forumux)

	// fmt.Println("Available templates:", temp.DefinedTemplates())
	fmt.Println(err)
}
