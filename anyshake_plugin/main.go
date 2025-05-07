package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func sendMessage(ipcObj *SeedLinkIPC, message *Message) error {
	if ipcObj == nil {
		return errors.New("IPC instance is nil")
	}
	if message == nil {
		return errors.New("message instance is nil")
	}

	station := fmt.Sprintf("%s.%s", message.Network, message.Station)
	channelId := fmt.Sprintf("CH%d", message.Index)
	return ipcObj.SendRaw3(station, channelId, message.Time, 0, 100, message.Data)
}

func main() {
	args := parseCommandLine()

	ipcObj := NewSeedlinkIPC()
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

	for reader := bufio.NewReader(conn); ; {
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
				message, err := NewMessage(line)
				if err != nil {
					log.Printf("error decoding message: %v", err)
					continue
				}

				if err := sendMessage(ipcObj, &message); err != nil {
					log.Printf("error sending message: %v", err)
					return
				}

				if args.verbose {
					log.Printf(
						"%s.%s.%s.%s, %d SPS - %s",
						message.Network, message.Station, message.Location, message.Channel,
						message.SampleRate,
						message.Time.Format(time.RFC3339Nano),
					)
				}
			}
		}
	}
}
