package client

import "io"

type (
	Cli struct {
		in       io.ReadCloser
		out, err io.Writer
	}
)

func (c *Cli) Out() io.Writer {
	return c.out
}

func (c *Cli) Err() io.Writer {
	return c.err
}

func (c *Cli) In() io.ReadCloser {
	return c.in
}
