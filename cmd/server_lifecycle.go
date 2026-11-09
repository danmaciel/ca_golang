package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

const shutdownTimeout = 10 * time.Second

// listenAndServe binds all addresses before starting any service. A startup
// error cannot leave another service running without its peers.
func listenAndServe(ctx context.Context, restAddr, graphqlAddr, grpcAddr string, restHandler, graphqlHandler http.Handler, grpcServer *grpc.Server) error {
	restListener, err := net.Listen("tcp", restAddr)
	if err != nil {
		return fmt.Errorf("listen REST: %w", err)
	}
	defer restListener.Close()

	graphqlListener, err := net.Listen("tcp", graphqlAddr)
	if err != nil {
		return fmt.Errorf("listen GraphQL: %w", err)
	}
	defer graphqlListener.Close()

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}
	defer grpcListener.Close()

	return serveListeners(ctx, restListener, graphqlListener, grpcListener, restHandler, graphqlHandler, grpcServer)
}

// serveListeners coordinates independent HTTP and gRPC services and returns
// when a stop signal arrives or any server fails.
func serveListeners(ctx context.Context, restListener, graphqlListener, grpcListener net.Listener, restHandler, graphqlHandler http.Handler, grpcServer *grpc.Server) error {
	rest := &http.Server{
		Handler: restHandler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	graphql := &http.Server{
		Handler: graphqlHandler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	failures := make(chan error, 3)
	go func() {
		if err := rest.Serve(restListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			failures <- fmt.Errorf("REST server: %w", err)
		}
	}()
	go func() {
		if err := graphql.Serve(graphqlListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			failures <- fmt.Errorf("GraphQL server: %w", err)
		}
	}()
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			failures <- fmt.Errorf("gRPC server: %w", err)
		}
	}()

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-failures:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	httpDone := make(chan error, 2)
	go func() { httpDone <- rest.Shutdown(shutdownCtx) }()
	go func() { httpDone <- graphql.Shutdown(shutdownCtx) }()
	for i := 0; i < 2; i++ {
		if err := <-httpDone; err != nil {
			runErr = errors.Join(runErr, fmt.Errorf("HTTP shutdown: %w", err))
			_ = rest.Close()
			_ = graphql.Close()
		}
	}

	grpcDone := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcDone)
	}()
	select {
	case <-grpcDone:
	case <-shutdownCtx.Done():
		grpcServer.Stop()
		<-grpcDone
		runErr = errors.Join(runErr, fmt.Errorf("gRPC graceful shutdown timed out: %w", shutdownCtx.Err()))
	}
	return runErr
}
