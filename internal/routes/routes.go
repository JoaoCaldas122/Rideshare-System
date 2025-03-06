package routes

import (
    "github.com/gorilla/mux"
    "rideshare-system/internal/handlers"
)

func SetupRoutes(router *mux.Router) {
    router.HandleFunc("/users/register", handlers.RegisterUser).Methods("POST")
    router.HandleFunc("/drivers/{id}/location", handlers.UpdateDriverLocation).Methods("POST")
    router.HandleFunc("/rides/request", handlers.RequestRide).Methods("POST")
    router.HandleFunc("/rides/find-driver/{rideRequestId}", handlers.FindNearestDriver).Methods("GET")
}