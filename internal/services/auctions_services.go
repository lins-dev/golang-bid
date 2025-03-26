package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota
)

type Message struct {
	Message string
	Kind MessageKind

}

type AuctionRoom struct {
	Id uuid.UUID
	Register chan *Client
	Unregister chan *Client
	Clients map[uuid.UUID]*Client
	Broadcast chan Message
	Context context.Context

	BidsService BidsService
}

type AuctionLobby struct {
	Rooms  map[uuid.UUID]*AuctionRoom
	sync.Mutex
}

type Client struct {
	Conn *websocket.Conn
	UserUuid uuid.UUID
	Send chan Message
	Room *AuctionRoom
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, BidService BidsService) * AuctionRoom {
	return &AuctionRoom{
		Id: id,
		Broadcast: make(chan Message),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Context: ctx,
		BidsService: BidService,
	}
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userUuid uuid.UUID) *Client {
	return &Client{
		Room: room,
		Conn: conn,
		Send: make(chan Message, 512),
		UserUuid: userUuid,
	}
}