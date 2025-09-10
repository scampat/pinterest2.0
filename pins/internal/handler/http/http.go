package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pinterest2.0/pins/internal/controller/pins"
	"pinterest2.0/pins/pkg"
)

type Handler struct {
	ctrl *pins.Controller
}

func NewHandler(ctrl *pins.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pins", h.handlePins)
}

func (h *Handler) handlePins(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createPin(w, r)
	case http.MethodGet:
		h.getPins(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createPin(w http.ResponseWriter, r *http.Request) {
	var req pkg.Pin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	p, err := h.ctrl.CreatePin(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to create pin", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) getPins(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	pins, err := h.ctrl.GetPinsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get pins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pins)
}
