# Go URL Shortener Microservice 🚀

A high-performance, containerized URL shortening service built with **Go (Golang)**, **Fiber**, **Redis**, and **PostgreSQL**. This project demonstrates a microservices architecture using the Cache-Aside pattern for optimal performance.

## 🏗️ Architecture & Tech Stack

* **Language:** Go (Golang) `v1.25.5`
* **Framework:** Fiber `v2`
* **Database:** PostgreSQL `15-alpine` (Docker Image)
* **Cache:** Redis `7-alpine` (Server) / `go-redis v9` (Client)
* **ORM:** GORM `v1.31.1`
* **Containerization:** Docker & Docker Compose

## ✨ Features

* **URL Shortening:** Generate short aliases for long URLs.
* **Custom Aliases:** Support for custom short codes.
* **Fast Redirection:** Uses Redis caching to serve redirects instantly (Cache-Aside Pattern).
* **Click Tracking:** Tracks the number of clicks for each link (Async updates).
* **Secure:** Automatic HTTP to HTTPS enforcement.
* **Dockerized:** Ready to deploy with a single command.
* **Automated Tests:** Integration tests included.

## 🚀 Getting Started

### Prerequisites

* [Docker](https://www.docker.com/) and Docker Compose
* [Go](https://golang.org/) (Optional, if running locally without Docker)

### Installation & Running

1.  **Clone the repository**
    ```bash
    git clone https://github.com/aliemreipek/go-url-shortener.git
    cd go-url-shortener
    ```

2.  **Environment Setup**
    Create a `.env` file in the root directory. You can copy the configuration below:
    ```env
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=password
    DB_NAME=shortener
    DB_PORT=5432
    REDIS_ADDR=cache:6379
    REDIS_PASSWORD=
    PORT=3000
    DOMAIN=localhost:3000
    ```

3.  **Run with Docker (Recommended)**
    ```bash
    docker-compose up -d --build
    ```
    *This will start the API, PostgreSQL, and Redis containers.*

4.  **Run Tests**
    To run the integration tests:
    ```bash
    go test ./cmd/... -v
    ```

## 🔌 API Endpoints

### 1. Shorten URL
**POST** `/api/v1`

**Request Body:**
```json
{
  "url": "https://github.com/aliemreipek",
  "short": "mygithub"
}
```
*(Note: "short" field is optional. If empty, a random ID is generated.)*

**Response:**
```json
{
  "url": "https://github.com/aliemreipek",
  "short": "localhost:3000/mygithub",
  "expiry": 24,
  "rate_limit": 10
}
```

### 2. Redirect
**GET** `/:url`

* Example: `http://localhost:3000/mygithub`
* Redirects to the original URL.

## 📂 Project Structure

```text
├── api         # API handlers and business logic
├── cmd         # Main entry point and tests
├── database    # Database and Redis connection setup
├── model       # Database structs and schemas
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 🗺️ Roadmap & Future Improvements

* [ ] **User Authentication:** Implement JWT authentication to allow users to manage their links.
* [ ] **Rate Limiting:** Add middleware to prevent abuse (e.g., 10 requests per minute).
* [ ] **Analytics Dashboard:** Create a simple frontend to visualize click statistics.
* [ ] **CI/CD Pipeline:** Add GitHub Actions for automated testing and deployment.

## 📄 License

This project is licensed under the MIT License.
