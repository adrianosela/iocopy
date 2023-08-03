package stitch

import (
	"context"
	"errors"
	"io"
	"net"

	"golang.org/x/sync/errgroup"
)

func shouldSupressError(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, io.ErrClosedPipe)
}

func setContextCloser(ctx context.Context, rws ...io.ReadWriteCloser) {
	go func() {
		defer func() {
			for _, rw := range rws {
				rw.Close()
			}
		}()
		<-ctx.Done()
	}()
}

// Stitch stitches N io.ReadWriteCloser implementations
// together, forwarding reads and writes between them.
func Stitch(ctx context.Context, rws ...io.ReadWriteCloser) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	setContextCloser(ctx, rws...)

	for i := 0; i < len(rws); i++ {
		src := rws[i]

		dsts := []io.Writer{}
		for j := 0; j < len(rws); j++ {
			if i != j {
				dsts = append(dsts, rws[j])
			}
		}

		dst := io.MultiWriter(dsts...)

		g.Go(func() error {
			defer cancel()

			if _, err := io.Copy(dst, src); err != nil {
				if shouldSupressError(err) {
					return nil
				}
				return err
			}

			return nil
		})
	}

	return g.Wait()
}
