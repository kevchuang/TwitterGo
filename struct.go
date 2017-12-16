package main

import (
	"time"
	"database/sql"
)

type ListData struct {
	Nicknames []string
	User_id string
}

type ProfileData struct {
	User UserData
	Follow string
	User_id string
}

type PostPageData struct {
	Comments []PostData
	Post PostData
}

type WebData struct {
	Error string
}

type FeedData struct {
	Posts []PostData
	User_id string
}

type UserData struct {
	Lastname string
	Firstname string
	Nickname string
	Mail string
	Nb_follow string
	User_id string
}



type PostData struct {
	Nickname string
	Post_id int
	Content string
	Date string
	Nb_of_likes int
	Ans_to_post string
	User_id int

}

type Friend struct {
	user_id string
	id_followed string
}

type Like struct {
	user_id string
	post_id string
}

type Post struct {
	post_id int
	content string
	date time.Time
	nb_of_likes int
	ans_to_post sql.NullString
	user_id int
}

type User struct {
	lastname string
	firstname string
	nickname string
	mail string
	login_username string
	password string
	user_id string
	nb_follow string
}
