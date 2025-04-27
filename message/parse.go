package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

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
