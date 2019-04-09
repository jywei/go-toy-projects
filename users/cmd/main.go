package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jywei/toy-projects/users"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.New("login").Parse(loginTemplate)
		t.Execute(w, nil)
	case "POST":
		user := r.FormValue("user")
		pass := r.FormValue("password")
		err := users.AuthenticateUser(user, pass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users.SetSession(w, user)
		w.Write([]byte("Signed in successfully"))
	}

}

func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetSession(w, r)
	w.Write([]byte(user))
}

func main() {
	username, password := "roywjy@gmail.com", "qwerty123"

	err := users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user: %s\n", err.Error())
		return
	}

	fmt.Printf("Succesfully created and authenticated user \033[32m%s\033[0m\n", username)

	http.HandleFunc("/", authHandler)
	http.HandleFunc("/restricted", restrictedHandler)

	http.ListenAndServe(":3000", nil)
}
