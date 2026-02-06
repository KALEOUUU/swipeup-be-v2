# Swipeup - School Canteen POS System

A Point of Sale (POS) system for school canteens built with Go, Gin-Gonic, GORM, and MySQL.

## Project Structure

```
swipeup-admin-v2/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go
├── internal/
│   ├── api/                 # API handlers
│   │   ├── admin/           # Admin handlers (users, products, categories, orders, transactions)
│   │   ├── siswa/           # Student handlers (profile, balance, orders, transactions)
│   │   └── auth/           # Authentication handlers (login, logout, refresh)
│   ├── app/                 # Core application logic
│   │   ├── database/       # Database connection and configuration
│   │   ├── models/          # GORM models (User, Product, Category, Order, OrderItem, Transaction)
│   │   └── utils/          # Utility functions
│   └── routes/              # Routing and middleware
│       ├── routes.go        # Route definitions
│       └── middleware.go    # Authentication and authorization middleware
├── configs/                # Configuration files
├── docs/                   # Documentation
├── .env.example           # Environment variables template
├── go.mod                # Go module definition
└── README.md             # This file
```

## Features

### Admin Features
- User management (CRUD operations)
- Product management (CRUD operations)
- Category management (CRUD operations)
- Order management (view all orders)
- Transaction management (view all transactions)
- Balance top-up for students

### Student Features
- View profile
- Check balance
- View order history
- View transaction history

### Authentication
- Login with student ID
- JWT-based authentication (TODO: implement proper JWT)
- Role-based access control (admin/student)

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd swipeup-admin-v2
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
mysql -u root -p
CREATE DATABASE swipeup;
```

4. Configure environment variables:
```bash
cp .env.example .env
```

Edit `.env` file with your database credentials:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=swipeup
SERVER_PORT=8080
GIN_MODE=debug
```

## Running the Application

1. Start the server:
```bash
go run cmd/server/main.go
```

Or build and run:
```bash
go build -o swipeup cmd/server/main.go
./swipeup
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /api/v1/health` - Health check endpoint

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh token

### Student Endpoints (Protected)
- `GET /api/v1/siswa/profile` - Get student profile
- `GET /api/v1/siswa/balance` - Get student balance
- `GET /api/v1/siswa/orders` - Get student orders
- `GET /api/v1/siswa/transactions` - Get student transactions

### Admin Endpoints (Protected + Admin Role)

#### Users
- `GET /api/v1/admin/users` - Get all users
- `GET /api/v1/admin/users/:id` - Get user by ID
- `POST /api/v1/admin/users` - Create new user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user
- `POST /api/v1/admin/users/:id/topup` - Top-up user balance

#### Products
- `GET /api/v1/admin/products` - Get all products
- `GET /api/v1/admin/products/:id` - Get product by ID
- `POST /api/v1/admin/products` - Create new product
- `PUT /api/v1/admin/products/:id` - Update product
- `DELETE /api/v1/admin/products/:id` - Delete product

#### Categories
- `GET /api/v1/admin/categories` - Get all categories
- `GET /api/v1/admin/categories/:id` - Get category by ID
- `POST /api/v1/admin/categories` - Create new category
- `PUT /api/v1/admin/categories/:id` - Update category
- `DELETE /api/v1/admin/categories/:id` - Delete category

#### Orders
- `GET /api/v1/admin/orders` - Get all orders
- `GET /api/v1/admin/orders/:id` - Get order by ID
- `POST /api/v1/admin/orders` - Create new order
- `PUT /api/v1/admin/orders/:id` - Update order
- `DELETE /api/v1/admin/orders/:id` - Delete order

#### Transactions
- `GET /api/v1/admin/transactions` - Get all transactions
- `GET /api/v1/admin/transactions/:id` - Get transaction by ID
- `GET /api/v1/admin/transactions/user/:user_id` - Get transactions by user

## Database Models

### User
- `id` - Primary key
- `student_id` - Unique student ID
- `name` - User name
- `email` - User email
- `phone` - User phone number
- `role` - User role (student, teacher, admin)
- `class` - User class
- `balance` - Account balance
- `is_active` - Account status
- `rfid_card` - RFID card number
- `password` - Hashed password

### Product
- `id` - Primary key
- `name` - Product name
- `description` - Product description
- `category_id` - Foreign key to Category
- `price` - Product price
- `stock` - Available stock
- `image_url` - Product image URL
- `is_active` - Product status

### Category
- `id` - Primary key
- `name` - Category name
- `description` - Category description
- `is_active` - Category status

### Order
- `id` - Primary key
- `order_number` - Unique order number
- `user_id` - Foreign key to User
- `total_amount` - Order total
- `status` - Order status (pending, completed, cancelled)
- `payment_method` - Payment method (card, cash)

### OrderItem
- `id` - Primary key
- `order_id` - Foreign key to Order
- `product_id` - Foreign key to Product
- `quantity` - Item quantity
- `price` - Item price
- `subtotal` - Item subtotal

### Transaction
- `id` - Primary key
- `transaction_number` - Unique transaction number
- `user_id` - Foreign key to User
- `type` - Transaction type (top_up, purchase, refund)
- `amount` - Transaction amount
- `balance_before` - Balance before transaction
- `balance_after` - Balance after transaction
- `description` - Transaction description
- `order_id` - Foreign key to Order (nullable)

## Development

### Adding New Features

1. Create or update models in `internal/app/models/`
2. Create handlers in `internal/api/admin/` or `internal/api/siswa/`
3. Add routes in `internal/routes/routes.go`
4. Add middleware if needed in `internal/routes/middleware.go`

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -ldflags="-s -w" -o swipeup cmd/server/main.go
```

## TODO

- [ ] Implement proper JWT token generation and validation
- [ ] Add password hashing with bcrypt
- [ ] Implement proper error handling
- [ ] Add input validation
- [ ] Add logging
- [ ] Add database migrations
- [ ] Add API documentation (Swagger)
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add Docker support
- [ ] Add CI/CD pipeline

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
