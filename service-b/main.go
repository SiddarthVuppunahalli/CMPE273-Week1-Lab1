package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var client = &http.Client{Timeout: 1 * time.Second}

func health(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	log.Printf("service=B endpoint=/health status=200 latency_ms=%d", time.Since(start).Milliseconds())
}

func callEcho(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	msg := r.URL.Query().Get("msg")

	serviceAURL := fmt.Sprintf("http://127.0.0.1:8080/echo?msg=%s", url.QueryEscape(msg))
	resp, err := client.Get(serviceAURL)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service_b": "ok",
			"service_a": "unavailable",
			"error":     err.Error(),
		})
		log.Printf("service=B endpoint=/call-echo status=503 error=%q latency_ms=%d", err.Error(), time.Since(start).Milliseconds())
		return
	}
	defer resp.Body.Close()

	var data map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&data)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"service_b": "ok",
		"service_a": data,
	})
	log.Printf("service=B endpoint=/call-echo status=200 latency_ms=%d", time.Since(start).Milliseconds())
}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/call-echo", callEcho)
	log.Println("service=B listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
