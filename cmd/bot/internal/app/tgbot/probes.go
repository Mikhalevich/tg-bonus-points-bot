package tgbot

import (
	"fmt"
	"net/http"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (t *TGBot) httpLivenessProbe() http.HandlerFunc {
	return doProbe(t.livenessProbe, "liveliness probe")
}

func (t *TGBot) httpReadinessProbe() http.HandlerFunc {
	return doProbe(t.readinessProbe, "rediness probe")
}

func (t *TGBot) SetLivenessProbe(probe Probe) {
	t.livenessProbe = probe
}

func (t *TGBot) SetReadinessProbe(probe Probe) {
	t.readinessProbe = probe
}

func doProbe(probe Probe, logDescription string) http.HandlerFunc {
	return http.HandlerFunc(
		//nolint:varnamelen
		func(w http.ResponseWriter, r *http.Request) {
			if probe == nil {
				logger.FromContext(r.Context()).
					WithField("probe", logDescription).
					Warn("no probe setup")

				return
			}

			if err := probe(r.Context()); err != nil {
				logger.FromContext(r.Context()).
					WithField("probe", logDescription).
					WithError(err).
					Info("probe failed")

				http.Error(w, "DOWN", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "UP")
		},
	)
}
