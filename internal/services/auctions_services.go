package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	// Request
	PlaceBid MessageKind = iota

	// Success
	SuccesfullyPlaceBid
	
	//Error
	FailedToPlaceBid

	//Info
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message string
	Kind MessageKind
	UserUuid uuid.UUID
	UserId int32
	Amount int32

}

type AuctionLobby struct {
	sync.Mutex
	Rooms  map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id uuid.UUID
	ProductId int32
	Register chan *Client
	Unregister chan *Client
	Clients map[uuid.UUID]*Client
	Broadcast chan Message
	Context context.Context

	BidsService BidsService
}

func NewAuctionRoom(ctx context.Context, productUuid uuid.UUID, productId int32,BidService BidsService) * AuctionRoom {
	return &AuctionRoom{
		Id: productUuid,
		ProductId: productId,
		Broadcast: make(chan Message),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[uuid.UUID]*Client),
		Context: ctx,
		BidsService: BidService,
	}
}

type Client struct {
	Room *AuctionRoom
	Conn *websocket.Conn
	Send chan Message
	UserUuid uuid.UUID
	UserId int32
}



func NewClient(room *AuctionRoom, conn *websocket.Conn, userUuid uuid.UUID, userId int32) *Client {
	return &Client{
		Room: room,
		Conn: conn,
		Send: make(chan Message, 512),
		UserUuid: userUuid,
		UserId: userId,
	}
}

func (r *AuctionRoom) registerClient(client *Client) {
	slog.Info("new user connected", "client", client)
	r.Clients[client.UserUuid] = client
}

func (r *AuctionRoom) unregisterClient(client *Client) {
	slog.Info("user has disconnected", "client", client)
	delete(r.Clients, client.UserUuid)
}

func (r *AuctionRoom) broadcastMessage(message Message) {
	slog.Info("new message recieved", "message", message, "room_id", r.Id, "user_uuid", message.UserUuid)
	switch message.Kind {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.ProductId, message.UserId, message.Amount)
		if err != nil {
			if errors.Is(err, ErrNewBidLowerThanBasePrice) || errors.Is(err, ErrNewBidLowerThanPrevious) {
				if client, ok := r.Clients[message.UserUuid]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrNewBidLowerThanBasePrice.Error()}
				}
				return
			}
		}
		if client, ok := r.Clients[message.UserUuid]; ok {
			client.Send <- Message{Kind: SuccesfullyPlaceBid, Message: "your bid was succesfully placed"}
		}
		
		for id, client := range r.Clients{
			newBidMessage := Message{Kind: NewBidPlaced, Message: "a new bid was placed", Amount: bid.BidAmount}
			if id == message.UserUuid {
				continue
			}
			client.Send <- newBidMessage
		}
	}

}



func (r *AuctionRoom) Run() {
	slog.Info("auction room started", "auction_id", r.Id)
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <- r.Register:
			//register client
			r.registerClient(client)
			continue
		case client := <- r.Unregister:
			//unregister client
			r.unregisterClient(client)
			continue
		case message := <- r.Broadcast:
			//send broadcast message
			r.broadcastMessage(message)
			continue
		case <- r.Context.Done():
			slog.Info("auction has ended", "auction_id", r.Id)
			for _, client := range r.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "auction has been finished"}
			}
			return
		}
	}
}