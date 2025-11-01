package helper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// GetClientIP extracts the client's IP address from the request.
func GetClientIP(ctx context.Context) string {
	req, ok := ctx.Value("request").(*http.Request)
	if !ok || req == nil {
		return ""
	}

	// Check for proxy headers (X-Forwarded-For, X-Real-IP)
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // First IP in list
	}

	xrip := req.Header.Get("X-Real-IP")
	if xrip != "" {
		return xrip
	}

	// Fallback to remote address
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}

	return ip
}

// GetUserAgent extracts the User-Agent from the request.
func GetUserAgent(ctx context.Context) string {
	req, ok := ctx.Value("request").(*http.Request)
	if !ok || req == nil {
		return ""
	}
	return req.Header.Get("User-Agent")
}

// Create HTTPS Client
func CreateHTTPSClient(certFile, keyFile, caCertFile string) *http.Client {
	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
		return nil
	}

	// Load system CA pool
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool() // Fallback to empty pool
	}

	// Load CA certificate if provided
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			log.Fatalf("Failed to read CA certificate: %v", err)
			return nil
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	// Create and return HTTPS client
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}
