package httpmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/handler"
)

func Start(
	ctx context.Context,
	port int,
	manager handler.Manager,
) error {
	var (
		mux = http.NewServeMux()
		h   = handler.New(manager)
	)

	mux.HandleFunc("GET /next-order-to-process", h.GetNextPendingOrderToProcess)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	//nolint:errcheck
	go srv.ListenAndServe()

	<-ctx.Done()

	//nolint:contextcheck
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
