package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

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
	InvalidJSON

	//Info
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message  string      `json:"message,omitempty"`
	Kind     MessageKind `json:"kind,omitempty"`
	UserUuid uuid.UUID   `json:"user_uuid,omitempty"`
	UserId   int32       `json:"user_id,omitempty"`
	Amount   int32       `json:"amount,omitempty"`

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

const (
	maxMessageSize = 512
	readDeadLine = 60 * time.Second
	writeWait = 10 * time.Second
	pingPeriod = (readDeadLine*9)/10 //90% da Read
)

func (c *Client) ReadEventLoop() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(readDeadLine))
		return nil
	})

	for {
		var m Message
		m.UserUuid = c.UserUuid
		m.UserId = c.UserId
		err := c.Conn.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Unexpected Close error", "error", err)
				return
			}

			c.Room.Broadcast <- Message{
				Kind: InvalidJSON,
				Message: "this message should be a valid json",
				UserUuid: m.UserUuid,
				UserId: m.UserId,
			}
			continue
		}

		c.Room.Broadcast <- m
	}
}

func(c *Client) WriteEventLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <- c.Send:
			if !ok {
				c.Conn.WriteJSON(Message{
					Kind: websocket.CloseMessage,
					Message: "closing websocket conn",
				})
				return
			}

			if message.Kind == AuctionFinished {
				close(c.Send)
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteJSON(message)
			if err != nil {
				c.Room.Unregister <- c
				return
			}
			
		case <- ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Unexpected write error", "error", err)
				return
			}
		}
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
		bid, err := r.BidsService.PlaceBid(r.Context, r.ProductId, message.UserId, message.Amount*100)
		if err != nil {
			slog.Info("Log broadcastMessage", "error", err)
			if errors.Is(err, ErrNewBidLowerThanBasePrice) || errors.Is(err, ErrNewBidLowerThanPrevious) {
				if client, ok := r.Clients[message.UserUuid]; ok {
					client.Send <- Message{
						Kind: FailedToPlaceBid, 
						Message: ErrNewBidLowerThanBasePrice.Error(), 
						UserUuid: message.UserUuid, 
						UserId: message.UserId,
					}
				}
				return
			}
		}
		if client, ok := r.Clients[message.UserUuid]; ok {
			client.Send <- Message{
				Kind: SuccesfullyPlaceBid, 
				Message: "your bid was succesfully placed", 
				UserUuid: message.UserUuid, 
				UserId: message.UserId,
			}
		}
		
		for id, client := range r.Clients{
			newBidMessage := Message{
				Kind: NewBidPlaced, 
				Message: "a new bid was placed", 
				Amount: bid.BidAmount,
				UserUuid: message.UserUuid, 
				UserId: message.UserId,
			}
			if id == message.UserUuid {
				continue
			}
			client.Send <- newBidMessage
		}
	case InvalidJSON:
		client, ok := r.Clients[message.UserUuid]
		if !ok {
			slog.Info("Client not found in hashmap", "user_uuid", message.UserUuid)
			return
		}
		client.Send <- message
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