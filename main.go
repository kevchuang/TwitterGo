package main

import (
	"html/template"
	"net/http"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

var db *sql.DB
var tpl *template.Template
var err error

func init()  {
	db, err = sql.Open("postgres", "host=localhost port=5435 user=postgres password=postgre dbname=BJTUitter sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database")
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func main()  {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/login", lgn)
	r.HandleFunc("/logout", lgout)
	r.HandleFunc("/registration", rgstr)
	r.HandleFunc("/settings", sttgs)
	r.HandleFunc("/feed", feed)
	r.HandleFunc("/add_post", add_post)
	r.HandleFunc("/like/{post_id}", like)
	r.HandleFunc("/like/{post_page}/{post_id}", like_post)
	r.HandleFunc("/post/{post_id}", post)
	r.HandleFunc("/comment/{post_id}", comment)
	r.HandleFunc("/edit/{post_page}/{post_id}", edit)
	r.HandleFunc("/delete/{post_page}/{post_id}", dlete)
	r.HandleFunc("/search", search)
	r.HandleFunc("/profile/{user_id}", profile)
	r.HandleFunc("/follow/{user_id}", follow)
	r.HandleFunc("/list_following/{user_id}", list_following)
	r.HandleFunc("/list_followers/{user_id}", list_followers)
	http.Handle("/", r)
	http.ListenAndServe(":8082", nil)
}

func index(w http.ResponseWriter, r *http.Request)  {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}