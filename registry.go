package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/distribution/distribution/v3/configuration"
	"github.com/distribution/distribution/v3/registry"
)

func startRegistry(ctx context.Context, backend StorageBackend, bucket string) (string, error) {
	port, err := findFreePort()
	if err != nil {
		return "", err
	}
	regAddr := fmt.Sprintf("localhost:%d", port)

	if err := backend.ValidateConfig(); err != nil {
		return "", err
	}

	storageDriverConfig := configuration.Storage{}
	storageDriverConfig[backend.Type()] = backend.GetStorageConfig(bucket)

	reg, err := registry.NewRegistry(ctx, &configuration.Configuration{
		Storage: storageDriverConfig,
		HTTP:    configuration.HTTP{Addr: regAddr},
		Log:     configuration.Log{Level: configuration.Loglevel("fatal"), AccessLog: configuration.AccessLog{Disabled: true}},
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
			slog.Error("Failed starting server", "error", err)
		}
	}()
	return regAddr, nil
}

// findFreePort listens on a random available TCP port (by specifying :0)
func findFreePort() (int, error) {
	// Listen on TCP port 0. The operating system will assign a free, ephemeral port.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to listen on a free port: %w", err)
	}
	defer func() { _ = listener.Close() }()

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("listener address is not a TCP address: %T", listener.Addr())
	}
	return tcpAddr.Port, nil
}
