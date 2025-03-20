package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/lins-dev/golang-bid.git/internal/services"
)

type Api struct {
	Router *chi.Mux
	UserService services.UserService
	Sessions *scs.SessionManager
	ProductService services.ProductService
	WsUpgrader websocket.Upgrader
}