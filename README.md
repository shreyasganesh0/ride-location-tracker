# Real-Time Ride Location Tracker (Go)

This is a high-performance backend system built in Go that simulates a real-time ride-sharing location tracker (like for Uber or DoorDash).

Drivers can publish their location, and clients can subscribe to a real-time feed of all driver movements within their vicinity.

This project was built to explore the challenges of real-time, stateful, and concurrent systems. The goal was to build the infrastructure for a "fan-out" system from scratch, using modern Go practices and a high-performance database (Redis).

## Core Features

* **Real-Time WebSocket API:** A WebSocket endpoint (`/ws`) allows clients to establish a persistent connection to receive live location updates.
* **In-Memory Pub/Sub Hub:** A custom, concurrent-safe Pub/Sub "Hub" was built from scratch using Go channels. This hub manages all active WebSocket clients, handling client registration, de-registration, and the broadcasting of messages.
* **Redis Geospatial Storage:** Driver locations are stored in **Redis** using the `GEOADD` command. This allows for incredibly fast geospatial queries.
* **Geofencing/Proximity Queries:** A REST endpoint (`/api/drivers`) that uses the Redis `GEORADIUS` command to find all drivers within a specified radius of a given latitude/longitude.
* **Containerized:** Fully containerized with a `docker-compose.yml` for easy local development, spinning up both the Go application and its Redis dependency.

## Technical Deep Dive

### 1. Concurrent Pub/Sub Hub

The core of this application is the `pubsub.Hub`. It runs as a single goroutine and uses channels to safely coordinate concurrent access from many different client connections. This avoids the need for complex mutexes.

* `register chan`: New clients are sent here.
* `unregister chan`: Disconnected clients are sent here.
* `broadcast chan`: Messages sent here are "fanned-out" to all connected clients.

This code from `internal/pubsub/hub.go` shows the main event loop of the hub:

```go
// From: internal/pubsub/hub.go
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Register a new client and add it to the map
			h.clients[client] = true
		
		case client := <-h.unregister:
			// Unregister a client, close its send channel
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		
		case message := <-h.broadcast:
			// Broadcast a message to all connected clients
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Failed to send, assume client is dead
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
```

2. WebSocket Connection Handling
The ws.go handler is responsible for upgrading an HTTP request to a persistent WebSocket connection. It then creates a Client struct and registers it with the Hub. This seamlessly connects the networking layer (HTTP/WebSockets) to the application's concurrency model (the Pub/Sub hub).

```Go

// From: internal/handler/ws.go
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		// ... error handling ...
		return
	}

	// Create a new client for this WebSocket connection
	client := pubsub.NewClient(h.Hub, conn)
	
	// Register the new client with the central hub
	h.Hub.Register(client)

	// Start the read/write goroutines for this client
	go client.WritePump()
	go client.ReadPump()
}
```

## How to Run
1. Make sure you have Docker and Docker Compose installed.
2. Clone the repository: ```git clone https://github.com/shreyasganesh0/ride-location-tracker.git```
3. Start the application and the Redis database:
```Bash
docker-compose up
```
4. The server will be running on http://localhost:8080.
  . REST API is at /api/...
  . WebSocket endpoint is at /ws
