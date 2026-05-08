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
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			println(err)
		}

		go handleClient(conn)
		
	}

}

func handleClient(conn net.Conn) {
	
	fmt.Println("connected")
	scanner := bufio.NewReader(conn)
	message, err := scanner.ReadString('\n')
	if err != nil {
		println(err)
	}
	fmt.Printf("recieved message: %v", message)
	conn.Write([]byte(fmt.Sprintf("sent message: %v", message)))
	conn.Close()
	
}