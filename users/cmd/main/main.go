package main

import (
	"fmt"
	"log"
	"net/http"

	usercontroller "pinterest2.0/users/internal/controller/user"
	userhttp "pinterest2.0/users/internal/handler/http"
)

func main() {
	ctrl := usercontroller.New(nil, nil)

	h := userhttp.NewHandler(ctrl)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	fmt.Println("Users service running at :8001")
	log.Fatal(http.ListenAndServe(":8001", mux))
}
