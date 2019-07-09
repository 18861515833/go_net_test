package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func login(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Println("method:", request.Method)
	if request.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(writer, nil)
	} else if request.Method == "POST" {
		fmt.Println("username:", request.Form["username"])
		fmt.Println("password:", request.Form["password"])
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.ListenAndServe(":9090", nil)
}
