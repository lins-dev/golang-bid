package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lins-dev/golang-bid.git/internal/jsonutils"
	"github.com/lins-dev/golang-bid.git/internal/services"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductUuid := chi.URLParam(r, "product_id")
	productUuid, err := uuid.Parse(rawProductUuid)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
	}
	
	_, err = api.ProductService.FindProductByUuid(r.Context(), productUuid)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any {
				"error": "product not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error in find user",
		})
		return
	}

	userUuid, ok := api.Sessions.Get(r.Context(), "AuthUserUuid").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error get user uuid",
		})
		return
	}
	// CODE WITH ERROR
	// RESOLVE THIS PROBLEM IS PRIORITY
	user, err := api.UserService.FindUserByUuid(r.Context(), userUuid)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "user not found",
		})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productUuid]
	api.AuctionLobby.Unlock()
	
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"msg" : "the auction has ended",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "could not upgrade connection to a websocket protocol",
			"error": err.Error(),
		})
		return
	}
	client := services.NewClient(room, conn, user.Uuid, user.ID)

	room.Register <- client
	// go client.ReadEventLoop()
	// go client.WriteEventLoop()

	for {

	}

}