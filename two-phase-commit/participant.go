package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
	"io/ioutil"
)

var participantAddress string

func main() {
	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	participantAddress = os.Args[1]

	// value records the proposed value which is later to be committed. It's
	// currently shared among all connections because we don't allow for more
	// than one coordinator in the network.
	var value string

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("error on accept connection: %v\n", err)
			continue
		}
		fmt.Println("accepted connection")

		request, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Errorf("error on read: %v\n", err)
			continue
		}
		fmt.Println("read request", request)

		requestParts := strings.Split(request, " ")
		cmd := requestParts[0]
		switch cmd {
		case "propose":
			fmt.Println("vote")
			vote(conn)
			value = requestParts[1]
		case "commit\n":
			fmt.Println("commit value", value)
			commit(value)
		case "abort\n":
			fmt.Println("abort value")
			abort(&value)
		default:
			fmt.Printf("invalid command: %s\n", cmd)
		}
	}
}

func vote(conn net.Conn) {
	fmt.Fprintf(conn, "Ack\n")
}

func commit(value string) {
	err := ioutil.WriteFile("./data/" + participantAddress, []byte(value), 0444)
	if err != nil {
		fmt.Errorf("commit: error writing file: %w", err)
	}
}

func abort(value *string) {
	*value = ""
}
