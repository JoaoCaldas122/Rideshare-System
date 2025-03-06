package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "rideshare-system/internal/kafka"
    "rideshare-system/internal/utils"
    "rideshare-system/internal/models"
    "strconv"

    "github.com/gorilla/mux"
)

// RegisterUser handles user registration.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    models.AddUser(user) // Add user to the database
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// UpdateDriverLocation handles updating the driver's location.
func UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    driverID, _ := strconv.Atoi(vars["id"])

    var driver models.User
    if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    driver.ID = driverID
    models.UpdateUserLocation(driver) // Update driver's location in the database

    // Publish driver location update to Kafka
    driverLocation, _ := json.Marshal(driver)
    if err := kafka.PublishDriverLocationUpdate(driverLocation); err != nil {
        log.Printf("Failed to publish driver location update: %v", err)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(driver)
}

// RequestRide handles ride requests from users.
func RequestRide(w http.ResponseWriter, r *http.Request) {
    var rideRequest models.RideRequest
    if err := json.NewDecoder(r.Body).Decode(&rideRequest); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    models.AddRideRequest(rideRequest) // Add ride request to the database

    // Publish ride request to Kafka
    rideRequestMessage, _ := json.Marshal(rideRequest)
    if err := kafka.PublishRideRequest(rideRequestMessage); err != nil {
        log.Printf("Failed to publish ride request: %v", err)
    }

    // Find the nearest driver
    driver, err := models.FindNearestDriver(rideRequest.InitialLat, rideRequest.InitialLon)
    if err != nil {
        http.Error(w, "No available drivers found", http.StatusNotFound)
        return
    }

    // Create a new ride
    ride := models.Ride{
        RiderID:    rideRequest.RiderID,
        DriverID:   driver.ID,
        InitialLat: rideRequest.InitialLat,
        InitialLon: rideRequest.InitialLon,
        FinalLat:   rideRequest.FinalLat,
        FinalLon:   rideRequest.FinalLon,
        Status:     "ongoing",
    }
    models.AddRide(ride)

    // Update driver's location to rider's initial location
    driver.Latitude = rideRequest.InitialLat
    driver.Longitude = rideRequest.InitialLon
    models.UpdateUserLocation(driver)

    // Publish driver location update to Kafka
    driverLocation, _ := json.Marshal(driver)
    if err := kafka.PublishDriverLocationUpdate(driverLocation); err != nil {
        log.Printf("Failed to publish driver location update: %v", err)
    }

    // Simulate the ride completion by updating the driver's location to the final location
    driver.Latitude = rideRequest.FinalLat
    driver.Longitude = rideRequest.FinalLon
    models.UpdateUserLocation(driver)

    // Publish driver location update to Kafka
    driverLocation, _ = json.Marshal(driver)
    if err := kafka.PublishDriverLocationUpdate(driverLocation); err != nil {
        log.Printf("Failed to publish driver location update: %v", err)
    }

    // Update rider's location to the final location
    rider := models.User{
        ID:        rideRequest.RiderID,
        Latitude:  rideRequest.FinalLat,
        Longitude: rideRequest.FinalLon,
    }
    models.UpdateUserLocation(rider)

    // Update ride status to "completed"
    ride.Status = "completed"
    models.UpdateRide(ride)

    // Publish ride completion notification to Kafka
    rideCompletionMessage, _ := json.Marshal(ride)
    if err := kafka.PublishNotification(rideCompletionMessage); err != nil {
        log.Printf("Failed to publish ride completion notification: %v", err)
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(ride)
}

// FindNearestDriver handles finding the nearest driver for a ride request and logs the action.
func FindNearestDriver(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    rideRequestID, _ := strconv.Atoi(vars["rideRequestId"])

    var rideRequest models.RideRequest
    if err := models.DB.QueryRow("SELECT id, rider_id, initial_latitude, initial_longitude, final_latitude, final_longitude FROM ride_requests WHERE id = ?", rideRequestID).Scan(&rideRequest.ID, &rideRequest.RiderID, &rideRequest.InitialLat, &rideRequest.InitialLon, &rideRequest.FinalLat, &rideRequest.FinalLon); err != nil {
        http.Error(w, "Ride request not found", http.StatusNotFound)
        return
    }

    driver, err := models.FindNearestDriver(rideRequest.InitialLat, rideRequest.InitialLon)
    if err != nil {
        http.Error(w, "No available drivers found", http.StatusNotFound)
        return
    }

    log.Printf("Nearest driver found for ride request %d: Driver ID %d, Distance %.2f km", rideRequestID, driver.ID, utils.Haversine(rideRequest.InitialLat, rideRequest.InitialLon, driver.Latitude, driver.Longitude))

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(driver)
}