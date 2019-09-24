package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	phonebook := make(map[string]string)
	scanner.Scan()
	num, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < num; i++ {
		scanner.Scan()
		record := strings.Split(scanner.Text(), " ")
		phonebook[record[0]] = record[1]
	}

	for scanner.Scan() {
		name := scanner.Text()
		if v, ok := phonebook[name]; ok {
			fmt.Printf("%v=%v\n", name, v)
		} else {
			fmt.Println("Not found")
		}
	}
}
