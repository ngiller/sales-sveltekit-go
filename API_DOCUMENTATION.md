# Magnum Sales API Documentation

This documentation provides details for all API endpoints in the backend.

## Base URL
`http://localhost:8080/api`

---

## 1. Authentication

### Login
Authenticates a user and returns a JWT token. Sets an HTTPOnly cookie named `token`.
*   **Method:** `POST`
*   **Path:** `/login`
*   **Auth Required:** No
*   **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
*   **Result (Success):**
    ```json
    {
      "data": {
        "token": "jwt_token_string",
        "user": { "id": 1, "name": "Admin", ... },
        "menus": [ { "id": 1, "menu_name": "Dashboard", ... } ]
      }
    }
    ```

### Logout
Clears the HTTPOnly session cookie.
*   **Method:** `POST`
*   **Path:** `/logout`
*   **Auth Required:** Yes

### Profile
Returns the current logged-in user profile and accessible menu tree.
*   **Method:** `GET`
*   **Path:** `/profile`
*   **Auth Required:** Yes

---

## 2. Users Management
**Policy Table:** `users`

### List All Users
*   **Method:** `GET`
*   **Path:** `/users`
*   **Policy Action:** `read`
*   **Query Parameters:**
    *   `search` (string): Search by name or email.
    *   `page` (int): Page number (default: 1).
    *   `limit` (int): Items per page (default: 50).
    *   `sort` (string): Column to sort by (default: "name").
    *   `order` (string): Sort direction ("asc" or "desc").

### Get User by ID
*   **Method:** `GET`
*   **Path:** `/users/:id`
*   **Policy Action:** `read`

### Create User
*   **Method:** `POST`
*   **Path:** `/users`
*   **Policy Action:** `create`
*   **Request Body:**
    ```json
    {
      "name": "Full Name",
      "email": "email@example.com",
      "password": "securepassword",
      "user_group_id": 1,
      "departement_id": 1,
      "phone_no": "08123456789",
      "enable": true
    }
    ```

### Update User
*   **Method:** `PUT`
*   **Path:** `/users/:id`
*   **Policy Action:** `update`

### Upload Avatar
*   **Method:** `POST`
*   **Path:** `/users/:id/avatar`
*   **Policy Action:** `update`
*   **Body:** `multipart/form-data` with field `avatar` (file)

### Upload Signature
*   **Method:** `POST`
*   **Path:** `/users/:id/signature`
*   **Policy Action:** `update`
*   **Body:** `multipart/form-data` with field `signature` (file)

---

## 3. Roles (User Groups)
**Policy Table:** `usergroups`

### List All Roles
*   **Method:** `GET`
*   **Path:** `/roles`
*   **Query Parameters:** `search`, `page`, `limit`, `sort`, `order`.

### Create/Update Role
*   **Method:** `POST` / `PUT`
*   **Request Body:**
    ```json
    {
      "name": "Manager Group"
    }
    ```

---

## 4. Departements
**Policy Table:** `departements`

### CRUD Operations
*   **List:** `GET /departements` (Query: `search`, `page`, `limit`, `sort`, `order`)
*   **Create:** `POST /departements` (Transaction)
*   **Update:** `PUT /departements/:id` (Transaction)
*   **Delete:** `DELETE /departements/:id` (Transaction)
*   **Request Body:**
    ```json
    {
      "name": "Accounting"
    }
    ```

---

## 5. Customer Categories
**Policy Table:** `customer category`

### CRUD Operations
*   **List:** `GET /customer-categories` (Query: `search`, `page`, `limit`, `sort`, `order`)
*   **Create:** `POST /customer-categories` (Transaction)
*   **Update:** `PUT /customer-categories/:id` (Transaction)
*   **Delete:** `DELETE /customer-categories/:id` (Transaction)
*   **Request Body:**
    ```json
    {
      "name": "Retail"
    }
    ```

---

## 6. Customers
**Policy Table:** `customers`

### List All Customers
*   **Method:** `GET`
*   **Path:** `/customers`
*   **Query Parameters:** `search`, `page`, `limit`, `sort`, `order`.

### Create Customer
*   **Method:** `POST`
*   **Path:** `/customers`
*   **Request Body:**
    ```json
    {
      "category_id": 1,
      "name": "PT. Example Jaya",
      "address": "Jl. Merdeka No. 1",
      "phone": "021-123456",
      "email": "info@example.com",
      "enable": true,
      "property_id": 1
    }
    ```

---

## 7. Customer Contacts
**Policy Table:** `customer_contact`

### List Contacts by Customer
*   **Method:** `GET`
*   **Path:** `/customer-contacts?customer_id=1`

### Create Contact
*   **Method:** `POST`
*   **Path:** `/customer-contacts`
*   **Request Body:**
    ```json
    {
      "customer_id": 1,
      "name": "John Doe",
      "phone": "08111222333",
      "email": "john@example.com",
      "position": "Procurement Manager"
    }
    ```

---

## 8. User Group Policies
**Policy Table:** `usergroupspolicies`

### List All Policies
*   **Method:** `GET`
*   **Path:** `/policies`

### Get Policies by Group ID
*   **Method:** `GET`
*   **Path:** `/policies/group/:groupID`

### Create Policy
*   **Method:** `POST`
*   **Path:** `/policies`
*   **Request Body:**
    ```json
    {
      "group_id": 1,
      "table_id": 3,
      "table_name": "usergroups",
      "action": "read",
      "property_id": 1
    }
    ```

---
 
## 9. Payment Terms
**Policy Table:** `payment term`
 
 ### CRUD Operations
 *   **List:** `GET /payment-terms` (Query: `search`, `page`, `limit`, `sort`, `order`)
 *   **Create:** `POST /payment-terms`
 *   **Update:** `PUT /payment-terms/:id`
 *   **Delete:** `DELETE /payment-terms/:id`
 *   **Request Body:**
     ```json
     {
       "name": "30 Days",
       "day": 30
     }
     ```
 
 ---
 
## 10. Project Levels
**Policy Table:** `project level`
 
 ### CRUD Operations
 *   **List:** `GET /project-levels`
 *   **Create:** `POST /project-levels`
 *   **Update:** `PUT /project-levels/:id`
 *   **Delete:** `DELETE /project-levels/:id`
 *   **Request Body:**
     ```json
     {
       "name": "Difficult"
     }
     ```
 
 ---
 
## 11. Project Priorities
**Policy Table:** `project priority`
 
 ### CRUD Operations
 *   **List:** `GET /project-priorities`
 *   **Create:** `POST /project-priorities`
 *   **Update:** `PUT /project-priorities/:id`
 *   **Delete:** `DELETE /project-priorities/:id`
 *   **Request Body:**
     ```json
     {
       "name": "Urgent"
     }
     ```
 
 ---

## Global Response Format

### Success Response (Single Item)
```json
{
  "data": { ... }
}
```

### Success Response (List with Pagination)
```json
{
  "data": {
    "items": [ ... ],
    "total": 100,
    "page": 1,
    "limit": 50
  }
}
```

### Error Response (400/401/403/500)
```json
{
  "message": "Error description message"
}
```
