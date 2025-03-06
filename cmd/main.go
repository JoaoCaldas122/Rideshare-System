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
    models.InitDB() // Initialize the database

    brokerAddress := "localhost:9092" // Replace with your Kafka broker address
    kafka.InitKafka(brokerAddress)    // Initialize Kafka

    router := mux.NewRouter()
    routes.SetupRoutes(router) // Setup API routes

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))

    defer kafka.Close() // Ensure Kafka writer is closed properly
}