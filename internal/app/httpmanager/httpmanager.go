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

const (
	readTimeout      = time.Second * 10
	writeTimeout     = time.Second * 10
	shoutdownTimeout = time.Second * 30
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

	httpManager := &HTTPManager{
		mux:     mux,
		humaAPI: humaAPI,
		logger:  logger,
	}

	httpManager.routes(handler.New(manager))

	return httpManager
}

func (m *HTTPManager) Start(
	ctx context.Context,
	port int,
) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      m.mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("service listen and serve error")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shoutdownTimeout)
	defer cancel()

	//nolint:contextcheck
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
