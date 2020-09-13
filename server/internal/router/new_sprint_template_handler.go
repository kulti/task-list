package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kulti/task-list/server/internal/models"
)

type sprintTempateStore interface {
	GetNewSprintTemplate(context.Context) (models.SprintTemplate, error)
	SetNewSprintTemplate(context.Context, models.SprintTemplate) error
}

type sprintTemplateHandler struct {
	sprintTempateStore sprintTempateStore
}

func newSprintTemplateHandler(sprintTempateStore sprintTempateStore) sprintTemplateHandler {
	return sprintTemplateHandler{
		sprintTempateStore: sprintTempateStore,
	}
}

func (h sprintTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p, _ := shiftPath(r.URL.Path); p != "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetSprintTemplate(w, r)
	case http.MethodPost:
		h.handleSetSprintTemplate(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h sprintTemplateHandler) handleGetSprintTemplate(w http.ResponseWriter, r *http.Request) {
	template, err := h.sprintTempateStore.GetNewSprintTemplate(r.Context())
	if err != nil {
		httpInternalServerError(w, "failed to get sprint template", err)
		return
	}

	httpJSON(w, &template)
}

func (h sprintTemplateHandler) handleSetSprintTemplate(w http.ResponseWriter, r *http.Request) {
	jsDecoder := json.NewDecoder(r.Body)

	var template models.SprintTemplate
	err := jsDecoder.Decode(&template)
	if err != nil {
		httpBadRequest(w, "failed to parse sprint template body", err)
		return
	}

	err = h.sprintTempateStore.SetNewSprintTemplate(r.Context(), template)
	if err != nil {
		httpInternalServerError(w, "failed to store sprint template in body", err)
		return
	}
}
