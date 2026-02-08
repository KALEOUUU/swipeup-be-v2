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
  - Create Order (direct/quick order)
  - Delete Order (cancel order before completion)
  - **Cart Management** (new enhanced flow)
    - Get Cart
    - Add Item to Cart
    - Update Cart Item
    - Remove Item from Cart
    - Clear Cart
    - Checkout (create orders from cart)
  - Get Transactions
  - **📊 Monthly Reporting** (new analytics features)
    - Get Orders by Month (monthly order history with summary)
    - Get Order Receipt (printable order receipt)

- **stand/** - Stand Admin endpoints (requires stand token)
  - Products Management
  - Orders Management
    - Get Orders
    - Get Pending Orders
    - Get Order by ID
    - Create Order
    - Update Order Status
    - **Delete Order** (cancel order before completion)
    - **📊 Monthly Reporting** (new analytics features)
      - Get Orders by Month (detailed monthly orders with summary)
      - Get Monthly Revenue Recap (annual revenue analytics)
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

## 📊 New Reporting Features

### Student Monthly Analytics

#### 1. Get Orders by Month
**Endpoint**: `GET /api/v1/siswa/orders/monthly?year=2024&month=2`

Returns student's order history for a specific month with summary statistics:
- Complete order details with product information
- Monthly totals: total orders, completed orders, pending orders, total spending
- Perfect for personal expense tracking and order history review

### Stand Admin Monthly Analytics

Two new endpoints have been added to provide comprehensive business intelligence for stand administrators:

#### 2. Get Orders by Month
**Endpoint**: `GET /api/v1/stand/orders/monthly?year=2024&month=2`

Returns detailed orders for a specific month along with summary statistics:
- Complete order details with user and product information
- Monthly totals: total orders, completed orders, pending orders, total revenue
- Perfect for detailed monthly performance analysis

#### 3. Get Monthly Revenue Recap
**Endpoint**: `GET /api/v1/stand/orders/revenue/monthly?year=2024`

Provides annual revenue analytics with monthly breakdown:
- 12-month revenue and order statistics
- Yearly summary totals
- Ideal for trend analysis, financial planning, and business reporting

### Student Receipt Printing

#### 4. Get Order Receipt
**Endpoint**: `GET /api/v1/siswa/orders/{id}/receipt`

Generates printable HTML receipt for any completed order:
- Professional receipt layout with order details
- Itemized list with quantities and prices
- Ready-to-print HTML format
- Perfect for record keeping and expense tracking

### Usage Examples

```bash
# Student: Get February 2024 order history
GET /api/v1/siswa/orders/monthly?year=2024&month=2

# Student: Get printable receipt for order #15
GET /api/v1/siswa/orders/15/receipt

# Stand Admin: Get February 2024 orders with details
GET /api/v1/stand/orders/monthly?year=2024&month=2

# Stand Admin: Get complete 2024 revenue analytics
GET /api/v1/stand/orders/revenue/monthly?year=2024
```

These features help both students and stand administrators track orders, manage finances, and maintain proper records for their canteen operations.
