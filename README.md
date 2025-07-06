# Vyapaar – Inventory Management System 🛒

Vyapaar is a backend service for managing shop inventories with user authentication, product management, email verification, and robust API design. This service is built using **Go**, with support for **JWT-based authentication**, **Redis**, **RabbitMQ**, and more.

---

## 🚀 Features

- 🔐 **User Authentication** (Signup, Login) with JWT
- 📧 **Email Verification** via SMTP using RabbitMQ
- 🏪 **Shop-based Product Management** using `shopID`
- 🧾 **CRUD APIs** for Products:
  - `GET /products`
  - `GET /products/:id`
  - `POST /products`
  - `PATCH /products/:id`
  - `DELETE /products/:id`
- 🔍 **Pagination, Sorting & Filtering** via query parameters
- 📦 **Redis** used for OTP storage and caching
- 🐇 **RabbitMQ** queues OTP emails
- 🐳 Docker support with images for Redis and RabbitMQ
- ⚙️ Configurable using a YAML file
- 🖥️ Frontend: Coming Soon (TBA)

---

## 🔐 Authentication Flow

1. User **signs up** with email and password.
2. Backend generates an **OTP**, stores in Redis, and pushes it to RabbitMQ.
3. Email service reads from queue and sends OTP via SMTP.
4. User **verifies** email using `/verify` endpoint.
5. On login, a **JWT token** is issued and required for protected routes.
6. If regenerate otp use `/generate-otp`

---
## To start the server 
```
  -> Make config file in the Backend_API name config/local.yaml 
  -> Make a env file in the Backend_API/internal/config/ and add the address of config file that we make above
  -> Make a env file and add email data in the /Microservices/sendMail 

  then: 
  
  docker-compose up --build
```

## 🧪 Example API Usage

### Signup
```http
POST /signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123"
}

POST /verify
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "12345"
}
POST /generate-otp
Content-Type: application/json

Make sure you login first

GET /products?limit=10&page=2&sort=name&filter=category:electronics
```

## 🛠️ Tech Stack

Language: Go (Golang)

Database: Sqlite 

Cache: Redis

Queue: RabbitMQ

Email: SMTP (Migrate to Google service later on)

Auth: JWT

Container: Docker


## 🚧 TODO

 Frontend (React/Vue) – Coming Soon!

 Rate limiting and abuse prevention

 API documentation with Swagger

## 🤝 Contributing
Pull requests are welcome. For major changes, open an issue first to discuss what you would like to change.


