package msghandlers

import "protocol"

//StartMessageHandlers will start n message handlers and return
//a list of channels that handlers will receive messages on. We
//use these channels to pass messages from network layer to the handlers.
func StartMessageHandlers(n int) []chan protocol.Message {
	handlerChannels := messageHandlerChannels(n)
	//each handler will read from a single handler channel
	startMessageHandlers(handlerChannels)

	return handlerChannels
}

//create n channels that we'll use when passing messages from network layer to message handling layer
func messageHandlerChannels(n int) []chan protocol.Message {
	channels := make([]chan protocol.Message, n)
	for i := range channels {
		//TODO: buffering?
		channels[i] = make(chan protocol.Message)
	}

	return channels
}

//starts a new goroutine for each channel. This goroutine is a message handler
//that listens to new messages received from the network and processes them
func startMessageHandlers(handlerChannels []chan protocol.Message) {
	for _, ch := range handlerChannels {
		handler := &PingPongHandler{messages: ch}
		handler.Start()
	}
}
