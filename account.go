package main

import (
	"log"
	"net/http"
	"database/sql"
	"strings"
)

func sttgs(w http.ResponseWriter, r *http.Request)  {
	session, _ := store.Get(r, "cookie-name")

	if r.Method == "POST" {
		_, err = db.Exec("UPDATE \"USER\" SET lastname = $1, firstname = $2, nickname = $3, mail = $4 WHERE user_id = $5",
			r.FormValue("lastname"), r.FormValue("firstname"), r.FormValue("nickname"), strings.Trim(strings.ToLower(r.FormValue("mail")), " "), session.Values["user_id"].(string))
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		if r.FormValue("password") != "" {
			_, err = db.Exec("UPDATE \"USER\" SET password = crypt($1, gen_salt('bf', 8)) WHERE user_id = $2",
				r.FormValue("password"), session.Values["user_id"].(string))
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

		}
		session.Values["lastname"] = r.FormValue("lastname")
		session.Values["firstname"] = r.FormValue("firstname")
		session.Values["nickname"] = r.FormValue("nickname")
		session.Values["mail"] = r.FormValue("mail")
		session.Save(r, w)
	}
	usr := UserData {}
	usr.Lastname = session.Values["lastname"].(string)
	usr.Firstname = session.Values["firstname"].(string)
	usr.Nickname = session.Values["nickname"].(string)
	usr.Mail = session.Values["mail"].(string)
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
	session, _ := store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func lgn(w http.ResponseWriter, r *http.Request) {
	wd := WebData {
		Error: "",
	}
	session, _ := store.Get(r, "cookie-name")

	if r.Method == "POST" {
		row := db.QueryRow("SELECT * from \"USER\" WHERE login_username = $1 AND password = crypt($2, password);",
			strings.Trim(strings.ToLower(r.FormValue("username")), " "), r.FormValue("password"))
		usr := User{}
		err := row.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.nb_follow)
		if err == sql.ErrNoRows {
			wd.Error = "Invalid username or password"
		} else {
			session.Values["authenticated"] = true
			session.Values["lastname"] = usr.lastname
			session.Values["firstname"] = usr.firstname
			session.Values["nickname"] = usr.nickname
			session.Values["mail"] = usr.mail
			session.Values["login_username"] = usr.login_username
			session.Values["user_id"] = usr.user_id
			session.Save(r, w)
			http.Redirect(w, r, "/feed", http.StatusSeeOther)
		}
	}
	err := tpl.ExecuteTemplate(w, "login.html", &wd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

