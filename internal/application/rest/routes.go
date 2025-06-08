package rest

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"fmt"
	"time"
)

func RegisterRoutes(
	r chi.Router,
	serverAdmin ServerAdmin,
) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           int((300 * time.Second).Seconds()),
	}))


	r.Route("/admin/v1", func(r chi.Router) {
		r.Post("/matrix", handler(serverAdmin.PostAdminV1CompetencyMatrix))
		r.Get("/matrix", handler(serverAdmin.GetAdminV1CompetencyMatrix))
	})
}

func handler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			fmt.Errorf("Handler error: %w", err)
		}
	}
}
