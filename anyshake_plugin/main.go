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
				sendMessage(message)
			}
		}
	}
}
