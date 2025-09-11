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
	reg := memory.NewRegistry()

	instanceID := registry.GenerateInstanceID(serviceName)

	if err := reg.Register(context.Background(), instanceID, serviceName, "localhost:8081"); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	pinsGateway := pins.New(reg)

	ctrl := board.New(pinsGateway)

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
