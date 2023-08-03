package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/adrianosela/stitch"
)

const (
	n = 3
)

func main() {
	l, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatalf("failed to start tcp listener: %v", err)
	}
	defer l.Close()

	conns := []io.ReadWriteCloser{}
	for i := 0; i < n; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("failed to accept conn %d: %v", i, err)
		}
		defer conn.Close()
		conns = append(conns, conn)
	}

	if err = stitch.Stitch(context.TODO(), conns...); err != nil {
		log.Fatalf("failed to stitch connections: %v", err)
	}
}
