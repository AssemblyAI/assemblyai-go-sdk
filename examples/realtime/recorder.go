package main

import (
	"bytes"
	"encoding/binary"

	"github.com/gordonklaus/portaudio"
)

type recorder struct {
	stream *portaudio.Stream
	buffer []int16
}

func newRecorder(sampleRate int, framesPerBuffer int) (*recorder, error) {
	buffer := make([]int16, framesPerBuffer)

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), framesPerBuffer, buffer)
	if err != nil {
		return nil, err
	}

	return &recorder{
		stream: stream,
		buffer: buffer,
	}, nil
}

func (r *recorder) Read() ([]byte, error) {
	if err := r.stream.Read(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, r.buffer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *recorder) Start() error {
	return r.stream.Start()
}

func (r *recorder) Stop() error {
	return r.stream.Stop()
}

func (r *recorder) Close() error {
	return r.stream.Close()
}
