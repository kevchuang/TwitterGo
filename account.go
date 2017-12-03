package main

import (
	"log"
	"net/http"
	"database/sql"
	"strings"
	"time"
)

func sttgs(w http.ResponseWriter, r *http.Request)  {
	user_id, _ := r.Cookie("user_id")
	firstname, _ := r.Cookie("firstname")
	lastname, _ := r.Cookie("lastname")
	mail, _ := r.Cookie("mail")
	nickname, _ := r.Cookie("nickname")

	if r.Method == "POST" {
		_, err = db.Exec("UPDATE \"USER\" SET lastname = $1, firstname = $2, nickname = $3, mail = $4 WHERE user_id = $5",
			r.FormValue("lastname"), r.FormValue("firstname"), r.FormValue("nickname"), strings.Trim(strings.ToLower(r.FormValue("mail")), " "), user_id.Value)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if r.FormValue("password") != "" {
			_, err = db.Exec("UPDATE \"USER\" SET password = crypt($1, gen_salt('bf', 8)) WHERE user_id = $2",
				r.FormValue("password"), user_id.Value)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

		}

		lastname.Value = r.FormValue("lastname")
		firstname.Value = r.FormValue("firstname")
		nickname.Value = r.FormValue("nickname")
		mail.Value = r.FormValue("mail")

		http.SetCookie(w, lastname)
		http.SetCookie(w, firstname)
		http.SetCookie(w, nickname)
		http.SetCookie(w, mail)
		http.SetCookie(w, user_id)
	}
	usr := UserData {}
	usr.Lastname = lastname.Value
	usr.Firstname = firstname.Value
	usr.Nickname = nickname.Value
	usr.Mail = mail.Value
	err := tpl.ExecuteTemplate(w, "settings.html", &usr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func rgstr(w http.ResponseWriter, r *http.Request) {
	wd := WebData{}
	if r.Method == "POST" {
		usr := User{}
		usr.lastname = r.FormValue("lastname")
		usr.firstname = r.FormValue("firstname")
		usr.nickname = r.FormValue("login")
		usr.mail = strings.Trim(strings.ToLower(r.FormValue("mail")), " ")
		usr.login_username = strings.Trim(strings.ToLower(r.FormValue("login")), " ")
		usr.password = r.FormValue("password")

		if usr.lastname == "" || usr.firstname == "" || usr.nickname == "" || usr.mail == "" || usr.login_username == "" || usr.password == "" {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		row := db.QueryRow("SELECT * FROM \"USER\" WHERE login_username = $1", usr.login_username)
		usrCheck := User{}
		check := row.Scan(&usrCheck.lastname, &usrCheck.firstname, &usrCheck.nickname, &usrCheck.mail, &usrCheck.login_username, &usrCheck.password, &usrCheck.user_id, &usrCheck.nb_follow)

		if check != sql.ErrNoRows {
			wd.Error = "Login already exists !"
		} else {

			_, err = db.Exec("INSERT INTO \"USER\" (lastname, firstname, nickname, mail, login_username, password) VALUES ($1, $2, $3, $4, $5, crypt($6, gen_salt('bf', 8)))",
				usr.lastname, usr.firstname, usr.login_username, usr.mail, usr.login_username, usr.password)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
	err := tpl.ExecuteTemplate(w, "registration.html", &wd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func lgout(w http.ResponseWriter, r *http.Request) {
	auth, _ := r.Cookie("authenticated")
	auth.Value = "false"
	http.SetCookie(w, auth)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func lgn(w http.ResponseWriter, r *http.Request) {
	wd := WebData {
		Error: "",
	}

	if r.Method == "POST" {
		row := db.QueryRow("SELECT * from \"USER\" WHERE login_username = $1 AND password = crypt($2, password);",
			strings.Trim(strings.ToLower(r.FormValue("username")), " "), r.FormValue("password"))
		usr := User{}
		err := row.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.nb_follow)
		if err == sql.ErrNoRows {
			wd.Error = "Invalid username or password"
		} else {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			auth := http.Cookie{Name: "authenticated", Value: "true", Expires: expiration}
			lastname := http.Cookie{Name: "lastname", Value: usr.lastname, Expires: expiration}
			firstname := http.Cookie{Name: "firstname", Value: usr.firstname, Expires: expiration}
			nickname := http.Cookie{Name: "nickname", Value: usr.nickname, Expires: expiration}
			mail := http.Cookie{Name: "mail", Value: usr.mail, Expires: expiration}
			login := http.Cookie{Name: "login_username", Value: usr.login_username, Expires: expiration}
			user_id := http.Cookie{Name: "user_id", Value: usr.user_id, Expires: expiration}
			http.SetCookie(w, &auth)
			http.SetCookie(w, &lastname)
			http.SetCookie(w, &firstname)
			http.SetCookie(w, &nickname)
			http.SetCookie(w, &mail)
			http.SetCookie(w, &login)
			http.SetCookie(w, &user_id)

			http.Redirect(w, r, "/feed", http.StatusSeeOther)
		}
	}
	err := tpl.ExecuteTemplate(w, "login.html", &wd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

