package main

import (
	"fmt"
	"os"
)

const (
	user, user2 = "jack", "roy"
	pwd, pwd2   = "1988", "1989"
	usgae       = "Usage: [username] [password]"
	errUser     = "Access denied for %q.\n"
	errPwd      = "Invalid passowrd for %q.\n"
	accessOK    = "Access granted to %q.\n"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println(usgae)
		return
	}

	u, p := args[1], args[2]
	if u != user && u != user2 {
		fmt.Printf(errUser, u)
	} else if (u == user && p == pwd) || (u == user2 && p == pwd2) {
		fmt.Printf(accessOK, u)
	} else {
		fmt.Printf(errPwd, u)
	}
}
