package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	// no_of_clients := 0
	clients := make([]net.Conn, 0)

	for {
		conn, err := listener.Accept()
		clients = append(clients, conn)
		if err != nil {
			println(err)
		}
		// no_of_clients += 1

		go handleClient(conn, clients)

	}

}

func handleClient(conn net.Conn, clients []net.Conn) {

	fmt.Printf("connected %v\n", len(clients))
	scanner := bufio.NewReader(conn)
	for {

		message, err := scanner.ReadString('\n')
		if err != nil {
			println(err)
		}
		fmt.Printf("recieved message: %v\n", message)
		// conn.Write([]byte(fmt.Sprintf("sent message: %v", message)))
		for i, client := range clients {
			client.Write([]byte(fmt.Sprintf("recieved message %v from another client %d\n", message, i)))
		}
	}

}
