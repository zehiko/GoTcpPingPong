# GoTcpPingPong

This is a simple minimalistic example of sending the same message back and forth (over a TCP connection) between a client and the server in Go. 

### Client:
* connect to the server
* create the message

[in loop]

* send the message to the server
* read expected reply 

### Server:
* start listening on predefined port
* start new connection handler for each new connection
* connection handler reads messages expecting them to be in predefined protocol format
* connection handler passes received message to message handling layer (we keep limited amount of message handling goroutines)
* message handling layer processes the message (do nothing right now) and sends back the 
same thing back to the sender/client

# TODO
* use protoocl on the client side
* remove logging everywhere, especially on the data path
* add argument parsing so that we can chose server port, etc. 
* add SBE encode/decode to do at least some work on the server side
* buffer - right now message handlers receive message over channels that don't have any buffering
so we can add that to avoid blocking network senders due to temporary "processing overload" on the 
handler side