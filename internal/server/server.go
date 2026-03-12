package server

import (
	"database/sql"
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/api"
	"github.com/LoomHubDev/loomhub/internal/auth"
	"github.com/LoomHubDev/loomhub/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	router chi.Router
	cfg    *config.Config
	db     *sql.DB
}

func New(cfg *config.Config, db *sql.DB) *Server {
	s := &Server{cfg: cfg, db: db}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(auth.Middleware(s.cfg.Auth.JWTSecret, s.db))

	h := api.NewHandler(s.cfg, s.db)

	// REST API
	r.Route("/api/v1", func(r chi.Router) {
		// Auth
		r.Post("/auth/login", h.Login)
		r.Post("/users", h.Register)

		// Authenticated user
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Get("/user", h.GetCurrentUser)
			r.Patch("/user", h.UpdateCurrentUser)
			r.Get("/user/looms", h.ListCurrentUserLooms)
			r.Get("/user/orgs", h.ListUserOrgs)
			r.Get("/user/feed", h.GetActivityFeed)
			r.Post("/looms", h.CreateLoom)
			r.Post("/orgs", h.CreateOrg)

			// Access tokens
			r.Get("/user/tokens", h.ListAccessTokens)
			r.Post("/user/tokens", h.CreateAccessToken)
			r.Delete("/user/tokens/{tokenID}", h.DeleteAccessToken)
		})

		// Public user routes
		r.Get("/users/{username}", h.GetUser)
		r.Get("/users/{username}/looms", h.ListUserLooms)
		r.Get("/users/{username}/activity", h.GetUserActivity)

		// Activity
		r.Get("/activity", h.GetGlobalActivity)

		// Organization routes
		r.Get("/orgs/{org}", h.GetOrg)
		r.Get("/orgs/{org}/members", h.ListOrgMembers)
		r.Get("/orgs/{org}/looms", h.ListOrgLooms)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Patch("/orgs/{org}", h.UpdateOrg)
			r.Post("/orgs/{org}/members", h.AddOrgMember)
			r.Delete("/orgs/{org}/members/{username}", h.RemoveOrgMember)
		})

		// Loom routes
		r.Get("/looms/{owner}/{loom}", h.GetLoom)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Patch("/looms/{owner}/{loom}", h.UpdateLoom)
			r.Delete("/looms/{owner}/{loom}", h.DeleteLoom)

			// Pin
			r.Put("/looms/{owner}/{loom}/pin", h.PinLoom)
			r.Delete("/looms/{owner}/{loom}/pin", h.UnpinLoom)
		})

		// Weave requests
		r.Get("/looms/{owner}/{loom}/weave-requests", h.ListWeaveRequests)
		r.Get("/looms/{owner}/{loom}/weave-requests/{number}", h.GetWeaveRequest)
		r.Get("/looms/{owner}/{loom}/weave-requests/{number}/comments", h.ListComments)
		r.Get("/looms/{owner}/{loom}/weave-requests/{number}/reviews", h.ListReviews)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Post("/looms/{owner}/{loom}/weave-requests", h.CreateWeaveRequest)
			r.Patch("/looms/{owner}/{loom}/weave-requests/{number}", h.UpdateWeaveRequest)
			r.Post("/looms/{owner}/{loom}/weave-requests/{number}/weave", h.WeaveWR)
			r.Post("/looms/{owner}/{loom}/weave-requests/{number}/close", h.CloseWR)
			r.Post("/looms/{owner}/{loom}/weave-requests/{number}/comments", h.CreateComment)
			r.Post("/looms/{owner}/{loom}/weave-requests/{number}/reviews", h.CreateReview)
		})

		// Labels
		r.Get("/looms/{owner}/{loom}/labels", h.ListLabels)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Post("/looms/{owner}/{loom}/labels", h.CreateLabel)
			r.Patch("/looms/{owner}/{loom}/labels/{labelID}", h.UpdateLabel)
			r.Delete("/looms/{owner}/{loom}/labels/{labelID}", h.DeleteLabel)
			r.Post("/looms/{owner}/{loom}/labels/{labelID}/assign", h.AddLabelToWR)
		})

		// Webhooks
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Get("/looms/{owner}/{loom}/webhooks", h.ListWebhooks)
			r.Post("/looms/{owner}/{loom}/webhooks", h.CreateWebhook)
			r.Patch("/looms/{owner}/{loom}/webhooks/{webhookID}", h.UpdateWebhook)
			r.Delete("/looms/{owner}/{loom}/webhooks/{webhookID}", h.DeleteWebhook)
			r.Get("/looms/{owner}/{loom}/webhooks/{webhookID}/deliveries", h.ListWebhookDeliveries)
		})

		// Collaborators
		r.Get("/looms/{owner}/{loom}/collaborators", h.ListCollaborators)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Post("/looms/{owner}/{loom}/collaborators", h.AddCollaborator)
			r.Delete("/looms/{owner}/{loom}/collaborators/{username}", h.RemoveCollaborator)
		})

		// Content browsing
		r.Get("/looms/{owner}/{loom}/streams", h.ListStreams)
		r.Get("/looms/{owner}/{loom}/entities", h.ListEntities)
		r.Get("/looms/{owner}/{loom}/checkpoints", h.ListCheckpoints)
		r.Get("/looms/{owner}/{loom}/objects/{ref}", h.GetEntityContent)
		r.Get("/looms/{owner}/{loom}/compare", h.CompareStreams)

		// Search
		r.Get("/search/looms", h.SearchLooms)
		r.Get("/search/users", h.SearchUsers)
		r.Get("/explore", h.Explore)

		// Admin
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Get("/admin/stats", h.AdminStats)
			r.Get("/admin/users", h.AdminListUsers)
			r.Get("/admin/looms", h.AdminListLooms)
		})
	})

	// Sync API (per-loom)
	r.Route("/{owner}/{loom}/api/v1", func(r chi.Router) {
		r.Post("/negotiate", h.HandleNegotiate)
		r.Post("/push", h.HandleSend)
		r.Post("/pull", h.HandleReceive)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	s.router = r
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
