package handlers

import (
	"context"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"go-yandex-practicum-metrics/internal/logger"
	"go-yandex-practicum-metrics/internal/usecases"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MetricUpdatePathUsecase interface {
	Execute(ctx context.Context, req *usecases.MetricUpdatePathRequest) (*usecases.MetricUpdatePathResponse, error)
}

func MetricUpdatePathHandler(uc MetricUpdatePathUsecase) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		metricType := utils.GetPathParam(ps, "type")
		metricName := utils.GetPathParam(ps, "name")
		metricValue := utils.GetPathParam(ps, "value")
		logger.Info("Received request to update metric", "type", metricType, "name", metricName, "value", metricValue)
		req := &usecases.MetricUpdatePathRequest{
			Type:  metricType,
			Name:  metricName,
			Value: metricValue,
		}
		resp, err := uc.Execute(r.Context(), req)
		if err != nil {
			logger.Error("Error executing MetricUpdatePath usecase", "error", err)
			metricUpdateHandleError(w, err)
			return
		}
		logger.Info("Metric updated successfully", "name", metricName)
		utils.MakeTextPlainResponse(w, []byte(resp.Message))
	}
}
