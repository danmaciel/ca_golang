package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestListenAndServeReturnsBindErrorBeforeStartingPeers(t *testing.T) {
	busy, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer busy.Close()

	restHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("REST must not start if another port cannot bind")
	})
	err = listenAndServe(context.Background(), "127.0.0.1:0",
		busy.Addr().String(), "127.0.0.1:0",
		restHandler, http.NotFoundHandler(), grpc.NewServer())
	if err == nil {
		t.Fatal("expected an error when the GraphQL port is occupied")
	}
}

func TestServersFinishInflightHTTPRequestOnCancellation(t *testing.T) {
	restListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer restListener.Close()
	graphqlListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer graphqlListener.Close()
	grpcListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer grpcListener.Close()

	started := make(chan struct{})
	release := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan error, 1)

	restHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "done")
	})
	go func() {
		finished <- serveListeners(ctx, restListener, graphqlListener, grpcListener,
			restHandler, http.NotFoundHandler(), grpc.NewServer())
	}()

	client := &http.Client{Timeout: 5 * time.Second}
	responseResult := make(chan error, 1)
	go func() {
		res, err := client.Get("http://" + restListener.Addr().String())
		if err != nil {
			responseResult <- err
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			responseResult <- err
			return
		}
		if string(body) != "done" || res.StatusCode != http.StatusOK {
			responseResult <- errors.New("request did not finish normally")
			return
		}
		responseResult <- nil
	}()

	select {
	case <-started:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("REST request did not reach the handler")
	}
	cancel()
	close(release)

	select {
	case err := <-responseResult:
		if err != nil {
			t.Fatalf("in-flight request failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("in-flight request stalled")
	}
	select {
	case err := <-finished:
		if err != nil {
			t.Fatalf("graceful shutdown failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("coordinated servers did not stop")
	}
}
