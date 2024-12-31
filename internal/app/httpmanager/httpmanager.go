package httpmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
)

type HTTPManager struct {
	mux     *http.ServeMux
	humaAPI huma.API
	logger  logger.Logger
}

func New(manager handler.Manager, logger logger.Logger) *HTTPManager {
	var (
		mux     = http.NewServeMux()
		humaAPI = humago.New(mux, huma.DefaultConfig("Bonus points", "1.0.0"))
	)

	m := &HTTPManager{
		mux:     mux,
		humaAPI: humaAPI,
		logger:  logger,
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
