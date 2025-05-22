# E-Money Backend

This is the backend for **E-Money**, a digital Monopoly companion that handles player state, room logic, properties, and real-time WebSocket communication. It's built using Go and [Gin](https://github.com/gin-gonic/gin) for routing.

---

## API Endpoints

All API routes are prefixed with `/v1`.

### 📡 WebSocket

- `GET /v1/ws/room/:code`  
  Establishes a WebSocket connection for real-time room updates.

---

### Rooms

- `POST /v1/rooms`  
  Create a new room.

- `GET /v1/rooms/:code/players`  
  Get a list of all players in a room.

- `GET /v1/rooms/:code/properties`  
  Get all available (unowned) properties in a room.

---

### Players

- `POST /v1/rooms/:code/players`  
  Join a room as a new player.

- `GET /v1/rooms/:code/players/:playerId`  
  Get full details about a player.

---

### Properties

- `POST /v1/rooms/:code/players/:playerId/properties/:propertyId`  
  Add a property to a player.

- `DELETE /v1/rooms/:code/players/:playerId/properties/:propertyId`  
  Remove a property from a player.

- `POST /v1/rooms/:code/players/:playerId/properties/:propertyId/mortgage`  
  Mortgage a property for a player.

---

## Local Development

### Prerequisites

- Go 1.20+
- A `.env` file (if needed for config)
- [Gin](https://github.com/gin-gonic/gin)

### Run the server

```bash
go run main.go
