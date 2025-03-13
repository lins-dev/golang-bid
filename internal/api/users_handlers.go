package api

import (
	"errors"
	"net/http"

	"github.com/lins-dev/golang-bid.git/internal/jsonutils"
	"github.com/lins-dev/golang-bid.git/internal/services"
	"github.com/lins-dev/golang-bid.git/internal/usecase/user"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	user, err := api.UserService.CreateUser(
		r.Context(),
		data.UserName,
		data.Email,
		data.Password,
		data.Bio,
	)

	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any {
				"error": "invalid email or user_name",
			})
			return
		}
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any {
		"user": user,
	})

}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	panic("NOT IMPLEMENTED")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("NOT IMPLEMENTED")
}