package main

import "io"

type ioStreams interface {
	In() io.Reader
	Out() io.Writer
	ErrOut() io.Writer
}

type basicStreams struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
}

func (s basicStreams) In() io.Reader {
	return s.in
}

func (s basicStreams) Out() io.Writer {
	return s.out
}

func (s basicStreams) ErrOut() io.Writer {
	return s.errOut
}
