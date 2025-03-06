package models

import (
    "database/sql"
    "log"
    "math"
    "rideshare-system/internal/utils"
    "sync"

    _ "github.com/mattn/go-sqlite3"
)

var (
    DB   *sql.DB
    once sync.Once
)

func InitDB() {
    once.Do(func() {
        var err error
        DB, err = sql.Open("sqlite3", "./rideshare.db")
        if err != nil {
            log.Fatalf("Failed to open database: %v", err)
        }

        createTables()
        log.Println("Database initialized")
    })
}

func createTables() {
    createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        role TEXT,
        latitude REAL,
        longitude REAL
    );`

    createRideRequestTable := `
    CREATE TABLE IF NOT EXISTS ride_requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        rider_id INTEGER,
        initial_latitude REAL,
        initial_longitude REAL,
        final_latitude REAL,
        final_longitude REAL,
        FOREIGN KEY(rider_id) REFERENCES users(id)
    );`

    createRideTable := `
    CREATE TABLE IF NOT EXISTS rides (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        rider_id INTEGER,
        driver_id INTEGER,
        initial_latitude REAL,
        initial_longitude REAL,
        final_latitude REAL,
        final_longitude REAL,
        status TEXT,
        FOREIGN KEY(rider_id) REFERENCES users(id),
        FOREIGN KEY(driver_id) REFERENCES users(id)
    );`

    log.Println("Creating users table...")
    if _, err := DB.Exec(createUserTable); err != nil {
        log.Fatalf("Failed to create users table: %v", err)
    }
    log.Println("Users table created successfully")

    log.Println("Creating ride_requests table...")
    if _, err := DB.Exec(createRideRequestTable); err != nil {
        log.Fatalf("Failed to create ride_requests table: %v", err)
    }
    log.Println("Ride requests table created successfully")

    log.Println("Creating rides table...")
    if _, err := DB.Exec(createRideTable); err != nil {
        log.Fatalf("Failed to create rides table: %v", err)
    }
    log.Println("Rides table created successfully")
}

func AddUser(user User) {
    stmt, err := DB.Prepare("INSERT INTO users(name, role, latitude, longitude) VALUES(?, ?, ?, ?)")
    if err != nil {
        log.Fatalf("Failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(user.Name, user.Role, user.Latitude, user.Longitude)
    if err != nil {
        log.Fatalf("Failed to execute statement: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Fatalf("Failed to retrieve last insert ID: %v", err)
    }
    user.ID = int(id)

    log.Printf("User added: %+v\n", user)
}

func AddRideRequest(rideRequest RideRequest) {
    stmt, err := DB.Prepare("INSERT INTO ride_requests(rider_id, initial_latitude, initial_longitude, final_latitude, final_longitude) VALUES(?, ?, ?, ?, ?)")
    if err != nil {
        log.Fatalf("Failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(rideRequest.RiderID, rideRequest.InitialLat, rideRequest.InitialLon, rideRequest.FinalLat, rideRequest.FinalLon)
    if err != nil {
        log.Fatalf("Failed to execute statement: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Fatalf("Failed to retrieve last insert ID: %v", err)
    }
    rideRequest.ID = int(id)

    log.Printf("Ride request added: %+v\n", rideRequest)
}

func AddRide(ride Ride) {
    stmt, err := DB.Prepare("INSERT INTO rides(rider_id, driver_id, initial_latitude, initial_longitude, final_latitude, final_longitude, status) VALUES(?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        log.Fatalf("Failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    result, err := stmt.Exec(ride.RiderID, ride.DriverID, ride.InitialLat, ride.InitialLon, ride.FinalLat, ride.FinalLon, ride.Status)
    if err != nil {
        log.Fatalf("Failed to execute statement: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Fatalf("Failed to retrieve last insert ID: %v", err)
    }
    ride.ID = int(id)

    log.Printf("Ride added: %+v\n", ride)
}

func UpdateRide(ride Ride) {
    stmt, err := DB.Prepare("UPDATE rides SET status = ? WHERE id = ?")
    if err != nil {
        log.Fatalf("Failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(ride.Status, ride.ID)
    if err != nil {
        log.Fatalf("Failed to execute statement: %v", err)
    }

    log.Printf("Ride updated: %+v\n", ride)
}

func FindNearestDriver(lat, lon float64) (User, error) {
    var nearestDriver User
    var minDistance float64 = math.MaxFloat64

    rows, err := DB.Query("SELECT id, name, role, latitude, longitude FROM users WHERE role = 'driver'")
    if err != nil {
        return nearestDriver, err
    }
    defer rows.Close()

    for rows.Next() {
        var driver User
        if err := rows.Scan(&driver.ID, &driver.Name, &driver.Role, &driver.Latitude, &driver.Longitude); err != nil {
            return nearestDriver, err
        }

        distance := utils.Haversine(lat, lon, driver.Latitude, driver.Longitude)
        if distance < minDistance {
            minDistance = distance
            nearestDriver = driver
        }
    }

    if err := rows.Err(); err != nil {
        return nearestDriver, err
    }

    return nearestDriver, nil
}

func UpdateUserLocation(user User) {
    stmt, err := DB.Prepare("UPDATE users SET latitude = ?, longitude = ? WHERE id = ?")
    if err != nil {
        log.Fatalf("Failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(user.Latitude, user.Longitude, user.ID)
    if err != nil {
        log.Fatalf("Failed to execute statement: %v", err)
    }

    log.Printf("User location updated: %+v\n", user)
}

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