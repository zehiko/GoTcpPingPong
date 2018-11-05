package protocol

import (
	"fmt"
)

//Message is what we pass between client and server. It contains the size of the payload
//and the payload (raw bytes)
type Message struct {
	Size    int
	Payload []byte
	//TODO: Not too happy about sending "writer" as part of the message
	//Maybe handlers should simply have access to separate writer and then that layer
	//would handle finding proper writer/connection for that specific handler?
	IOWriter IOWrite
}

func (m Message) String() string {
	return fmt.Sprintf("size: %d, content: %s", m.Size, m.Payload)
}

//IORead supports reading of Messages from the network
type IORead interface {
	ReadMessage() (Message, error)
}

//IOWrite supports writing Messages to the network
type IOWrite interface {
	WriteMessage(Message) error
}

//IOCodec is a simple interface that will hide away interactions with the
//network (i.e. network connection) from upper layers of the application.
//Implementation is expected to support specific protocol (encode / decode
//raw bytes from the network into/from Messages)
type IOCodec interface {
	IORead
	IOWrite
}
