package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lins-dev/golang-bid.git/internal/api"
	"github.com/lins-dev/golang-bid.git/internal/services"
)

func main()  {
	
	gob.Register(uuid.UUID{})

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	sessions :=  scs.New()
	sessions.Store = pgxstore.New(pool)
	sessions.Lifetime = 12 * time.Hour
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode

	api := api.Api {
		Router: chi.NewMux(),
		UserService: services.NewUserService(pool),
		Sessions: sessions,
		ProductService: services.NewProductService(pool),
		BidService: services.NewBidService(pool),
		WsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		AuctionLobby: services.AuctionLobby{
			Rooms: make(map[uuid.UUID]*services.AuctionRoom),
		},
	}

	api.BindRoutes()

	fmt.Println("starting server on port :3080")

	if err := http.ListenAndServe("localhost:3080", api.Router); err != nil {
		panic(err)
	}
}

// run Air
// air --build.cmd "go build -o ./bin/api ./cmd/api" --build.bin "./bin/api"

// create migrations
// in folder: internal/store/pgstore/migrations
// tern new create_users_table

// generate SQLc files
// sqlc generate -f ./internal/store/pgstore/sqlc.yml

// run migration
// go run ./cmd/terndotenv