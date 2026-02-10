package main

import (
	"fmt"
	"net"
	"testing"
)

func TestFindFreePort(t *testing.T) {
	// Test that findFreePort returns a valid port
	port, err := findFreePort()
	if err != nil {
		t.Fatalf("findFreePort() error = %v", err)
	}

	if port <= 0 || port > 65535 {
		t.Errorf("findFreePort() = %d, want valid port range (1-65535)", port)
	}

	// Test that the port is actually free
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Errorf("Port %d returned by findFreePort is not actually free: %v", port, err)
	} else {
		listener.Close()
	}
}

func TestStartRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping registry test in short mode")
	}

	t.Skip("Registry test requires Docker and complex setup - skipping for now")
}

func TestStartRegistryCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping registry cancellation test in short mode")
	}

	t.Skip("Registry cancellation test requires complex setup - skipping for now")
}
