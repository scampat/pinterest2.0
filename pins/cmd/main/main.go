package main

import (
	"fmt"
	"log"
	"net/http"

	"pinterest2.0/pins/internal/controller/pins"
	pinHandler "pinterest2.0/pins/internal/handler/http"
	"pinterest2.0/pins/internal/repository/memory"
)

func main() {
	repo := memory.NewRepository()
	ctrl := pins.NewController(repo)
	handler := pinHandler.NewHandler(ctrl)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	fmt.Println("Pins service running at :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
