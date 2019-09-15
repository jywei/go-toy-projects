package main

import (
	"fmt"
	"os"
)

const (
	user     = "jack"
	pwd      = "1988"
	usgae    = "Usage: [username] [password]"
	errUser  = "Access denied for %q.\n"
	errPwd   = "Invalid passowrd for %q.\n"
	accessOK = "Access granted to %q.\n"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println(usgae)
		return
	}

	u, p := args[1], args[2]
	if u != user {
		fmt.Printf(errUser, u)
	} else if p != pwd {
		fmt.Printf(errPwd, u)
	} else {
		fmt.Printf(accessOK, u)
	}
}
