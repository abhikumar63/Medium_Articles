# Idempotent Orders – Go + Gin + Redis demo

A bite-sized example that shows how to build **idempotent REST endpoints** in Go.  
Hit `/orders` as many times as you like with the same **Idempotency-Key** and body – you’ll always get the exact same response, and the order is created only once.

---

## 🌟 What is an idempotent API?

*An operation is idempotent when performing it one time or many times has the **same side-effect** on the server.*

| Safe by default | Can be idempotent with care |
|-----------------|-----------------------------|
| `GET /resource` | `PUT /resource/123`  |
|                 | `DELETE /resource/123` |
|                 | `POST /orders` **with an Idempotency-Key** |

Why you want this:

* **Network retries** (mobile drops, load balancer timeouts) won’t duplicate payments or orders.  
* **At-least-once delivery** queues can re-deliver safely.  
* **Predictability** makes SDKs, workflows, and audits simpler.

---

## 🏃‍♂️ Quick start

```bash
# clone
git clone https://github.com/<your-handle>/idempotent-orders-go.git
cd idempotent-orders-go

# deps
go mod tidy

# start Redis (if you don’t already have one)
redis-server &          # default port 6379

# run the API
go run main.go
```

---

## 🖥️ How the demo works

1. Client sends:
    - JSON body describing the order
    - Idempotency-Key: <any-unique-string>

2. Server computes:
```
hash = SHA256(body || idempotency-key)
```

3. Cache lookup
    - Key exists in Redis → return the cached JSON + original status (200 OK).
    - Key missing → create the order (fake business logic), store the response & status with a TTL (24 h), reply with 201 Created.

4. Result: identical key/body combos never create duplicates.

The core logic lives in main.go (~140 lines including boilerplate).

---
## Test it Yourself
```
# 1️⃣ One request, one order
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: abc123" \
  -d '{"item":"Apple Watch Ultra","qty":1,"price":799.00}'

# 2️⃣ Retry the exact same request
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: abc123" \
  -d '{"item":"Apple Watch Ultra","qty":1,"price":799.00}'

# 3️⃣ Same key, different body → new order
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: abc123" \
  -d '{"item":"AirPods Pro","qty":2,"price":249.00}'

# 4️⃣ Different key, same body → new order
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: xyz789" \
  -d '{"item":"Apple Watch Ultra","qty":1,"price":799.00}'

# 5️⃣ Missing key → 400 error
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"item":"MacBook Air M3","qty":1,"price":1299.00}'

# 6️⃣ Malformed JSON → 400 error
curl -i -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: badjson1" \
  -d '{"item":"iPad","qty":2,}'
```
| Test | Expected status     | Note                    |
| ---- | ------------------- | ----------------------- |
| 1    | **201 Created**     | Order created & cached  |
| 2    | **200 OK**          | Cache hit, no duplicate |
| 3    | **201 Created**     | Body changed → new hash |
| 4    | **201 Created**     | Key changed → new hash  |
| 5    | **400 Bad Request** | Header missing          |
| 6    | **400 Bad Request** | Invalid JSON            |

---
## 🛠️ Extending the example
Swap in PostgreSQL or MongoDB for durable orders.

Move the hashing & lookup into Gin middleware for reuse.

Parameterize TTL via env var.

Add unit tests (cache hit/miss, validation).

Protect secrets/config with Docker + Compose for easier spin-up.
