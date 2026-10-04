package main

import (
	"fmt"
	"time"

	"github.com/bclswl0827/mseedio"
)

const MINISEED_CHUNK_SAMPLES = 100

type MiniSeedData struct {
	Station    string
	Network    string
	Location   string
	Channel    string
	Timestamp  int64
	SampleRate int
	Data       []int32
}

func NewMiniSeedData(timestamp time.Time, station, network, location, channel string, sampleRate int, data []int32) MiniSeedData {
	return MiniSeedData{
		Timestamp:  timestamp.UnixMilli(),
		Station:    station,
		Network:    network,
		Location:   location,
		Channel:    channel,
		SampleRate: sampleRate,
		Data:       data,
	}
}

func (m *MiniSeedData) chunkInt32Slice(data []int32, chunkSamples int) [][]int32 {
	var chunks [][]int32

	for i := 0; i < len(data); i += chunkSamples {
		end := min(i+chunkSamples, len(data))
		chunks = append(chunks, data[i:end])
	}

	return chunks
}

func (m *MiniSeedData) EncodeChunk(sequenceNumber int) ([][]byte, error) {
	if m.SampleRate <= 0 {
		return nil, fmt.Errorf("invalid sample rate: %d", m.SampleRate)
	}

	startTime := time.UnixMilli(m.Timestamp).UTC()
	var buf [][]byte

	for i, c := range m.chunkInt32Slice(m.Data, MINISEED_CHUNK_SAMPLES) {
		var miniseed mseedio.MiniSeedData
		if err := miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST); err != nil {
			return nil, err
		}

		// Calculate from the total sample offset to avoid accumulating rounding errors.
		offset := time.Duration(i) * MINISEED_CHUNK_SAMPLES * time.Second / time.Duration(m.SampleRate)
		if err := miniseed.Append(c, &mseedio.AppendOptions{
			ChannelCode:    m.Channel,
			StationCode:    m.Station,
			LocationCode:   m.Location,
			NetworkCode:    m.Network,
			SampleRate:     float64(m.SampleRate),
			SequenceNumber: fmt.Sprintf("%06d", sequenceNumber),
			StartTime:      startTime.Add(offset),
		}); err != nil {
			return nil, err
		}

		for i := 0; i < len(miniseed.Series); i++ {
			miniseed.Series[i].BlocketteSection.RecordLength = 9
		}

		msData, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
		if err != nil {
			return nil, err
		}

		buf = append(buf, msData)
	}

	return buf, nil
}
