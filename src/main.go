package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	server := "example.com:443"

	// Dial the server
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", server, &tls.Config{
		InsecureSkipVerify: true, // Temporarily skip verification; we'll handle it manually.
	})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Retrieve the certificate chain
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		log.Fatalf("No certificates presented by the server")
	}

	// Extract the leaf certificate
	leafCert := certs[0]

	// Display expiration date of the leaf certificate
	fmt.Printf("TLS Certificate is valid until: %s\n", leafCert.NotAfter)

	// Check if the certificate is expired
	if time.Now().After(leafCert.NotAfter) {
		log.Fatalf("TLS Certificate has expired on: %s", leafCert.NotAfter)
	}

	// Load system root CAs
	roots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Failed to load system root CAs: %v", err)
	}
	if roots == nil {
		log.Fatal("No system root CAs available")
	}

	// Create a verification options structure
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
	}

	// Add all but the leaf certificate to intermediates pool
	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}

	// Verify the leaf certificate
	if _, err := leafCert.Verify(opts); err != nil {
		log.Fatalf("Certificate verification failed: %v", err)
	}

	fmt.Println("Certificate verification successful!")
}
