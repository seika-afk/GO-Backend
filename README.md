# Restaurant Management GO Backend
REST API for managing users, tables, menus, foods, orders, order items, and invoices for a restaurant

---
## Tech Stack

- Go
- Gin
- MongoDB
- JWT auth

---
The API connects to the `restaurant` database and uses these collections:

- `user`
- `table`
- `menu`
- `food`
- `order`
- `orderItem`
- `invoice`
---

## API Overview

### Users

Public routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/users` | List users with pagination |
| GET | `/users/:user_id` | Get a single user by `user_id` |
| POST | `/users/signup` | Create a new user |
| POST | `/users/login` | Log in and return user details plus tokens |

#### `POST /users/signup`

Creates a user and generates both access and refresh tokens.

Example body:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "password": "secret123",
  "email": "john@example.com",
  "avatar": "https://example.com/avatar.png",
  "phone": "9876543210"
}
```

#### `POST /users/login`

Authenticates with email and password.

Example body:

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

#### `GET /users`

Optional query params:

- `page`
- `recordPerPage`

Example:

```text
/users?page=1&recordPerPage=10
```

### Tables

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/tables` | List all tables |
| GET | `/tables/:table_id` | Get a table by `table_id` |
| POST | `/tables` | Create a table |
| PATCH | `/tables/:table_id` | Update a table |

#### `POST /tables`

Example body:

```json
{
  "number_of_guests": 4,
  "table_number": 12
}
```

#### `PATCH /tables/:table_id`

Send only the fields you want to update.

Example body:

```json
{
  "number_of_guests": 6,
  "table_number": 15
}
```

### Menus

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/menus` | Menu lookup route, currently wired to the same handler as `GET /menus/:menu_id` |
| GET | `/menus/:menu_id` | Get a menu item by `menu_id` |
| POST | `/menus` | Create a menu |
| PATCH | `/menus/:menu_id` | Update a menu |

#### `POST /menus`

Example body:

```json
{
  "name": "Lunch Special",
  "category": "Main Course",
  "start_date": "2026-05-17T12:00:00Z",
  "end_date": "2026-05-17T15:00:00Z"
}
```

#### `PATCH /menus/:menu_id`

You can update:

- `name`
- `category`
- `start_date`
- `end_date`

When both dates are sent, the API checks that the time window is valid.

### Foods

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/foods` | List foods with pagination |
| GET | `/foods/:food_id` | Get a food item by `food_id` |
| POST | `/foods` | Create a food item |
| POST | `/foods/:food_id` | Update a food item |

#### `GET /foods`

Optional query params:

- `page`
- `recordPerPage`
- `startIndex`

Example:

```text
/foods?page=1&recordPerPage=10
```

#### `POST /foods`

Example body:

```json
{
  "name": "Paneer Tikka",
  "price": 180,
  "food_image": "https://example.com/paneer-tikka.png",
  "menu_id": "menu-id-here"
}
```

`menu_id` must exist before the food can be created.

#### `POST /foods/:food_id`

Update supports partial fields:

- `name`
- `price`
- `food_image`
- `menu_id`

### Orders

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/orders` | List all orders |
| GET | `/orders/:order_id` | Get a single order |
| POST | `/orders` | Create an order |
| PATCH | `/orders/:order_id` | Update an order |

#### `POST /orders`

Example body:

```json
{
  "order_date": "2026-05-17T10:30:00Z",
  "table_id": "table-id-here"
}
```

`table_id` must exist before the order is created.

#### `PATCH /orders/:order_id`

Currently supports updating the linked `table_id`.

### Order Items

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/orderItems` | List all order items |
| GET | `/orderItems/:order_item_id` | Get one order item |
| GET | `/orderItems-order/:order_id` | Get all items for one order, grouped with table and total due |
| POST | `/orderItems` | Create an order plus its items |
| PATCH | `/orderItems/:order_item_id` | Update an order item |

#### `POST /orderItems`

This endpoint creates:

- one order
- multiple order items attached to that order

Example body:

```json
{
  "table_id": "table-id-here",
  "order_items": [
    {
      "food_id": "food-id-1",
      "quantity": "S",
      "unit_price": 120
    },
    {
      "food_id": "food-id-2",
      "quantity": "M",
      "unit_price": 150
    }
  ]
}
```

`quantity` is validated as a size label.

#### `GET /orderItems-order/:order_id`

Returns aggregated order details for:

- food name and image
- quantity
- table number
- total payment due

This is the data source used by invoice details.

### Invoices

Protected routes.

| Method | Path | Description |
| --- | --- | --- |
| GET | `/invoices` | List all invoices |
| GET | `/invoices/:invoice_id` | Get invoice details with order summary |
| POST | `/invoices` | Create an invoice |
| PATCH | `/invoices/:invoice_id` | Update invoice payment details |

#### `POST /invoices`

Example body:

```json
{
  "order_id": "order-id-here",
  "payment_method": "CASH",
  "payment_status": "PENDING"
}
```

If `payment_status` is omitted, it defaults to `PENDING`.

`order_id` must exist before the invoice can be created.

#### `GET /invoices/:invoice_id`

Returns an invoice view with:

- invoice metadata
- order id
- payment due
- table number
- order item breakdown

## Common ID Fields

The server generates string IDs after insert:

- `user_id`
- `table_id`
- `menu_id`
- `food_id`
- `order_id`
- `order_item_id`
- `invoice_id`
