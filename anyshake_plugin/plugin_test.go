package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

func capturePluginOutput(t *testing.T, send func(*SeedLinkPluginIPC) error) ([]byte, error) {
	t.Helper()
	fd, err := os.CreateTemp(t.TempDir(), "ipc")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { fd.Close() })
	ipc := SeedLinkPluginIPC{fd: fd}
	sendErr := send(&ipc)
	data, err := os.ReadFile(fd.Name())
	if err != nil {
		t.Fatal(err)
	}
	return data, sendErr
}

func TestRawPluginPackets(t *testing.T) {
	// Local March 1 is still February 29 UTC (day 60 of a leap year).
	timestamp := time.Date(2024, 3, 1, 0, 30, 12, 123456000, time.FixedZone("UTC+8", 8*60*60))
	samples := make([]int32, 1001)
	for i := range samples {
		samples[i] = int32(i) - 500
	}
	data, err := capturePluginOutput(t, func(ipc *SeedLinkPluginIPC) error {
		if err := ipc.SendRaw3("XX.TEST", "BHZ", timestamp, -25, 90, nil); err != nil {
			return err
		}
		if err := ipc.SendRaw3("XX.TEST", "BHZ", timestamp, -25, 90, samples); err != nil {
			return err
		}
		if err := ipc.SendRaw3("XX.TEST", "BHZ", time.Time{}, 0, 0, samples[:2]); err != nil {
			return err
		}
		if err := ipc.SendGap3("XX.TEST", "BHZ", -50, -1, 1800); err != nil {
			return err
		}
		return ipc.SendFlush3("XX.TEST", "BHZ")
	})
	if err != nil {
		t.Fatal(err)
	}

	clock := PluginPacketHeader{
		PackType: 8, Year: 2024, Yday: 60, Hour: 16, Minute: 30, Second: 12,
		Usec: 123456, UsecCorrection: -25, TimingQuality: 90,
	}
	timedData := clock
	timedData.DataSize = 1000
	want := []struct {
		header  PluginPacketHeader
		samples []int32
	}{
		{clock, nil},
		{timedData, samples[:1000]},
		{PluginPacketHeader{PackType: 9, DataSize: 1}, samples[1000:]},
		{PluginPacketHeader{PackType: 9, DataSize: 2}, samples[:2]},
		{PluginPacketHeader{PackType: 10, UsecCorrection: -50, TimingQuality: -1, DataSize: 1800}, nil},
		{PluginPacketHeader{PackType: 11}, nil},
	}

	// Decode the stream as SeedLink does: native header, then samples only for
	// raw data packets. A gap's data_size is a count, not a payload length.
	reader := bytes.NewReader(data)
	for i, packet := range want {
		copy(packet.header.Station[:], "XX.TEST")
		copy(packet.header.Channel[:], "BHZ")
		var header PluginPacketHeader
		if err := binary.Read(reader, binary.NativeEndian, &header); err != nil {
			t.Fatal(err)
		}
		if header != packet.header {
			t.Fatalf("packet %d: header = %+v, want %+v", i, header, packet.header)
		}
		if header.PackType == 8 || header.PackType == 9 {
			got := make([]int32, header.DataSize)
			if err := binary.Read(reader, binary.NativeEndian, got); err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, packet.samples) {
				t.Fatalf("packet %d: sample payload mismatch", i)
			}
		}
	}
	if reader.Len() != 0 {
		t.Fatalf("unexpected trailing data: %d bytes", reader.Len())
	}
}

func TestSendLog3(t *testing.T) {
	for _, tc := range []struct{ name, message, want string }{
		{"empty", "", "\r\n"},
		{"newlines", "ready\r\n\n", "ready\r\n"},
		{"oversized", strings.Repeat("x", 5000), strings.Repeat("x", 445) + "\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := capturePluginOutput(t, func(ipc *SeedLinkPluginIPC) error {
				return ipc.SendLog3("XX.TEST", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), tc.message)
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(data) != 60+len(tc.want) {
				t.Fatalf("packet length = %d, want %d", len(data), 60+len(tc.want))
			}
			var header PluginPacketHeader
			if err := binary.Read(bytes.NewReader(data), binary.NativeEndian, &header); err != nil {
				t.Fatal(err)
			}
			if header.PackType != 12 || header.Year != 2025 || header.Yday != 1 || header.DataSize != int32(len(tc.want)) {
				t.Fatalf("unexpected LOG header: %+v", header)
			}
			if string(data[60:]) != tc.want {
				t.Fatalf("LOG payload = %q, want %q", data[60:], tc.want)
			}
		})
	}
}

func TestMiniSeedPluginPackets(t *testing.T) {
	for _, withChannel := range []bool{false, true} {
		for _, size := range []int{0, 511, 512, 513, 1024} {
			t.Run(fmt.Sprintf("channel=%t/size=%d", withChannel, size), func(t *testing.T) {
				payload := bytes.Repeat([]byte{0xa5}, size)
				data, err := capturePluginOutput(t, func(ipc *SeedLinkPluginIPC) error {
					if withChannel {
						return ipc.SendMSeed2("XX.TEST", "BHZ", 123, payload)
					}
					return ipc.SendMSeed("XX.TEST", payload)
				})
				if size != 512 {
					if err == nil || len(data) != 0 {
						t.Fatalf("invalid size: err = %v, wrote %d bytes", err, len(data))
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if len(data) != 572 {
					t.Fatalf("packet length = %d, want 572", len(data))
				}
				want := PluginPacketHeader{PackType: 13, DataSize: 512}
				copy(want.Station[:], "XX.TEST")
				if withChannel {
					copy(want.Channel[:], "BHZ")
					want.TimingQuality = 123
				}
				var header PluginPacketHeader
				if err := binary.Read(bytes.NewReader(data), binary.NativeEndian, &header); err != nil {
					t.Fatal(err)
				}
				if header != want || !bytes.Equal(data[60:], payload) {
					t.Fatalf("miniSEED packet mismatch: header = %+v", header)
				}
			})
		}
	}
}
