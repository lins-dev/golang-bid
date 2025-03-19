package api

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/lins-dev/golang-bid.git/internal/jsonutils"
)

func (api *Api) AuthMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthUserUuid") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"msg": "must be logged",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *Api) HandleGetCsrfTokenMiddleware (w http.ResponseWriter, r *http.Request)  {
	token := csrf.Token(r)
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}