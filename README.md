# Simple Payment Gateway Backend (Go)

This is a simple backend project functioning as a payment gateway, integrated with **Midtrans**. It is built using **Go (Golang)** with a clean, production-ready architecture.

---

## Features

- User registration and login (JWT authentication)
- Create payment transactions (Midtrans Snap integration)
- Check transaction status
- View transaction history (for logged-in users)
- Webhook to receive payment status notifications from Midtrans
- Database schema for users, transactions, and transaction items

---

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Payment Gateway**: Midtrans
- **Authentication**: JWT
- **Containerization**: Docker

---

## How to Run Locally

1. **Clone the Repository**
   ```bash
   git clone https://github.com/bagussubagja/backend-payment-gateway-go.git
   cd backend-payment-gateway-go
   ```

2. **Configure Environment Variables**
   Copy the `.env.example` file to `.env` and fill in your credentials.
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Run Database Migration (Once)**
   Temporarily uncomment the `AutoMigrate` block in `storage/postgres.go`, then run the server to create tables. Re-comment it afterward.

5. **Start the Application**
   ```bash
   go run main.go
   ```
   The server will run at `http://localhost:8080`.

---

## Environment Variables (.env)

Make sure the following variables are set in your `.env` file:

- `PORT`: App server port (e.g., `8080`)
- `DB_HOST`: PostgreSQL host (e.g., `localhost`)
- `DB_PORT`: PostgreSQL port (e.g., `5432`)
- `DB_USER`: PostgreSQL username (e.g., `postgres`)
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: Database name (e.g., `payment_db`)
- `JWT_SECRET_KEY`: Secret key for signing JWT
- `JWT_EXPIRATION_HOURS`: Token expiration in hours (e.g., `24h`)
- `MIDTRANS_SERVER_KEY`: Midtrans server key
- `MIDTRANS_CLIENT_KEY`: Midtrans client key
- `MIDTRANS_ENVIRONMENT`: Midtrans environment (`sandbox` or `production`)

---

## API Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `POST /api/v1/payments/notification` - Midtrans webhook notification
- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/payments/create` - Create a new payment transaction
- `POST /api/v1/payments/qris` - Create a QRIS transaction
- `GET /api/v1/payments/status/:orderID` - Get transaction status by order ID
- `GET /api/v1/payments/history` - Get user transaction history
