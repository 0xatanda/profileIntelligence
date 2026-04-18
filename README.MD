# Profile Intelligence API

A backend service that enriches a user-provided name using external APIs, classifies the data, persists it in PostgreSQL, and exposes RESTful endpoints to manage profiles.

---

## 🚀 Overview

This service:

* Accepts a name input
* Fetches enriched data from 3 external APIs:

  * Genderize
  * Agify
  * Nationalize
* Processes and classifies the data
* Stores it in PostgreSQL
* Exposes REST endpoints to create, retrieve, filter, and delete profiles

---

## 🧠 External APIs Used

* [https://api.genderize.io?name={name}](https://api.genderize.io?name={name})
* [https://api.agify.io?name={name}](https://api.agify.io?name={name})
* [https://api.nationalize.io?name={name}](https://api.nationalize.io?name={name})

---

## ⚙️ Features

* Multi-API data aggregation
* Age classification (child, teenager, adult, senior)
* Country inference from probability ranking
* Idempotent profile creation (no duplicates by name)
* PostgreSQL persistence
* Filtering support (gender, country, age group)
* Full CRUD API
* Structured JSON responses
* CORS enabled

---

## 📦 Tech Stack

* Go (net/http)
* PostgreSQL
* SQL Driver: lib/pq
* Systemd (Linux service management)
* Nginx (reverse proxy)
* UUID v7 for IDs

---

## 🗄️ Database Schema

```sql
CREATE TABLE profiles (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE,
    gender TEXT,
    gender_probability FLOAT,
    sample_size INT,
    age INT,
    age_group TEXT,
    country_id TEXT,
    country_probability FLOAT,
    created_at TIMESTAMP
);
```

---

## 📡 API Endpoints

---

### 1. Create Profile

**POST** `/api/profiles`

Request:

```json
{
  "name": "ella"
}
```

Response (201):

```json
{
  "status": "success",
  "data": {
    "id": "uuid-v7",
    "name": "ella",
    "gender": "female",
    "gender_probability": 0.99,
    "sample_size": 1234,
    "age": 46,
    "age_group": "adult",
    "country_id": "NG",
    "country_probability": 0.85,
    "created_at": "2026-04-18T12:00:00Z"
  }
}
```

### Idempotency

If name already exists:

```json
{
  "status": "success",
  "message": "Profile already exists",
  "data": { ...existing profile }
}
```

---

### 2. Get Profile by ID

**GET** `/api/profiles/{id}`

Response:

```json
{
  "status": "success",
  "data": { ...profile }
}
```

---

### 3. Get All Profiles

**GET** `/api/profiles`

Optional filters:

* `gender`
* `country_id`
* `age_group`

Example:

```
/api/profiles?gender=male&country_id=NG
```

Response:

```json
{
  "status": "success",
  "count": 2,
  "data": []
}
```

---

### 4. Delete Profile

**DELETE** `/api/profiles/{id}`

Response:

* `204 No Content`

---

## ⚠️ Error Handling

Standard format:

```json
{
  "status": "error",
  "message": "error description"
}
```

### Error Codes

* 400 → Missing name
* 422 → Invalid request body
* 404 → Profile not found
* 502 → External API failure
* 500 → Server error

---

## ⚠️ Edge Cases

* Gender = null or count = 0 → 502 error
* Age = null → 502 error
* No country returned → 502 error

---

## 🌍 Deployment (EC2)

### Run service

```bash
sudo systemctl start profile-api
sudo systemctl status profile-api
```

### Nginx config

```nginx
server {
    listen 80;
    server_name YOUR_IP_OR_DOMAIN;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Restart nginx:

```bash
sudo systemctl restart nginx
```

---

## 🧪 Testing

### Create profile

```bash
curl -X POST http://localhost:8080/api/profiles \
-H "Content-Type: application/json" \
-d '{"name":"ella"}'
```

### Get all profiles

```bash
curl http://localhost:8080/api/profiles
```

---

## 🧱 Build

```bash
go mod tidy
go build -o app ./cmd/server
```


