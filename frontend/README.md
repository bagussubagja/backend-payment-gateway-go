# Payment Gateway Frontend

Simple React frontend for the Payment Gateway backend.

## Features

- User registration and login
- Dashboard with profile information
- Create payment transactions (regular and QRIS)
- View transaction history
- Responsive design with inline styles

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm start
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

## API Configuration

The frontend is configured to connect to the backend at `http://192.168.1.8:8080`. 

To change the API URL, edit the `API_BASE_URL` in `src/services/api.js`.

## Pages

- **Home** (`/`) - Landing page
- **Login** (`/login`) - User login
- **Register** (`/register`) - User registration
- **Dashboard** (`/dashboard`) - User profile and navigation
- **Payment** (`/payment`) - Create new payments
- **History** (`/history`) - View transaction history

## Authentication

The app uses JWT tokens stored in localStorage for authentication. Protected routes automatically redirect to login if no token is present.