package main

import (
    "log"
    "net/http"
    "rideshare-system/internal/kafka"
    "rideshare-system/internal/models"
    "rideshare-system/internal/routes"
    "github.com/gorilla/mux"
)

func main() {
    models.InitDB()

    brokerAddress := "localhost:9092"
    kafka.InitKafka(brokerAddress)

    router := mux.NewRouter()
    routes.SetupRoutes(router)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))

    defer kafka.Close()
}