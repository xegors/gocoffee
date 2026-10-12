# GoCoffee Evolution Roadmap & Career Alignment 🚀

This roadmap outlines the phased evolution of the **GoCoffee** project. Each stage is intentionally mapped to critical industry requirements for a **Junior/Middle Go Developer**, transforming a simple script into a robust, cloud-native microservice architecture.

---

## Stage 1: Foundation & Clean Architecture (Current Phase 📍)
* **Project Scope:** Core domain logic, file-based persistence (`encoding/json`), repository interfaces, FIFO cart management, and payment strategy patterns.
* **Skills Mastered:** Go syntax, pointers, slices, interfaces, error handling, file I/O, and software design principles (SOLID, Separation of Concerns).
* **Resume Impact:** Demonstrates strong command of idiomatic Go and clean architectural separation.

## Stage 2: Terminal-First SaaS & SSH Authentication
* **Project Scope:** Implementing an SSH server (`charmbracelet/wish`) where users connect directly via terminal (`ssh coffee@gocoffee.local`), authenticating via SSH public keys and managing user sessions automatically.
* **Skills Mastered:** Network protocols, SSH handling, key-based identity management, and terminal application workflows.
* **Resume Impact:** **Network technologies, Linux console mastery**, and standout portfolio uniqueness.

## Stage 3: Interactive TUI (Terminal User Interface)
* **Project Scope:** Building a modern, interactive terminal UI using the Charm ecosystem (`Bubble Tea`, `Lip Gloss`), allowing users to navigate catalogs, manage carts, and checkout seamlessly within their SSH or local session.
* **Skills Mastered:** State management (The Elm Architecture), event loops, and command-line aesthetics.
* **Resume Impact:** Advanced command-line proficiency and UX design in terminal environments.

## Stage 4: HTTP REST API
* **Project Scope:** Exposing the business logic over HTTP using Go's standard `net/http` or lightweight routers (`chi`), handling JSON payloads, status codes, and middleware.
* **Skills Mastered:** HTTP protocol mechanics, RESTful design principles, routing, and request validation.
* **Resume Impact:** **Working with HTTP REST, gRPC**.

## Stage 5: Relational Database Migration (SQL & PostgreSQL)
* **Project Scope:** Transitioning from JSON text files to PostgreSQL, designing normalized schemas (`users`, `items`, `orders`), writing database migrations, and ensuring ACID compliance via transactions.
* **Skills Mastered:** SQL queries, `pgx`/`sqlx` drivers, transaction management (ensuring atomicity during checkout), and database migrations.
* **Resume Impact:** **SQL/PostgreSQL**, data consistency, and transactional safety.

## Stage 6: Testing & Quality Assurance
* **Project Scope:** Writing comprehensive unit and integration tests for `ShopService` using Go's built-in `testing` package and mock generation tools. Practicing structured Git workflows (Git Flow / GitLab Flow).
* **Skills Mastered:** Table-driven testing, mocking dependencies via interfaces, code reviews, and branching strategies.
* **Resume Impact:** **Writing tests for your own code**, **Git Flow / GitLab Flow**.

## Stage 7: Containerization & Infrastructure (Docker & Linux)
* **Project Scope:** Packaging the application with multi-stage Docker builds, orchestrating multi-container environments using `docker-compose` alongside PostgreSQL.
* **Skills Mastered:** Dockerfiles, service isolation, environment configuration, and Linux environment management.
* **Resume Impact:** **Confident use of the Linux console**, containerization, and deployment readiness.

## Stage 8: Enterprise-Grade Architecture (Redis, Kafka, gRPC)
* **Project Scope:** 
  * Integrating **Redis** for active cart caching and session management.
  * Introducing **gRPC** for high-performance internal service communication.
  * Implementing **Kafka** for event-driven architecture (emitting `OrderCreated` events for asynchronous worker notification processing).
* **Skills Mastered:** Caching strategies, event-driven microservices, message brokers, and high-load patterns.
* **Resume Impact:** **Understanding microservice application architecture**, Redis, and Kafka experience.
