package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akashkncse/xotp/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewServer()

	log.Printf("xotp API listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Handler()))
}
