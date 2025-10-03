package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/savageking-io/ogbuser/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Client struct {
	hostname       string
	port           uint16
	ErrorChan      chan error
	conn           *grpc.ClientConn
	client         proto.UserServiceClient
	maxFailedPings int
	maxReconnects  int
}

func NewClient(hostname string, port uint16) *Client {
	return &Client{
		hostname: hostname,
		port:     port,
	}
}

func (c *Client) Run() error {
	if c.hostname == "" {
		return fmt.Errorf("hostname is empty")
	}
	if c.port == 0 {
		return fmt.Errorf("port is empty")
	}
	c.ErrorChan = make(chan error)

	c.maxReconnects = 5
	c.maxFailedPings = 5

	if err := c.Connect(); err != nil {
		log.Errorf("Failed to connect to server: %v", err)
		c.Reconnect()
	}

	lastPing := time.Unix(0, 0)
	failedPings := 0
	for {
		if failedPings >= c.maxFailedPings {
			// Too many pings failed - reconnect
			c.Reconnect()
			continue
		}
		if time.Since(lastPing) > time.Second*5 {
			if err := c.ping(); err != nil {
				log.Errorf("Ping to user microservice failed: %s", err.Error())
				if errors.Is(err, context.DeadlineExceeded) {
					lastPing = time.Unix(0, 0)
					failedPings++
					continue // Send another ping immediately
				}
				if errors.Is(err, grpc.ErrClientConnClosing) {
					// Server disconnected - reconnect
					c.Reconnect()
				}
			}
			lastPing = time.Now()
			failedPings = 0
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (c *Client) Reconnect() {
	if err := c.Disconnect(); err != nil {
		log.Errorf("Failed to close connection: %s", err.Error())
	}
	c.ConnectWithRetry()
}

func (c *Client) Disconnect() error {
	return c.conn.Close()
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = grpc.NewClient(fmt.Sprintf("%s:%d", c.hostname, c.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.client = proto.NewUserServiceClient(c.conn)
	return nil
}

func (c *Client) ConnectWithRetry() {
	timeout := time.Unix(0, 0)
	reconnectsNum := 0
	for {
		if reconnectsNum >= c.maxReconnects {
			log.Errorf("Failed to connect to server after %d reconnects", c.maxReconnects)
			return
		}
		if time.Since(timeout) > time.Millisecond*1000 {
			timeout = time.Now()
			if err := c.Connect(); err != nil {
				reconnectsNum++
				log.Errorf("Failed to connect: %s", err.Error())
			} else {
				return
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (c *Client) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.client.Ping(ctx, &proto.PingMessage{SentAt: timestamppb.New(time.Now())})
	if err != nil {
		return err
	}
	sentAt := resp.SentAt.AsTime()
	repliedAt := resp.RepliedAt.AsTime()
	diff := repliedAt.Sub(sentAt)
	log.Tracef("Ping to user microservice replied in %s", diff.String())
	return nil
}

// ValidateToken is a helper function for easy token validation
func (c *Client) ValidateToken(ctx context.Context, token string) (bool, int32, error) {
	log.Traceln("User::Client::ValidateToken")
	if c.conn == nil {
		return false, -1, fmt.Errorf("connection is not initialized")
	}
	if c.client == nil {
		return false, -1, fmt.Errorf("client is not initialized")
	}

	result, err := c.client.ValidateToken(ctx, &proto.ValidateTokenRequest{Token: token})
	if err != nil {
		return false, -1, err
	}
	if result.Code != 0 {
		return false, -1, fmt.Errorf("validation failed with code %d: %s", result.Code, result.Error)
	}
	return result.IsValid, result.UserId, nil
}
