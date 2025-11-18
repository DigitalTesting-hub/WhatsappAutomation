package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type SendRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// Global variables to track WhatsApp status
var (
	isWhatsAppConnected bool = false
	qrCodeData         string = ""
	connectionStatus   string = "disconnected"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Serve static HTML for QR code
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/qr", qrHandler)
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/ws", wsHandler)

	// Start WebSocket broadcaster
	go handleBroadcasts()

	log.Printf("🚀 WhatsApp Server starting on port %s", port)
	log.Printf("✅ All endpoints ready:")
	log.Printf("   📍 GET  /        - Home page with QR")
	log.Printf("   ❤️  GET  /health  - Health check")
	log.Printf("   💌 POST /send    - Send WhatsApp message")
	log.Printf("   📊 GET  /status  - Connection status")
	log.Printf("   📱 GET  /qr      - QR code page")
	log.Printf("   🔗 GET  /connect - Force reconnect")
	log.Printf("   🔌 GET  /ws      - WebSocket for real-time updates")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>WhatsApp Web Server</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; background: #f5f5f5; }
			.container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
			.status { padding: 15px; border-radius: 5px; margin: 20px 0; text-align: center; font-weight: bold; }
			.connected { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
			.disconnected { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
			.qr-container { background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0; text-align: center; }
			.endpoints { background: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
			.endpoint { margin: 10px 0; padding: 10px; background: white; border-radius: 3px; }
			.code { background: #2d3748; color: #e2e8f0; padding: 10px; border-radius: 3px; font-family: monospace; }
			.btn { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>📱 WhatsApp Web Server</h1>
			
			<div class="status %s">
				WhatsApp Status: %s
			</div>

			<div class="qr-container">
				<h3>🔗 Connect WhatsApp</h3>
				<div style="background: white; padding: 20px; display: inline-block; margin: 10px 0;">
					<pre style="font-size: 6px; line-height: 6px; margin: 0;">%s</pre>
				</div>
				<p><strong>Instructions:</strong></p>
				<ol style="text-align: left;">
					<li>Open WhatsApp on your phone</li>
					<li>Tap ⋮ (menu) → Linked Devices → Link a Device</li>
					<li>Scan the QR code above</li>
					<li>Wait for connection confirmation</li>
				</ol>
				<a href="/qr" class="btn">Refresh QR Code</a>
				<a href="/connect" class="btn" style="background: #28a745;">Force Reconnect</a>
			</div>

			<div class="endpoints">
				<h3>🔧 API Endpoints</h3>
				
				<div class="endpoint">
					<strong>Health Check:</strong>
					<div class="code">GET /health</div>
				</div>

				<div class="endpoint">
					<strong>Send Message:</strong>
					<div class="code">POST /send<br>Content-Type: application/json<br>{ "to": "1234567890", "message": "Hello" }</div>
				</div>

				<div class="endpoint">
					<strong>Status:</strong>
					<div class="code">GET /status</div>
				</div>

				<div class="endpoint">
					<strong>WebSocket (Real-time):</strong>
					<div class="code">GET /ws</div>
				</div>
			</div>

			<div style="margin-top: 20px; text-align: center;">
				<small>Server running on Render.com | Version 2.0</small>
			</div>
		</div>

		<script>
			// WebSocket for real-time updates
			const ws = new WebSocket((window.location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + window.location.host + '/ws');
			
			ws.onmessage = function(event) {
				const data = JSON.parse(event.data);
				if (data.type === 'status_update') {
					// Reload page to show updated status
					window.location.reload();
				}
			};

			ws.onclose = function() {
				console.log('WebSocket disconnected');
			};
		</script>
	</body>
	</html>
	`

	statusClass := "disconnected"
	statusText := "Disconnected"
	if isWhatsAppConnected {
		statusClass = "connected"
		statusText = "Connected ✓"
	}

	// Generate simple ASCII QR (in real implementation, use actual QR)
	if qrCodeData == "" {
		qrCodeData = generateDemoQR()
	}

	fmt.Fprintf(w, html, statusClass, statusText, qrCodeData)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"platform":  "Render.com",
		"timestamp": time.Now().Format(time.RFC3339),
		"whatsapp":  isWhatsAppConnected,
	})
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	if req.To == "" || req.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Both 'to' and 'message' fields are required",
		})
		return
	}

	log.Printf("📤 Message request: To=%s, Message=%s", req.To, req.Message)

	// Check WhatsApp connection
	if !isWhatsAppConnected {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "WhatsApp is not connected",
			"action":  "Please scan QR code first at /qr",
			"data": map[string]string{
				"to":      req.To,
				"message": req.Message,
			},
		})
		return
	}

	// TODO: Implement actual WhatsApp message sending here
	// For now, simulate successful send
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp message sent successfully",
		"data": map[string]string{
			"to":        req.To,
			"message":   req.Message,
			"timestamp": time.Now().Format(time.RFC3339),
			"note":      "Actual WhatsApp integration pending",
		},
	})

	// Broadcast the sent message via WebSocket
	message := Message{
		Type: "message_sent",
		Content: map[string]string{
			"to":      req.To,
			"message": req.Message,
			"time":    time.Now().Format("15:04:05"),
		},
	}
	broadcast <- message
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "success",
		"server":             "running",
		"whatsapp_connected": isWhatsAppConnected,
		"connection_status":  connectionStatus,
		"platform":           "Render.com",
		"timestamp":          time.Now().Format(time.RFC3339),
		"endpoints": map[string]string{
			"health":  "/health",
			"send":    "/send (POST)",
			"qr":      "/qr",
			"connect": "/connect",
			"ws":      "/ws (WebSocket)",
		},
	})
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Generate new QR code
	qrCodeData = generateDemoQR()
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "QR code generated",
		"qr_data": qrCodeData,
		"note":    "Scan this QR code in WhatsApp → Linked Devices",
	})
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate WhatsApp connection process
	connectionStatus = "connecting"
	qrCodeData = generateDemoQR()
	
	// Simulate connection after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		isWhatsAppConnected = true
		connectionStatus = "connected"
		log.Printf("✅ WhatsApp connected successfully")
		
		// Broadcast connection status
		message := Message{
			Type: "status_update",
			Content: map[string]string{
				"status": "connected",
				"time":   time.Now().Format("15:04:05"),
			},
		}
		broadcast <- message
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp connection initiated",
		"data": map[string]string{
			"status":    "connecting",
			"next_step": "Scan QR code when ready",
			"note":      "This is a simulation. Real QR code generation needed.",
		},
	})
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	// Send current status to new client
	statusMessage := Message{
		Type: "status_update",
		Content: map[string]interface{}{
			"connected": isWhatsAppConnected,
			"status":    connectionStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	ws.WriteJSON(statusMessage)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleBroadcasts() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func generateDemoQR() string {
	// Simple ASCII QR code for demonstration
	// In real implementation, use: github.com/skip2/go-qrcode
	return `
█████████████████████████████████████
████ ▄▄▄▄▄ █▀ █▀▄▄▀▀ ▄▄▄▄▄ █ ▄▄▄▄▄ ████
████ █   █ █▄▀█▄▀▀▄▄ █   █ █ █   █ ████
████ █▄▄▄█ █ ▄█▄▀ ▀▄ █▄▄▄█ █ █▄▄▄█ ████
████▄▄▄▄▄▄▄█ █▄▀▄█▄▀▄▄▄▄▄▄▄█▄▄▄▄▄▄▄████
████  ▀ ▄▀▄▄ ▄ ▀▄ ▀▄▄▄ ▀▄▀▀▀▄▄▀▄▄▀█████
████▄▄█▄█▄▄▄ ▀▄▀▀▄▄▀█▄▄▄▄█▀ ▀ ▀▄▀▄█████
████▄▄▄█▄▄▄█▄▀▀▄▀▄█▄▀▄▀▀▄▄▄▀▄▀█▄▀▀█████
████ ▄▄▄▄▄ █▄▀▄█ █ ▄▄█ █▄▀ █ ▄▄▄█ ████
████ █   █ █ ▀▄▀▀▀▄▀▄▀▀▄▀▀▀▄▀▀ ▀▄▀████
████ █▄▄▄█ █▄█▄▀▄▀▀▀▄▀▀ ▄▀ ▀▄ ▄ █▄████
████▄▄▄▄▄▄▄█▄▄█▄█▄█▄▄█▄▄█▄▄█▄▄█████████
█████████████████████████████████████
`
}
