package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gilevskaya/distributed/lib"
)

func main() {
	dbpath := flag.String("db", "", "location to store your db")
	port := flag.Uint("port", 8080, "port to listen on")
	peersStr := flag.String("peers", "", "comma separated list of peers to connect to")
	flag.Parse()

	if *dbpath == "" {
		log.Fatalln("dbpath is mandatory")
	}

	peers := strings.Split(*peersStr, ",")
	if len(peers) == 0 {
		log.Fatalln("Your system is not decentralized enough.")
	}

	parsedPort := uint32(*port)
	server := lib.NewServer(parsedPort, *dbpath)

	server.Wg.Add(1)
	go server.Listen()

	for _, peer := range peers {
		go server.ConnectToPeer(peer)
	}

	server.REPL()

	server.Wait()
}
