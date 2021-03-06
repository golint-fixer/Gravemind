package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// We don't respond to pings so that unused connections die naturally
func NewConn(addr, username, token string, messages chan string) {
	log.Printf("Connecting to IRC as %q...\n", username)
	defer log.Printf("Disconnecting from IRC as %q.\n", username)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func() { logErr(conn.Close()) }()

	// Connect
	_, err = fmt.Fprintf(conn, "PASS oauth:%s\r\n", token)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(conn, "NICK %s\r\n", username)
	if err != nil {
		return
	}

	// Send messages
	b := make([]byte, 4096)
	for m := range messages {
		// Read all pending data to help detect dead connections
		logErr(conn.SetReadDeadline(time.Now().Add(10 * time.Microsecond)))
		for {
			if _, err := conn.Read(b); err != nil {
				if err == io.EOF {
					messages <- m
					return
				}
				break
			}
		}
		// Try to send it
		log.Println(m)
		sent, err := fmt.Fprintf(conn, "%s\r\n", m)
		// Die if something seemed to go wrong
		if err != nil || sent != len(m)+2 {
			messages <- m
			return
		}
		time.Sleep(1 * time.Second)
	}
}
