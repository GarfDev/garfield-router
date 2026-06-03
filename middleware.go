package main

import (
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var apiToken = os.Getenv("ROUTER_API_TOKEN")

// extractHeaders pulls Garfield-specific routing metadata from request headers.
func extractHeaders(r *http.Request) RouteRequest {
	tier := 0
	if t := r.Header.Get("X-Garfield-Tier"); t != "" {
		fmt.Sscanf(t, "%d", &tier)
	}

	priority := r.Header.Get("X-Garfield-Priority")
	if priority == "" {
		priority = "normal"
	}

	return RouteRequest{
		Service:        r.Header.Get("X-Garfield-Service"),
		CallType:       r.Header.Get("X-Garfield-CallType"),
		Priority:       priority,
		Tier:           tier,
		PersonaID:      r.Header.Get("X-Garfield-PersonaID"),
		ResponseSchema: r.Header.Get("X-Garfield-Response-Schema"),
		Reflect:        isTruthyHeader(r.Header.Get("X-Garfield-Reflect")),
		Consensus:      isTruthyHeader(r.Header.Get("X-Garfield-Consensus")),
	}
}

// isTruthyHeader treats "1"/"true"/"yes" (case-insensitive) as enabled.
func isTruthyHeader(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

// loggingMiddleware logs each request with method, path, status, and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		duration := time.Since(start)

		// Don't log health checks to reduce noise
		if r.URL.Path == "/health" {
			return
		}

		logger.Printf("%s %s %d %s [%s %s]",
			r.Method, r.URL.Path, sw.status, duration.Round(time.Millisecond),
			r.Header.Get("X-Garfield-Service"),
			r.Header.Get("X-Garfield-CallType"),
		)
	})
}

// corsMiddleware adds CORS headers for API access.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Garfield-Service, X-Garfield-CallType, X-Garfield-Priority, X-Garfield-Tier, X-Garfield-PersonaID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Powered-By, X-Garfield-Router-Version, X-Garfield-Backend, X-Garfield-Rule, X-Garfield-Request-Cost")

		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// statusWriter wraps ResponseWriter to capture the status code.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

// ValidateExternalURL checks a URL is not targeting internal/private networks (SSRF prevention).
func ValidateExternalURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got %s", parsed.Scheme)
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL has no host")
	}
	// Block metadata endpoints
	if host == "169.254.169.254" || host == "metadata.google.internal" {
		return fmt.Errorf("URL targets cloud metadata endpoint")
	}
	// Block localhost
	if host == "localhost" || host == "127.0.0.1" || host == "::1" || host == "0.0.0.0" {
		return fmt.Errorf("URL targets localhost")
	}
	// Block private networks
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("URL targets private/internal network")
		}
	}
	return nil
}

// Flush implements http.Flusher so streaming SSE works through the middleware.
func (sw *statusWriter) Flush() {
	if f, ok := sw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// authMiddleware protects /api/* endpoints with a bearer token.
// If ROUTER_API_TOKEN is unset, all requests are allowed (open access).
// The UI (/) and proxy endpoint (/v1/) are never gated.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for non-API paths (UI, health, proxy)
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		// No token configured: open access
		if apiToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Check Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			writeErrorJSON(w, 401, "authentication required: set Authorization: Bearer <token>")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if subtle.ConstantTimeCompare([]byte(token), []byte(apiToken)) != 1 {
			writeErrorJSON(w, 403, "invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
