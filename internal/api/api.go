package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/lins-dev/golang-bid.git/internal/services"
)

type Api struct {
	Router *chi.Mux
	UserService services.UserService
}