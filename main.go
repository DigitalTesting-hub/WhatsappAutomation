package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "WhatsApp Server is running on Render!",
			"version": "1.0.0",
		})
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		if r.Method != "POST" {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Only POST method allowed",
			})
			return
		}

		var request struct {
			To      string `json:"to"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid JSON: " + err.Error(),
			})
			return
		}

		log.Printf("Message to %s: %s", request.To, request.Message)

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Message received - WhatsApp integration pending",
			"to":      request.To,
			"text":    request.Message,
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy", 
			"platform": "Render.com",
		})
	})

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
