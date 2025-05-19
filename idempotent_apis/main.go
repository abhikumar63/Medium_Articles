package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Order struct {
	ID    string  `json:"id"`
	Item  string  `json:"item"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ttl = 24 * time.Hour // how long we keep idempotency keys
)

func hash(body []byte, key string) string {
	h := sha256.Sum256(append(body, []byte(key)...))
	return hex.EncodeToString(h[:])
}

func createOrder(c *gin.Context) {
	idKey := c.GetHeader("Idempotency-Key")
	if idKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
		return
	}

	body, _ := c.GetRawData()
	cacheKey := "idem:" + hash(body, idKey)

	// Was this request (body+key) seen before?
	if cached, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var resp gin.H
		_ = json.Unmarshal([]byte(cached), &resp)
		c.JSON(http.StatusOK, resp)
		return
	}

	// Fake business logic
	var o Order
	if err := json.Unmarshal(body, &o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.ID = time.Now().UTC().Format("20060102150405")

	resp := gin.H{"order": o, "status": "created"}
	respBytes, _ := json.Marshal(resp)

	// Store response against the hash for future retries
	_ = rdb.Set(ctx, cacheKey, respBytes, ttl).Err()

	c.JSON(http.StatusCreated, resp)
}

func main() {
	r := gin.Default()
	r.POST("/orders", createOrder)
	_ = r.Run(":8080")
}
