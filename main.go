package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

func getCode(msg string) string {
	h := hmac.New(sha256.New, []byte("ovo je key"))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w,r,"/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w,r,"/", http.StatusSeeOther)
		return
	}

	code := getCode(email)

	c := http.Cookie{
		Name: "session",
		Value: code + "|" + email,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	message := "Not logged in!"

	xs := strings.SplitN(c.Value, "|", 2)
	if len(xs) == 2 {
		cCode := xs[0]
		cEmail := xs[1]

		code := getCode(cEmail)

		if hmac.Equal([]byte(cCode), []byte(code)) {
			message = "Logged in!"
		} else {
			message = "Not logged in!"
		}
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
				<input type="email" name="email"/>
				<input type="submit" name="Submit"/>
			</form>
		</body>
	</html>`

	io.WriteString(w, html)
}
