package main

import (
	"context"
	"log"
	"net/http"

	"pinterest2.0/boards/internal/controller/board"
	"pinterest2.0/boards/internal/gateway/pins"
	handler "pinterest2.0/boards/internal/handler/http"
	"pinterest2.0/pkg/discovery/memory"
	"pinterest2.0/pkg/registry"
)

const serviceName = "boards"

func main() {
	// Crear el registry en memoria
	reg := memory.NewRegistry()

	// Generar un ID único para esta instancia
	instanceID := registry.GenerateInstanceID(serviceName)

	// Registrar este servicio
	if err := reg.Register(context.Background(), instanceID, serviceName, "localhost:8081"); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	// Gateway a pins
	pinsGateway := pins.New(reg)

	// Controller
	ctrl := board.New(pinsGateway)

	// HTTP handler
	h := handler.New(ctrl)

	http.HandleFunc("/boards", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateBoard(w, r)
		} else if r.Method == http.MethodGet {
			h.GetBoardsByUser(w, r)
		}
	})

	http.HandleFunc("/boards/pins", h.GetBoardPins)

	log.Println("Boards service running on :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
