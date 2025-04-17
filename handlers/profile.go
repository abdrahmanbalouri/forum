package handlers

import (
	"fmt"
	"forum/config"
	"forum/database"
	"forum/security"
	"net/http"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(userIDKey).(int)
	user := database.GetUserInfo(user_id)
	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "profile.html", user) // when u excute 2 template the get concatinated one in top of the other
}

func UpddateProfile(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("value")
	switch value {
	case "username":
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", map[string]bool{"username": true})
		return
	case "email":
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", map[string]bool{"email": true})
		return
	case "password":
		config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", map[string]bool{"password": true})
		return
	default:
		http.Error(w, "bad req", http.StatusBadRequest)
	}
}

func SaveChanges(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(userIDKey).(int)
	switch r.PathValue("value") {
	case "username":
		new_username := r.FormValue("username")
		password := r.FormValue("current")
		if !config.ValidUsername(new_username) {
			http.Redirect(w, r, "/profile/update/username", http.StatusSeeOther)
			return
		}
		hash := database.GetUserHashById(user_id)
		if !security.CheckPassword(password, hash) {
			http.Redirect(w, r, "/profile/update/username", http.StatusSeeOther)
			return
		}
		database.UpdateUsernmae(user_id, new_username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case "email":
		new_email := r.FormValue("email")
		password := r.FormValue("current")
		if !config.ValidEmail(new_email) {
			http.Redirect(w, r, "/profile/update/email", http.StatusSeeOther)
			return
		}
		hash := database.GetUserHashById(user_id)
		if !security.CheckPassword(password, hash) {
			http.Redirect(w, r, "/profile/update/email", http.StatusSeeOther)
			return
		}
		database.UpdateEmail(user_id, new_email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case "password":
		current := r.FormValue("current")
		new := r.FormValue("new")
		confirm := r.FormValue("confirm")
		hash := database.GetUserHashById(user_id)
		if !security.CheckPassword(current, hash) {
			http.Redirect(w, r, "/profile/update/password", http.StatusSeeOther)
			return
		}
		if new != confirm {
			http.Redirect(w, r, "/profile/update/password", http.StatusSeeOther)
			return
		}
		new_hash, err := security.HashPassword(new)
		if err != nil {
			fmt.Println(err)
		}
		database.UpdatePassword(user_id, new_hash)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		http.Error(w, "bad req", http.StatusBadRequest)
		return
	}
}

func ServeDelete(w http.ResponseWriter, r *http.Request) {
	config.GLOBAL_TEMPLATE.ExecuteTemplate(w, "delete.html", nil)
}

func DeleteConfirmation(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(userIDKey).(int)
	password := r.FormValue("password")
	hash := database.GetUserHashById(user_id)
	if !security.CheckPassword(password, hash) {
		http.Redirect(w, r, "/profile/delete", http.StatusSeeOther)
		return
	}
	database.DeleteUser(user_id)
	LogoutHandler(w, r)
}
