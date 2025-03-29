package utils

import (
	"go-yandex-practicum-metrics/internal/logger"

	"github.com/julienschmidt/httprouter"
)

func GetPathParam(ps httprouter.Params, paramName string) string {
	paramValue := ps.ByName(paramName)
	logger.Info("Retrieved path parameter", "param", paramName, "value", paramValue)
	return paramValue
}
