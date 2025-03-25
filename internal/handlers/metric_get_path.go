package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"go-yandex-practicum-metrics/internal/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricGetPathUsecase interface {
	Execute(
		ctx context.Context, req *types.MetricGetByTypeAndIDPathRequest,
	) (*types.MetricGetByTypeAndIDPathResponse, error)
}

func MetricGetPathHandler(uc MetricGetPathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		metricType := chi.URLParam(r, "type")

		req := &types.MetricGetByTypeAndIDPathRequest{
			Name: name,
			Type: metricType,
		}

		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			MetricErrorResponse(w, err)
			return
		}

		utils.SendTextResponse(w, string(*resp))
	}
}
