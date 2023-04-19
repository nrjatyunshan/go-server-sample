package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

func test(w http.ResponseWriter, req *http.Request) {
	numIn := runtime.NumGoroutine()
	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	val, err := redisClient.Set("key", "value", time.Second).Result()
	if err != nil {
		panic(err.Error())
	}

	val, err = redisClient.Get("key").Result()
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "============== Redis =============\n")
	fmt.Fprintf(w, "value: %s\n", val)
	fmt.Fprintf(w, "==================================\n\n")

	// MySQL
	db, err := sql.Open("mysql", "root@tcp(db)/information_schema")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT TABLE_NAME FROM tables limit 3")
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "============== MySQL =============\n")
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "table name: %s\n", tableName)
	}
	fmt.Fprintf(w, "==================================\n\n")

	// HTTP
	// 客户端默认会启用连接池,会重用连接和协程,导致追踪失败.
	// 关闭 HTTP KeepAlive
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true
	c := &http.Client{Transport: t}
	resp, err := c.Get("http://web")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "============== HTTP =============\n")
	fmt.Fprintf(w, "StatusCode: %d\n", resp.StatusCode)
	fmt.Fprintf(w, "==================================\n\n")

	// Go
	numOut := runtime.NumGoroutine()
	fmt.Fprintf(w, "============== Go =============\n")
	fmt.Fprintf(w, "NumGoroutine: %d %d\n", numIn, numOut)
	fmt.Fprintf(w, "==================================\n\n")
}

func main() {
	// HTTP handler
	http.HandleFunc("/", test)

	// Start HTTP server
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
