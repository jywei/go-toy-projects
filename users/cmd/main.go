package main

import (
	"fmt"
	"net/http"

	"github.com/jywei/toy-projects/users"
)

func main() {
	username, password := "roywjy@gmail.com", "qwerty123"

	err := users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user: %s\n", err.Error())
		return
	}

	err = users.AuthenticateUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't authenticate user: %s\n", err.Error())
		return
	}

	fmt.Printf("Succesfully created and authenticated user \033[32m%s\033[0m\n", username)

	// Send reset email
	err = users.SendPasswordResetEmail(username)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/reset", users.ResetPassword)
	http.ListenAndServe(":3000", nil)
}
