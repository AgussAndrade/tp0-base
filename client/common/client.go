package common

import (
	"bufio"
	"net"
	"time"
	"context"
	"os"
	"strings"
	"os/signal"
	"syscall"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	bets := getBets(c.config.ID)
	go func() {
		<-sigs
		cancel()
	}()

	for _, bet := range bets {
		select {
		case <-ctx.Done():
			log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
			c.conn.Close()
			return	
		default:
			// Create the connection the server in every loop iteration. Send an
			c.createClientSocket()

			msg := formatBetMessage(bet)

			totalSent := c.SendMessage(msg)

			if totalSent != len(msg) {
				c.conn.Close()
				return
			}
			response := c.ReceiveMessage()
			if response == "" {
				log.Errorf("action: receive_response | result: fail | client_id: %v | error: Not Ok msg",
					c.config.ID,
				)
			}
			c.conn.Close()

			if !isOkMsg(response) {
				log.Errorf("action: receive_response | result: fail | client_id: %v | error: Not Ok msg",
					c.config.ID,
				)
				return
			}

			log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s",
				bet.Document,
				bet.Number,
			)

			// Wait a time between sending one message and the next one
			time.Sleep(c.config.LoopPeriod)
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) SendMessage(msg string) int {
	totalSent := 0
	for totalSent < len(msg) {
		n, err := c.conn.Write([]byte(msg)[totalSent:])
		if err != nil {
			log.Errorf("action: send_bet | result: fail | error: %v", err)
			return totalSent
		}
		totalSent += n
	}
	return totalSent
}

func (c *Client) ReceiveMessage() string {

	var response strings.Builder
	reader := bufio.NewReader(c.conn)
	
	for {
		chunk, err := reader.ReadByte()
		if err != nil {
			return ""
		}
		response.WriteByte(chunk)
		if isEndOfMsg(chunk) {
			break
		}
	}
	return response.String()
}