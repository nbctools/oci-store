package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/distribution/distribution/v3/configuration"
	"github.com/distribution/distribution/v3/registry"
)

func startRegistry(ctx context.Context, bucket string) (string, error) {
	port, err := findFreePort()
	if err != nil {
		return "", err
	}
	regAddr := fmt.Sprintf("localhost:%d", port)
	loglevel := "error"
	if verbose {
		loglevel = "info"
	}
	reg, err := registry.NewRegistry(ctx, &configuration.Configuration{
		Storage: configuration.Storage{"s3": {"bucket": bucket,
			"region":         region,
			"regionendpoint": endpoint,
			"loglevel":       loglevel}},
		HTTP: configuration.HTTP{Addr: regAddr},
		Log:  configuration.Log{Level: configuration.Loglevel(loglevel), AccessLog: configuration.AccessLog{Disabled: true}},
	})
	if err != nil {
		return "", err
	}

	go func() {
		<-ctx.Done()
		slog.Debug("Stopping registry", "addr", regAddr)
		err := reg.Shutdown(ctx)
		if err != nil {
			slog.Error("Failed stopping server", "error", err)
		}
	}()

	go func() {
		slog.Debug("Starting registry", "addr", regAddr)
		err := reg.ListenAndServe()
		if err != nil {
			// slog.Error("Failed starting server", "error", err)
		}
	}()
	return regAddr, nil
}

// findFreePort listens on a random available TCP port (by specifying :0)
// and returns the assigned port number. It then closes the listener.
func findFreePort() (int, error) {
	// Listen on TCP port 0. The operating system will assign a free, ephemeral port.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to listen on a free port: %w", err)
	}
	defer listener.Close() // Ensure the listener is closed after we've found the port

	// Get the actual assigned address and extract the port number.
	// We expect a *net.TCPAddr from a "tcp" listener.
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("listener address is not a TCP address: %T", listener.Addr())
	}

	return tcpAddr.Port, nil
}
