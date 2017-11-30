package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/url"
	"github.com/gorilla/sessions"
	"database/sql"
)

var (
	keyTest = []byte("super-secret-key")
	storeTest = sessions.NewCookieStore(keyTest)
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/login", lgn)
	router.HandleFunc("/registration", rgstr)
	router.HandleFunc("/feed", feed)
	return router
}

func TestLogin(t *testing.T)  {
	request, _ := http.NewRequest("GET", "/login", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Ok response is expected")
}


func TestFeedNoLog(t *testing.T)  {
	request, _ := http.NewRequest("GET", "/feed", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "forbidden response is expected")
}

func TestRegister(t *testing.T) {
	request, _ := http.NewRequest("GET", "/registration", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Ok response is expected")
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
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusSeeOther, response.Code, "Redirection response is expected")
}

func TestLoginWrongUser(t *testing.T) {
	request, _ := http.NewRequest("POST", "/login", nil)
	form := url.Values{}
	form.Add("username", "wrong")
	form.Add("password", "wrong")
	request.Form = form
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code, "Wrong user response is expected")
}