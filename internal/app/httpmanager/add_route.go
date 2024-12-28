package httpmanager

import (
	"encoding/json"
	"net/http"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/httpmanager/internal/httperror"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/tracing"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) *httperror.ErrorHTTPResponse

type errorBody struct {
	Message string `json:"message"`
}

func (m *HTTPManager) addRoute(pattern string, hf handlerFunc) {
	m.mux.HandleFunc(pattern, m.makeHTTPHandler(pattern, hf))
}

func (m *HTTPManager) makeHTTPHandler(pattern string, hf handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.StartSpanName(r.Context(), pattern)
		defer span.End()

		var (
			log    = m.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog = logger.WithLogger(ctx, log)
		)

		if err := hf(w, r.WithContext(ctxLog)); err != nil {
			if err.Err != nil {
				log.WithError(err.Err).
					Error("http handler error")
			}

			w.WriteHeader(err.Code)

			if err := json.NewEncoder(w).Encode(
				errorBody{
					Message: err.Message,
				},
			); err != nil {
				log.WithError(err).Error("encode json body error")
			}

			return
		}
	}
}
