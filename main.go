package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Global variables to track WhatsApp status
var (
	isWhatsAppConnected bool   = false
	connectionStatus    string = "disconnected"
	whatsappSession     string = ""
)

type SendRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Register handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/qr", qrHandler)
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/disconnect", disconnectHandler)

	log.Printf("🚀 WhatsApp Server starting on port %s", port)
	log.Printf("✅ All endpoints ready:")
	log.Printf("   📍 GET  /           - Home page")
	log.Printf("   ❤️  GET  /health     - Health check")
	log.Printf("   💌 POST /send       - Send WhatsApp message")
	log.Printf("   📊 GET  /status     - Connection status")
	log.Printf("   📱 GET  /qr         - QR code page")
	log.Printf("   🔗 GET  /connect    - Connect WhatsApp")
	log.Printf("   ❌ GET  /disconnect - Disconnect WhatsApp")

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
			.connecting { background: #fff3cd; color: #856404; border: 1px solid #ffeaa7; }
			.qr-container { background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0; text-align: center; }
			.endpoints { background: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
			.endpoint { margin: 10px 0; padding: 10px; background: white; border-radius: 3px; }
			.code { background: #2d3748; color: #e2e8f0; padding: 10px; border-radius: 3px; font-family: monospace; font-size: 12px; }
			.btn { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; margin: 5px; }
			.btn-success { background: #28a745; }
			.btn-danger { background: #dc3545; }
			.console { background: #000; color: #0f0; padding: 15px; border-radius: 5px; font-family: monospace; font-size: 12px; margin: 10px 0; max-height: 200px; overflow-y: auto; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>📱 WhatsApp Web Server</h1>
			
			<div class="status %s">
				🔔 WhatsApp Status: %s
			</div>

			<div class="qr-container">
				<h3>🔗 Connect WhatsApp</h3>
				<div style="background: white; padding: 20px; margin: 10px 0; border-radius: 5px;">
					<div style="font-size: 80px; margin: 20px;">📱</div>
					<p><strong>Simulation Mode Active</strong></p>
					<p>Click "Connect WhatsApp" below to simulate connection.</p>
					<p>In production, this would show a real QR code.</p>
				</div>
				<p><strong>📋 Instructions:</strong></p>
				<ol style="text-align: left;">
					<li>Click <strong>"Connect WhatsApp"</strong> below</li>
					<li>Wait 5 seconds for auto-connection</li>
					<li>Status will change to <strong>"Connected"</strong></li>
					<li>Then you can send test messages</li>
				</ol>
				<div>
					<a href="/connect" class="btn btn-success">🔗 Connect WhatsApp</a>
					<a href="/disconnect" class="btn btn-danger">❌ Disconnect</a>
					<a href="/status" class="btn">📊 Status</a>
				</div>
			</div>

			<div class="endpoints">
				<h3>🔧 Send Message API</h3>
				<div class="endpoint">
					<strong>💌 Send WhatsApp Message:</strong>
					<div class="code">
POST /send<br>
Content-Type: application/json<br>
{<br>
  "to": "1234567890",<br>
  "message": "Hello from WhatsApp API!"<br>
}
					</div>
				</div>

				<div style="margin-top: 15px;">
					<strong>Try it now:</strong>
					<div style="background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 10px 0;">
						<input type="text" id="phoneNumber" placeholder="Phone number" style="width: 200px; padding: 8px; margin: 5px; border: 1px solid #ddd; border-radius: 3px;" value="1234567890">
						<input type="text" id="messageText" placeholder="Your message" style="width: 300px; padding: 8px; margin: 5px; border: 1px solid #ddd; border-radius: 3px;" value="Hello from API Dashboard!">
						<button onclick="sendMessage()" class="btn btn-success">📤 Send Message</button>
					</div>
					<div id="response" class="console">Response will appear here...</div>
				</div>
			</div>

			<div style="margin-top: 20px; text-align: center;">
				<small>🚀 Server running on Render.com | Version 1.0 | Simulation Mode</small>
			</div>
		</div>

		<script>
			function sendMessage() {
				const phone = document.getElementById('phoneNumber').value;
				const message = document.getElementById('messageText').value;
				const responseDiv = document.getElementById('response');
				
				responseDiv.innerHTML = '📤 Sending message...';
				
				fetch('/send', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						to: phone,
						message: message
					})
				})
				.then(response => response.json())
				.then(data => {
					responseDiv.innerHTML = '✅ Response:\\n' + JSON.stringify(data, null, 2);
				})
				.catch(error => {
					responseDiv.innerHTML = '❌ Error:\\n' + error.toString();
				});
			}

			// Auto-refresh status every 5 seconds
			setInterval(() => {
				fetch('/status')
					.then(response => response.json())
					.then(data => {
						if (data.whatsapp_connected !== %t) {
							location.reload();
						}
					});
			}, 5000);
		</script>
	</body>
	</html>
	`

	statusClass := "disconnected"
	statusText := "Disconnected"
	if isWhatsAppConnected {
		statusClass = "connected"
		statusText = "Connected ✓"
	} else if connectionStatus == "connecting" {
		statusClass = "connecting"
		statusText = "Connecting..."
	}

	fmt.Fprintf(w, html, statusClass, statusText, isWhatsAppConnected)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"platform":  "Render.com",
		"timestamp": time.Now().Format(time.RFC3339),
		"whatsapp":  isWhatsAppConnected,
		"version":   "1.0",
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

	if !isWhatsAppConnected {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "WhatsApp is not connected",
			"action":  "Please connect WhatsApp first",
			"data": map[string]string{
				"to":      req.To,
				"message": req.Message,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp message delivered successfully",
		"data": map[string]string{
			"to":        req.To,
			"message":   req.Message,
			"timestamp": time.Now().Format(time.RFC3339),
			"message_id": fmt.Sprintf("WA_%d", time.Now().Unix()),
			"status":    "delivered",
		},
	})
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
		"version":            "1.0",
	})
}

func qrHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	connectionStatus = "connecting"
	isWhatsAppConnected = false

	log.Printf("🔄 Starting WhatsApp connection simulation")

	go func() {
		time.Sleep(5 * time.Second)
		isWhatsAppConnected = true
		connectionStatus = "connected"
		log.Printf("✅ WhatsApp connected successfully (simulation)")
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp connection started",
		"data": map[string]string{
			"status":    "connecting",
			"wait_time": "5 seconds",
			"note":      "Simulation mode - will auto-connect",
		},
	})
}

func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	isWhatsAppConnected = false
	connectionStatus = "disconnected"
	whatsappSession = ""

	log.Printf("🔌 WhatsApp disconnected")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp disconnected successfully",
	})
}
