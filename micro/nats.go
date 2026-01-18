package micro

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Config struct {
	NatsUrl            string
	NatsServiceName    string
	NatsServiceVersion string
	Timeout            time.Duration
}

type NatsClient interface {
	GetInstance() *natsClient
	Disconnect()
}

type natsClient struct {
	Conn    *nats.Conn
	Service micro.Service
	Timeout time.Duration
}

func (n *natsClient) GetInstance() *natsClient {
	return n
}

func (n *natsClient) Disconnect() {
	fmt.Println("disconnecting nats..")
	n.Conn.Close()
	fmt.Println("disconnected nats")
}

func NewNatsClient(config *Config) NatsClient {
	fmt.Println("connecting to nats..")

	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		panic(err)
	}

	srv, err := micro.AddService(nc, micro.Config{
		Name:    config.NatsServiceName,
		Version: config.NatsServiceVersion,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to nats")

	return &natsClient{
		Conn:    nc,
		Service: srv,
		Timeout: config.Timeout,
	}
}
