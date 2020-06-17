package router

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func httpJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		httpInternalServerError(w, "failed to encode response body", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("failed to write response")
	}
}

func httpBadRequest(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusBadRequest, msg, err)
}

func httpInternalServerError(w http.ResponseWriter, msg string, err error) {
	httpError(w, http.StatusInternalServerError, msg, err)
}

func httpError(w http.ResponseWriter, httpStatus int, msg string, err error) {
	http.Error(w, msg, httpStatus)
	zap.L().Warn(msg, zap.Error(err))
}
