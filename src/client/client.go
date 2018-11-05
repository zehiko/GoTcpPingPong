package main

import (
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Print("Error while establishing connection to the server", err)
	}

	var msg []byte
	msg = make([]byte, 12)
	msg[0] = 4
	msg[4] = 1
	content := []byte("test")
	copy(msg[8:], content)
	//TODO: use binary protocol impl on the client side as well and
	//not just make this assumption that server really sent back exactly
	//what client send initially
	//TODO: also remove restraint of just sending 10 messages
	for i := 0; i < 10; i++ {
		conn.Write(msg)
		io.ReadFull(conn, msg[0:])
		log.Printf("Received back: %s", msg[8:])
	}
	conn.Close()

}
