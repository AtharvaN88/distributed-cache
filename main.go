package main

import (
	"cache-project/cache"
	"cache-project/handlers"
	"cache-project/hashring"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Get environment variables
	selfAddr := os.Getenv("SELF_ADDR") // e.g. "cache-a:8081"
	peers := strings.Split(os.Getenv("PEERS"), ",") // e.g. "cache-a:8081,cache-b:8082,cache-c:8083"

	// Setup cache + ring
	c := cache.NewCache(60)
	hr := hashring.New(peers)

	// Setup handlers
	h := &handlers.Handler{
		Cache:    c,
		HashRing: hr,
		SelfAddr: selfAddr,
	}

	http.HandleFunc("/set", h.MakeSetHandler())
	http.HandleFunc("/get/", h.MakeGetHandler())

	http.ListenAndServe(":"+strings.Split(selfAddr, ":")[1], nil)
}