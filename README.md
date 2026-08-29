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

This document provides an overview of the backend microservices architecture for the **Writeful** project, detailing the services, flow diagrams, environment configuration, and operation instructions.

---

## 📌 System Overview

The **Writeful** backend is built on an API-driven **Microservices** architecture designed for independence, security, and scalability. 

The system consists of the following key components:
1. **API Gateway (KrakenD)**: Centralized entry point that routes HTTP requests from clients to their corresponding downstream microservices.
2. **Cloudflare Tunnels**: Secure tunnels that expose the API Gateway (HTTP) and the Chat Service (WebSocket) to the internet without opening public ports on the host.
3. **Auth Service (Go - REST)**: Manages users, sign-ups, sign-ins, and issues/validates JWT tokens.
4. **Content Service (Go - REST)**: Manages articles (posts), draft stories, search tags, and attached music.
5. **Media Service (Go - REST)**: Uploads and manages media files (images/videos) using the Cloudinary cloud storage service.
6. **Chat Service (Go - WebSockets/REST)**: Enables real-time private messaging between users over WebSocket connections.

---

## 🏗️ System Architecture

![System Architecture Diagram](docs/images/system_architecture.png)

Below is the general system data flow for the Writeful backend:

```mermaid
graph TD
    Client["Client / FE App"] -->|HTTP/HTTPS Requests| CFTunnel["Cloudflare Tunnel - HTTP"]
    Client -->|WebSocket Connections| CFTunnelWS["Cloudflare Tunnel - WS"]
    
    subgraph API Gateway Layer
        CFTunnel -->|Port 8080| Gateway["KrakenD API Gateway"]
    end
    
    subgraph Backend Microservices
        Gateway -->|"/v1/auth/*"| AuthService["Auth Service<br/>Port 8004"]
        Gateway -->|"/v1/posts/*, /v1/stories/*"| ContentService["Content Service<br/>Port 8003"]
        Gateway -->|"/v1/media/*"| MediaService["Media Service<br/>Port 8005"]
        
        CFTunnelWS -->|Port 8006| ChatService["Chat Service<br/>Port 8006"]
        
        ContentService -->|"HTTP JWT Validate<br/>Port 8004"| AuthService
        ChatService -->|"HTTP JWT Validate<br/>Port 8004"| AuthService
    end

    subgraph External Resources
        MediaService -->|Upload Images| Cloudinary["Cloudinary Cloud Storage"]
    end

    subgraph Database Layer
        AuthService -->|"Schema: auth_service"| DB[("PostgreSQL<br/>Port 5438")]
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

## 🔄 Core Flows (Flow Diagrams)

### 1. User Authentication & Token Verification (HTTP)
When a client requests a protected endpoint (e.g. creating a post in the `Content Service`), the target service makes a synchronous HTTP call to the `Auth Service` to verify the user token:

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Gateway as KrakenD API Gateway
    participant Service as Content Service
    participant Auth as Auth Service (HTTP)
    participant DB as PostgreSQL (auth_service)

    Client->>Gateway: HTTP Request (Header: Authorization: Bearer JWT)
    Gateway->>Service: Route request with Token
    Note over Service: User details required
    Service->>Auth: HTTP Call: VerifyToken(Token)
    Auth->>DB: Query user / check session details (if needed)
    DB-->>Auth: User Record
    Note over Auth: Verify signature & expiration
    Auth-->>Service: HTTP Response: VerifyTokenResponse (User ID, Role, Valid: true)
    Note over Service: Process request with user permissions
    Service-->>Gateway: HTTP Response
    Gateway-->>Client: Return response to Client
```

### 2. Real-Time WebSocket Chat Flow
Since API Gateway (KrakenD) does not natively support robust, persistent WebSocket proxying, the real-time chat traffic is isolated through a dedicated Cloudflare Tunnel pointing directly to the `Chat Service`:

