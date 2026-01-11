# Helpdesk Ticketing API
Clean Architecture REST API menggunakan Golang

---

## Project Overview

Project ini adalah Helpdesk Ticketing System yang dirancang untuk organisasi/perusahaan yang membutuhkan sistem untuk mengelola permintaan bantuan teknis (technical support) secara terstruktur.

Project ini adibangun menggunakan **Clean Architecture**.  
Tujuan utama arsitektur ini adalah menjaga **business logic tetap terisolasi** dari framework, database, dan delivery mechanism.




---

## API Information

- **Base URL**:  http://localhost:8080
- **Format**: JSON  
- **Auth**: JWT Bearer Token 

## API Endpoints


### 1.  Login


**Request:**
```json
POST /api/v1/login
Content-Type: application/json

{
    "username": "foo",
    "password": "1234567" 
}
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgxNTg4OTQsImlhdCI6MTc2ODExNTY5NCwicm9sZSI6IlVTRVIiLCJ1c2VyX2lkIjoiVVNFUi0yNDE1ZjA0Mi1mZDZlLTRlNjEtYjQ2Ny1iNTBlMGZkOWUxMjIifQ.K6IMJk5tNDeueI7iK5O6mI2I5M1BFhE9Gbs-zng2ys0"
}
```
**Failed Response:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "error": "invalid credentials: username tidak ditemukan"
}

``` 


### 2.  Register


**Request:**
```json
POST /api/v1/register
Content-Type: application/json

{
    "email":"bananajoe@nope.com",
    "username":"bananajoe",
    "password":"pass123",
    "department":"HR department",
    "role":"USER", // role ADMIN or USER
    "is_remote":"false"
}

```

**Successful Response:**
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
    "message": "User registered successfully",
    "user": {
        "user_id": "USER-48050fe4-aeb5-403c-b121-91bf116c43be",
        "username": "bananajoe",
        "email": "bananajoe@nope.com",
        "role": "USER",
        "profile_pict": "",
        "department": "HR department",
        "remote": false,
        "created_at": "2026-01-11T14:21:21.11+07:00",
        "updated_at": "2026-01-11T14:21:21.11+07:00"
    }
}
```
**Failed Response:**
```json
HTTP/1.1 400 Bad request
Content-Type: application/json

{
   "error": "User sudah ada"
}

``` 

```json

{
    "error": "Role must be either 'ADMIN' or 'USER'"
}

``` 


### 3.  Create ticket


**Request:**
```json
POST /api/v1/tickets
Headers: Authorization: Bearer <token>
Content-Type: application/json

{
    "title":"Printer Rusak",
    "description":"Printer di ruang HR tidak bisa connect ke pc ..."
}
```

**Successful Response:**
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
    "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
    "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
    "title": "Printer Rusak",
    "description": "Printer di ruang HR tidak bisa connect ke pc ...",
    "status": "OPEN",
    "created_at": "2026-01-11T14:29:24.299+07:00",
    "updated_at": "2026-01-11T14:29:24.299+07:00",
    "deleted_at": null,
    "owner": null
}
```
**Failed Response:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "error": "Invalid token"
}

``` 
```json
HTTP/1.1 400 Bad request
Content-Type: application/json

{
   "error": "descripton is empty"
}

{
    "error": "title is required"
}

``` 


### 4.  Get My ticket


**Request:**
```json
GET /api/v1/tickets/my-tickets
Headers: Authorization: Bearer <token>
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
        "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
        "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
        "title": "Printer Rusak",
        "description": "Printer di ruang HR tidak bisa connect ke pc ...",
        "status": "OPEN",
        "created_at": "2026-01-11T14:29:24.299+07:00",
        "updated_at": "2026-01-11T14:29:24.299+07:00",
        "deleted_at": null,
        "owner": null
    },
    {
        "ticket_id": "tick-3bc2047a-999b-4dd4-9bdf-e8e25094b177",
        "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
        "title": "Speed Please i need this",
        "description": "My PM kinda project less",
        "status": "OPEN",
        "created_at": "2026-01-11T14:28:08.778+07:00",
        "updated_at": "2026-01-11T14:28:08.778+07:00",
        "deleted_at": null,
        "owner": null
    },
    {
        "ticket_id": "tick-6d7bd2a0-9f6c-4319-8046-2bcdf4b9a1b2",
        "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
        "title": "Speed Please i need this",
        "description": "My PM kinda project less",
        "status": "OPEN",
        "created_at": "2026-01-06T21:25:22.458+07:00",
        "updated_at": "2026-01-06T21:25:22.458+07:00",
        "deleted_at": null,
        "owner": null
    }
]
```
**Failed Response:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "error": "Invalid token"
}

