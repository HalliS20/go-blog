package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"go-blog/config"
	"go-blog/internal/router"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/quic-go/quic-go/http3"
)

var e *gin.Engine

func main() {
	// Initialize the database

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	initializeServer()

	router.Init(e)

	cert, err := generateDevCert()
	if err != nil {
		log.Fatalf("Failed to generate development certificate: %v", err)
	}

	// Configure TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3", "h2", "http/1.1"},
	}

	// HTTP/3 Server
	http3Server := &http3.Server{
		Addr:      ":443",
		Handler:   e,
		TLSConfig: tlsConfig,
	}

	// Start HTTP/3 server
	go func() {
		log.Println("Starting HTTP/3 server on :443")
		if err := http3Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP/3 server error: %v", err)
		}
	}()

	//======== Shutdown the all layers when the server is closed
	defer router.Shutdown()
	//======== Run the server
	err = e.Run(":8080")
	if err != nil {
		return
	}
}

func initializeServer() {
	e = gin.Default()
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // use gzip for text compression
	e.LoadHTMLGlob("templates/*")
	e.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/public/miniStyles/total.min.css" || c.Request.URL.Path == "/public/scripts/main.min.js" {
			c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
		}
		c.Next()
	})
}

func generateDevCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Development Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180), // Valid for 180 days
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}
