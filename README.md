# Exercise for Backend Software Engineer: Simulating a Rideshare System with Apache Kafka Integration (Go)

This is a simplified rideshare system that integrates Apache Kafka using Go for managing real-time events such as ride requests, driver location updates, and ride matching notifications.

## Project Structure

```
rideshare-system
├── cmd
│   └── main.go                # Entry point of the application
├── internal
│   ├── consumer
│   │   └── consumer.go        # Kafka consumer implementation
│   ├── handlers
│   │   └── handlers.go        # HTTP handler functions for API endpoints
│   ├── kafka
│   │   └── kafka.go           # Kafka producer setup and event publishing
│   ├── models
│   │   └── models.go          # Data models for users, drivers, and rides
│   ├── routes
│   │   └── routes.go          # API routes and routing logic
│   └── utils
│       └── utils.go           # Utility functions for various operations
├── go.mod                      # Go module configuration file
└── README.md                   # Project documentation and explanation
```

## Setup Instructions

1. *Clone the repository:*
   
   git clone <repository-url>
   cd rideshare-system
   

2. *Install dependencies:*
   
   go mod tidy
   

3. *Run the application:*
   
   go run cmd/main.go
   

## API Usage

- *User Registration*
  - Endpoint: POST /users/register
  - Description: Register a new user.

- *Update Driver Location*
  - Endpoint: POST /drivers/{id}/location
  - Description: Update the location of a driver.

- *Request Ride*
  - Endpoint: POST /rides/request
  - Description: Request a ride.

- *Fetch Nearest Driver*
  - Endpoint: GET /rides/nearest-driver
  - Description: Find the nearest available driver.

## Application Development

The first thing that was made to start developing the rideshare system was to understand the problem. Understand that there are two agents: a rider and a driver, and a rider can't ride another user to any place. After that, A small drawing was made to simulate an environment where there were riders and drivers in different places(coordinates) and write up all the different events that could and needed to happen. If user A(rider) requested for a ride to go from location A(0,0) to location B(1,2), a ride request should be done, then the nearest driver should be found, by calculating the distance from each available driver's location to the rider that made the request. After the driver was found, a ride was created to take user A from location A to B, driven by user C with an "ongoing" status. Followed by that, the driver(User C) should update his location and move to user's A location to pick him up and then both of them should travel to location B and have their location updated again. Once they got to the final location, the ride was completed and its status was updated to "complete".

Given this, I thought of creating 3 main structures: User, Ride and Ride Request. The User would have an id(which was his primary key), a name, a role(either rider or driver), and the user location(latitude and longitude). The Ride Request would have its id too, the rider id, the matched the rider that called, the initial latitude and longitude and the final latitude and longitude. The Ride would have an id, the rider id, the driver id, the initial and final location and the ride status, that can be ongoing or complete.
Below there are the 3 structs created in the code:

```
type User struct {
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    Role      string  `json:"role"` // "driver" or "rider"
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

type RideRequest struct {
    ID             int     `json:"id"`
    RiderID        int     `json:"rider_id"`
    InitialLat     float64 `json:"initial_latitude"`
    InitialLon     float64 `json:"initial_longitude"`
    FinalLat       float64 `json:"final_latitude"`
    FinalLon       float64 `json:"final_longitude"`
}

type Ride struct {
    ID             int     `json:"id"`
    RiderID        int     `json:"rider_id"`
    DriverID       int     `json:"driver_id"`
    InitialLat     float64 `json:"initial_latitude"`
    InitialLon     float64 `json:"initial_longitude"`
    FinalLat       float64 `json:"final_latitude"`
    FinalLon       float64 `json:"final_longitude"`
    Status         string  `json:"status"` // "ongoing", "completed"
}
```

## Consumer

The consumer file was made to initialize the Kafka consumer and process the messages sent to it.

## Handlers

Here were developed the main functions to react to the posts made to the servers, such as registering a user, update a drivers location, make a ride request and the function to find the nearest driver. All these handlers add the respective request to the database.

## Kafka

This file writes and sends all the Kafka messages of any ride request or notifications, and sends a log to the terminal to check if it happened

## Models 

This file is where the database is created and initialized, such as its tables that are users, rides and ride_requests. There are also queries and functions to populate and insert the data into the tables when a request is made, and by the end the 3 structures that I have explained before.

## Routes

This file has the routes to handle the requests to the server and their format according to the REST API.

## Utils 

This file has a function that calculates the distance between two points on the Earth specified by latitude and longitude(Haversine formula) and is used to find the nearest driver.