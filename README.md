# 🏪 Vyapaar – Shop Inventory Management System

**Vyapaar** is a full-stack **Shop Inventory Management** application built to help small businesses manage their products, users, and sales efficiently.

The project is currently under active development, with the backend implemented in **Go (Golang)** and JWT-based authentication already in place. The frontend development is underway.

---

## 📦 Tech Stack

### 🚀 Backend
- **Language**: Go (Golang)
- **Authentication**: JWT (JSON Web Tokens)
- **Routing**: `net/http`
- **Database**: sqlite
- **Key Features**:
  - User Signup/Login
  - JWT-based Auth Middleware
  - Product CRUD APIs

### 💻 Frontend *(in development)*
- **Framework**: [To be added, e.g., React.js, Vue.js]
- **Design**: Responsive UI for inventory, orders, and user dashboard

---


---

## 🔐 Authentication Flow

- On login/signup, the backend issues a **JWT token**.
- The token must be included in the `token` cookies (`Bearer <token>`) for protected routes like product operations.
- Middleware verifies the token and extracts user info from it.

---

## 🧪 API Overview

### 🔑 Auth Routes
- `POST /signup` – Register a new user
- `POST /login` – Authenticate user and receive JWT

### 📦 Product Routes (Protected)
- `GET /products` – List all products
- `POST /products` – Add a new product
- `PUT /products/:id` – Update a product
- `DELETE /products/:id` – Delete a product

---

## 🛠️ Setup Instructions

```bash
# Clone the repo
git clone https://github.com/yourusername/vyapaar.git
cd vyapaar

# Initialize Go modules
go mod tidy

# Run the server
go run cmd/main.go

