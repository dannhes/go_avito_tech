package e2e

import (
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := waitForPort("localhost:8081", 30*time.Second); err != nil {
		log.Fatalf("server not ready: %v", err)
	}
	os.Exit(m.Run())
}

func waitForPort(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return &ServerNotReadyError{}
}

type ServerNotReadyError struct{}

func (e *ServerNotReadyError) Error() string {
	return "server did not become ready in time"
}
