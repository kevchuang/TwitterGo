package main

import (
	"net/http"
	"log"
	"time"
	"github.com/gorilla/mux"
	"database/sql"
	"github.com/lib/pq"
	"strconv"
)

func getAllUserId(w http.ResponseWriter, user_id string) []int {
	valueUser, errorLog := strconv.Atoi(user_id)
	if errorLog != nil {
		log.Println(errorLog)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	buffer := []int{valueUser}

	rows, err := db.Query("SELECT * FROM \"FRIENDS\" WHERE user_id = $1", user_id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer rows.Close()
	if rows != nil {
		friend := Friend{}
		for rows.Next() {
			rows.Scan(&friend.user_id, &friend.id_followed)
			value, err := strconv.Atoi(friend.id_followed)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			buffer = append(buffer, value)
		}
	}
	return buffer
}

func feed(w http.ResponseWriter, r *http.Request) {
	auth, _ := r.Cookie("authenticated")
	user_id, _ := r.Cookie("user_id")

	posts := make([]PostData, 0)
	if auth == nil || auth.Value != "true" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	rows, err := db.Query("SELECT * FROM \"POSTS\" WHERE ans_to_post IS NULL AND user_id = ANY($1) ORDER BY date DESC", pq.Array(getAllUserId(w, user_id.Value)))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	if rows != nil {
		postScan := Post{}
		for rows.Next() {
			rows.Scan(&postScan.post_id, &postScan.content, &postScan.date, &postScan.nb_of_likes, &postScan.ans_to_post, &postScan.user_id)
			post := PostData{}
			post.Post_id = postScan.post_id
			post.Content = postScan.content
			post.Date = postScan.date.Format("Mon Jan _2 15:04:05 2006")
			post.Nb_of_likes = postScan.nb_of_likes
			post.Ans_to_post = postScan.ans_to_post.String
			post.User_id = postScan.user_id
			nicknameQuery := db.QueryRow("SELECT nickname FROM \"USER\" WHERE user_id = $1", post.User_id)
			nickname := string("")
			nicknameQuery.Scan(&nickname)
			post.Nickname = nickname
			posts = append(posts, post)
		}
	}
	data := FeedData {}
	data.Posts = posts
	data.User_id = user_id.Value
	errorLog := tpl.ExecuteTemplate(w, "feed.html", &data)
	if errorLog != nil {
		log.Println(errorLog)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func add_post(w http.ResponseWriter, r *http.Request) {
	user_id, _ := r.Cookie("user_id")
	if r.Method == "POST" {
		if r.FormValue("post-content") != "" {
			_, err := db.Exec("INSERT INTO \"POSTS\" (content, date, user_id) VALUES ($1, $2, $3)",
				r.FormValue("post-content"), time.Now(), user_id.Value)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/feed", http.StatusNotModified)
			return
		}
		http.Redirect(w, r, "/feed", http.StatusSeeOther)
		return
	}
}

func like_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	result := db.QueryRow("SELECT * from \"LIKES\" WHERE user_id = $1 AND post_id = $2", user_id.Value, vars["post_id"])

	lk := Like{}
	err := result.Scan(&lk.user_id, &lk.post_id)
	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO \"LIKES\" (user_id, post_id) VALUES ($1, $2)", user_id.Value, vars["post_id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		_, errorLog := db.Exec("UPDATE \"POSTS\" SET nb_of_likes = nb_of_likes + 1 WHERE post_id = $1", vars["post_id"])
		if errorLog != nil {
			log.Println(errorLog)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		_, err := db.Exec("DELETE FROM \"LIKES\" WHERE user_id = $1 AND post_id = $2", user_id.Value, vars["post_id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		_, errorLog := db.Exec("UPDATE \"POSTS\" SET nb_of_likes = nb_of_likes - 1 WHERE post_id = $1", vars["post_id"])
		if errorLog != nil {
			log.Println(errorLog)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/post/" + vars["post_page"], http.StatusSeeOther)
}

func like(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	result := db.QueryRow("SELECT * from \"LIKES\" WHERE user_id = $1 AND post_id = $2", user_id.Value, vars["post_id"])

	lk := Like{}
	err := result.Scan(&lk.user_id, &lk.post_id)
	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO \"LIKES\" (user_id, post_id) VALUES ($1, $2)", user_id.Value, vars["post_id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_, errorLog := db.Exec("UPDATE \"POSTS\" SET nb_of_likes = nb_of_likes + 1 WHERE post_id = $1", vars["post_id"])
		if errorLog != nil {
			log.Println(errorLog)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		_, err := db.Exec("DELETE FROM \"LIKES\" WHERE user_id = $1 AND post_id = $2", user_id.Value, vars["post_id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_, errorLog := db.Exec("UPDATE \"POSTS\" SET nb_of_likes = nb_of_likes - 1 WHERE post_id = $1", vars["post_id"])
		if errorLog != nil {
			log.Println(errorLog)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
	return
}

func post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := PostPageData{}
	result := db.QueryRow("SELECT * FROM \"POSTS\" WHERE post_id = $1", vars["post_id"])

	post := Post{}
	err := result.Scan(&post.post_id, &post.content, &post.date, &post.nb_of_likes, &post.ans_to_post, &post.user_id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	postData := PostData{}
	postData.Post_id = post.post_id
	postData.Content = post.content
	postData.Date = post.date.Format("Mon Jan _2 15:04:05 2006")
	postData.Nb_of_likes = post.nb_of_likes
	postData.Ans_to_post = post.ans_to_post.String
	postData.User_id = post.user_id
	nicknameQuery := db.QueryRow("SELECT nickname FROM \"USER\" WHERE user_id = $1", postData.User_id)
	nickname := string("")
	nicknameQuery.Scan(&nickname)
	postData.Nickname = nickname
	data.Post = postData

	comments := make([]PostData, 0)
	rows, err := db.Query("SELECT * FROM \"POSTS\" WHERE ans_to_post = $1", vars["post_id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer rows.Close()
	if rows != nil {
		postScan := Post{}
		for rows.Next() {
			rows.Scan(&postScan.post_id, &postScan.content, &postScan.date, &postScan.nb_of_likes, &postScan.ans_to_post, &postScan.user_id)
			post := PostData{}
			post.Post_id = postScan.post_id
			post.Content = postScan.content
			post.Date = postScan.date.Format("Mon Jan _2 15:04:05 2006")
			post.Nb_of_likes = postScan.nb_of_likes
			post.Ans_to_post = postScan.ans_to_post.String
			post.User_id = postScan.user_id
			nicknameQuery := db.QueryRow("SELECT nickname FROM \"USER\" WHERE user_id = $1", post.User_id)
			nickname := string("")
			nicknameQuery.Scan(&nickname)
			post.Nickname = nickname
			comments = append(comments, post)
		}
		data.Comments = comments
	}
	errorLog := tpl.ExecuteTemplate(w, "post.html", &data)
	if errorLog != nil {
		log.Println(errorLog)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func comment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	if r.Method == "POST" {
		if r.FormValue("comment") != "" {
			_, err := db.Exec("INSERT INTO \"POSTS\" (content, date, ans_to_post, user_id) VALUES ($1, $2, $3, $4)",
				r.FormValue("comment"), time.Now(), vars["post_id"], user_id.Value)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/post/" + vars["post_id"], http.StatusNotModified)
			return
		}
		http.Redirect(w, r, "/post/" + vars["post_id"], http.StatusSeeOther)
		return
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if r.Method == "POST" {
		if r.FormValue("editPostText" ) != "" {
			_, err := db.Exec("UPDATE \"POSTS\" SET content = $1 WHERE post_id = $2",
				r.FormValue("editPostText"), vars["post_id"])
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/post/" + vars["post_page"], http.StatusNotModified)
			return
		}
		http.Redirect(w, r, "/post/" + vars["post_page"], http.StatusSeeOther)
		return
	}
}

func dlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := r.Cookie("user_id")

	row := db.QueryRow("SELECT * FROM \"POSTS\" WHERE user_id = $1 AND post_id = $2 AND ans_to_post IS NULL",
		user_id.Value, vars["post_id"])
	post := Post{}
	result := row.Scan(&post.post_id, &post.content, &post.date, &post.nb_of_likes, &post.ans_to_post, &post.user_id)
	_, err := db.Exec("DELETE FROM \"POSTS\" WHERE user_id = $1 AND post_id = $2", user_id.Value, vars["post_id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	if result != sql.ErrNoRows {
		rows, err := db.Query("SELECT * FROM \"POSTS\" WHERE ans_to_post = $1", vars["post_id"])
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		defer rows.Close()
		if rows != nil {
			comment := Post{}
			for rows.Next() {
				rows.Scan(&comment.post_id, &comment.content, &comment.date, &comment.nb_of_likes, &comment.ans_to_post, &comment.user_id)
				_, err := db.Exec("DELETE FROM \"POSTS\" WHERE post_id = $1", comment.post_id)
				if err != nil {
					log.Println(err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}
		http.Redirect(w, r, "/feed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/post/"+vars["post_page"], http.StatusSeeOther)
	return
}