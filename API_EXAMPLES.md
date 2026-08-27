# API Examples & Testing Guide

## Curl Examples

### 1. Create a Product

```bash
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Pro 15",
    "description": "High-performance laptop with 16GB RAM",
    "price": 1299.99,
    "sale_price": 999.99
  }' \
  | jq
```

**Response:**
```json
{
  "successful": true,
  "error_code": "",
  "data": {
    "data1": 1,
    "data2": "Laptop Pro 15",
    "price": 1299.99,
    "sale_price": 999.99
  }
}
```

---

### 2. Get a Product by ID

```bash
curl -X GET http://localhost:8080/product/1 \
  | jq
```

**Response:**
```json
{
  "successful": true,
  "error_code": "",
  "data": {
    "data1": 1,
    "data2": "Laptop Pro 15",
    "price": 1299.99,
    "sale_price": 999.99
  }
}
```

---

### 3. Create Multiple Products

```bash
# Product 2
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "USB Mouse",
    "description": "Wireless mouse with USB receiver",
    "price": 49.99
  }' | jq

# Product 3
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard",
    "description": "RGB backlit mechanical keyboard",
    "price": 149.99,
    "sale_price": 99.99
  }' | jq

# Product 4
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monitor 27 inch",
    "price": 349.99,
    "sale_price": 299.99
  }' | jq
```

---

### 4. Update a Product (Partial Update)

**Update only the price:**
```bash
curl -X PATCH http://localhost:8080/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 1199.99
  }' | jq
```

**Update only the name:**
```bash
curl -X PATCH http://localhost:8080/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Pro 15 Inch - Updated"
  }' | jq
```

**Update multiple fields:**
```bash
curl -X PATCH http://localhost:8080/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Pro 16",
    "price": 1499.99,
    "sale_price": 1199.99
  }' | jq
```

**Clear sale_price (set to null):**
```bash
curl -X PATCH http://localhost:8080/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "sale_price": null
  }' | jq
```

**Response:**
```json
{
  "successful": true,
  "error_code": ""
}
```

---

### 5. Delete a Product

```bash
curl -X DELETE http://localhost:8080/product/1 | jq
```

**Response:**
```json
{
  "successful": true,
  "error_code": ""
}
```

---

### 6. Error Scenarios

#### Invalid Request (Missing required field)
```bash
curl -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Missing name and price"
  }' | jq
```

**Response (400):**
```json
{
  "successful": false,
  "error_code": "INVALID_REQUEST",
  "data": null
}
```

#### Product Not Found
```bash
curl -X GET http://localhost:8080/product/99999 | jq
```

**Response (404):**
```json
{
  "successful": false,
  "error_code": "NOT_FOUND",
  "data": null
}
```

#### Invalid ID Format
```bash
curl -X GET http://localhost:8080/product/invalid | jq
```

**Response (400):**
```json
{
  "successful": false,
  "error_code": "INVALID_ID",
  "data": null
}
```

---

## HTTPie Examples

If you prefer HTTPie over curl:

### Create Product
```bash
http POST localhost:8080/product \
  name="Laptop Pro 15" \
  description="High-performance laptop" \
  price:=1299.99 \
  sale_price:=999.99
```

### Get Product
```bash
http GET localhost:8080/product/1
```

### Update Product
```bash
http PATCH localhost:8080/product/1 \
  price:=1199.99
```

### Delete Product
```bash
http DELETE localhost:8080/product/1
```

---

## Postman Collection

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Product API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Create Product",
      "request": {
        "method": "POST",
        "url": "{{base_url}}/product",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "body": {
          "mode": "raw",
          "raw": "{\"name\": \"Laptop\", \"price\": 1299.99, \"sale_price\": 999.99}"
        }
      }
    },
    {
      "name": "Get Product",
      "request": {
        "method": "GET",
        "url": "{{base_url}}/product/1"
      }
    },
    {
      "name": "Update Product",
      "request": {
        "method": "PATCH",
        "url": "{{base_url}}/product/1",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "body": {
          "mode": "raw",
          "raw": "{\"price\": 1199.99}"
        }
      }
    },
    {
      "name": "Delete Product",
      "request": {
        "method": "DELETE",
        "url": "{{base_url}}/product/1"
      }
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    }
  ]
}
```

---

## JavaScript/Fetch Examples

### Create Product
```javascript
fetch('http://localhost:8080/product', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'Laptop Pro 15',
    description: 'High-performance laptop',
    price: 1299.99,
    sale_price: 999.99
  })
})
  .then(res => res.json())
  .then(data => console.log(data));
