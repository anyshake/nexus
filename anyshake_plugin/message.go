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

type Message struct {
	RawMessage string

	Time       time.Time
	SampleRate int

	Network  string
	Station  string
	Channel  string
	Location string

	Data     []int
	Checksum uint8
}

func NewMessage(msg string) (Message, error) {
	messageObj := Message{RawMessage: msg}
	if err := messageObj.Parse(); err != nil {
		return Message{}, err
	}

	return messageObj, nil
}

func (m *Message) Parse() error {
	if err := m.Validate(); err != nil {
		return fmt.Errorf("failed to validate message: %w", err)
	}

	fields := strings.Split(m.RawMessage, ",")

	timestamp, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}
	m.SampleRate, err = strconv.Atoi(fields[5])
	if err != nil {
		return fmt.Errorf("failed to parse sample rate: %w", err)
	}

	m.Network = fields[0][1:] // without $ flag
	m.Station = fields[1]
	m.Location = fields[2]
	m.Channel = fields[3]
	m.Time = time.UnixMilli(timestamp).UTC()
	m.Data = lo.Map(fields[6:len(fields)-2], func(field string, _ int) int {
		data, _ := strconv.ParseInt(field, 10, 32)
		return int(data)
	})

	return nil
}

func (m *Message) Validate() (err error) {
	fields := strings.Split(m.RawMessage, ",")

	// Minimum message fields length is 7 (only 1 sample)
	if len(fields) < 7 {
		return errors.New("message fields length is less than 7")
	}

	// Convert data fields to int32
	var dataArr []int32
	for _, field := range fields[6 : len(fields)-1] {
		data, err := strconv.Atoi(field)
		if err != nil {
			return err
		}
		dataArr = append(dataArr, int32(data))
	}

	// Get message checksum by XOR operation
	var calcChecksum uint8
	for _, data := range dataArr {
		bytes := (*[4]uint8)(unsafe.Pointer(&data))[:]
		for j := 0; j < int(unsafe.Sizeof(int32(0))); j++ {
			calcChecksum ^= bytes[j]
		}
	}

	recvChecksum, err := m.getChecksum()
	if err != nil {
		return err
	}
	if calcChecksum != recvChecksum {
		return errors.New("checksum is not equal")
	}

	m.Checksum = recvChecksum
	return nil
}

func (m *Message) getChecksum() (uint8, error) {
	for idx := len(m.RawMessage) - 1; idx >= 0; idx-- {
		ch := m.RawMessage[idx]
		if ch == '*' {
			checksum, err := strconv.ParseUint(m.RawMessage[idx+1:idx+3], 16, 8)
			if err != nil {
				return 0, err
			}
			return uint8(checksum), nil
		}
	}

	return 0, errors.New("checksum not found")
}
