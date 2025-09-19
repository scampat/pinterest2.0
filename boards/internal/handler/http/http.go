package http

import (
	"encoding/json"
	"log"
	"net/http"

	"pinterest2.0/boards/internal/controller/board"
)

type Handler struct {
	ctrl *board.Controller
}

func New(ctrl *board.Controller) *Handler {
	return &Handler{ctrl}
}

// POST /boards → crear board
func (h *Handler) CreateBoard(w http.ResponseWriter, req *http.Request) {
	var input struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
		Title  string `json:"title"`
	}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	b := h.ctrl.CreateBoard(input.ID, input.UserID, input.Title)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}

// GET /boards?user_id=123 → obtener boards de un usuario
func (h *Handler) GetBoardsByUser(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user_id")
	boards, err := h.ctrl.GetBoardsByUser(userID)
	if err != nil {
		http.Error(w, "boards not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(boards); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}

// GET /boards/pins?board_id=123 → obtener pins de un board
func (h *Handler) GetBoardPins(w http.ResponseWriter, req *http.Request) {
	boardID := req.URL.Query().Get("board_id")
	pins, err := h.ctrl.GetBoardPins(req.Context(), boardID)
	if err != nil {
		http.Error(w, "pins not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(pins); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
