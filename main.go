package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from Render environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Home endpoint - FIXED: Proper routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "WhatsApp Server is running!",
			"version": "1.0",
			"endpoints": map[string]string{
				"health": "/health",
				"send":   "/send (POST)",
				"status": "/status",
			},
		})
	})

	// Health endpoint - FIXED: Exact path
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "healthy",
			"platform": "Render.com",
		})
	})

	// Send endpoint - FIXED: Exact path
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/send" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid JSON: " + err.Error(),
			})
			return
		}

		// Basic validation
		if request.To == "" || request.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Both 'to' and 'message' fields are required",
			})
			return
		}

		log.Printf("Message to %s: %s", request.To, request.Message)

		// Success response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Message received successfully",
			"data": map[string]string{
				"to":      request.To,
				"message": request.Message,
				"note":    "WhatsApp integration pending",
			},
		})
	})

	// Status endpoint - FIXED: Exact path
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "running",
			"platform":  "Render.com",
			"whatsapp":  "not_connected",
			"timestamp": "now",
		})
	})

	// Simple QR endpoint - FIXED: Exact path
	http.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/qr" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "QR code endpoint",
			"note":    "QR functionality will be implemented here",
		})
	})

	// Start server - FIXED: Better logging
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Endpoints available:")
	log.Printf("   📍 GET  /        - Home")
	log.Printf("   ❤️  GET  /health  - Health check")
	log.Printf("   💌 POST /send    - Send message")
	log.Printf("   📊 GET  /status  - Status")
	log.Printf("   📱 GET  /qr      - QR code")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
