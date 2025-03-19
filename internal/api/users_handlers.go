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
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
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
	data, problems, err := jsonutils.DecodeValidJson[user.LoginUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
	}
	user, err :=  api.UserService.AuthUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any {
				"error": "invalid email or password",
			})
			return
		}
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any {
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any {
			"error": "unexpected internal server error",
		})
		return
	}
	api.Sessions.Put(r.Context(), "AuthUserUuid", user.Uuid)
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"msg": "login successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthUserUuid")
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any {
		"msg": "logout successfully",
	})
}