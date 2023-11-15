package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/gordonklaus/portaudio"
)

type realtimeHandler struct{}

func (h *realtimeHandler) SessionBegins(event assemblyai.SessionBegins) {
	fmt.Println("session begins")
}

func (h *realtimeHandler) SessionTerminated(event assemblyai.SessionTerminated) {
	fmt.Println("session terminated")
}

func (h *realtimeHandler) FinalTranscript(transcript assemblyai.FinalTranscript) {
	fmt.Println(transcript.Text)
}

func (h *realtimeHandler) PartialTranscript(transcript assemblyai.PartialTranscript) {
	fmt.Printf("%s\r", transcript.Text)
}

func (h *realtimeHandler) Error(err error) {
	fmt.Println(err)
}

func main() {
	logger := log.New(os.Stderr, "", log.Lshortfile)

	var (
		apiKey = os.Getenv("ASSEMBLYAI_API_KEY")

		// Number of samples per seconds.
		sampleRate = 16000

		// Number of samples to send at once.
		framesPerBuffer = 3200
	)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var h realtimeHandler

	client := assemblyai.NewRealTimeClient(apiKey, &h)

	ctx := context.Background()

	if err := client.Connect(ctx); err != nil {
		logger.Fatal(err)
	}

	// We need portaudio to record the microphone.
	portaudio.Initialize()
	defer portaudio.Terminate()

	rec, err := newRecorder(sampleRate, framesPerBuffer)
	if err != nil {
		logger.Fatal(err)
	}

	if err := rec.Start(); err != nil {
		logger.Fatal(err)
	}

	for {

		select {
		case <-sigs:
			fmt.Println("stopping recording...")

			if err := rec.Stop(); err != nil {
				log.Fatal(err)
			}

			if err := client.Disconnect(ctx, true); err != nil {
				log.Fatal(err)
			}

			os.Exit(0)
		default:
			b, err := rec.Read()
			if err != nil {
				logger.Fatal(err)
			}

			// Send partial audio samples.
			if err := client.Send(ctx, b); err != nil {
				logger.Fatal(err)
			}
		}
	}
}
