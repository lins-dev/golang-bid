package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/lins-dev/golang-bid.git/internal/jsonutils"
	"github.com/lins-dev/golang-bid.git/internal/services"
	"github.com/lins-dev/golang-bid.git/internal/usecase/product"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)
	slog.Info("log info", "dataReq", data)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userUuid, ok := api.Sessions.Get(r.Context(), "AuthUserUuid").(uuid.UUID)
	slog.Info("log info", "userUuid", userUuid)
	slog.Info("log info", "user uuid ok", ok)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error",
		})
		return
	}

	user, err := api.UserService.FindUserByUuid(r.Context(), userUuid)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
			"error": "user not found",
		})
		return
	}
	slog.Info("log info", "user", user)
	product, err := api.ProductService.CreateProduct(
		r.Context(),
		user.ID,
		data.ProductName,
		data.Description,
		int32(data.Price*100),
		data.AuctionEnd,
	)

	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	auctionRoom := services.NewAuctionRoom(ctx, product.Uuid, product.ID, api.BidService)
	
	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[product.Uuid] = auctionRoom
	api.AuctionLobby.Unlock()

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"msg": "product created and auction has started",
		"product": product,
	})
}