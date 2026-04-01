package model

import (
	"net"
	"net/url"
	"time"
)

// HealthStatus represents the health state of a registered service.
type HealthStatus string

const (
	StatusUnknown   HealthStatus = "unknown"
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// Service is the core domain object stored in the registry.
type Service struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	URL           string       `json:"url"`
	DisplayURL    string       `json:"display_url"`
	ExternalURL   string       `json:"external_url,omitempty"`
	Description   string       `json:"description,omitempty"`
	RemoteIP      string       `json:"remote_ip"`
	Status        HealthStatus `json:"status"`
	RegisteredAt  time.Time    `json:"registered_at"`
	LastCheckedAt *time.Time   `json:"last_checked_at,omitempty"`
}

// RegisterRequest is the payload for registering a new service.
type RegisterRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// isLoopback returns true if ip is a loopback address (127.x.x.x, ::1).
func isLoopback(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ip == "localhost"
	}
	return parsed.IsLoopback()
}

// urlHasLoopbackHost returns true if the URL's host part is localhost / 127.x / ::1.
func urlHasLoopbackHost(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := u.Hostname()
	return isLoopback(host)
}

// replaceURLHost returns a new URL with the host replaced by newHost,
// preserving the original port, path, and query.
func replaceURLHost(rawURL, newHost string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	port := u.Port()
	if port != "" {
		u.Host = net.JoinHostPort(newHost, port)
	} else {
		u.Host = newHost
	}
	return u.String()
}

// ResolveDisplayURLs computes DisplayURL and ExternalURL based on the
// registered URL and the caller's remote IP.
//
// Rules:
//   - Remote IP is loopback + URL host is loopback:
//     display_url = original URL (works on this machine)
//     external_url = URL with host replaced by localIP (for other devices)
//   - Remote IP is external + URL host is loopback:
//     display_url = URL with host replaced by remote IP
//     external_url = empty
//   - Otherwise:
//     display_url = original URL
//     external_url = empty
func (s *Service) ResolveDisplayURLs(localIP string) {
	// If remote IP is unknown (legacy data), just show the original URL
	// and generate an external URL if the URL host is loopback.
	if s.RemoteIP == "" {
		s.DisplayURL = s.URL
		if urlHasLoopbackHost(s.URL) && localIP != "" {
			s.ExternalURL = replaceURLHost(s.URL, localIP)
		}
		return
	}

	remoteIsLoopback := isLoopback(s.RemoteIP)
	urlIsLoopback := urlHasLoopbackHost(s.URL)

	switch {
	case remoteIsLoopback && urlIsLoopback:
		s.DisplayURL = s.URL
		if localIP != "" {
			s.ExternalURL = replaceURLHost(s.URL, localIP)
		}
	case !remoteIsLoopback && urlIsLoopback:
		s.DisplayURL = replaceURLHost(s.URL, s.RemoteIP)
		s.ExternalURL = ""
	default:
		s.DisplayURL = s.URL
		s.ExternalURL = ""
	}
}