``` 


### 5.  Get all ticket (admin)


**Request:**
```json
GET /api/v1/tickets/admin/all
Headers: Authorization: Bearer <token>
```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json

[
    {
    "Jumlah": 3,
    "tickets": [
        {
            "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
            "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
            "title": "Printer Rusak",
            "description": "Printer di ruang HR tidak bisa connect ke pc ...",
            "status": "OPEN",
            "created_at": "2026-01-11T14:29:24.299+07:00",
            "updated_at": "2026-01-11T14:29:24.299+07:00",
            "deleted_at": null,
            "owner": {
                "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
                "username": "abdan",
                "email": "abdan@ninga.com",
                "role": "USER",
                "profile_pict": "",
                "department": "IT department",
                "remote": false,
                "created_at": "2026-01-06T21:04:04.802+07:00",
                "updated_at": "2026-01-06T21:04:04.802+07:00"
            }
        },
        {
            "ticket_id": "tick-6d7bd2a0-9f6c-4319-8046-2bcdf4b9a1b2",
            "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
            "title": "Speed Please i need this",
            "description": "My PM kinda project less",
            "status": "OPEN",
            "created_at": "2026-01-06T21:25:22.458+07:00",
            "updated_at": "2026-01-06T21:25:22.458+07:00",
            "deleted_at": null,
            "owner": {
                "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
                "username": "abdan",
                "email": "abdan@ninga.com",
                "role": "USER",
                "profile_pict": "",
                "department": "IT department",
                "remote": false,
                "created_at": "2026-01-06T21:04:04.802+07:00",
                "updated_at": "2026-01-06T21:04:04.802+07:00"
            }
        },
        {
            "ticket_id": "tick-d1f559f4-59d8-4e61-b357-9f92d21c5ca9",
            "user_id": "USER-63a17844-6562-40dc-994a-4240b74287d8",
            "title": "Butuh akun Microsoft baru kak",
            "description": "Ini akun aku yang lama dikira scam , jadi di ban oleh community , help meee",
            "status": "OPEN",
            "created_at": "2025-12-31T14:15:23.212+07:00",
            "updated_at": "2025-12-31T14:15:23.212+07:00",
            "deleted_at": null,
            "owner": {
                "user_id": "USER-63a17844-6562-40dc-994a-4240b74287d8",
                "username": "huan",
                "email": "huan@hual.com",
                "role": "USER",
                "profile_pict": "",
                "department": "",
                "remote": false,
                "created_at": "2025-12-31T12:11:32.951+07:00",
                "updated_at": "2025-12-31T12:11:32.951+07:00"
            }
        }
    ]
}
]
```
**Failed Response:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "error": "Invalid token"
}

``` 



### 6.  Change status ticket (admin)


**Request:**
```json
/api/v1/tickets/admin/{{idTicket}}/status
Headers: Authorization: Bearer <token>

{
    "status":"DONE" //DONE , OPEN , or IN_PROGRESS
}

```

**Successful Response:**
```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "message": "update success",
    "ticket": {
        "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
        "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
        "title": "Printer Rusak",
        "description": "Printer di ruang HR tidak bisa connect ke pc ...",
        "status": "DONE",
        "created_at": "2026-01-11T14:29:24.299+07:00",
        "updated_at": "2026-01-11T14:47:14.184+07:00",
        "deleted_at": null,
        "owner": {
            "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
            "username": "abdan",
            "email": "abdan@ninga.com",
            "role": "USER",
            "profile_pict": "",
            "department": "IT department",
            "remote": false,
            "created_at": "2026-01-06T21:04:04.802+07:00",
            "updated_at": "2026-01-06T21:04:04.802+07:00"
        }
    }
}

```
**Failed Response:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "error": "Invalid token"
}

``` 
```json
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
    "error": "Status must be  'DONE' or 'IN_PROGRESS'"
}

``` 

