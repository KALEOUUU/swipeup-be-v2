# Stand Order Management API

This directory contains API documentation for stand administrator order management endpoints.

## Available Endpoints

### Basic Order Operations

1. **Get Orders** (`get-orders.bru`)
   - Get all orders for the authenticated stand
   - Returns orders with user and product details

2. **Get Pending Orders** (`get-pending-orders.bru`)
   - Get orders that are waiting for processing
   - Useful for real-time order queue management

3. **Get Order by ID** (`get-order-by-id.bru`)
   - Get detailed information for a specific order
   - Includes full order items and customer details

4. **Create Order** (`create-order.bru`)
   - Create a new order (typically for walk-in customers)
   - Manual order creation for stand staff

5. **Update Order Status** (`update-order-status.bru`)
   - Update order status (request → cooking → done)
   - Status progression management

6. **Delete Order** (`delete-order.bru`)
   - Cancel/delete an order before completion
   - Soft delete with audit trail

### 📊 Advanced Reporting Features

7. **Get Orders by Month** (`get-orders-monthly.bru`)
   - **NEW**: Get detailed orders filtered by specific month and year
   - Includes monthly summary statistics
   - Query parameters: `year`, `month`
   - Returns: orders array + summary (total_orders, completed_orders, pending_orders, total_revenue)

8. **Get Monthly Revenue Recap** (`get-monthly-revenue-recap.bru`)
   - **NEW**: Annual revenue analytics with monthly breakdown
   - Shows 12-month performance data
   - Query parameter: `year` (optional, defaults to current year)
   - Returns: monthly_data array + yearly_summary

## Authentication

All endpoints require **Stand Admin** authentication. Use a valid stand token obtained from `auth/login.bru`.

## Common Query Parameters

### Get Orders by Month
```
GET /api/v1/stand/orders/monthly?year=2024&month=2
```
- `year` (required): Year to filter (e.g., 2024)
- `month` (required): Month number 1-12 (e.g., 2 for February)

### Get Monthly Revenue Recap
```
GET /api/v1/stand/orders/revenue/monthly?year=2024
```
- `year` (optional): Year to analyze (defaults to current year)

## Response Examples

### Monthly Orders Response
```json
{
  "orders": [...],
  "summary": {
    "year": "2024",
    "month": "2",
    "total_orders": 15,
    "completed_orders": 12,
    "pending_orders": 3,
    "total_revenue": 352800
  }
}
```

### Revenue Recap Response
```json
{
  "year": "2024",
  "monthly_data": [
    {
      "month": 1,
      "month_name": "January",
      "total_orders": 45,
      "completed_orders": 42,
      "total_revenue": 1058400
    }
    // ... 11 more months
  ],
  "yearly_summary": {
    "total_orders": 520,
    "completed_orders": 485,
    "total_revenue": 12241600
  }
}
```

## Use Cases

- **Daily Operations**: Get Orders, Get Pending Orders, Update Order Status
- **Order Management**: Create Order, Delete Order, Get Order by ID
- **Business Analytics**: Get Orders by Month, Get Monthly Revenue Recap
- **Financial Reporting**: Monthly revenue tracking and trend analysis
- **Performance Monitoring**: Order completion rates and revenue analytics

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (missing/invalid parameters)
- `401`: Unauthorized (invalid or missing token)
- `404`: Not Found (order doesn't exist)
- `500`: Internal Server Error

Error responses include a JSON object with an `error` field describing the issue.