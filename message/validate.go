package message

import (
	"errors"
	"strconv"
	"strings"
	"unsafe"
)

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
