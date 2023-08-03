# Stitch

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/stitch)](https://goreportcard.com/report/github.com/adrianosela/stitch)
[![Documentation](https://godoc.org/github.com/adrianosela/stitch?status.svg)](https://godoc.org/github.com/adrianosela/stitch)
[![license](https://img.shields.io/github/license/adrianosela/stitch.svg)](https://github.com/adrianosela/stitch/blob/master/LICENSE)

Stitch stitches N io.ReadWriteCloser implementations together, forwarding reads and writes between them.

### Example

```
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
```