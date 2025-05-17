# 7-Solutions API

A RESTful API built using **Go**, **Gin**, and **MongoDB**, implementing JWT-based authentication for user registration, login, and CRUD operations.

---

## 🚀 Project Setup

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

### 🔧 Run the Project

```bash
git clone https://github.com/your-username/7-solutions.git
cd 7-solutions

cp .env.example .env

docker-compose up --build

```

# Sample API requests/responses

```bash
1. Register 
URL : http://localhost:8080/register
Method : POST
Requests :
{
  "name": "Wasawat Test",
  "email": "Yean@example.com",
  "password": "password123"
}

Responses :
{
    "message": "registered"
}

2. Login
URL : http://localhost:8080/login
Method : POST
Requests :
{
  "email": "Yean@example.com",
  "password": "password123"
}
Responses :
{
   Authorization: Bearer <token>
}

3. Get All Users
URL : http://localhost:8080/users
Method : GET
Headers Requests :
{
   Authorization: Bearer <token>
}
Responses :
[
    {
        "id": "6825f072ad10a50069b84d46",
        "name": "Wasawat Test",
        "email": "Yean@example.com",
        "created_at": "2025-05-15T13:47:30.819Z"
    }
]

4. Get User By ID
URL : http://localhost:8080/users/<id>
Method : GET
Headers Requests :
{
   Authorization: Bearer <token>
}
Responses :
{
    "id": "6825f072ad10a50069b84d46",
    "name": "Wasawat Test",
    "email": "Yean@example.com",
    "created_at": "2025-05-15T13:47:30.819Z"
}

5. Update User 
URL : http://localhost:8080/users/<id>
Method : PUT
Headers Requests :
{
   Authorization: Bearer <token>
}
Requests :
{
  "name": "Wasawat Updated",
  "email": "Yean@example.com"
}
Responses :
{
    "message": "updated"
}

6. Delete User 
URL : http://localhost:8080/users/<id>
Method : DELETE
Headers Requests :
{
   Authorization: Bearer <token>
}
Responses :
{
    "message": "deleted"
}
```

# TODO / Improvements

1. Add Swagger/OpenAPI documentation
2. Implement refresh tokens
3. Use Redis for tokens

# Project Structure
```bash
├── main.go
├── Dockerfile
├── docker-compose.yml
├── .env
├── .env.example
├── config/
│   └── config.go
├── handler/
├── middleware/
├── repository/
├── usecase/
├── utils/
├── model/
├── proto/
├── grpc/
```