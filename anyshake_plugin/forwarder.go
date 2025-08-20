package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/samber/lo"
)

type ForwarderMessage struct {
	RawMessage string

	Time       time.Time
	SampleRate int

	Index    int
	Network  string
	Station  string
	Channel  string
	Location string

	Data     []int32
	Checksum uint8
}

func NewForwarderMessage(msg string) (ForwarderMessage, error) {
	m := ForwarderMessage{RawMessage: msg}
	if err := m.Parse(); err != nil {
		return ForwarderMessage{}, fmt.Errorf("failed to parse message: %w", err)
	}
	return m, nil
}

func (m *ForwarderMessage) Parse() error {
	fields := strings.Split(m.RawMessage, ",")
	if len(fields) < 8 {
		return errors.New("message fields length is less than 8")
	}

	// Parse index without leading '$'
	index, err := strconv.Atoi(fields[0][1:])
	if err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}
	m.Index = index

	m.Network = fields[1]
	m.Station = fields[2]
	m.Location = fields[3]
	m.Channel = fields[4]

	timestamp, err := strconv.ParseInt(fields[5], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}
	m.Time = time.UnixMilli(timestamp).UTC()

	m.SampleRate, err = strconv.Atoi(fields[6])
	if err != nil {
		return fmt.Errorf("failed to parse sample rate: %w", err)
	}

	// Extract and validate data
	m.Data = m.extractDataArr(fields)
	calcChecksum := m.calculateChecksum()
	recvChecksum, err := m.extractChecksum()
	if err != nil {
		return fmt.Errorf("failed to extract checksum: %w", err)
	}

	if calcChecksum != recvChecksum {
		return errors.New("checksum mismatch")
	}

	m.Checksum = recvChecksum
	return nil
}

func (m *ForwarderMessage) extractDataArr(fields []string) []int32 {
	return lo.Map(fields[7:len(fields)-1], func(field string, _ int) int32 {
		data, _ := strconv.ParseInt(field, 10, 32)
		return int32(data)
	})
}

func (m *ForwarderMessage) calculateChecksum() uint8 {
	var checksum uint8
	for _, data := range m.Data {
		bytes := (*[4]uint8)(unsafe.Pointer(&data))[:]
		for _, b := range bytes {
			checksum ^= b
		}
	}
	return checksum
}

func (m *ForwarderMessage) extractChecksum() (uint8, error) {
	idx := strings.LastIndex(m.RawMessage, "*")
	if idx == -1 || idx+2 >= len(m.RawMessage) {
		return 0, errors.New("checksum not found")
	}

	checksum, err := strconv.ParseUint(m.RawMessage[idx+1:idx+3], 16, 8)
	if err != nil {
		return 0, fmt.Errorf("failed to parse checksum: %w", err)
	}

	return uint8(checksum), nil
}
