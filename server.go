package main

import (
	"bufio"
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
	useRealMode  = false
)

func main() {
	if _, err := os.Stat("./wabot"); err == nil {
		useRealMode = true
		log.Println("✅ REAL MODE - wabot found")
		go startRealWabot()
	} else {
		log.Println("🔷 SIMULATION MODE - wabot not found")
		go simulateMode()
	}

	http.HandleFunc("/", serveLanding)
	http.HandleFunc("/qr", getQR)
	http.HandleFunc("/status", getStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s [%s]", port, map[bool]string{true: "REAL", false: "SIMULATION"}[useRealMode])
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func startRealWabot() {
	cmd := exec.Command("./wabot")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("wabot pipe error: %v", err)
		simulateMode()
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("wabot start error: %v", err)
		simulateMode()
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("[wabot] %s", line)

		switch {
		case strings.HasPrefix(line, "QR:"):
			b64 := strings.TrimPrefix(line, "QR:")
			png, err := base64.StdEncoding.DecodeString(b64)
			if err == nil {
				mu.Lock()
				currentQRPNG = png
				mu.Unlock()
			}
		case strings.HasPrefix(line, "CODE:"):
			mu.Lock()
			currentCode = strings.TrimPrefix(line, "CODE:")
			mu.Unlock()
		case strings.HasPrefix(line, "LOGGEDIN:"):
			jid := strings.TrimPrefix(line, "LOGGEDIN:")
			mu.Lock()
			sessions = append(sessions, jid)
			mu.Unlock()
			log.Printf("✅ SESSION: %s", jid)
		case strings.HasPrefix(line, "MSG:"):
			parts := strings.SplitN(strings.TrimPrefix(line, "MSG:"), "|", 3)
			if len(parts) == 3 {
				log.Printf("📩 %s -> %s: %s", parts[0], parts[1], parts[2])
			}
		}
	}
	cmd.Wait()
	log.Println("wabot exited → fallback to simulation")
	simulateMode()
}

func simulateMode() {
	for {
		time.Sleep(25 * time.Second)
		qrContent := fmt.Sprintf("whatsapp-link-sim-%d", time.Now().Unix())
		png, _ := qrcode.Encode(qrContent, qrcode.Medium, 256)
		mu.Lock()
		currentQRPNG = png
		currentCode = fmt.Sprintf("%04d-%04d", time.Now().Unix()%10000, (time.Now().Unix()/1000)%10000)
		mu.Unlock()
	}
}

func serveLanding(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func getQR(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	resp := map[string]interface{}{
		"qr":        "",
		"code":      currentCode,
		"time":      time.Now().Unix(),
		"real_mode": useRealMode,
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uptime":      time.Since(startTime).String(),
		"sessions":    len(sessions),
		"session_ids": sessions,
		"has_code":    currentCode != "",
		"mode":        map[bool]string{true: "REAL", false: "SIM"}[useRealMode],
	})
}
