package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/url"
)

var (
	clt = &http.Client{}
	router = Router()
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", lgn)
	router.HandleFunc("/registration", rgstr)
	router.HandleFunc("/feed", feed)
	router.HandleFunc("/add_post", add_post)
	router.HandleFunc("/search", search)
	return router
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
//
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
	clt.Do(request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
}

func TestFeedLog(t *testing.T)  {
	request, _ := http.NewRequest("GET", "/feed", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "forbidden response is expected")
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

//func TestAddPost(t *testing.T) {
//	request, _ := http.NewRequest("POST", "/add_post", nil)
//	form := url.Values{}
//	form.Add("post-content", "Unit test post")
//	request.Form = form
//	response := httptest.NewRecorder()
//	Router().ServeHTTP(response, request)
//	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirect response is expected")
//}