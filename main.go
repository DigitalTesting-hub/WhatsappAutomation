package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// WhatsApp Business API configuration
type WhatsAppConfig struct {
	APIURL      string
	AccessToken string
	PhoneID     string
}

// Global variables
var (
	whatsappConfig *WhatsAppConfig
	isConnected    bool   = false
	qrCodeData     string = ""
	sessionID      string = ""
)

type SendRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type WhatsAppMessage struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize WhatsApp configuration
	initWhatsAppConfig()

	// Register handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/disconnect", disconnectHandler)
	http.HandleFunc("/api/setup", setupHandler)
	http.HandleFunc("/api/send-template", sendTemplateHandler)

	log.Printf("🚀 WhatsApp Business API Server starting on port %s", port)
	log.Printf("📱 Ready for REAL WhatsApp integration")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initWhatsAppConfig() {
	whatsappConfig = &WhatsAppConfig{
		APIURL:      "https://graph.facebook.com/v17.0",
		AccessToken: os.Getenv("WHATSAPP_TOKEN"),
		PhoneID:     os.Getenv("WHATSAPP_PHONE_ID"),
	}

	if whatsappConfig.AccessToken != "" && whatsappConfig.PhoneID != "" {
		isConnected = true
		log.Printf("✅ WhatsApp Business API configured")
	} else {
		log.Printf("⚠️  WhatsApp credentials not set - using simulation mode")
		log.Printf("💡 Set WHATSAPP_TOKEN and WHATSAPP_PHONE_ID environment variables")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>WhatsApp Business API - REAL</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style>
			body { font-family: Arial, sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; background: #f0f2f5; }
			.container { background: white; padding: 30px; border-radius: 15px; box-shadow: 0 2px 20px rgba(0,0,0,0.1); }
			.status { padding: 20px; border-radius: 10px; margin: 25px 0; text-align: center; font-size: 18px; font-weight: bold; }
			.connected { background: linear-gradient(135deg, #d4edda, #c3e6cb); color: #155724; border: 2px solid #28a745; }
			.disconnected { background: linear-gradient(135deg, #f8d7da, #f5c6cb); color: #721c24; border: 2px solid #dc3545; }
			.config { background: linear-gradient(135deg, #fff3cd, #ffeaa7); color: #856404; border: 2px solid #ffc107; }
			.btn { background: #007bff; color: white; padding: 12px 25px; border: none; border-radius: 8px; cursor: pointer; margin: 8px; font-size: 16px; font-weight: bold; transition: all 0.3s; }
			.btn:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,123,255,0.3); }
			.btn-success { background: #28a745; }
			.btn-success:hover { box-shadow: 0 4px 12px rgba(40,167,69,0.3); }
			.btn-danger { background: #dc3545; }
			.btn-warning { background: #ffc107; color: #000; }
			.console { background: #1a1a1a; color: #00ff00; padding: 20px; border-radius: 10px; font-family: 'Courier New', monospace; font-size: 14px; margin: 20px 0; border: 1px solid #333; }
			.input-group { margin: 15px 0; }
			.input-group label { display: block; margin-bottom: 8px; font-weight: bold; color: #333; }
			.input-group input, .input-group textarea { width: 100%; padding: 12px; border: 2px solid #ddd; border-radius: 8px; font-size: 16px; transition: border 0.3s; }
			.input-group input:focus, .input-group textarea:focus { border-color: #007bff; outline: none; }
			.section { background: #f8f9fa; padding: 20px; border-radius: 10px; margin: 20px 0; border-left: 4px solid #007bff; }
			.code { background: #2d3748; color: #e2e8f0; padding: 15px; border-radius: 8px; font-family: monospace; font-size: 14px; margin: 10px 0; }
			.step { background: white; padding: 15px; margin: 10px 0; border-radius: 8px; border-left: 4px solid #28a745; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1 style="text-align: center; color: #25D366; margin-bottom: 10px;">📱 WhatsApp Business API</h1>
			<p style="text-align: center; color: #666; font-size: 18px; margin-bottom: 30px;">Real WhatsApp Integration - Send Actual Messages</p>
			
			<div class="status %s">%s</div>

			<div class="section">
				<h3>🚀 Quick Start</h3>
				<div class="step">
					<strong>Step 1:</strong> Get WhatsApp Business API credentials from 
					<a href="https://developers.facebook.com/docs/whatsapp/cloud-api/get-started" target="_blank">Facebook Developer Portal</a>
				</div>
				<div class="step">
					<strong>Step 2:</strong> Set environment variables in Render:
					<div class="code">
WHATSAPP_TOKEN=your_access_token_here<br>
WHATSAPP_PHONE_ID=your_phone_number_id_here
					</div>
				</div>
				<div class="step">
					<strong>Step 3:</strong> Test with the form below!
				</div>
			</div>

			<div class="section">
				<h3>📤 Send Real WhatsApp Message</h3>
				<div class="input-group">
					<label>📞 Phone Number (with country code, no +):</label>
					<input type="text" id="phoneNumber" value="14155552671" placeholder="e.g., 14155552671 for US number">
				</div>
				<div class="input-group">
					<label>💬 Message:</label>
					<textarea id="messageText" rows="4" placeholder="Enter your WhatsApp message here">Hello from REAL WhatsApp Business API! This is an actual WhatsApp message.</textarea>
				</div>
				<button onclick="sendRealMessage()" class="btn btn-success">📤 Send Real WhatsApp Message</button>
			</div>

			<div class="section">
				<h3>🔧 API Testing</h3>
				<button onclick="testConnection()" class="btn">🔍 Test Connection</button>
				<button onclick="sendTemplateMessage()" class="btn btn-warning">📝 Send Template Message</button>
				<button onclick="setupWebhook()" class="btn">🔄 Setup Webhook</button>
			</div>

			<div id="response" class="console">👆 Response will appear here after API calls...</div>

			<div class="section">
				<h3>📚 API Documentation</h3>
				<div class="code">
// Send Message Endpoint<br>
POST /send<br>
Content-Type: application/json<br>
{<br>
  "to": "1234567890",<br>
  "message": "Your message here"<br>
}<br><br>
// Check Status<br>
GET /status<br><br>
// Setup Webhook<br>
POST /api/setup
				</div>
			</div>
		</div>

		<script>
			function sendRealMessage() {
				const phone = document.getElementById('phoneNumber').value;
				const message = document.getElementById('messageText').value;
				const responseDiv = document.getElementById('response');
				
				responseDiv.innerHTML = '🔄 Sending REAL WhatsApp message...';
				
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
					responseDiv.innerHTML = '📱 API Response:\\n' + JSON.stringify(data, null, 2);
					if(data.status === 'success') {
						responseDiv.style.borderColor = '#28a745';
					} else {
						responseDiv.style.borderColor = '#dc3545';
					}
				})
				.catch(error => {
					responseDiv.innerHTML = '❌ Error:\\n' + error.toString();
					responseDiv.style.borderColor = '#dc3545';
				});
			}

			function testConnection() {
				const responseDiv = document.getElementById('response');
				responseDiv.innerHTML = '🔍 Testing WhatsApp connection...';
				
				fetch('/status')
					.then(response => response.json())
					.then(data => {
						responseDiv.innerHTML = '📊 Connection Status:\\n' + JSON.stringify(data, null, 2);
					})
					.catch(error => {
						responseDiv.innerHTML = '❌ Error:\\n' + error.toString();
					});
			}

			function sendTemplateMessage() {
				const responseDiv = document.getElementById('response');
				responseDiv.innerHTML = '📝 Sending template message...';
				
				fetch('/api/send-template', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					}
				})
				.then(response => response.json())
				.then(data => {
					responseDiv.innerHTML = '📋 Template Response:\\n' + JSON.stringify(data, null, 2);
				})
				.catch(error => {
					responseDiv.innerHTML = '❌ Error:\\n' + error.toString();
				});
			}

			function setupWebhook() {
				const responseDiv = document.getElementById('response');
				responseDiv.innerHTML = '🔄 Setting up webhook...';
				
				fetch('/api/setup', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					}
				})
				.then(response => response.json())
				.then(data => {
					responseDiv.innerHTML = '🌐 Webhook Setup:\\n' + JSON.stringify(data, null, 2);
				})
				.catch(error => {
					responseDiv.innerHTML = '❌ Error:\\n' + error.toString();
				});
			}

			// Auto-refresh status
			setInterval(() => {
				fetch('/status')
					.then(response => response.json())
					.then(data => {
						const statusDiv = document.querySelector('.status');
						if(data.connected) {
							statusDiv.className = 'status connected';
							statusDiv.innerHTML = '✅ CONNECTED - Ready to send REAL WhatsApp messages';
						} else if(data.configured) {
							statusDiv.className = 'status config';
							statusDiv.innerHTML = '⚙️  CONFIGURED - Set WHATSAPP_TOKEN & WHATSAPP_PHONE_ID in environment variables';
						} else {
							statusDiv.className = 'status disconnected';
							statusDiv.innerHTML = '❌ NOT CONFIGURED - Setup required for real WhatsApp messages';
						}
					});
			}, 5000);
		</script>
	</body>
	</html>
	`

	statusClass := "disconnected"
	statusText := "❌ NOT CONFIGURED - Setup required for real WhatsApp messages"

	if isConnected {
		statusClass = "connected"
		statusText = "✅ CONNECTED - Ready to send REAL WhatsApp messages"
	} else if whatsappConfig != nil {
		statusClass = "config"
		statusText = "⚙️  CONFIGURED - Set WHATSAPP_TOKEN & WHATSAPP_PHONE_ID in environment variables"
	}

	fmt.Fprintf(w, html, statusClass, statusText)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Only POST method allowed",
		})
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Invalid JSON: " + err.Error(),
		})
		return
	}

	if req.To == "" || req.Message == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Both 'to' and 'message' fields are required",
		})
		return
	}

	log.Printf("📤 REAL WhatsApp message request: To=%s, Message=%s", req.To, req.Message)

	// Send via REAL WhatsApp Business API
	if isConnected {
		result, err := sendWhatsAppMessage(req.To, req.Message)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Failed to send via WhatsApp: " + err.Error(),
				"details": result,
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "REAL WhatsApp message sent successfully",
			"data": map[string]interface{}{
				"to":        req.To,
				"message":   req.Message,
				"timestamp": time.Now().Format(time.RFC3339),
				"api_response": result,
				"integration": "whatsapp_business_api",
			},
		})
	} else {
		// Simulation mode
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "simulation",
			"message": "WhatsApp not configured - running in simulation mode",
			"data": map[string]interface{}{
				"to":        req.To,
				"message":   req.Message,
				"timestamp": time.Now().Format(time.RFC3339),
				"note":      "Set WHATSAPP_TOKEN and WHATSAPP_PHONE_ID environment variables for real messages",
			},
		})
	}
}

func sendWhatsAppMessage(to, message string) (map[string]interface{}, error) {
	// Format phone number (remove any non-digit characters)
	phone := ""
	for _, char := range to {
		if char >= '0' && char <= '9' {
			phone += string(char)
		}
	}

	// Create WhatsApp message payload
	whatsappMsg := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		To:               phone,
	}
	whatsappMsg.Text.Body = message

	// Convert to JSON
	jsonData, err := json.Marshal(whatsappMsg)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	apiURL := fmt.Sprintf("%s/%s/messages", whatsappConfig.APIURL, whatsappConfig.PhoneID)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+whatsappConfig.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	log.Printf("✅ REAL WhatsApp message sent to %s", to)
	return result, nil
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connected":   isConnected,
		"configured":  whatsappConfig != nil,
		"has_token":   whatsappConfig != nil && whatsappConfig.AccessToken != "",
		"has_phone_id": whatsappConfig != nil && whatsappConfig.PhoneID != "",
		"timestamp":   time.Now().Format(time.RFC3339),
		"service":     "whatsapp_business_api",
	})
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	initWhatsAppConfig()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"connected": isConnected,
		"message":   "WhatsApp configuration reloaded",
	})
}

func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	isConnected = false
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "WhatsApp disconnected",
	})
}

func setupHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Webhook setup instructions",
		"steps": []string{
			"1. Go to Facebook Developer Portal",
			"2. Create WhatsApp Business App",
			"3. Get Access Token and Phone Number ID",
			"4. Set environment variables in Render",
			"5. Test with /send endpoint",
		},
		"environment_variables": map[string]string{
			"WHATSAPP_TOKEN":   "Your WhatsApp Business API Access Token",
			"WHATSAPP_PHONE_ID": "Your Phone Number ID from Facebook",
		},
	})
}

func sendTemplateHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "info",
		"message": "Template messages require business approval",
		"note":    "Use regular messages for testing",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "whatsapp_business_api",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "2.0",
	})
}