```

### Get Product
```javascript
fetch('http://localhost:8080/product/1')
  .then(res => res.json())
  .then(data => console.log(data));
```

### Update Product
```javascript
fetch('http://localhost:8080/product/1', {
  method: 'PATCH',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    price: 1199.99
  })
})
  .then(res => res.json())
  .then(data => console.log(data));
```

### Delete Product
```javascript
fetch('http://localhost:8080/product/1', {
  method: 'DELETE'
})
  .then(res => res.json())
  .then(data => console.log(data));
```

---

## Python Examples

### Using Requests Library

```python
import requests
import json

BASE_URL = 'http://localhost:8080'

# Create Product
response = requests.post(
    f'{BASE_URL}/product',
    json={
        'name': 'Laptop Pro 15',
        'description': 'High-performance laptop',
        'price': 1299.99,
        'sale_price': 999.99
    }
)
print(response.json())

# Get Product
response = requests.get(f'{BASE_URL}/product/1')
print(response.json())

# Update Product
response = requests.patch(
    f'{BASE_URL}/product/1',
    json={'price': 1199.99}
)
print(response.json())

# Delete Product
response = requests.delete(f'{BASE_URL}/product/1')
print(response.json())
```

---

## Testing Workflow

### Setup
```bash
# Start server
go run cmd/main.go

# In another terminal, run this script
```

### Test Script (Bash)
```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Testing Product API ==="

# Create
echo -e "\n1. Create Product"
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "price": 99.99
  }')
echo $CREATE_RESPONSE | jq

# Get product ID from response
PRODUCT_ID=$(echo $CREATE_RESPONSE | jq '.data.data1')
echo "Created product ID: $PRODUCT_ID"

# Get
echo -e "\n2. Get Product"
curl -s -X GET $BASE_URL/product/$PRODUCT_ID | jq

# Update
echo -e "\n3. Update Product"
curl -s -X PATCH $BASE_URL/product/$PRODUCT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "price": 149.99
  }' | jq

# Get again to verify update
echo -e "\n4. Get Updated Product"
curl -s -X GET $BASE_URL/product/$PRODUCT_ID | jq

# Delete
echo -e "\n5. Delete Product"
curl -s -X DELETE $BASE_URL/product/$PRODUCT_ID | jq

# Try to get deleted product
echo -e "\n6. Get Deleted Product (should be not found)"
curl -s -X GET $BASE_URL/product/$PRODUCT_ID | jq
```

---

## Swagger UI Testing

1. Start the server: `go run cmd/main.go`
2. Open Swagger UI: `http://localhost:8080/api-docs/index.html`
3. Click on each endpoint
4. Click "Try it out"
5. Fill in request body
6. Click "Execute"
7. See response

---

## Load Testing (Using Apache Bench)

```bash
# Create 1000 products
ab -n 1000 -c 10 -p payload.json -T application/json http://localhost:8080/product

# Get product (simple request)
ab -n 1000 -c 10 http://localhost:8080/product/1

# Payload file (payload.json)
{
  "name": "Load Test Product",
  "price": 99.99
}
```

---

## Common Test Scenarios

### Scenario 1: Complete CRUD
1. Create a product
2. Get the product
3. Update its price
4. Get it again to verify update
5. Delete the product
6. Try to get deleted product (expect not found)

### Scenario 2: Partial Updates
1. Create product with full data
2. Update only name
3. Update only price
4. Update only sale_price
5. Set sale_price to null
6. Verify each change

### Scenario 3: Validation
1. Create without name (should fail)
2. Create with invalid price (should fail)
3. Create with negative price (should fail)
4. Get with invalid ID (should fail)
5. Get non-existent ID (should fail)

### Scenario 4: Edge Cases
1. Create with very long name
2. Create with very large price
3. Create with description containing special characters
4. Update with empty optional fields
5. Multiple updates in quick succession

---

## Monitoring Requests

Using `curl -v` for verbose output:

```bash
curl -v -X POST http://localhost:8080/product \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "price": 99.99}'
```

This shows:
- Request headers
- Response headers
- Response body
- Connection details
- Timing information

---

## Environment for Testing

Set these environment variables before testing:

```bash
export BASE_URL="http://localhost:8080"
export DB_HOST="localhost"
export DB_PORT="5432"
```

Then use in curl:
```bash
curl -X GET $BASE_URL/product/1
```

---

**Happy Testing!** 🚀

For more information, see [README.md](README.md) or [QUICKSTART.md](QUICKSTART.md)
