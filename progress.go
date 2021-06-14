package main

import (
	"fmt"
	"io"
)

type ProgressWriter struct {
	now   int64
	Total int64
	Event func(int64, int64)
	w     io.Writer
	r     io.Reader
}

func NewProgressWriter(w io.Writer, t int64) *ProgressWriter {
	var writer ProgressWriter
	writer.Event = func(now, total int64) {
		fmt.Printf("\r%d/%d", now, total)
	}
	writer.w = w
	writer.now = 0
	writer.Total = t
	return &writer
}

func (w *ProgressWriter) Copy(r io.Reader) (int64, error) {
	return io.Copy(w.w, io.TeeReader(r, w))
}

func (w *ProgressWriter) Write(b []byte) (int, error) {
	w.now += int64(len(b))
	w.Event(w.now, w.Total)
	return w.w.Write(b)
}
