# Training Plan: Go Backend — Finance Tracker

## Learner Profile
- **Topic**: Go backend development with real-time features
- **Skill level**: Beginner (Go) — experienced with React/Next.js and Java Spring Boot
- **Language/Stack**: Go + Next.js + PostgreSQL + sqlc + JWT + WebSockets
- **Project style**: Capstone

---

## Project: PocketLedger

> A personal finance tracker where users can log expenses by category, set monthly budget limits, and receive real-time alerts and dashboard updates when spending approaches or exceeds those limits.

**Repo structure**: Monorepo (`/frontend` → Vercel, `/backend` → Railway)

## Requirements

- **Relational Database** (PostgreSQL) — stores users, categories, expenses, budgets
- **Database Migrations** (`golang-migrate`) — manages schema evolution across environments
- **Authentication** (JWT + bcrypt) — secures all endpoints, scopes data per user
- **REST API** (Go `net/http` / `chi`) — exposes CRUD and summary endpoints
- **Real-time** (WebSockets via `gorilla/websocket`) — live dashboard updates and budget alerts
- **Containerisation** (Docker + docker-compose) — local dev environment and Railway deployment
- **Testing** (unit + integration) — verifies business logic and API contracts

---

## Tasks

### Task 1: Go Project Skeleton
- [ ] **Deliverable**: A Go project inside `/backend` with a working HTTP server that responds to `GET /health`, structured with separate packages for `handler`, `service`, `repository`, and `config`

**Objective**: Establish a clean, layered project structure in Go that mirrors the separation of concerns you're familiar with from Spring Boot — so every subsequent task has a consistent home for its code.

---

### Task 2: Docker Dev Environment
- [ ] **Deliverable**: A `docker-compose.yml` at the repo root that spins up PostgreSQL and the Go backend together; running `docker compose up` starts both services and `GET /health` returns `200 OK`

**Objective**: Eliminate manual database setup so the entire dev environment is reproducible with a single command — a habit that pays off when onboarding teammates or deploying to Railway.

---

### Task 3: Database Schema + Migrations
- [ ] **Deliverable**: Migration files (using `golang-migrate`) that create tables for `users`, `categories`, `expenses`, and `budgets`; migrations run automatically on app startup and are idempotent

**Objective**: Design a normalized schema that captures the relationships between users, their categories, their expenses, and their budget limits — and learn how Go projects manage schema evolution without an ORM.

---

### Task 4: Type-Safe DB Queries with sqlc
- [ ] **Deliverable**: `sqlc`-generated Go code for CRUD operations on all four tables; a `sqlc.yaml` config in `/backend`; no raw string queries anywhere in the codebase

**Objective**: Learn how `sqlc` turns SQL into type-safe Go functions, replacing the boilerplate you'd write manually — and understand why Go projects often prefer this over ORMs like GORM.

---

### Task 5: User Registration and Login
- [ ] **Deliverable**: `POST /auth/register` and `POST /auth/login` endpoints; passwords hashed with bcrypt; login returns a signed JWT; a middleware function that validates the JWT and attaches the user ID to the request context

**Objective**: Implement stateless authentication from scratch in Go — understanding how JWT validation middleware works at the HTTP layer so every subsequent endpoint can trust the identity of the caller.

---

### Task 6: Categories API
- [ ] **Deliverable**: Authenticated CRUD endpoints (`GET`, `POST`, `PUT`, `DELETE`) for `/categories`; each category belongs to the authenticated user; a seed of default categories (Food, Transport, Entertainment, Health, Other) is created on first login

**Objective**: Build your first fully authenticated resource in Go, enforcing that users can only read and modify their own data — establishing the ownership pattern all subsequent resources will follow.

---

### Task 7: Expenses API
- [ ] **Deliverable**: Authenticated CRUD endpoints for `/expenses`; each expense has an amount, date, note, and category; `GET /expenses` supports filtering by category and date range via query params

**Objective**: Implement the core data entry feature of PocketLedger — and practice parsing, validating, and responding to query parameters in Go's HTTP layer.

---

### Task 8: Budgets API + Spending Summary
- [ ] **Deliverable**: Authenticated CRUD endpoints for `/budgets` (a budget links a category to a monthly spending limit); a `GET /budgets/summary` endpoint that returns each budget with its current month's total spend, remaining amount, and a percentage used

**Objective**: Write a non-trivial SQL aggregation query via `sqlc` that joins expenses to budgets — producing the summary data that will drive both the dashboard and the real-time alert logic.

---

### Task 9: WebSocket Server
- [ ] **Deliverable**: A `GET /ws` endpoint that upgrades to a WebSocket connection (authenticated via JWT passed as a query param); the server maintains a registry of active connections per user; a helper function exists to broadcast a JSON message to all connections for a given user ID

**Objective**: Learn how Go handles long-lived concurrent connections using goroutines and channels — the core concurrency model that makes Go well-suited for real-time features.

---

### Task 10: Real-Time Dashboard Updates
- [ ] **Deliverable**: When an expense is created or deleted via the REST API, the server broadcasts an updated spending summary payload to all active WebSocket connections for that user — no page refresh needed to see the change

**Objective**: Wire the REST layer to the WebSocket broadcast system so that state changes propagate in real time — connecting the two systems you built independently in Tasks 7–9.

---

### Task 11: Budget Alert Notifications
- [ ] **Deliverable**: After every expense creation, the server checks if any budget has crossed the 80% or 100% threshold; if so, it broadcasts a distinct alert message (with category name, percentage used, and severity) to the user's WebSocket connections

**Objective**: Implement the core value proposition of PocketLedger — proactive alerts — by adding business logic that runs as a side effect of expense writes, without blocking the HTTP response.

---

### Task 12: Next.js Frontend
- [ ] **Deliverable**: A Next.js app in `/frontend` with: an auth flow (register/login), a dashboard showing spending summaries per category, an expense entry form, a budgets management page, and a toast notification system that reacts to WebSocket alert messages in real time

**Objective**: Connect your existing Next.js skills to the Go backend you built — focusing on the WebSocket client integration and how to keep the UI in sync with server-pushed events without polling.

---

### Task 13: Unit + Integration Tests
- [ ] **Deliverable**: Unit tests for the budget alert threshold logic and spending summary calculation; integration tests for the auth, expenses, and budgets endpoints using a real test PostgreSQL database (not mocks); all tests pass in CI

**Objective**: Learn Go's built-in `testing` package and establish the habit of testing against a real database — avoiding the class of bugs where mocked tests pass but production queries fail.

---

### Task 14: Deployment
- [ ] **Deliverable**: Go backend and PostgreSQL deployed to Railway with environment variables for secrets; Next.js frontend deployed to Vercel pointing at the Railway API URL; a live public URL in the GitHub README with a screenshot of the dashboard

**Objective**: Ship PocketLedger as a real, publicly accessible application — because a portfolio project with a live URL is worth ten times one that only runs locally.

---

## Suggested Learning Resources

- [Go Tour](https://go.dev/tour) — complete this before Task 1
- [Go by Example](https://gobyexample.com) — reference as you build each task
- [`sqlc` docs](https://docs.sqlc.dev) — read the quickstart before Task 4
- [`golang-migrate` docs](https://github.com/golang-migrate/migrate) — read before Task 3
- [`gorilla/websocket` examples](https://github.com/gorilla/websocket/tree/main/examples) — read before Task 9
