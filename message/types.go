package message

import "time"

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
