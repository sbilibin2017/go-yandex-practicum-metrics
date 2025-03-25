package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"go-yandex-practicum-metrics/internal/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MetricUpdatePathUsecase interface {
	Execute(
		ctx context.Context, req *types.MetricUpdatePathRequest,
	) (*types.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameParam := chi.URLParam(r, "name")
		typeParam := chi.URLParam(r, "type")
		valueParam := chi.URLParam(r, "value")

		req := types.MetricUpdatePathRequest{
			Type:  typeParam,
			Name:  nameParam,
			Value: valueParam,
		}

		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			MetricErrorResponse(w, err)
			return
		}

		utils.SendTextResponse(w, string(*resp))

	}
}
