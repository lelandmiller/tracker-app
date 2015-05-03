package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type SessionStorer interface {
	Login(int) (string, error)
	Logout(string) error
	Check(string) (int, error)
}

type Handler func(http.ResponseWriter, *http.Request)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!, %s!", r.URL.Path[1:])
}

func user(users map[string]User) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["username"]
		switch r.Method {
		case "GET":
			fmt.Fprintf(w, "Hi there!, your name is %s! Here is your info: \n%+v", name, users[name])
		case "POST":
			var u User
			b, _ := ioutil.ReadAll(r.Body)
			log.Printf("%s", b)
			json.Unmarshal(b, &u)
			log.Print(u)
			users[name] = u
		}
	}
}

func loginHander() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		name, pass, ok := r.BasicAuth()
		if !ok {
			log.Printf("login request ignored")
		}
		GenAuth(w, r)
		log.Print(r)

		log.Printf("token request for %s with pass %s", name, pass)
	}
}

func main() {
	Init()
	users := make(map[string]User)
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/auth/", loginHander())
	r.HandleFunc("/user/{username}/", user(users))
	//r.HandleFunc("/login/", login)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "HELLO!")
	})
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
