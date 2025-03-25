package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"net/http"
)

type MetricListHTMLUsecase interface {
	Execute(ctx context.Context) (string, error)
}

func MetricListHTMLHandler(uc MetricListHTMLUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := uc.Execute(r.Context())
		if err != nil {
			MetricErrorResponse(w, err)
			return
		}

		err = utils.RenderHTML(w, resp)
		if err != nil {
			return
		}
	}
}
