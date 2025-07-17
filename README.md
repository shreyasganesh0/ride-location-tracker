# Real-Time Location Tracker
A real-time location tracking service backend built with Go and WebSockets, inspired by the architecture of ride-sharing apps like Uber. This project is designed to handle multiple concurrent clients (drivers) and broadcast their location data in real-time.

## Current Status
The service is currently capable of:

- Managing Multiple Clients: Accepts and manages numerous concurrent WebSocket connections.

- Structured Data Broadcast: Receives location updates from any client as a JSON object.

- Real-Time Broadcasting: Instantly broadcasts the received location data to all other connected clients.

- Resilient Concurrency: Uses a non-blocking, channel-based architecture to ensure a single slow or disconnected client cannot impact the entire system.

## Getting Started
Follow these instructions to get the project running on your local machine.

## Prerequisites
- Go (Version 1.18+ recommended)
- websocat (A command-line WebSocket client for testing)

## Installation & Running
- Clone the repository:

``` 

git clone https://github.com/your-username/location-tracker.git
cd location-tracker
```

- Run the server:

```
go run .
The server will start on localhost:8080.
```

## Usage & API
The server exposes a single WebSocket endpoint for real-time communication.

WebSocket Endpoint: /ws
This is the primary endpoint for clients to connect and send/receive location updates.

1. Connecting a Client:

Use websocat to establish a persistent connection. Open a new terminal for each client.

```
websocat ws://localhost:8080/ws
```

2. Sending a Location Update:

To simulate a driver sending their location, send a JSON object to the server. The server expects the following format:

JSON

{
  "latitude": 12.345,
  "longitude": 67.890,
  "driverId": "driver-xyz-789"
}
You can send this from a new terminal using echo:


```
echo '{"latitude": 12.345, "longitude": 67.890, "driverId": "driver-xyz-789"}' | websocat -n1 ws://localhost:8080/ws
```

Any client connected (from step 1) will instantly receive this JSON object.

## Architectural Decisions
This section documents key engineering decisions to showcase the thought process behind the design.

Concurrency Model (Channels over Mutexes): We use Go's native channels and goroutines (CSP) to manage state changes within the Hub. This avoids the need for explicit mutex locks, preventing potential deadlocks and making the concurrent code simpler and easier to reason about. All state mutations are handled by a single hub.run() goroutine, ensuring safe, sequential access.

Decoupled Client Pumps: Each client connection spawns two dedicated goroutines (readPump and writePump). This decouples network I/O from the central Hub. The Hub communicates with a client via a buffered channel, so a slow network connection on one client will not block the Hub from broadcasting to others. This makes the system highly resilient.

Structured Data Contract (JSON): We transitioned from raw byte streams to strongly-typed Go structs with json tags. This enforces a clear and versionable API contract, making the system robust against malformed data and easier for new clients (like a future frontend) to integrate with.

Roadmap & Future Milestones
[ ] Persistence: Store the last known location of each driver in a database (e.g., Redis or PostgreSQL).

[ ] Frontend Client: Build a simple web frontend with Leaflet.js or Mapbox to visualize drivers on a map.

[ ] Geospatial Indexing: Implement a method to efficiently query for drivers near a given point.

[ ] Authentication: Secure the WebSocket endpoint so only authenticated drivers can connect and send updates.

[ ] Containerization: Dockerize the application and create a k8s deployment configuration.
