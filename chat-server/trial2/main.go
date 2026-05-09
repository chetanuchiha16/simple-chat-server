package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	// no_of_clients := 0
	// clients := make([]net.Conn, 0)
	var clients []net.Conn // let it not be a pointer, pass by reference instead
	var mu sync.Mutex

	for {
		conn, err := listener.Accept()
		mu.Lock()
		clients = append(clients, conn)
		mu.Unlock()
		if err != nil {
			println(err)
		}
		// no_of_clients += 1

		go handleClient(conn, &clients, &mu)

	}

}

func handleClient(conn net.Conn, clients *[]net.Conn, mu *sync.Mutex) {
	
	mu.Lock()
	fmt.Printf("connected %v\n", len(*clients))
	mu.Unlock()
	scanner := bufio.NewReader(conn)
	for {

		message, err := scanner.ReadString('\n')
		if err != nil {
			println(err)
		}
		fmt.Printf("recieved message: %v\n", message)
		// conn.Write([]byte(fmt.Sprintf("sent message: %v", message)))
		mu.Lock()
		for i, client := range *clients {
			client.Write([]byte(fmt.Sprintf("recieved message %v from another client %d\n", message, i)))
		}
		mu.Unlock()
	}

}
