package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pinterest2.0/users/internal/controller/user"
	"pinterest2.0/users/pkg/model"
)

type Handler struct {
	ctrl *user.Controller
}

func NewHandler(ctrl *user.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/users/boards", h.getUserBoards)
	mux.HandleFunc("/users/pins", h.getUserPins)
}

// ---- USERS ----

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodGet:
		h.getUser(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	u := h.ctrl.CreateUser(r.Context(), req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u.ToResponse())
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	u, err := h.ctrl.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u.ToResponse())
}

// ---- INTEGRACIÓN CON BOARDS ----

func (h *Handler) getUserBoards(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	boards, err := h.ctrl.GetUserBoards(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// ---- INTEGRACIÓN CON PINS ----

func (h *Handler) getUserPins(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	pins, err := h.ctrl.GetUserPins(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pins)
}
