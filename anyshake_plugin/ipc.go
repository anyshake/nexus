package main

import (
	"encoding/binary"
	"os"
	"time"
)

type SeedLinkIPC struct {
	fd *os.File
}

func NewSeedlinkIPC() *SeedLinkIPC {
	return &SeedLinkIPC{fd: os.NewFile(63, "seedlink")}
}

func (s *SeedLinkIPC) Close() {
	s.fd.Close()
}

func (s *SeedLinkIPC) SendRaw3(station, channel string, t time.Time, usecCorr, timingQuality int, data []int32) error {
	header := make([]byte, 4+10+10+9*4)
	binary.LittleEndian.PutUint32(header[0:4], 8)

	copy(header[4:14], station)
	copy(header[14:24], channel)

	binary.LittleEndian.PutUint32(header[24:28], uint32(t.Year()))
	binary.LittleEndian.PutUint32(header[28:32], uint32(mdy2dy(int(t.Month()), t.Day(), t.Year())))
	binary.LittleEndian.PutUint32(header[32:36], uint32(t.Hour()))
	binary.LittleEndian.PutUint32(header[36:40], uint32(t.Minute()))
	binary.LittleEndian.PutUint32(header[40:44], uint32(t.Second()))
	binary.LittleEndian.PutUint32(header[44:48], uint32(t.Nanosecond()/1000))
	binary.LittleEndian.PutUint32(header[48:52], uint32(usecCorr))
	binary.LittleEndian.PutUint32(header[52:56], uint32(timingQuality))
	binary.LittleEndian.PutUint32(header[56:60], uint32(len(data)))

	if _, err := s.fd.Write(header); err != nil {
		return err
	}

	dataBytes := make([]byte, len(data)*4)
	for i, v := range data {
		binary.LittleEndian.PutUint32(dataBytes[i*4:(i+1)*4], uint32(v))
	}

	if _, err := s.fd.Write(dataBytes); err != nil {
		return err
	}

	_ = s.fd.Sync() // just ignore the error
	return nil
}
