package handlers

import (
	"context"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PingHandler struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPingHandler(pool *pgxpool.Pool, ctx context.Context) *PingHandler {
	return &PingHandler{
		pool: pool,
		ctx:  ctx,
	}
}

func (h *PingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if h.pool == nil {
		log.Error("Database not accessable")
		applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)

		return
	}
	status := h.pool.Ping(h.ctx)
	if status != nil {
		log.Error("Database not accessable: ", status)
		applicationerrors.SwitchError(status, &w)

		return
	}

	w.WriteHeader(http.StatusOK)
}
