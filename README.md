# ChronoCache

ChronoCache is a smart, AI-powered caching system built using Go for performance and Python for intelligent TTL (time-to-live) prediction. It demonstrates modern backend practices including microservice communication, intelligent caching, and frontend observability.

## Features:

- LRU-based in-memory cache with auto-expiry
- Dynamic TTL prediction via FastAPI ML model
- RESTful API for cache interactions
- Frontend dashboard to monitor keys and TTLs
- Auto-refreshing UI with live countdowns

---

## How It Works

This project is split into two key components:

### 1. **Go Backend (Core Cache Layer + Web Server)**

- I implemented the caching layer using an LRU cache from HashiCorp’s `golang-lru` package.
- A `ChronoCache` struct wraps the LRU cache, tracking expiry times using TTL.
- The Go HTTP server exposes `/set` and `/cache` routes.
- When `/set` is called, the Go backend sends the key to a Python-based ML model to receive a TTL prediction, and caches the key-value pair accordingly.

### 2. **Python FastAPI Service (TTL Predictor)**

- The Python microservice is built with FastAPI and uses a simple ML model to predict an appropriate TTL.
- It receives keys via POST requests and returns a TTL in seconds.
- This simulates an intelligent decision-making service that could be upgraded with real-world features like key category detection or ML-based expiration strategy.

### 3. **Frontend (Live Dashboard)**

- A static HTML + JavaScript UI is served by the Go backend.
- It fetches live cache contents from `/cache` every few seconds.
- Users can insert keys through a simple form and view their TTLs in real-time.
- The UI auto-updates and visually distinguishes expired vs. active keys.

---

## Setup Instructions

### Prerequisites

- Go (>=1.20)
- Python (>=3.9)
- Git

---

### Step-by-step (PowerShell)

#### 1. Clone the Repository

```powershell
git clone https://github.com/LiciaSaikia/chronocache.git
cd chronocache
```

---

#### 2. Run the Go Cache Server

```powershell
cd go-cache
go mod tidy
go run main.go
```

Runs at: [http://localhost:8080](http://localhost:8080)

---

#### 3. Run the Python ML Predictor

In a new terminal:

```powershell
cd ml-predictor
python -m venv venv
.\venv\Scripts\activate
pip install fastapi uvicorn scikit-learn
uvicorn main:app --reload
```

Runs at: [http://127.0.0.1:8000](http://127.0.0.1:8000)

---

#### 4. Open the Frontend Dashboard

Visit:

```
http://localhost:8080
```

- Add key-value entries via form
- TTL is dynamically predicted and applied
- Entries are removed once expired

---

## API Endpoints

### POST `/set`

```json
{
  "key": "username",
  "value": "darkmatter"
}
```

Response:

```
Key username set with TTL 44s
```

---

### GET `/cache`

Returns currently cached (non-expired) entries:

```json
[
  {
    "key": "username",
    "value": "darkmatter",
    "ttl": 32
  }
]
```

---

## Use Cases

- Demonstrates AI-powered caching strategies
- Showcases Go + Python microservice orchestration

- Basis for extending to gRPC, Redis, or persistent stores




