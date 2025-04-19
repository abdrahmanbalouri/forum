package handlers

import (
	"context"
	"fmt"
	"forum/config"
	"forum/database"
	"forum/security"
	"log"
	"net/http"
	"time"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMidleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil || sessionCookie.Value == "" {
			ServLogin(w, r)
			return
		}

		user_id, exist := database.GetUserBySession(sessionCookie.Value)
		if !exist {
			ServLogin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, user_id) 
		next(w, r.WithContext(ctx))
	}
}

func SwitchLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ServLogin(w, r)
	case http.MethodPost:
		SubmitLogin(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func SwitchRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ServRegister(w, r)
	case http.MethodPost:
		SubmitRegister(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func ServLogin(w http.ResponseWriter, r *http.Request) {
	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "login.html", nil)
}

func ServRegister(w http.ResponseWriter, r *http.Request) {
	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "register.html", nil)
}

func SubmitRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")
	fmt.Println(confirm_password)


	hash, err := security.HashPassword(password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = database.AddNewUser(username, email, hash)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" || err.Error() == "UNIQUE constraint failed: users.email" {
			http.Error(w, "Username or email already exists.", http.StatusConflict)
		} else {
			log.Printf("Error adding user: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func SubmitLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	

	user_id, hash := database.GetUserHash(username)
	

	if !security.CheckPassword(password, hash) {
		http.Error(w, "wrong passwrod", http.StatusUnauthorized)
		return
	}

	session := security.GenerateToken(32) // TODO: UUID bonus csrf implementation genrate csrf read it in front end js and match it with server go

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
	})

	database.SetSessionToken(user_id, session)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	hasSession := database.DeleteUserBySession(sessionCookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	if !hasSession {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
