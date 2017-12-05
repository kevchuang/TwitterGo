package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/url"
	"time"
)

var (
	router = Router()
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/login", lgn)
	router.HandleFunc("/logout", lgout)
	router.HandleFunc("/registration", rgstr)
	router.HandleFunc("/settings", sttgs)
	router.HandleFunc("/feed", feed)
	router.HandleFunc("/add_post", add_post)
	router.HandleFunc("/like/{post_id}", like)
	router.HandleFunc("/like/{post_page}/{post_id}", like_post)
	router.HandleFunc("/post/{post_id}", post)
	router.HandleFunc("/comment/{post_id}", comment)
	router.HandleFunc("/edit/{post_page}/{post_id}", edit)
	router.HandleFunc("/delete/{post_page}/{post_id}", dlete)
	router.HandleFunc("/search", search)
	router.HandleFunc("/profile/{user_id}", profile)
	router.HandleFunc("/follow/{user_id}", follow)
	router.HandleFunc("/list_following/{user_id}", list_following)
	router.HandleFunc("/list_followers/{user_id}", list_followers)
	return router
}

func SetCookies(w http.ResponseWriter) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	auth := http.Cookie{Name: "authenticated", Value: "true", Expires: expiration}
	lastname := http.Cookie{Name: "lastname", Value: "unittest", Expires: expiration}
	firstname := http.Cookie{Name: "firstname", Value: "unittest", Expires: expiration}
	nickname := http.Cookie{Name: "nickname", Value: "unittest", Expires: expiration}
	mail := http.Cookie{Name: "mail", Value: "unittest", Expires: expiration}
	login := http.Cookie{Name: "login_username", Value: "unittest", Expires: expiration}
	user_id := http.Cookie{Name: "user_id", Value: "1", Expires: expiration}
	http.SetCookie(w, &auth)
	http.SetCookie(w, &lastname)
	http.SetCookie(w, &firstname)
	http.SetCookie(w, &nickname)
	http.SetCookie(w, &mail)
	http.SetCookie(w, &login)
	http.SetCookie(w, &user_id)
}

func TestFeedNoLog(t *testing.T)  {
	request, _ := http.NewRequest("GET", "/feed", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "forbidden response is expected")
}

func TestLoginWrongUser(t *testing.T) {
	request, _ := http.NewRequest("POST", "/login", nil)
	form := url.Values{}
	form.Add("username", "wrong")
	form.Add("password", "wrong")
	request.Form = form
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Wrong user response is expected")
}

//func TestRegisterUser(t *testing.T) {
//	request, _ := http.NewRequest("POST", "/registration", nil)
//	form := url.Values{}
//	form.Add("lastname", "unittest")
//	form.Add("firstname", "unittest")
//	form.Add("login", "unittest")
//	form.Add("mail", "unittest")
//	form.Add("password", "unittest")
//	request.Form = form
//	response := httptest.NewRecorder()
//	Router().ServeHTTP(response, request)
//	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirection response is expected")
//}

func TestLoginUser(t *testing.T) {
	request, _ := http.NewRequest("POST", "/login", nil)
	form := url.Values{}
	form.Add("username", "   uNitTeSt           ")
	form.Add("password", "unittest")
	request.Form = form
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
}

func TestLogout(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("GET", "/logout", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
}

func TestLogoutNolog(t *testing.T) {
	request, _ := http.NewRequest("GET", "/logout", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusForbidden, response.Code, "Forbidden response is expected")
}

func TestSettings(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/settings", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	form := url.Values{}
	form.Add("mail", "unittest@test.com")
	form.Add("lastname", "unittest")
	form.Add("firstname", "unittest")
	form.Add("nickname", "unittest")
	request.Form = form
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Ok response is expected")
}

func TestFeedLog(t *testing.T)  {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("GET", "/feed", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Ok response is expected")
}

func TestSearchWrongUser(t *testing.T) {
	request, _ := http.NewRequest("POST", "/search", nil)
	form := url.Values{}
	form.Add("search-content", "undefined user")
	request.Form = form
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad request response is expected")
}

func TestSearchUser(t *testing.T) {
	request, _ := http.NewRequest("POST", "/search", nil)
	form := url.Values{}
	form.Add("search-content", "unittest")
	request.Form = form
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
}

func TestAddPost(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/add_post", nil)
	form := url.Values{}
	form.Add("post-content", "Unit test post")
	request.Form = form
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}

	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
}

func TestAddEmptyPost(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/add_post", nil)
	form := url.Values{}
	form.Add("post-content", "")
	request.Form = form
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}

	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad request response is expected")
}

func TestLike(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/like/1", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Ok response is expected")
}

func TestPost(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("GET", "/post/1", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Ok response is expected")
}

func TestComment(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/comment/1", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	form := url.Values{}
	form.Add("comment", "Comment unit test")
	request.Form = form
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Ok response is expected")
}

func TestNoComment(t *testing.T) {
	response := httptest.NewRecorder()
	SetCookies(response)
	request, _ := http.NewRequest("POST", "/comment/1", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	form := url.Values{}
	form.Add("comment", "")
	request.Form = form
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad request response is expected")
}

