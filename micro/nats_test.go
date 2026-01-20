package micro

import (
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
)

// RunNatsServerOnPort is a test helper to run a NATS server on a specific port.
func RunNatsServerOnPort(t *testing.T, port int) *server.Server {
	t.Helper()
	opts := test.DefaultTestOptions
	opts.Port = port
	return test.RunServer(&opts)
}

func TestNewNatsClient(t *testing.T) {
	t.Run("should panic with invalid NATS URL", func(t *testing.T) {
		config := &Config{
			NatsUrl:            "nats://localhost:4223", // Invalid/non-running port
			NatsServiceName:    "test-service",
			NatsServiceVersion: "1.0.0",
			Timeout:            2 * time.Second,
		}

		assert.Panics(t, func() {
			NewNatsClient(config)
		})
	})

	t.Run("should connect successfully with valid config", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		config := &Config{
			NatsUrl:            s.ClientURL(),
			NatsServiceName:    "test-service",
			NatsServiceVersion: "1.0.0",
			Timeout:            2 * time.Second,
		}

		client := NewNatsClient(config)
		assert.NotNil(t, client)
		assert.NotNil(t, client.GetInstance().Conn)
		assert.NotNil(t, client.GetInstance().Service)
		assert.Equal(t, config.Timeout, client.GetInstance().Timeout)
		
		// Cleanup
		client.Disconnect()
	})

	t.Run("should set timeout correctly", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		timeout := 5 * time.Second
		config := &Config{
			NatsUrl:            s.ClientURL(),
			NatsServiceName:    "timeout-test-service",
			NatsServiceVersion: "1.0.0",
			Timeout:            timeout,
		}

		client := NewNatsClient(config)
		assert.NotNil(t, client)
		assert.Equal(t, timeout, client.GetInstance().Timeout)
		
		client.Disconnect()
	})
}

func TestNatsClient_GetInstance(t *testing.T) {
	t.Run("should return the same instance", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		config := &Config{
			NatsUrl:            s.ClientURL(),
			NatsServiceName:    "test-service",
			NatsServiceVersion: "1.0.0",
			Timeout:            2 * time.Second,
		}

		client := NewNatsClient(config)
		instance1 := client.GetInstance()
		instance2 := client.GetInstance()

		assert.Equal(t, instance1, instance2)
		assert.NotNil(t, instance1.Conn)
		assert.NotNil(t, instance1.Service)
		
		client.Disconnect()
	})
}

func TestNatsClient_Disconnect(t *testing.T) {
	t.Run("should disconnect and close connection", func(t *testing.T) {
		s := RunNatsServerOnPort(t, -1)
		defer s.Shutdown()

		config := &Config{
			NatsUrl:            s.ClientURL(),
			NatsServiceName:    "disconnect-test-service",
			NatsServiceVersion: "1.0.0",
			Timeout:            2 * time.Second,
		}

		client := NewNatsClient(config)
		instance := client.GetInstance()
		
		assert.False(t, instance.Conn.IsClosed())
		
		client.Disconnect()
		
		assert.True(t, instance.Conn.IsClosed())
	})
}
