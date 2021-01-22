package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	p1 := person{
		First: "Vuk",
	}
	p2 := person{
		First: "knez Lazar",
	}

	xp := []person{p1, p2}
	bs, err := json.Marshal(&xp)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("PRINT JSON: " + string(bs))

	xp2 := []person{}
	err = json.Unmarshal(bs, &xp2)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Back in to go struct: ", xp2)

	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	http.ListenAndServe(":8080", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

}

func bar(w http.ResponseWriter, r *http.Request) {

}
