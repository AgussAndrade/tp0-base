package common

import (
	"context"
	"encoding/csv"
	"os"
	"io"
)

type BatchStream struct {
	reader       *csv.Reader
	clientID     string
	file         *os.File
	maxAmount    int
	maxBytes     int
}

func NewBatchStream(path string, clientID string, maxAmount int, maxBytes int) (*BatchStream, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &BatchStream{
		reader:    csv.NewReader(file),
		clientID:  clientID,
		file:      file,
		maxAmount: maxAmount,
		maxBytes:  maxBytes,
	}, nil
}

func (bs *BatchStream) Close() {
	bs.file.Close()
}

func (bs *BatchStream) NextBatch(ctx context.Context) ([]Bet, error) {
	batch := make([]Bet, 0, bs.maxAmount)
	currentSize := 0

	for len(batch) < bs.maxAmount && currentSize <= bs.maxBytes {
		if ctx.Err() != nil {
			break
		}

		record, err := bs.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		bet, valid := ParseBet(record, bs.clientID)
		if !valid {
			continue
		}

		msg := formatBetMessage(bet)
		msgLen := len(msg)

		if msgLen > bs.maxBytes {
			log.Warningf("action: bet_too_large | result: fail | size: %d", msgLen)
			continue
		}
		batch = append(batch, bet)
		currentSize += msgLen
		// add tolerance to avoid sending max len than default.
		if currentSize > bs.maxBytes-msgLen {
			break
		}
	}

	if len(batch) > 0 {
		return batch, nil
	}
	return nil, nil
}


