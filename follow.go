package main

import (
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
)

func search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.FormValue("search-content") != "" {
			row := db.QueryRow("SELECT * FROM \"USER\" WHERE login_username = $1 OR mail = $2",
				r.FormValue("search-content"), r.FormValue("search-content"))
			usr := User{}
			result := row.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.nb_follow)
			if result != sql.ErrNoRows {
				http.Redirect(w, r, "/profile/" + usr.user_id, http.StatusSeeOther)
				return
			}
		}
	}
	http.Redirect(w, r, "/feed", http.StatusBadRequest)
}

func profile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	row := db.QueryRow("SELECT * FROM \"USER\" WHERE user_id = $1", vars["user_id"])
	data := ProfileData{}
	data.User_id = user_id.Value
	usr := User{}
	result := row.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.nb_follow)
	data.User.User_id = usr.user_id
	data.User.Lastname = usr.lastname
	data.User.Mail = usr.mail
	data.User.Firstname = usr.firstname
	data.User.Nickname = usr.nickname
	data.User.Nb_follow = usr.nb_follow
	rowFriend := db.QueryRow("SELECT * FROM \"FRIENDS\" WHERE user_id = $1 and id_followed = $2",
		user_id.Value, vars["user_id"])

	friend := Friend{}
	result = rowFriend.Scan(&friend.user_id, &friend.id_followed)
	if vars["user_id"] == user_id.Value {
		data.Follow = ""
	} else if result == sql.ErrNoRows {
		data.Follow = "Follow"
	} else {
		data.Follow = "Unfollow"
	}
	err := tpl.ExecuteTemplate(w, "profile.html", &data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func follow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	row := db.QueryRow("SELECT * FROM \"FRIENDS\" WHERE user_id = $1 AND id_followed = $2",
		user_id.Value, vars["user_id"])
	friend := Friend{}

	result := row.Scan(&friend.user_id, &friend.id_followed)
	if result == sql.ErrNoRows {
		db.Exec("INSERT INTO \"FRIENDS\" (user_id, id_followed) VALUES ($1, $2)",
			user_id.Value, vars["user_id"])
		db.Exec("UPDATE \"USER\" SET nb_follow = nb_follow + 1 WHERE user_id = $1", user_id.Value)
	} else {
		db.Exec("DELETE FROM \"FRIENDS\" WHERE user_id = $1 AND id_followed = $2",
			user_id.Value, vars["user_id"])
		db.Exec("UPDATE \"USER\" SET nb_follow = nb_follow - 1 WHERE user_id = $1", user_id.Value)
	}
	http.Redirect(w, r, "/profile/" + vars["user_id"], http.StatusOK)
	return
}

func list_following(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	rows, err := db.Query("SELECT * FROM \"FRIENDS\" WHERE user_id = $1", vars["user_id"])

	data := ListData{}
	data.User_id = user_id.Value
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	if rows != nil {
		friend := Friend{}
		for rows.Next() {
			rows.Scan(&friend.user_id, &friend.id_followed)
			rowsUser := db.QueryRow("SELECT * FROM \"USER\" WHERE user_id = $1", friend.id_followed)
			usr := User{}
			rowsUser.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.password)
			data.Nicknames = append(data.Nicknames, usr.nickname)
		}
	}
	err = tpl.ExecuteTemplate(w, "list_following.html", &data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func list_followers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	rows, err := db.Query("SELECT * FROM \"FRIENDS\" WHERE id_followed = $1", vars["user_id"])

	data := ListData{}
	data.User_id = user_id.Value
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	if rows != nil {
		friend := Friend{}
		for rows.Next() {
			rows.Scan(&friend.user_id, &friend.id_followed)
			rowsUser := db.QueryRow("SELECT * FROM \"USER\" WHERE user_id = $1", friend.user_id)
			usr := User{}
			rowsUser.Scan(&usr.lastname, &usr.firstname, &usr.nickname, &usr.mail, &usr.login_username, &usr.password, &usr.user_id, &usr.password)
			data.Nicknames = append(data.Nicknames, usr.nickname)
		}
	}
	err = tpl.ExecuteTemplate(w, "list_following.html", &data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}