package httpmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

type HTTPManager struct {
	mux    *http.ServeMux
	logger logger.Logger
}

func New(manager handler.Manager, logger logger.Logger) *HTTPManager {
	m := &HTTPManager{
		mux:    http.NewServeMux(),
		logger: logger,
	}

	m.routes(handler.New(manager))

	return m
}

func (m *HTTPManager) Start(
	ctx context.Context,
	port int,
) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      m.mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("service listen and serve error")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	//nolint:contextcheck
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
