package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	//p1 := person{
	//	First: "Vuk",
	//}
	//p2 := person{
	//	First: "knez Lazar",
	//}
	//
	//xp := []person{p1, p2}
	//bs, err := json.Marshal(&xp)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//fmt.Println("PRINT JSON: " + string(bs))
	//
	//xp2 := []person{}
	//err = json.Unmarshal(bs, &xp2)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//fmt.Println("Back in to go struct: ", xp2)
	//
	//http.HandleFunc("/encode", foo)
	//http.HandleFunc("/decode", bar)
	//http.ListenAndServe(":8080", nil)

	pass := "123456789"

	hashedPass, err := hashPassword(pass)
	if err != nil {
		panic(err)
	}

	err = comparePassword(pass, hashedPass)
	if err != nil {
		log.Fatalln("Not logged in")
	}

	log.Println("Logged in!")
}

func foo(w http.ResponseWriter, r *http.Request) {
	p1 := person{
		First: "Vuk",
	}

	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("encoded bad data: ", err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	var p1 person
	err := json.NewDecoder(r.Body).Decode(&p1)
	if err != nil {
		log.Println("encoded bad data: ", err)
	}

	log.Println("Person decoded: ", p1)
}

func hashPassword(password string) ([]byte, error) {
	bs, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while generationg bcrypt from password: %w", err)
	}
	return bs, nil
}

func comparePassword(password string, hashedPass []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid password: %w", err)
	}
	return nil
}
