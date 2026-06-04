package handlers

import (
	"cache-project/cache"
	"cache-project/hashring"
	"io"
	"encoding/json"
	"net/http"
	"log"
	"strings"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
}

type Handler struct {
	Cache    *cache.Cache
	HashRing *hashring.HashRing
	SelfAddr string
}

func (h *Handler) MakeSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close() 

		// Find responsible node
		node := h.HashRing.GetNode(req.Key)
		log.Printf("Incoming SET: key=%s owner=%s self=%s", req.Key, node.Address, h.SelfAddr)

		if node.Address == h.SelfAddr {
			// Store locally
			h.Cache.Set(req.Key, req.Value, req.TTL)
			json.NewEncoder(w).Encode(map[string]string{"status": "stored locally"})
		} else {
			// Forward request to correct node
			forwardBody, err := json.Marshal(req)
			if err != nil {
                http.Error(w, "failed to marshal request body", http.StatusInternalServerError)
                return
            }

			resp, err := http.Post("http://"+node.Address+"/set", "application/json", strings.NewReader(string(forwardBody)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
			io.Copy(w, resp.Body)
		}
	}
}

func (h *Handler) MakeGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/get/")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		// Find responsible node
		node := h.HashRing.GetNode(key)

		if node.Address == h.SelfAddr {
			if val, ok := h.Cache.Get(key); ok {
				json.NewEncoder(w).Encode(map[string]string{"value": val})
			} else {
				json.NewEncoder(w).Encode(map[string]string{"value": ""})
			}
		} else {
			// Proxy to correct node
			resp, err := http.Get("http://" + node.Address + "/get/" + key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
			io.Copy(w, resp.Body)
		}
	}
}