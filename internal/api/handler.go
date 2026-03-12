package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LoomHubDev/loomhub/internal/config"
	"github.com/LoomHubDev/loomhub/internal/store"
	"github.com/LoomHubDev/loomhub/internal/sync"
	"gorm.io/gorm"
)

type Handler struct {
	cfg           *config.Config
	db            *gorm.DB
	owners        *store.OwnerStore
	users         *store.UserStore
	looms         *store.LoomStore
	collabs       *store.CollaboratorStore
	weaveRequests *store.WeaveRequestStore
	comments      *store.CommentStore
	reviews       *store.ReviewStore
	pins          *store.PinStore
	activities    *store.ActivityStore
	orgs          *store.OrgStore
	webhooks      *store.WebhookStore
	labels        *store.LabelStore
	search        *store.SearchStore
	tokens        *store.TokenStore
	sync          *sync.Handler
}

func NewHandler(cfg *config.Config, db *gorm.DB) *Handler {
	// Extract raw *sql.DB for sync handler (per-loom SQLite databases)
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get underlying sql.DB: %v", err))
	}

	return &Handler{
		cfg:           cfg,
		db:            db,
		owners:        store.NewOwnerStore(db),
		users:         store.NewUserStore(db),
		looms:         store.NewLoomStore(db),
		collabs:       store.NewCollaboratorStore(db),
		weaveRequests: store.NewWeaveRequestStore(db),
		comments:      store.NewCommentStore(db),
		reviews:       store.NewReviewStore(db),
		pins:          store.NewPinStore(db),
		activities:    store.NewActivityStore(db),
		orgs:          store.NewOrgStore(db),
		webhooks:      store.NewWebhookStore(db),
		labels:        store.NewLabelStore(db),
		search:        store.NewSearchStore(db),
		tokens:        store.NewTokenStore(db),
		sync:          sync.NewHandler(sqlDB, cfg.Server.DataDir),
	}
}

func (h *Handler) Close() {
	h.sync.Close()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
