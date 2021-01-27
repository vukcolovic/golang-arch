package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

type myClaims struct {
	jwt.StandardClaims
	Email string
}

const myKey = "ovo je key"

func getJWT(msg string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				},
		Email: msg,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	ss, err := token.SignedString([]byte(myKey))
	if err != nil {
		return "", fmt.Errorf("couldn.t get JWT signed string in NewWithClaims %w", err)
	}

	return ss, nil
}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w,r,"/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("emailThing")
	if email == "" {
		http.Redirect(w,r,"/", http.StatusSeeOther)
		return
	}

	ss, err := getJWT(email)
	if err != nil {
		http.Error(w, "couldn't get JWT token", http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name: "session",
		Value: ss,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	ss := c.Value
	afterVerificationToken, err := jwt.ParseWithClaims(ss, &myClaims{}, func(beforeVerificationToken *jwt.Token) (interface{}, error) {
		if beforeVerificationToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
 			return nil, fmt.Errorf("Someone tried to change sining metod")
		}
		return []byte(myKey), nil
	})

	//standard claims has valid method and on ParseWithClaims Valid() is called and on token Valid field is set
	isEqual := err == nil && afterVerificationToken.Valid

	message := "Not logged in!"
	if isEqual {
		message = "Logged in"
	}

		html := `<!DOCTYPE html>
	<html lang="en"
		<head>
			<meta charset="UTF-8"> 
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>HMAC Example</title>
		</head>
		<body>
			<p>Hi Vuk, cookie is ` + c.Value + `</p>
			<p>` + message + `</p>
			<form action="/submit" method="post">	
				<input type="email" name="emailThing"/>
				<input type="submit" name="Submit"/>
			</form>
		</body>
	</html>`

	io.WriteString(w, html)
}
