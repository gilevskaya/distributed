package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const fileName = "map.db"

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

// FindRecord does stuff with records
func (server *Server) FindRecord(key string) (*string, error) {
	r := getFileRead(server.dbpath)
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

// AddRecord ...
func (server *Server) AddRecord(key string, value string) {
	w := getFileWrite(server.dbpath)
	defer w.Close()

	// TODO: handle re-assignment
	fmt.Fprintf(w, "%s=%s\n", key, value)
}