```mermaid
sequenceDiagram
    autonumber
    actor Client A
    actor Client B
    participant WS as Chat Service (WS Hub)
    participant Auth as Auth Service (HTTP)
    participant DB as PostgreSQL (chat_service)

    Client A->>WS: Establish WS Connection (Query Param: token=JWT)
    WS->>Auth: HTTP VerifyToken(JWT)
    Auth-->>WS: Token confirmed valid (User ID: A)
    WS-->>Client A: Connection upgraded & accepted (Joined Hub)
    
    Note over Client A, Client B: Client B is already connected & registered to the WS Hub
    
    Client A->>WS: Send WS Message {to: User B, content: "Hello"}
    WS->>DB: Save Message to database
    DB-->>WS: Message saved successfully
    WS->>WS: Search active connection for Client B in Hub
    WS->>Client B: Dispatch WS Event: New Message {from: User A, content: "Hello"}
```

---

## 🛠️ Microservice Directory Details

### 1. Gateway Service
* **Technology**: KrakenD API Gateway.
* **Responsibilities**: Central entry routing, rate limiting, request/response transformation, and endpoint aggregation.
* **Configurations**: Located in [gateway-service](file:///Users/haopham/go-playground/writeful/gateway-service).

### 2. Auth Service
* **Technology**: Go (Gin Web Framework, SQLBoiler/gorm).
* **Responsibilities**: User signup/signin, authorization, session management, and exposing token verification APIs.
* **Environment Configuration**: Setup details in [auth-service/.env.example](file:///Users/haopham/go-playground/writeful/auth-service/.env.example).

### 3. Content Service
* **Technology**: Go (Gin Web Framework, HTTP Client).
* **Responsibilities**: Core posting engine, draft/publish workflows for stories, tags indexing, and background music associations.
* **Environment Configuration**: Setup details in [content-service/.env.example](file:///Users/haopham/go-playground/writeful/content-service/.env.example).

### 4. Media Service
* **Technology**: Go (Gin Web Framework).
* **Responsibilities**: Uploads file streams to Cloudinary, stores references in database, and returns static asset URLs.
* **Environment Configuration**: Setup details in [media-service/.env.example](file:///Users/haopham/go-playground/writeful/media-service/.env.example).

### 5. Chat Service
* **Technology**: Go (Gorilla Websocket, HTTP Client).
* **Responsibilities**: Handles raw WebSocket connections, maintains connection hub state, publishes messages, and tracks unread message statistics.
* **Environment Configuration**: Setup details in [chat-service/.env.example](file:///Users/haopham/go-playground/writeful/chat-service/.env.example).

---

## ⚙️ Environment Configuration

For each backend microservice, copy the `.env.example` file to `.env` and fill in the target variables:

```bash
# Execute within the corresponding microservice directories
cp auth-service/.env.example auth-service/.env
cp content-service/.env.example content-service/.env
cp media-service/.env.example media-service/.env
cp chat-service/.env.example chat-service/.env
```

### Key Configurations:
* **`PORT`**: The HTTP port for the service (Defaults: Auth - `8004`, Content - `8003`, Media - `8005`, Chat - `8006`).
* **`DB.URL`**: PostgreSQL connection string. Format: `postgres://<username>:<password>@<host>:<port>/<db_name>?sslmode=disable&search_path=<schema>`.
* **`JWT.SECRET`**: Signature secret key for JWT tokens. Replace with a secure key in production.
* **`CLOUDINARY`**: Credentials (Cloud Name, API Key, API Secret) for the Cloudinary integration in `media-service`.
* **`AUTH_SERVICE_GRPC_ADDR`**: The host:port for gRPC connections (Optional / Used for testing only).

---

## 🚀 Operation & Running Guide

A root-level `Makefile` is provided to simplify orchestration, container builds, and database utilities.

### Key Make Targets:

* **Start the ecosystem** (in background mode):
  ```bash
  make up
  ```
* **Stop all running containers**:
  ```bash
  make down
  ```
* **Stream container logs**:
  ```bash
  make logs             # Logs for all services
  make logs-auth        # Auth Service logs only
  make logs-content     # Content Service logs only
  make logs-chat        # Chat Service logs only
  ```
* **Rebuild all local Docker images**:
  ```bash
  make build-all
  ```
