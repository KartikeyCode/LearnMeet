package server

import (
	"log"
	"net/http"
	"server/internal/server/objects"
	"server/pkg/packets"
)

// A structure for a state machine to process the client's messages
type ClientStateHandler interface {
	Name() string

	//Inject the client into the state handler
	SetClient(client ClientInterfacer)

	OnEnter()
	HandleMessage(senderId uint64, message packets.Msg)

	//cleanup the state handler and perform any last actions
	OnExit()
}

type ClientInterfacer interface {

	//return client id
	Id() uint64

	//Process the message
	ProcessMessage(senderId uint64, message packets.Msg)

	//set client Id and other things that need to be initialized
	Initialize(id uint64)

	//sets the state
	SetState(newState ClientStateHandler)

	// puts data from client into write pump
	SocketSend(message packets.Msg)

	//puts data from another client into the write pump
	SocketSendAs(message packets.Msg, senderId uint64)

	// Forward message to another client for processing
	PassToPeer(message packets.Msg, peerId uint64)

	// Forward message to all other clients for processing
	Broadcast(message packets.Msg)

	// close client conncetion and cleanup
	Close(reason string)

	// pump data from socket to client
	ReadPump()

	// pump data from client to socket
	WritePump()
}

// the hub is the central point of communication between all connected clients
type Hub struct {
	Clients *objects.SharedCollection[ClientInterfacer]

	//packets in this channel will be processed by all connected clients except the sender
	BroadcastChan chan *packets.Packet

	//clients in this channel will be registered to the hub
	RegisterChan chan ClientInterfacer

	// clients in this channel will be unregistered from the hub
	UnregisterChan chan ClientInterfacer
}

func NewHub() *Hub {
	return &Hub{
		Clients:        objects.NewSharedCollection[ClientInterfacer](),
		BroadcastChan:  make(chan *packets.Packet),
		RegisterChan:   make(chan ClientInterfacer),
		UnregisterChan: make(chan ClientInterfacer),
	}
}

func (h *Hub) Run() {
	log.Println("Awaiting client registration")
	for {
		select {
		case client := <-h.RegisterChan:
			client.Initialize(h.Clients.Add(client))
		case client := <-h.UnregisterChan:
			h.Clients.Remove(client.Id())
		case packet := <-h.BroadcastChan:
			h.Clients.ForEach(func(clientId uint64, client ClientInterfacer) {
				if clientId != packet.SenderId {
					client.ProcessMessage(packet.SenderId, packet.Msg)
				}
			})
		}
	}
}

func (h *Hub) Serve(getNewClient func(*Hub, http.ResponseWriter, *http.Request) (ClientInterfacer, error), writer http.ResponseWriter, request *http.Request) {
	log.Println("New Client connected from", request.RemoteAddr)
	client, err := getNewClient(h, writer, request)

	if err != nil {
		log.Printf("Error Occurred: %v", err)
		return
	}

	h.RegisterChan <- client
	//using go keyword allows us to do a process in background thread
	go client.WritePump()
	go client.ReadPump()
}
