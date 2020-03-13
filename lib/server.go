package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Server struct
type Server struct {
	running  uint32
	port     uint32
	outgoing map[string]net.Conn
	incoming map[string]net.Conn
	Wg       sync.WaitGroup
	dbpath   string
}

// NewServer ...
func NewServer(port uint32, dbpath string) Server {
	return Server{
		incoming: make(map[string]net.Conn, 0),
		outgoing: make(map[string]net.Conn, 0),
		Wg:       sync.WaitGroup{},
		port:     port,
		running:  0,
		dbpath:   dbpath,
	}
}

func (server *Server) handlePeerConn(conn net.Conn) {
	// Make a buffer to hold incoming data.
	// buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	// reqLen, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }
	// Send a response back to person contacting us.``
	conn.Write([]byte("Message received."))

	// fmt.Println(reqLen)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	fmt.Println(rw)

	// Close the connection when you're done with it.
	defer conn.Close()
	for {
		cmd, err := rw.ReadString('\n')
		log.Print("Receive command")

		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		log.Println(cmd)
	}
}

// ConnectToPeer connects to peers
func (server *Server) ConnectToPeer(peer string) error {
	server.Wg.Add(1)

	var conn net.Conn
	var err error
	for attempt := 0; attempt < 5; attempt++ {
		log.Printf("Attempt %d to connect to %s", attempt, peer)
		conn, err = net.Dial("tcp", peer)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if conn != nil {
			break
		}
	}
	if conn == nil {
		return fmt.Errorf("all attempts failed, last error: %v", err)
	}

	server.outgoing[conn.RemoteAddr().String()] = conn
	go server.handlePeerConn(conn)
	return nil
}

// Listen listens
func (server *Server) Listen() error {
	atomic.StoreUint32(&server.running, 1)
	server.Wg.Add(1)

	iface := fmt.Sprintf("0.0.0.0:%d", server.port)

	ln, err := net.Listen("tcp", iface)
	if err != nil {
		return err
	}
	log.Printf("Listening on interface %s", iface)
	defer ln.Close()

	for atomic.LoadUint32(&server.running) > 0 {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		remoteEndpoint := conn.RemoteAddr().String()
		server.incoming[remoteEndpoint] = conn

		// Handle connections in a new goroutine.
		go server.handlePeerConn(conn)
	}

	return nil
}

// Wait waits until server is done handling all connections including listening
func (server *Server) Wait() {
	server.Wg.Wait()
}

// REPL is the repl loop
func (server *Server) REPL() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for input.Scan() {
		input := strings.Fields(input.Text())
		if err := server.processInput(input); err != nil {
			fmt.Println("ERROR HAPPENED:", err.Error())
		}
		fmt.Print("> ")
	}
}

func (server *Server) processInput(words []string) error {
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
		server.AddRecord(pair[0], pair[1])
	case "get":
		value, err := server.FindRecord(arg)
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
