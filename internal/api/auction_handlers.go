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
	
	_, err = api.UserService.FindUserByUuid(r.Context(), productUuid)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any {
				"error": "product not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error",
		})
		return
	}

	userUuid, ok := api.Sessions.Get(r.Context(), "AuthUserUuid").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal error",
		})
	}
	// CODE WITH ERROR
	// RESOLVE THIS PROBLEM IS PRIORITY
	user, err := api.UserService.FindUserByUuid(r.Context(), userUuid)
	connection, err := api.WsUpgrader.Upgrade(w, r, nil)
	defer close(connection)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "could not upgrade connection to a websocket",
		})
		return
	}



}