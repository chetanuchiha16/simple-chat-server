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
	no_of_clients := 0
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			println(err)
		}
		no_of_clients += 1

		go handleClient(conn, no_of_clients)
		
	}

}

func handleClient(conn net.Conn, no_of_clients int) {
	
	fmt.Printf("connected %v\n", no_of_clients)
	scanner := bufio.NewReader(conn)
	message, err := scanner.ReadString('\n')
	if err != nil {
		println(err)
	}
	fmt.Printf("recieved message: %v", message)
	conn.Write([]byte(fmt.Sprintf("sent message: %v", message)))
	conn.Close()
	
}