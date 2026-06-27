# Complaint Portal — Go REST API

A concurrency-safe HTTP JSON API built with **pure Go** (zero third-party packages).
Users can submit complaints; administrators can review and resolve them.

---

## Project Structure

```
complaint-portal/
├── main.go                   # Entry point — server bootstrap & route registration
├── go.mod                    # Module definition
├── models/
│   └── models.go             # Structs: User, Complaint, all request/response types
├── store/
│   └── store.go              # In-memory data store (sync.RWMutex — concurrency safe)
└── handlers/
    └── handlers.go           # HTTP handler for every route
```

### Layer Responsibilities

```
HTTP Request
     │
     ▼
handlers/      ← decode JSON, validate input, write JSON response
     │
     ▼
store/         ← business logic + thread-safe in-memory state
     │
     ▼
models/        ← shared data types used across layers
```

---

## Prerequisites & Setup

### 1. Install Go

**Windows**  
Download the installer from https://go.dev/dl/ and run it.

**macOS**

```bash
brew install go
```

**Ubuntu / Debian**

```bash
sudo apt update && sudo apt install -y golang-go
```

Verify:

```bash
go version
```

### 2. Clone / Download the project

```bash
# If using git
git clone https://github.com/Shreya-Ankolnerkar/Complaint-portal
cd complaint-portal

# Or just unzip and enter the folder
cd complaint-portal
```

### 3. Run the server

```bash
go run main.go
```

You should see:

```
Complaint Portal API running on http://localhost:8080
Admin Secret Code: ADMIN-SECRET-2024
─────────────────────────────────────────
```

---

## Authentication Model

- Every user gets a **unique `secret_code`** generated at registration.
- All protected endpoints require `secret_code` in the request body.
- The **admin** uses a fixed secret (`ADMIN-SECRET-2024`) configured in `main.go`.
- Regular users can only access their own complaints; the admin can access everything.

---

## API Reference

All endpoints accept and return **JSON**.  
Method: **POST** for all routes.

---

### `POST /register`

Create a new user.

**Request**

```json
{
  "name": "Shreya Sharma",
  "email": "shreya@example.com"
}
```

**Response**

```json
{
  "success": true,
  "data": {
    "id": "a1b2c3d4e5f6...",
    "secret_code": "deadbeef1234...",
    "name": "Shreya Sharma",
    "email": "shreya@example.com",
    "complaints": []
  }
}
```

---

### `POST /login`

Login with secret code.

**Request**

```json
{ "secret_code": "deadbeef1234..." }
```

**Response** — full user record (same shape as register).

---

### `POST /submitComplaint`

Submit a complaint (user only).

**Request**

```json
{
  "secret_code": "deadbeef1234...",
  "title": "AC not working in Lab 3",
  "summary": "The air conditioning has been broken for 2 weeks.",
  "severity": 4
}
```

> `severity` is an integer **1–5** (1 = low, 5 = critical).

**Response**

```json
{
  "success": true,
  "data": {
    "id": "complaint-id-here",
    "title": "AC not working in Lab 3",
    "summary": "...",
    "severity": 4,
    "status": "open",
    "user_id": "..."
  }
}
```

---

### `POST /getAllComplaintsForUser`

List all complaints submitted by the authenticated user.

**Request**

```json
{ "secret_code": "deadbeef1234..." }
```

**Response** — array of complaint objects (without `user_name`).

---

### `POST /getAllComplaintsForAdmin`

List **all** complaints across all users (admin only).

**Request**

```json
{ "secret_code": "ADMIN-SECRET-2024" }
```

**Response** — array of complaint objects with `user_name` populated.

---

### `POST /viewComplaint`

View a single complaint. Accessible by its owner or the admin.

**Request**

```json
{
  "secret_code": "deadbeef1234...",
  "complaint_id": "complaint-id-here"
}
```

---

### `POST /resolveComplaint`

Mark a complaint as resolved (admin only).

**Request**

```json
{
  "secret_code": "ADMIN-SECRET-2024",
  "complaint_id": "complaint-id-here"
}
```

**Response** — updated complaint object with `"status": "resolved"`.

---

## Quick Test with curl

```bash
# 1. Register
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Shreya","email":"shreya@test.com"}' | jq

# 2. Copy the secret_code from above, then login
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"secret_code":"<YOUR_SECRET>"}' | jq

# 3. Submit a complaint
curl -s -X POST http://localhost:8080/submitComplaint \
  -H "Content-Type: application/json" \
  -d '{"secret_code":"<YOUR_SECRET>","title":"WiFi Issue","summary":"No internet in hostel","severity":3}' | jq

# 4. Admin: list all complaints
curl -s -X POST http://localhost:8080/getAllComplaintsForAdmin \
  -H "Content-Type: application/json" \
  -d '{"secret_code":"ADMIN-SECRET-2024"}' | jq

# 5. Admin: resolve a complaint
curl -s -X POST http://localhost:8080/resolveComplaint \
  -H "Content-Type: application/json" \
  -d '{"secret_code":"ADMIN-SECRET-2024","complaint_id":"<COMPLAINT_ID>"}' | jq
```

---

## Error Responses

All errors follow this shape:

```json
{
  "success": false,
  "message": "human-readable error"
}
```

| HTTP Status | Meaning                            |
| ----------- | ---------------------------------- |
| 400         | Missing / invalid fields           |
| 401         | Wrong secret code or access denied |
| 404         | User or complaint not found        |
| 405         | Wrong HTTP method (must be POST)   |
| 409         | Email already registered           |
| 500         | Internal server error              |

---

## Concurrency Safety

The store uses `sync.RWMutex`:

- **Multiple readers** can query simultaneously (read lock).
- **Writers** (register, submit, resolve) acquire an exclusive lock.
- No data race is possible even under high concurrent load.

---

## Design Decisions

| Decision                  | Reason                                                                        |
| ------------------------- | ----------------------------------------------------------------------------- |
| Zero third-party packages | Requirement — only Go stdlib                                                  |
| POST for all routes       | All routes need a body; avoids query-param leakage of secret codes            |
| In-memory store           | Keeps scope simple; swap with a DB by implementing the same method signatures |
| Admin secret in `main.go` | Easy to spot and change; would be an env var in production                    |
| `sync.RWMutex`            | Allows concurrent reads without blocking, serialises writes                   |
