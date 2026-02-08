# Swipeup API Documentation

This directory contains Bruno API documentation for testing the Swipeup School Canteen POS system.

## Bruno Setup

1. Install Bruno from https://www.usebruno.com/downloads
2. Open Bruno and import this `docs` folder
3. Set your environment variables in Bruno (optional)
4. Run the requests to test the API

## Service Structure

- **auth/** - Authentication endpoints (public)
  - Login
  - Logout
  - Refresh Token
  - Register (create new user)

- **student/** - Student endpoints (requires student token)
  - Get Profile
  - Get Balance
  - Get Orders
  - Create Order
  - **Delete Order** (cancel order before completion)
  - Get Transactions

- **stand/** - Stand Admin endpoints (requires stand token)
  - Products Management
  - Orders Management
    - Get Orders
    - Get Pending Orders
    - Get Order by ID
    - Create Order
    - Update Order Status
    - **Delete Order** (cancel order before completion)
  - Settings Management

- **admin/** - Admin endpoints (requires admin token)
  - Users Management
  - Categories Management
  - Stand Canteens Management
  - Global Settings Management

## Authentication

Most endpoints require authentication. Use the `auth/login.bru` request to get a token for the appropriate role:

- **Student**: Login with student credentials to get student token
- **Stand Admin**: Login with stand credentials to get stand token
- **Admin**: Login with admin credentials to get admin token

**Public Registration**: Use `auth/register.bru` to create a new user account without authentication.

Tokens are automatically used in subsequent requests via Bruno's environment variables.

## Environment Variables

Set these in Bruno environment:

```
BASE_URL=http://localhost:8080
STUDENT_TOKEN=<your_student_token>
STAND_TOKEN=<your_stand_token>
ADMIN_TOKEN=<your_admin_token>
```

## API Base URL

```
http://localhost:8080/api/v1
```
