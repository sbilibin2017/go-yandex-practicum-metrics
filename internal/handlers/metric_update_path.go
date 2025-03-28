package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"go-yandex-practicum-metrics/internal/requests"
	"go-yandex-practicum-metrics/internal/responses"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MetricUpdatePathUsecase interface {
	Execute(ctx context.Context, req *requests.MetricUpdatePathRequest) (*responses.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		metricType := utils.GetPathParam(ps, "type")
		metricName := utils.GetPathParam(ps, "name")
		metricValue := utils.GetPathParam(ps, "value")
		req := &requests.MetricUpdatePathRequest{
			Type:  metricType,
			Name:  metricName,
			Value: metricValue,
		}
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			metricUpdateHandleError(w, err)
			return
		}
		utils.MakeTextPlainResponse(w, *resp.ToResponse())
	}
}
