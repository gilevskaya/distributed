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
		processInput(input)
		fmt.Print("> ")
	}
}

func processInput(words []string) {
	if len(words) == 2 && (words[0] == "get" || words[0] == "set") {
		action := words[0]
		arg := words[1]

		switch action {
		case "set":
			pair := strings.Split(arg, "=")
			if len(pair) == 2 {
				addRecord(pair[0], pair[1])
				return
			}
		case "get":
			value := findRecord(arg)
			if value != "" {
				fmt.Println(arg, " => ", value)
				return
			}
		}
	}
	fmt.Println("hmm...")
}

func getFileRead(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func getFileWrite(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func findRecord(key string) (value string) {
	r := getFileRead(fileName)
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if key == parts[0] {
			return parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ""
}

func addRecord(key string, value string) {
	w := getFileWrite(fileName)
	defer w.Close()
	// TODO: handle re-assignment
	fmt.Fprintln(w, key+"="+value)
}
