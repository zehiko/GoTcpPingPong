package binary

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"protocol"
)

const (
	headerLength    = 8
	protocolVersion = 1
)

//SimpleIOCodec is an implementation of our little simple protocol.
type SimpleIOCodec struct {
	conn net.Conn //connection to receive from / write to
	buff []byte   //byte buffer used for reading. we create 1 upfront to avoid creating on each new message
}

//ReadMessage will read a Message from the underyling network connection. It knows
//the expected format of the protocol. It will extract the payload from the passed
//message and return a new Message
func (p *SimpleIOCodec) ReadMessage() (protocol.Message, error) {
	var n int
	var err error
	//read header which is 8 bytes, 4 bytes for msg size and 4 bytes protocol version (this could obviously be smaller)
	_, err = io.ReadFull(p.conn, p.buff[:headerLength])
	//handle any error
	if err != nil {
		if err == io.EOF {
			log.Printf("%s disconnected", p.conn.RemoteAddr())
		} else {
			log.Println("Error on connection, closing connection", err)
		}
		p.conn.Close()
		return protocol.Message{}, err
	}

	//TODO: check protocol version, check msg size is not over limit

	//extract expected message size
	msgContentSize := int(binary.LittleEndian.Uint32(p.buff[:4]))
	log.Printf("Receiving a message of size %d", msgContentSize)

	//try to read message size bytes
	n, err = io.ReadFull(p.conn, p.buff[headerLength:headerLength+msgContentSize])
	//again handle errors
	if err != nil {
		if err == io.EOF {
			log.Printf("%s disconnected while reading full message", p.conn.RemoteAddr())
		} else {
			log.Println("Error on connection while reading full message, closing connection", err)
		}
		p.conn.Close()
		return protocol.Message{}, err
	}

	if n < msgContentSize {
		log.Printf("Failed to read entire message. Expected %d, read %d", msgContentSize, n)
	}
	return protocol.Message{Size: msgContentSize,
		Payload: p.buff[headerLength : headerLength+msgContentSize], IOWriter: p}, nil
}

//WriteMessage will extract payload/bytes from the Message, wrap it into expected formatting
//and send it to the other side of the connection
func (p *SimpleIOCodec) WriteMessage(message protocol.Message) error {
	var msg []byte
	msg = make([]byte, headerLength+message.Size)
	assignIntToSliceStart(message.Size, msg[0:])
	assignIntToSliceStart(protocolVersion, msg[4:])
	copy(msg[8:], message.Payload)

	p.conn.Write(msg)
	return nil
}

//there's got to be a simpler i.e. builtin solution for this
func assignIntToSliceStart(n int, sl []byte) {
	sl[0] = (byte(n & 0xFF))
	sl[1] = (byte(n>>8) & 0xFF)
	sl[2] = (byte(n>>8) & 0xFF)
	sl[3] = (byte(n>>8) & 0xFF)
}

//NewSimpleIOCodec is a small helper function for creating binary protocol "instances"
func NewSimpleIOCodec(conn net.Conn) *SimpleIOCodec {
	return &SimpleIOCodec{conn: conn, buff: make([]byte, 4096)}
}
