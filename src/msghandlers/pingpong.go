package msghandlers

import (
	"log"
	"protocol"
)

//PingPongHandler handles received Message from the network layer and sends back the response
type PingPongHandler struct {
	messages <-chan protocol.Message
}

//PingPongHandler doesn't actually do anything with the message, it just writes back
//the same thing to the sender
//TODO: as part of this example we need to also use SBE encode/decode, so client should
//encode the message and send it, server (i.e. this handler) should decode it, "read" the payload
//encode the answer (same message) and send it back. In this case server will at least do 1 decode/encode
//per message.
func (h *PingPongHandler) handleMessages() {
	for {
		msg := <-h.messages
		log.Printf("Received message: %s", msg)
		msg.IOWriter.WriteMessage(msg)
	}
}

//Start will start the handler
func (h *PingPongHandler) Start() {
	go h.handleMessages()
}
