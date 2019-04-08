package main

import (
	"fmt"

	"github.com/jywei/toy-projects/users"
)

func main() {
	username, password := "roy", "qwerty123"

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

	fmt.Printf("Successfully created and authenticated user %s", username)

	err = users.NewUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't create user: %s\n", err.Error())
		return
	}

	err = users.AuthenticateUser(username, password)
	if err != nil {
		fmt.Printf("Couldn't authenticate user: %s\n", err.Error())
		return
	}
}
