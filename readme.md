# Taskflow Backend

## 1. Overview

Taskflow is a backend service for managing users, projects, and tasks. It provides authentication, task management, and basic observability features .

### Tech Stack
- Go (Golang)
- PostgreSQL
- Docker & Docker Compose

---

## 2. Architecture Decisions

- **Docker-first approach**: Entire stack (DB + API) runs via Docker Compose for reproducibility.
- **Multi-stage Docker build**: Keeps runtime image small and production-ready.
- **Koanf-based config**: Flexible environment-based configuration.
- **Auto migrations**: Database schema and seed data are applied automatically at container startup.
- **Separation of concerns**: Database, cache, and API are separate services.

---

## 3. Running Locally

Assuming only Docker is installed:

```bash
git clone https://github.com/CrossStack-Q/taskflow-assignment.git
cd backend
docker compose up --build
```

Server will start at:

http://localhost:8080

---

## 4. Running Migrations

Migrations run automatically when the API container starts.

No manual steps required.

---

## 5. Test Credentials

Use the following credentials:

Email: test@example.com  
Password: password123

---

## 6. API Reference

### Auth Login

POST /api/v1/auth/login

Request:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "user": {
    "id": "uuid",
    "name": "Anurag Sharma",
    "email": "test@example.com",
    "created_at": "2026-04-12T17:50:00.862954Z"
  },
  "token": "<token>"
}
```

---

## 7. All protected API's in backend readme file

You can find all protected APIs here:  
[Backend API Documentation](https://github.com/CrossStack-Q/taskflow-assignment/blob/main/backend/README.md)