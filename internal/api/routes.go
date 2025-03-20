package api

import (

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID)
	api.Router.Use(middleware.Recoverer)
	api.Router.Use(middleware.Logger)
	api.Router.Use(api.Sessions.LoadAndSave)

	// csrfSecure, _ := strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("CSRF_KEY")),
	// 	csrf.Path("/"),
	// 	csrf.Secure(csrfSecure), //false only in dev
	// )

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// r.Get("/csrftoken", api.HandleGetCsrfTokenMiddleware)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/login", api.handleLoginUser)
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					
					r.Post("/logout", api.handleLogoutUser)
				})
			})
			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)

					r.Post("/", api.handleCreateProduct)
				})
			})
		})
	})
}
