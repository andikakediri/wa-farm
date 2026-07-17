package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
)

var (
	currentQRPNG []byte
	currentCode  string
	mu           sync.Mutex
	sessions     []string
	startTime    = time.Now()
)

func main() {
	// Start whatsmeow in background
	go startWhatsmeow()

	// Routes
	http.HandleFunc("/", serveLanding)
	http.HandleFunc("/qr", getQR)
	http.HandleFunc("/status", getStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startWhatsmeow() {
	// Check if wabot binary exists
	wabotPath := "./wabot"
	if _, err := os.Stat(wabotPath); os.IsNotExist(err) {
		log.Printf("wabot not found at %s, waiting for build...", wabotPath)
		// Try alternate path
		altPath := "/app/wabot"
		if _, err := os.Stat(altPath); err == nil {
			wabotPath = altPath
		} else {
			log.Println("wabot binary not available. Using simulated mode.")
			simulateMode()
			return
		}
	}

	cmd := exec.Command(wabotPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		simulateMode()
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error creating stderr pipe: %v", err)
		simulateMode()
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting wabot: %v", err)
		simulateMode()
		return
	}

	// Read stdout
	buf := make([]byte, 4096)
	go func() {
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			line := string(buf[:n])
			log.Printf("[wabot] %s", line)

			if strings.Contains(line, "QR code:") || strings.Contains(line, "QR:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					qrData := strings.TrimSpace(parts[1])
					png, _ := qrcode.Encode(qrData, qrcode.Medium, 256)
					mu.Lock()
					currentQRPNG = png
					mu.Unlock()
				}
			}
			if strings.Contains(line, "Pairing code") || strings.Contains(line, "code:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					mu.Lock()
					currentCode = strings.TrimSpace(parts[1])
					mu.Unlock()
				}
			}
			if strings.Contains(line, "Logged in") || strings.Contains(line, "success") {
				mu.Lock()
				sessions = append(sessions, line)
				mu.Unlock()
				log.Printf("✅ NEW SESSION CAPTURED: %s", line)
			}
		}
	}()

	// Read stderr
	go func() {
		buf2 := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf2)
			if err != nil {
				break
			}
			log.Printf("[wabot-err] %s", string(buf2[:n]))
		}
	}()

	cmd.Wait()
}

func simulateMode() {
	log.Println("Running in simulation mode - generating test QR codes")
	// Generate QR codes periodically
	go func() {
		for {
			time.Sleep(25 * time.Second)
			qrContent := fmt.Sprintf("whatsapp-link-sim-%d", time.Now().Unix())
			png, _ := qrcode.Encode(qrContent, qrcode.Medium, 256)
			mu.Lock()
			currentQRPNG = png
			currentCode = fmt.Sprintf("%04d-%04d",
				time.Now().Unix()%10000,
				(time.Now().Unix()/1000)%10000)
			mu.Unlock()
			log.Printf("New QR + pairing code generated: %s", currentCode)
		}
	}()
}

func serveLanding(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Landing page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func getQR(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	resp := map[string]interface{}{
		"qr":   "",
		"code": currentCode,
		"time": time.Now().Unix(),
	}

	if len(currentQRPNG) > 0 {
		resp["qr"] = base64.StdEncoding.EncodeToString(currentQRPNG)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	resp := map[string]interface{}{
		"uptime":      time.Since(startTime).String(),
		"sessions":    len(sessions),
		"session_ids": sessions,
		"has_code":    currentCode != "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
