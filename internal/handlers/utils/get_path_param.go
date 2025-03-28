package utils

import "github.com/julienschmidt/httprouter"

func GetPathParam(ps httprouter.Params, paramName string) string {
	return ps.ByName(paramName)
}
