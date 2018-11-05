package main

import (
	"log"
	"math/rand"
	"msghandlers"
	"net"
	"protocol"
	"protocol/binary"
	"runtime"
)

const (
	headerLength = 8
)

func main() {
	//TODO args
	log.Println("starting listener on port 8080")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listen failed", err)
	}
	log.Println("listener started")

	//gorouting accepting new connection will send them to this channel
	//our main goroutine listens on this channel and starts per-channel handling
	//goroutines
	conns := make(chan net.Conn)

	//start goroutine that will accept new connections from the client
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("failed to accept new connection", err)
			}
			log.Printf("received new connection from %s", conn.RemoteAddr())
			conns <- conn
		}
	}()

	//we'll have twice as many message handlers (i.e. message handling goroutines) as we have CPUs
	//unlike connection handlers where we have 1 per connection, we keep number of message handlers
	//limited. This reasoning comes from the assumption that since we're doing blocking I/O, we want
	//1 goroutine per connection and these will be scheduled off quickly once we block waiting for I/O.
	//For message handler - these will handle our app logic, constantly receving messages from many
	//different connections and doing some "cpu intesive work" (not really in this light example),
	//so we keep number of these small (2x num of CPUs of the machine)
	handlerChannels := msghandlers.StartMessageHandlers(2 * runtime.NumCPU())

	//main goroutine receives new connections from the goroutine started above
	//and starts a new handling goroutine for every connection
	for {
		conn := <-conns
		log.Printf("preparing to handle connection from %s", conn.RemoteAddr())
		protocolCodec := binary.NewSimpleIOCodec(conn)
		//some really basic multiplexing where we push messages from many different connections
		//onto 1 message handler
		connHandlerPipe := handlerChannels[rand.Intn(2*runtime.NumCPU())]
		go handleConnection(protocolCodec, connHandlerPipe)
	}
}

func handleConnection(p protocol.IOCodec, handlersPipe chan<- protocol.Message) {
	for {
		msg, err := p.ReadMessage()
		//stop reading the messages if there's any error
		if err != nil {
			log.Println(err)
			return
		}
		handlersPipe <- msg
	}
}
