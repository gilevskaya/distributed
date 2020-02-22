package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const fileName = "map.db"

func main() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for input.Scan() {
		input := strings.Fields(input.Text())
		if err := processInput(input); err != nil {
			fmt.Println("ERROR HAPPENED:", err.Error())
		}
		fmt.Print("> ")
	}
}

func processInput(words []string) error {
	if len(words) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	action, arg := words[0], words[1]
	switch action {
	case "set":
		pair := strings.Split(arg, "=")
		if len(pair) != 2 {
			return fmt.Errorf("invalid syntax for set")
		}
		addRecord(pair[0], pair[1])
	case "get":
		value, err := findRecord(arg)
		if err != nil {
			return fmt.Errorf("failed to find record: %s", err.Error())
		}
		if value != nil {
			fmt.Println("Your value was:", *value)
		} else {
			fmt.Println("Unable to find your value")
		}

	}
	return nil
}

func getFileRead(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func getFileWrite(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func findRecord(key string) (*string, error) {
	r := getFileRead(fileName)
	defer r.Close()

	scanner := bufio.NewScanner(r)
	var res *string
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) != 2 {
			continue
		}
		if key == parts[0] {
			res = &parts[1]
		}
	}
	return res, scanner.Err()
}

func addRecord(key string, value string) {
	w := getFileWrite(fileName)
	defer w.Close()

	// TODO: handle re-assignment
	fmt.Fprintf(w, "%s=%s\n", key, value)
}
