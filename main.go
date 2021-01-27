package main

import (
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/url"
)

var db = map[string][]byte{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	errMsg := r.FormValue("errorMsg")

	html := `<!DOCTYPE html>
	<html lang="en"
		<head>
			<meta charset="UTF-8"> 
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>HMAC Example</title>
		</head>
		<body>
			<p>` + errMsg + `</p>
			<h1>REGISTER</h1>
			<form action="/register" method="POST">	
				<input type="text" name="name"/>
				<input type="password" name="password"/>
				<input type="submit" name="Register"/>
			</form>
			<h1>LOG IN</h1>
			<form action="/login" method="POST">	
				<input type="text" name="name"/>
				<input type="password" name="password"/>
				<input type="submit" name="Login"/>
			</form>
		</body>
	</html>`

	io.WriteString(w, html)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("your method was not post")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")

	if name == "" || password == "" {
		errorMsg := url.QueryEscape("name or password is empty")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	bsp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		errorMsg := url.QueryEscape("error encrpting password")
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	db[name] = bsp
	http.Redirect(w, r, "/hello?name=" + name, http.StatusSeeOther)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("your method was not post")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")

	if name == "" || password == "" {
		errorMsg := url.QueryEscape("name or password is empty")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	if _, ok := db[name]; !ok {
		errorMsg := url.QueryEscape("your name and password don't match")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(db[name]), []byte(password))
	if err != nil {
		errorMsg := url.QueryEscape("wrong password")
		http.Redirect(w, r, "/?errorMsg=" + errorMsg, http.StatusSeeOther)
		return
	}

	errorMsg := url.QueryEscape("you logged in " + name)
	http.Redirect(w, r, "?errorMsg=" + errorMsg, http.StatusSeeOther)
}

