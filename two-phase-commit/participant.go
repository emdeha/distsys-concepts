package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

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
		case "commit\n":
			fmt.Println("commit value")
		case "abort\n":
			fmt.Println("abort value")
		default:
			fmt.Printf("invalid command: %s\n", cmd)
		}
	}
}

func vote(conn net.Conn) {
	fmt.Fprintf(conn, "Ack\n")
}
