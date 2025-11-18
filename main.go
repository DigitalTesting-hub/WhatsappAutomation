package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// WhatsApp Web implementation
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

// Global variables
var (
	client          *whatsmeow.Client
	isConnected     bool = false
	qrCodeData      string = ""
	connectionStatus string = "disconnected"
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

	// Initialize WhatsApp client
	initWhatsApp()

	// Register handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/disconnect", disconnectHandler)
	http.HandleFunc("/qr", getQRHandler)

	log.Printf("🚀 WhatsApp Server with REAL implementation starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initWhatsApp() {
	// Initialize database
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	// Get the first device
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Printf("Failed to get device: %v", err)
		return
	}

	// Create client
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	// Add event handler for connection
	client.AddEventHandler(eventHandler)

	log.Printf("✅ WhatsApp client initialized - ready for connection")
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		log.Printf("✅ WhatsApp connected successfully")
		isConnected = true
		connectionStatus = "connected"
	case *events.Disconnected:
		log.Printf("❌ WhatsApp disconnected")
		isConnected = false
		connectionStatus = "disconnected"
	case *events.LoggedOut:
		log.Printf("🔒 WhatsApp logged out")
		isConnected = false
		connectionStatus = "disconnected"
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>WhatsApp - REAL Implementation</title>
		<meta charset="utf-8">
		<style>
			body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
			.status { padding: 15px; border-radius: 5px; margin: 20px 0; text-align: center; font-weight: bold; }
			.connected { background: #d4edda; color: #155724; }
			.disconnected { background: #f8d7da; color: #721c24; }
			.connecting { background: #fff3cd; color: #856404; }
			.btn { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; margin: 5px; }
			.btn-success { background: #28a745; }
			.btn-danger { background: #dc3545; }
			.console { background: #000; color: #0f0; padding: 15px; border-radius: 5px; font-family: monospace; margin: 10px 0; }
		</style>
	</head>
	<body>
		<h1>📱 WhatsApp - REAL Implementation</h1>
		
		<div class="status %s">Status: %s</div>

		<div>
			<button onclick="connectWhatsApp()" class="btn btn-success">🔗 Connect WhatsApp</button>
			<button onclick="getQRCode()" class="btn">📱 Get QR Code</button>
			<button onclick="disconnectWhatsApp()" class="btn btn-danger">❌ Disconnect</button>
		</div>

		<div id="qrContainer" style="margin: 20px 0;"></div>

		<div>
			<h3>Send Message</h3>
			<input type="text" id="phone" placeholder="Phone (with country code)" value="1234567890" style="width: 200px; padding: 8px; margin: 5px;">
			<input type="text" id="message" placeholder="Message" value="Hello from REAL WhatsApp!" style="width: 300px; padding: 8px; margin: 5px;">
			<button onclick="sendMessage()" class="btn btn-success">📤 Send Real Message</button>
		</div>

		<div id="response" class="console">Response will appear here...</div>

		<script>
			function connectWhatsApp() {
				fetch('/connect', { method: 'POST' })
					.then(r => r.json())
					.then(data => {
						document.getElementById('response').innerText = JSON.stringify(data, null, 2);
						if(data.qr_code) {
							document.getElementById('qrContainer').innerHTML = 
								'<h4>Scan QR Code:</h4><div style="background:#fff;padding:20px;">' + 
								data.qr_code + '</div>';
						}
					});
			}

			function getQRCode() {
				fetch('/qr')
					.then(r => r.json())
					.then(data => {
						document.getElementById('response').innerText = JSON.stringify(data, null, 2);
					});
			}

			function disconnectWhatsApp() {
				fetch('/disconnect', { method: 'POST' })
					.then(r => r.json())
					.then(data => {
						document.getElementById('response').innerText = JSON.stringify(data, null, 2);
					});
			}

			function sendMessage() {
				const phone = document.getElementById('phone').value;
				const message = document.getElementById('message').value;
				
				fetch('/send', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ to: phone, message: message })
				})
				.then(r => r.json())
				.then(data => {
					document.getElementById('response').innerText = JSON.stringify(data, null, 2);
				});
			}

			setInterval(() => {
				fetch('/status')
					.then(r => r.json())
					.then(data => {
						const statusDiv = document.querySelector('.status');
						if(data.connected) {
							statusDiv.className = 'status connected';
							statusDiv.innerHTML = 'Status: ✅ Connected to WhatsApp';
						} else {
							statusDiv.className = 'status disconnected';
							statusDiv.innerHTML = 'Status: ❌ Disconnected';
						}
					});
			}, 3000);
		</script>
	</body>
	</html>
	`

	statusClass := "disconnected"
	statusText := "❌ Disconnected"
	if isConnected {
		statusClass = "connected"
		statusText = "✅ Connected to WhatsApp"
	} else if connectionStatus == "connecting" {
		statusClass = "connecting"
		statusText = "⏳ Connecting..."
	}

	fmt.Fprintf(w, html, statusClass, statusText)
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if client == nil {
		initWhatsApp()
	}

	if client.Store.ID != nil {
		// Already logged in
		isConnected = true
		connectionStatus = "connected"
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Already connected to WhatsApp",
			"connected": true,
		})
		return
	}

	// Generate QR code
	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	connectionStatus = "connecting"

	// Wait for QR code
	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrCodeData = evt.Code
				log.Printf("📱 QR code generated: %s", evt.Code)
			}
		}
	}()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"message":  "QR code generation started",
		"qr_code":  qrCodeData,
		"connected": false,
	})
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isConnected || client == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "WhatsApp not connected",
		})
		return
	}

	// Send real WhatsApp message
	jid := types.NewJID(req.To, types.DefaultUserServer)
	_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: &req.Message,
	})

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Failed to send message: " + err.Error(),
		})
		return
	}

	log.Printf("✅ Real WhatsApp message sent to %s", req.To)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Real WhatsApp message sent successfully",
		"to":      req.To,
		"message": req.Message,
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connected": isConnected,
		"status":    connectionStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	if client != nil {
		client.Disconnect()
	}
	isConnected = false
	connectionStatus = "disconnected"
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp disconnected",
	})
}

func getQRHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"qr_code":   qrCodeData,
		"has_qr":    qrCodeData != "",
		"connected": isConnected,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"whatsapp":  isConnected,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
