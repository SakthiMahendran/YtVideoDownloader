package yt

import (
	"io"
)

func NewVideoReader(reader io.Reader) VideoReader {
	vr := VideoReader{}
	vr.Reader = reader

	return vr
}

type VideoReader struct {
	io.Reader
	total float64
}

func (vr *VideoReader) Read(p []byte) (int, error) {
	n, err := vr.Reader.Read(p)
	vr.total += float64(n)

	return n, err
}
