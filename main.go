package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/anyshake/nexus/fifo"
	"github.com/anyshake/nexus/message"
)

func main() {
	args := parseCommandLine()

	conn, err := net.Dial("tcp", args.address)
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}
	log.Printf("connected to server at %s", args.address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fifoBuf := fifo.New[string](128)

	go func(readInterval time.Duration) {
		timer := time.NewTimer(readInterval)

		for reader := bufio.NewReader(conn); ; {
			timer.Reset(readInterval)

			select {
			case <-timer.C:
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("error reading from connection: %v", err)
					cancel()
					return
				}
				if strings.HasPrefix(line, "$") {
					_, _ = fifoBuf.Write(line)
				}
			case <-ctx.Done():
				log.Println("exiting from data packet reader")
				timer.Stop()
				conn.Close()
				return
			}
		}
	}(100 * time.Millisecond)

	go func(decodeInterval time.Duration) {
		timer := time.NewTimer(decodeInterval)

		for {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				messages, err := fifoBuf.Read(3)
				if err != nil {
					continue
				}
				for _, msg := range messages {
					message, err := message.New(msg)
					if err != nil {
						log.Printf("error decoding message: %v", err)
						continue
					}
					seisComp3DaemonCallback(message)
				}

			case <-ctx.Done():
				log.Println("exiting from data packet decoder")
				timer.Stop()
				return
			}
		}
	}(500 * time.Millisecond)

	<-ctx.Done()
	log.Println("interrupt signal received, exiting...")
}
