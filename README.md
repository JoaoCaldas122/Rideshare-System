# Rideshare System

This is a simplified rideshare system that integrates Apache Kafka for managing real-time events such as ride requests, driver location updates, and ride matching notifications.

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
└── README.md                   # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd rideshare-system
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

## API Usage

- **User Registration**
  - Endpoint: `POST /users/register`
  - Description: Register a new user.

- **Update Driver Location**
  - Endpoint: `POST /drivers/{id}/location`
  - Description: Update the location of a driver.

- **Request Ride**
  - Endpoint: `POST /rides/request`
  - Description: Request a ride.

- **Fetch Nearest Driver**
  - Endpoint: `GET /rides/nearest-driver`
  - Description: Find the nearest available driver.

## System Architecture

The system is designed to handle real-time events using Apache Kafka. It listens for events related to ride requests and driver updates, processes them, and notifies the relevant parties. The architecture is modular, allowing for easy maintenance and scalability.