package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"go-yandex-practicum-metrics/internal/types"
	"net/http"
)

type MetricUpdatesBodyUsecase interface {
	Execute(
		ctx context.Context, req *types.MetricUpdatesBodyRequest,
	) (*types.MetricUpdatesBodyResponse, error)
}

func MetricUpdatesBodyHandler(uc MetricUpdatesBodyUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MetricUpdatesBodyRequest

		err := utils.DecodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			MetricErrorResponse(w, err)
			return
		}

		err = utils.EncodeResponseBody(w, resp)
		if err != nil {
			return
		}
	}
}
