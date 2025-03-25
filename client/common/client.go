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
	"fmt"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchAmount   int
	BatchMaxBytes int
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

func (c *Client) StartClientLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	err := c.createClientSocket()
	if err != nil {
		log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}
	defer func() {
		c.conn.Close()
		log.Infof("action: close_socket | result: success | client_id: %v", c.config.ID)
		//sleep based on https://campusgrado.fi.uba.ar/mod/forum/discuss.php?d=29739#p52493
		time.Sleep(200 * time.Millisecond)
	}()

	err = c.SendBatchesFromCSV(ctx, "/data/agency.csv")
	if err != nil {
		log.Errorf("action: batch_loop | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	select {
	case <-ctx.Done():
		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
	default:
		log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	}
	_, err = c.conn.Write([]byte{'\t'})
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

func (c *Client) SendBatchesFromCSV(ctx context.Context, path string) error {
	stream, err := NewBatchStream(path, c.config.ID, c.config.BatchAmount, c.config.BatchMaxBytes)
	if err != nil {
		return fmt.Errorf("failed to create batch stream: %w", err)
	}
	defer stream.Close()

	for {
		batch, err := stream.NextBatch(ctx)
		if err != nil {
			return err
		}
		if batch == nil {
			break
		}
		if err := c.retryBatchUntilSuccess(batch, 3); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) retryBatchUntilSuccess(batch []Bet, maxRetries int) error {
	msg := formatBatchMessageWithEnd(batch)

	for i := 0; i < maxRetries; i++ {
		totalSent := c.SendMessage(msg)
		if totalSent != len(msg) {
			log.Errorf("action: send_batch | result: fail | error: short_write")
			return fmt.Errorf("short write")
		}

		response := c.ReceiveMessage()
		// sleep period between messages either fail or success
		time.Sleep(c.config.LoopPeriod)
		if isOkMsg(response) {
			log.Infof("action: send_batch | result: success | cantidad: %d", len(batch))
			return nil
		}

		log.Warningf("action: batch_response | result: fail | attempt: %d/%d", i+1, maxRetries)
		
	}

	log.Warning("action: batch_response | result: fail | reason: max_retries_exceeded")
	return fmt.Errorf("max retries exceeded")
}
