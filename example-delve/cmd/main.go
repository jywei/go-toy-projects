package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type user struct {
	Name string
	Age  int
}

func main() {
	r := strings.NewReader(`
      {
        "Name": "Roy",
        "Age": 30
      }
  `)
	var v user

	json.NewDecoder(r).Decode(&v)
	fmt.Printf("%v\n", v)

}
