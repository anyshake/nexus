package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestSetForwarderRealtime(t *testing.T) {
	for _, tc := range []struct {
		name     string
		realtime bool
		want     string
	}{
		{"disabled", false, "AT+REALTIME=0\r\n"},
		{"enabled", true, "AT+REALTIME=1\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, server := net.Pipe()
			t.Cleanup(func() {
				client.Close()
				server.Close()
			})
			if err := server.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}

			result := make(chan error, 1)
			go func() {
				result <- setForwarderRealtime(client, tc.realtime, 5*time.Second)
			}()

			command, err := bufio.NewReader(server).ReadString('\n')
			if err != nil {
				t.Fatal(err)
			}
			if command != tc.want {
				t.Fatalf("command = %q, want %q", command, tc.want)
			}

			// Observer does not acknowledge mode changes.
			select {
			case err := <-result:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("mode command did not complete without an acknowledgement")
			}
		})
	}
}
