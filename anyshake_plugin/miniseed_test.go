package main

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestMiniSeedChunkTimes(t *testing.T) {
	start := time.Date(2024, 12, 31, 23, 59, 59, 900000000, time.UTC)
	for _, tc := range []struct {
		rate    int
		offsets []time.Duration
	}{
		{100, []time.Duration{0, time.Second, 2 * time.Second}},
		// BTIME stores fractions of a second at 100-microsecond resolution.
		{150, []time.Duration{0, 666600 * time.Microsecond, 1333300 * time.Microsecond}},
		{300, []time.Duration{0, 333300 * time.Microsecond, 666600 * time.Microsecond}},
		{1500, []time.Duration{0, 66600 * time.Microsecond, 133300 * time.Microsecond}},
	} {
		t.Run(fmt.Sprint(tc.rate), func(t *testing.T) {
			packet := NewMiniSeedData(start, "TEST", "XX", "00", "BHZ", tc.rate, make([]int32, 201))
			records, err := packet.EncodeChunk(1)
			if err != nil {
				t.Fatal(err)
			}
			if len(records) != len(tc.offsets) {
				t.Fatalf("record count = %d, want %d", len(records), len(tc.offsets))
			}
			for i, record := range records {
				if len(record) != 512 {
					t.Fatalf("record %d: size = %d, want 512", i, len(record))
				}
				// Read the on-wire BTIME fields directly, including its 0.0001 s units.
				btime := record[20:30]
				got := time.Date(
					int(binary.BigEndian.Uint16(btime[0:2])), time.January,
					int(binary.BigEndian.Uint16(btime[2:4])),
					int(btime[4]), int(btime[5]), int(btime[6]),
					int(binary.BigEndian.Uint16(btime[8:10]))*100000, time.UTC,
				)
				want := start.Add(tc.offsets[i])
				if !got.Equal(want) {
					t.Fatalf("record %d: start = %s, want %s", i, got, want)
				}
			}
		})
	}
}

func TestMiniSeedRejectsInvalidSampleRate(t *testing.T) {
	for _, rate := range []int{0, -1} {
		packet := NewMiniSeedData(time.Now(), "TEST", "XX", "00", "BHZ", rate, []int32{1})
		if _, err := packet.EncodeChunk(1); err == nil {
			t.Errorf("sample rate %d: expected an error", rate)
		}
	}
}
