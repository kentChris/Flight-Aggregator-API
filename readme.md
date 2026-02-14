Flight Aggregator API ✈️

A Go-based flight search engine that aggregates data from Garuda, AirAsia, Lion Air, and Batik Air. It features Redis caching, advanced filtering (price, duration, stops, time), and a "Best Value" ranking algorithm.

📋 Prerequisites
Before starting, ensure you have the following installed:

- Go 1.24+ (Note: 1.26 is not yet a stable release; 1.24 is recommended).
- Docker & Docker Desktop (for Containerized run).
- Redis (if running locally without Docker).

🛠️ Option 1: Run with Docker Compose (Fastest)
This method sets up the API, the mock data paths, and the Redis database automatically with a single command.

- Ensure main.go is set for Docker: In cmd/app/main.go, ensure the Redis address is set to the container name:
"redisService := redis.NewRedisService("redis:6379", "", 0)"

- Build and Run:
"docker compose up --build"


💻 Option 2: Run Locally (Development Mode)
Use this if you want to make quick code changes and run via go run.

Step 1: Start Redis
docker run -d --name flight-redis -p 6379:6379 redis:alpine


Step 2: Configure Redis Address
In cmd/app/main.go, switch the address to localhost:
// Comment out the docker one, uncomment localhost
redisService := redis.NewRedisService("localhost:6379", "", 0)


Step 3: Install Dependencies
Download the required Go modules:
go mod tidy
go mod download


Step 4: Run the Application
go run cmd/app/main.go


⚙️ How to Modify the Search (Mocking)
Since the current version uses a mock trigger in the controller, you can change the search criteria (filters/sorting) by editing: internal/controller/flight.go
