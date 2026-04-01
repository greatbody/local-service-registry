package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"net/http"
	"strings"
)

// generateID produces a short random hex ID (8 bytes = 16 hex chars).
func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// extractIP returns the IP address of the HTTP client.
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For / X-Real-IP first (reverse proxy).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// detectLocalIP returns the first non-loopback IPv4 address of this machine.
func detectLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("warning: cannot detect local IP: %v", err)
		return ""
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok &&
			!ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}
