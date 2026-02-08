# Student Order Management API

This directory contains API documentation for student order management endpoints.

## Available Endpoints

### Basic Order Operations

1. **Get Orders** (`get-orders.bru`)
   - Get all orders for the authenticated student
   - Returns orders with product details

2. **Create Order** (`create-order.bru`)
   - Create a new order (direct/quick order)
   - Manual order creation for students

### 📊 Advanced Reporting Features

3. **Get Orders by Month** (`get-orders-monthly.bru`)
   - **NEW**: Get order history filtered by specific month and year
   - Includes monthly summary statistics
   - Query parameters: `year`, `month`
   - Returns: orders array + summary (total_orders, completed_orders, pending_orders, total_spending)

4. **Get Order Receipt** (`get-order-receipt.bru`)
   - **NEW**: Generate printable HTML receipt for any order
   - Path parameter: `id` (order ID)
   - Returns: HTML content ready for printing
   - Perfect for expense tracking and record keeping

## Authentication

All endpoints require **Student** authentication. Use a valid student token obtained from `auth/login.bru`.

## Common Query Parameters

### Get Orders by Month
```
GET /api/v1/siswa/orders/monthly?year=2024&month=2
```
- `year` (required): Year to filter (e.g., 2024)
- `month` (required): Month number 1-12 (e.g., 2 for February)

### Get Order Receipt
```
GET /api/v1/siswa/orders/15/receipt
```
- `id` (required): Order ID in URL path

## Response Examples

### Monthly Orders Response
```json
{
  "orders": [
    {
      "id": 15,
      "order_number": "ORD-51-51-1770581708",
      "user_id": 51,
      "total_amount": 23520,
      "status": "done",
      "payment_method": "cash",
      "created_at": "2024-02-09T03:15:08.698Z",
      "order_items": [
        {
          "product_id": 5,
          "quantity": 1,
          "price": 23520,
          "subtotal": 23520,
          "product": {
            "name": "Montblanc",
            "price": 23520
          }
        }
      ]
    }
  ],
  "summary": {
    "year": "2024",
    "month": "2",
    "total_orders": 8,
    "completed_orders": 7,
    "pending_orders": 1,
    "total_spending": 187360
  }
}
```

### Receipt Response
Returns HTML content that can be opened in a browser and printed. The HTML includes:
- Order header with order number
- Customer information
- Order date and payment method
- Itemized table with products, quantities, and prices
- Total amount
- Professional styling for printing

## Use Cases

- **Order History**: Get Orders, Get Orders by Month
- **Expense Tracking**: Get Orders by Month with spending summary
- **Record Keeping**: Get Order Receipt for printing/saving
- **Order Management**: Create Order for quick purchases
- **Financial Planning**: Monthly spending analysis

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (missing/invalid parameters)
- `401`: Unauthorized (invalid or missing token)
- `403`: Forbidden (accessing other student's orders)
- `404`: Not Found (order doesn't exist)
- `500`: Internal Server Error

Error responses include a JSON object with an `error` field describing the issue.