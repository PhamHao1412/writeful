<div align="center">

# ✍️ Writeful Backend Microservices (`writeful`)

**The high-performance, event-driven Microservices platform for blogging and real-time messaging.**

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![API Gateway](https://img.shields.io/badge/KrakenD-API%20Gateway-1D70B8?style=flat-square&logo=krakend)](https://www.krakend.io/)
[![WebSockets](https://img.shields.io/badge/WebSockets-Real--Time%20Chat-010101?style=flat-square&logo=socketdotio)](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)
[![Database](https://img.shields.io/badge/PostgreSQL-Multi--Schema-4169E1?style=flat-square&logo=postgresql)](https://www.postgresql.org)
[![Media Storage](https://img.shields.io/badge/Cloudinary-SDK-3448C5?style=flat-square&logo=cloudinary)](https://cloudinary.com)
[![Cloudflare](https://img.shields.io/badge/Cloudflare-Tunnels-F38020?style=flat-square&logo=cloudflare)](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
[![Deployment](https://img.shields.io/badge/Deploy-Docker%20Compose-2496ED?style=flat-square&logo=docker)](https://www.docker.com)

</div>

---

**Writeful** is a modular, decoupled microservices platform designed for modern blogging, content authoring, and real-time private messaging. 

The system leverages **KrakenD** as a high-throughput API Gateway, isolated **Cloudflare Tunnels** for secure zero-trust network ingress (HTTP and WebSockets), independent **Go backend services**, a multi-schema **PostgreSQL** database, and **Cloudinary** for media transformations.

---

## System Architecture

```mermaid
graph TD
    Client["Client / Frontend App"] -->|HTTP / HTTPS REST| CFTunnel["Cloudflare Tunnel (HTTP)"]
    Client -->|WebSocket Connections| CFTunnelWS["Cloudflare Tunnel (WS)"]
    
    subgraph API Gateway Layer
        CFTunnel -->|Port 8080| Gateway["KrakenD API Gateway"]
    end
    
    subgraph Microservices Layer
        Gateway -->|"/v1/auth/*"| AuthService["Auth Service<br/>Port 8004"]
        Gateway -->|"/v1/posts/*, /v1/stories/*"| ContentService["Content Service<br/>Port 8003"]
        Gateway -->|"/v1/media/*"| MediaService["Media Service<br/>Port 8005"]
        
        CFTunnelWS -->|Port 8006| ChatService["Chat Service (WS Hub)<br/>Port 8006"]
        
        ContentService -->|"Token Verification (HTTP)"| AuthService
        ChatService -->|"Token Verification (HTTP)"| AuthService
    end

    subgraph Storage & External Services
        MediaService -->|Direct Upload| Cloudinary["Cloudinary CDN Storage"]
    end

    subgraph Multi-Schema Database Layer
        AuthService -->|"Schema: auth_service"| DB[("PostgreSQL")]
        ContentService -->|"Schema: content_service"| DB
        MediaService -->|"Schema: media_service"| DB
        ChatService -->|"Schema: chat_service"| DB
    end
    
    style Gateway fill:#f9f,stroke:#333,stroke-width:2px
    style AuthService fill:#bbf,stroke:#333,stroke-width:1px
    style ContentService fill:#bbf,stroke:#333,stroke-width:1px
    style MediaService fill:#bbf,stroke:#333,stroke-width:1px
    style ChatService fill:#bbf,stroke:#333,stroke-width:1px
    style DB fill:#bfb,stroke:#333,stroke-width:2px
    style Cloudinary fill:#ffb,stroke:#333,stroke-width:1px
```

---

## Microservices Breakdown

| Service | Port | Gateway Route Prefix | Database Schema | Primary Responsibilities |
|---|---|---|---|---|
| **`gateway-service`** | `8080` | `/v1/*` | - | Central routing, CORS handling, rate limiting, request proxying |
| **`auth-service`** | `8004` | `/v1/auth/*` | `auth_service` | User registration, login, JWT token issuance, token verification |
| **`content-service`** | `8003` | `/v1/posts/*`, `/v1/stories/*` | `content_service` | Blog posts, draft stories, tagging, category feeds, background music |
| **`media-service`** | `8005` | `/v1/media/*` | `media_service` | Media file uploads, Cloudinary CDN integrations, image metadata |
| **`chat-service`** | `8006` | WebSocket Direct | `chat_service` | Real-time WebSocket message hub, private messaging, chat history |

---

## Core Interaction Flows

### 1. Synchronous Token Verification Flow
When a client requests a protected endpoint in a downstream microservice (e.g. creating a post in `content-service`), the downstream service verifies the JWT directly against `auth-service`:

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Gateway as KrakenD Gateway (:8080)
    participant Content as Content Service (:8003)
    participant Auth as Auth Service (:8004)
    participant DB as PostgreSQL

    Client->>Gateway: HTTP Request (Authorization: Bearer <JWT>)
    Gateway->>Content: Proxy Request with Token
    Content->>Auth: HTTP POST /internal/v1/auth/verify-token
    Auth->>DB: Validate user status & token signature
    DB-->>Auth: User Record
    Auth-->>Content: 200 OK {valid: true, user_id: "...", role: "..."}
    Content->>DB: Execute post creation in content_service schema
    Content-->>Gateway: 201 Created
    Gateway-->>Client: JSON Response
```

### 2. Real-Time WebSocket Chat Flow
WebSocket connections bypass KrakenD and route through a dedicated Cloudflare Tunnel directly into `chat-service`'s in-memory connection hub:

```mermaid
sequenceDiagram
    autonumber
    actor Client A
    actor Client B
    participant WS as Chat Service Hub (:8006)
    participant Auth as Auth Service (:8004)
    participant DB as PostgreSQL (chat_service)

    Client A->>WS: Establish WebSocket (query: ?token=<JWT>)
    WS->>Auth: Verify JWT Token
    Auth-->>WS: Token valid (User ID: A)
    WS-->>Client A: Connection upgraded (Registered in Hub)
    Note over Client A, Client B: Client B is already connected and active in Hub

    Client A->>WS: Send message payload {to: "User B", content: "Hello"}
    WS->>DB: Persist message record
    WS->>WS: Lookup Client B socket in active connection pool
    WS->>Client B: Push WebSocket Event {from: "User A", content: "Hello"}
```

---

## Repository Structure

```
.
├── auth-service/                # Authentication & User Management (Go)
│   ├── cmd/serverd/
│   ├── internal/
│   └── .env.example
├── content-service/             # Articles, Stories, Tags & Music (Go)
│   ├── cmd/serverd/
│   ├── internal/
│   └── .env.example
├── media-service/               # Cloudinary Uploads & Media Assets (Go)
│   ├── cmd/serverd/
│   ├── internal/
│   └── .env.example
├── chat-service/                # Real-Time WebSocket Chat Engine (Go)
│   ├── cmd/serverd/
│   ├── internal/
│   └── .env.example
├── gateway-service/             # KrakenD API Gateway configuration & Dockerfile
│   └── krakend.json
├── docker-compose.yml           # Multi-container orchestration specification
└── Makefile                     # Unified management and operations interface
```

---

## Getting Started

### Prerequisites

- **Docker & Docker Compose** (v2+)
- **Go**: `1.22+` (for local standalone development)
- **PostgreSQL**: PostgreSQL 14+ running on port `5438` (or configured host)

### 1. Configure Environment Files

Create `.env` files for each microservice from their corresponding `.env.example`:

```bash
cp auth-service/.env.example auth-service/.env
cp content-service/.env.example content-service/.env
cp media-service/.env.example media-service/.env
cp chat-service/.env.example chat-service/.env
```

### 2. Start All Microservices

Run the complete microservices ecosystem in background mode:

```bash
make up
```

Verify services status:

```bash
make ps
```

Once running, the API Gateway is available at:
- **API Gateway**: `http://localhost:8080`
- **Health Endpoint**: `http://localhost:8080/health`

---

## Makefile Reference

| Command | Description |
|---|---|
| `make up` | Start all microservices in detached mode |
| `make down` | Stop all microservices containers |
| `make restart` | Restart all running microservices |
| `make restart-gateway` | Recreate and restart `gateway-service` |
| `make restart-auth` | Recreate and restart `auth-service` |
| `make restart-content` | Recreate and restart `content-service` |
| `make restart-media` | Recreate and restart `media-service` |
| `make restart-chat` | Recreate and restart `chat-service` |
| `make logs` | Stream live logs for all microservices |
| `make logs-gateway` | Stream logs for `gateway-service` |
| `make logs-auth` | Stream logs for `auth-service` |
| `make logs-content` | Stream logs for `content-service` |
| `make logs-media` | Stream logs for `media-service` |
| `make logs-chat` | Stream logs for `chat-service` |
| `make health` | Check health status of the KrakenD Gateway |
| `make rebuild` | Pull latest images and force-recreate all containers |
| `make clean` | Stop and remove all containers and networks |

---

## License

This project is licensed under the [MIT License](LICENSE).
