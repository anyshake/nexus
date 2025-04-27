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

	connTimeout := time.Duration(args.timeout) * time.Second
	_, connCancel := context.WithTimeout(context.Background(), connTimeout)
	defer connCancel()

	conn, err := net.DialTimeout("tcp", args.address, connTimeout)
	if err != nil {
		log.Printf("error connecting to server: %v", err)
		return
	}
	log.Printf("connected to server at %s", args.address)

	shutdownCtx, shutdownCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer shutdownCancel()

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
					shutdownCancel()
					return
				}
				if strings.HasPrefix(line, "$") {
					_, _ = fifoBuf.Write(line)
				}
			case <-shutdownCtx.Done():
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
					seisCompDaemonCallback(message)
				}

			case <-shutdownCtx.Done():
				log.Println("exiting from data packet decoder")
				timer.Stop()
				return
			}
		}
	}(500 * time.Millisecond)

	<-shutdownCtx.Done()
	log.Println("interrupt signal received, exiting...")
}
