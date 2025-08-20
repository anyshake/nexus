package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	args := parseCommandLine()

	ipcObj := NewSeedlinkPluginIPC()
	defer ipcObj.Close()

	connTimeout := time.Duration(args.timeout) * time.Second
	_, connCancel := context.WithTimeout(context.Background(), connTimeout)
	defer connCancel()

	conn, err := net.DialTimeout("tcp", args.address, connTimeout)
	if err != nil {
		log.Printf("error connecting to server: %v", err)
		return
	}
	log.Printf("connected to server at %s", args.address)
	defer conn.Close()

	shutdownCtx, shutdownCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer shutdownCancel()

	reader := bufio.NewReader(conn)
	for sequenceNumber := 0; ; sequenceNumber++ {
		select {
		case <-shutdownCtx.Done():
			log.Println("interrupt signal received, exiting...")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			line, err := reader.ReadString('\n')
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					log.Println("timed out when reading from connection")
					continue
				}
				log.Printf("error reading from connection: %v", err)
				return
			}
			if strings.HasPrefix(line, "$") {
				message, err := NewForwarderMessage(line)
				if err != nil {
					log.Printf("error decoding message from forwarder: %v", err)
					continue
				}

				if err := sendPluginMessage(ipcObj, message, sequenceNumber); err != nil {
					log.Printf("error sending forwarder message to plugin: %v", err)
					return
				}

				if args.verbose {
					log.Printf(
						"[%s] - %s.%s.%s.%s (%d SPS)",
						message.Time.Format(time.RFC3339Nano),
						message.Network, message.Station, message.Location, message.Channel,
						message.SampleRate,
					)
				}
			}
		}
	}
}

func sendPluginMessage(ipcObj SeedLinkPluginIPC, message ForwarderMessage, sequenceNumber int) error {
	time, station, network, location, channel, sampleRate, data := message.Time, message.Station, message.Network, message.Location, message.Channel, message.SampleRate, message.Data
	packet := NewMiniSeedData(time, station, network, location, channel, sampleRate, data)
	chunk, err := packet.EncodeChunk(sequenceNumber)
	if err != nil {
		return err
	}

	pluginStationName := fmt.Sprintf("%s.%s", message.Network, message.Station)
	for _, data := range chunk {
		if err := ipcObj.SendMSeed(pluginStationName, data); err != nil {
			return err
		}
	}

	return nil
}
