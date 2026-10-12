# GoCoffee ☕

> A minimalist, terminal-first specialty coffee shop backend built in idiomatic Go, following Clean Architecture principles.

## 🚀 Project Overview

**GoCoffee** is an educational yet production-grade backend application simulating a specialty coffee shop. Designed as a growing pet project, it serves as a practical sandbox for mastering modern Go development, software design patterns, concurrent programming, and infrastructure technologies.

Rather than relying on over-engineered frameworks or monolithic scripts, GoCoffee is structured from the ground up to reflect professional software engineering standards—featuring clear boundaries between business logic, storage layers, and user interfaces.

---

## 🏗️ Architectural Principles

- **Clean Architecture:** Strict separation of concerns across layers (`domain`, `storage`, `service`, and presentation).
- **Idiomatic Go:** Minimalist structures, pointer receivers, slice-based data manipulation (FIFO carts), and anonymous field embedding (field promotion) for clean, flat data models.
- **Dependency Inversion:** Use of interfaces (such as `Payer` and repositories) to decouple business logic from external implementations.
- **Pluggable Interfaces:** Extensible payment strategy pattern allowing seamless integration of internal wallets, cryptocurrencies, and bank cards.

---

## 📂 Project Structure

```text
gocoffee/
├── cmd/
│   └── gocoffee/
│       └── main.go          # Application entry point & wiring
├── internal/
│   ├── domain/              # Core business entities, models & interfaces
│   │   ├── user.go          # User, Wallet, Cart (FIFO), Discount
│   │   ├── item.go          # Coffee products, categories, ItemRepository
│   │   └── payment.go       # Payer interface & wallet payment logic
│   ├── storage/             # Data persistence layers
│   │   ├── json_user.go     # JSON file-based user storage
│   │   └── json_item.go     # JSON file-based item catalog storage
│   └── service/             # Business coordination & orchestration
│       └── shop.go          # ShopService (AddToCart, RemoveFromCart, Checkout)
├── data/                    # Local persistent JSON state
│   ├── users.json
│   └── items.json
├── go.mod
└── README.md
```

---

## 🛠️ Core Domain Models

1. **User & Wallet (Anonymous Embedding):**
   `User` embeds `Wallet` anonymously, enabling flat field access (`user.Balance`, `user.BonusBalance`) while maintaining clean data structures.
2. **FIFO Cart:**
   `Cart` is implemented as a slice of `CartItem` (`[]CartItem`), guaranteeing strict first-in, first-out ordering with built-in deduplication and removal methods.
3. **Payment Strategy (`Payer`):**
   The `Checkout` process accepts any type satisfying the `Payer` interface, abstracting away the underlying transaction mechanism.
