package main

import (
	"encoding/binary"
	"os"
	"time"
)

const (
	PLUGIN_FD             = 63
	PLUGIN_MSEED_SIZE     = 512
	PLUGIN_MAX_MSG_SIZE   = 448
	PLUGIN_MAX_DATA_BYTES = 4000
)

const (
	PLUGIN_RAW_TIME_PACKET  = 8
	PLUGIN_RAW_PACKET       = 9
	PLUGIN_RAW_GAP_PACKET   = 10
	PLUGIN_RAW_FLUSH_PACKET = 11
	PLUGIN_LOG_PACKET       = 12
	PLUGIN_MSEED_PACKET     = 13
)

type PluginPacketHeader struct {
	PackType       uint32
	Station        [10]byte
	Channel        [10]byte
	Year           uint32
	Yday           uint32
	Hour           uint32
	Minute         uint32
	Second         uint32
	Usec           uint32
	UsecCorrection int32
	TimingQuality  int32
	DataSize       int32
}

type SeedLinkPluginIPC struct {
	fd *os.File
}

func (s *SeedLinkPluginIPC) sendPacket(head *PluginPacketHeader, data []byte) error {
	headerBuf := make([]byte, 60)

	binary.LittleEndian.PutUint32(headerBuf[0:4], head.PackType)
	copy(headerBuf[4:14], head.Station[:])
	copy(headerBuf[14:24], head.Channel[:])
	binary.LittleEndian.PutUint32(headerBuf[24:28], head.Year)
	binary.LittleEndian.PutUint32(headerBuf[28:32], head.Yday)
	binary.LittleEndian.PutUint32(headerBuf[32:36], head.Hour)
	binary.LittleEndian.PutUint32(headerBuf[36:40], head.Minute)
	binary.LittleEndian.PutUint32(headerBuf[40:44], head.Second)
	binary.LittleEndian.PutUint32(headerBuf[44:48], head.Usec)
	binary.LittleEndian.PutUint32(headerBuf[48:52], uint32(head.UsecCorrection))
	binary.LittleEndian.PutUint32(headerBuf[52:56], uint32(head.TimingQuality))
	binary.LittleEndian.PutUint32(headerBuf[56:60], uint32(head.DataSize))

	if _, err := s.fd.Write(headerBuf); err != nil {
		return err
	}
	if data != nil {
		if _, err := s.fd.Write(data); err != nil {
			return err
		}
	}

	return nil
}

func (s *SeedLinkPluginIPC) isLeap(y int) bool {
	return (y%400 == 0) || (y%4 == 0 && y%100 != 0)
}

func (s *SeedLinkPluginIPC) ldoy(y, m int) int {
	doy := [...]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365}
	if s.isLeap(y) && m >= 3 {
		return doy[m-1] + 1
	}
	return doy[m-1]
}

func (s *SeedLinkPluginIPC) mdy2dy(month, day, year int) int {
	return s.ldoy(year, month) + day - 1
}

func (s *SeedLinkPluginIPC) Close() {
	_ = s.fd.Close()
}

func (s *SeedLinkPluginIPC) SendRaw3(station, channel string, t time.Time, usecCorr, timingQuality int, data []int32) error {
	const maxSamplesPerPacket = PLUGIN_MAX_DATA_BYTES / 4 // 4000 / 4 = 1000

	sent := 0
	total := len(data)
	first := true

	for sent < total {
		end := sent + maxSamplesPerPacket
		if end > total {
			end = total
		}
		chunk := data[sent:end]

		var head PluginPacketHeader
		copy(head.Station[:], station)
		copy(head.Channel[:], channel)

		if first {
			head.PackType = PLUGIN_RAW_TIME_PACKET
			head.Year = uint32(t.Year())
			head.Yday = uint32(s.mdy2dy(int(t.Month()), t.Day(), t.Year()))
			head.Hour = uint32(t.Hour())
			head.Minute = uint32(t.Minute())
			head.Second = uint32(t.Second())
			head.Usec = uint32(t.Nanosecond() / 1000)
			head.UsecCorrection = int32(usecCorr)
			head.TimingQuality = int32(timingQuality)
			first = false
		} else {
			head.PackType = PLUGIN_RAW_PACKET
		}

		head.DataSize = int32(len(chunk))

		dataBytes := make([]byte, len(chunk)*4)
		for i, v := range chunk {
			binary.LittleEndian.PutUint32(dataBytes[i*4:(i+1)*4], uint32(v))
		}

		if err := s.sendPacket(&head, dataBytes); err != nil {
			return err
		}

		sent = end
	}

	return nil
}

func (s *SeedLinkPluginIPC) SendFlush3(station, channel string) error {
	var head PluginPacketHeader
	copy(head.Station[:], station)
	copy(head.Channel[:], channel)
	head.PackType = PLUGIN_RAW_FLUSH_PACKET
	head.DataSize = 0

	return s.sendPacket(&head, nil)
}

func (s *SeedLinkPluginIPC) SendMSeed(station string, data []byte) error {
	if len(data) != PLUGIN_MSEED_SIZE {
		return nil
	}

	var head PluginPacketHeader
	copy(head.Station[:], station)
	head.PackType = PLUGIN_MSEED_PACKET
	head.DataSize = int32(len(data))

	return s.sendPacket(&head, data)
}

func (s *SeedLinkPluginIPC) SendMSeed2(station, channel string, seq int, data []byte) error {
	if len(data) != PLUGIN_MSEED_SIZE {
		return nil
	}

	var head PluginPacketHeader
	copy(head.Station[:], station)
	copy(head.Channel[:], channel)
	head.PackType = PLUGIN_MSEED_PACKET
	head.TimingQuality = int32(seq)
	head.DataSize = int32(len(data))

	return s.sendPacket(&head, data)
}

func (s *SeedLinkPluginIPC) SendLog3(station string, t time.Time, msg string) error {
	var head PluginPacketHeader
	copy(head.Station[:], station)
	head.PackType = PLUGIN_LOG_PACKET

	head.Year = uint32(t.Year())
	head.Yday = uint32(s.mdy2dy(int(t.Month()), t.Day(), t.Year()))
	head.Hour = uint32(t.Hour())
	head.Minute = uint32(t.Minute())
	head.Second = uint32(t.Second())
	head.Usec = uint32(t.Nanosecond() / 1000)

	data := []byte(msg)
	head.DataSize = int32(len(data))

	return s.sendPacket(&head, data)
}

func (s *SeedLinkPluginIPC) SendRawDepoch(station, channel string, depoch float64, usecCorr, timingQuality int, data []int32) error {
	sec := int64(depoch)
	usec := int((depoch - float64(sec)) * 1e6)
	t := time.Unix(sec, int64(usec)*1000).UTC()
	return s.SendRaw3(station, channel, t, usecCorr, timingQuality, data)
}

func NewSeedlinkPluginIPC() SeedLinkPluginIPC {
	return SeedLinkPluginIPC{
		fd: os.NewFile(PLUGIN_FD, "seedlink"),
	}
}
