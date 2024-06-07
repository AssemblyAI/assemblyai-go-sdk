package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/gordonklaus/portaudio"
)

type realtimeHandler struct{}

func (h *realtimeHandler) SessionBegins(event assemblyai.SessionBegins) {
	slog.Info("session begins")
}

func (h *realtimeHandler) SessionTerminated(event assemblyai.SessionTerminated) {
	slog.Info("session terminated")
}

func (h *realtimeHandler) FinalTranscript(transcript assemblyai.FinalTranscript) {
	fmt.Println(transcript.Text)
}

func (h *realtimeHandler) PartialTranscript(transcript assemblyai.PartialTranscript) {
	fmt.Printf("%s\r", transcript.Text)
}

func (h *realtimeHandler) Error(err error) {
	slog.Error("Something bad happened", "err", err)
}

func main() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// We need portaudio to record the microphone.
	err := portaudio.Initialize()
	checkErr(err)
	defer portaudio.Terminate()

	var (
		// Number of samples per seconds.
		sampleRate = 16_000

		// Number of samples to send at once.
		framesPerBuffer = 3_200
	)

	var h realtimeHandler

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")

	client := assemblyai.NewRealTimeClientWithOptions(
		assemblyai.WithRealTimeAPIKey(apiKey),
		assemblyai.WithRealTimeSampleRate(int(sampleRate)),
		assemblyai.WithHandler(&h),
	)

	ctx := context.Background()

	err = client.Connect(ctx)
	checkErr(err)

	slog.Info("connected to real-time API", "sample_rate", sampleRate, "frames_per_buffer", framesPerBuffer)

	rec, err := newRecorder(sampleRate, framesPerBuffer)
	checkErr(err)

	err = rec.Start()
	checkErr(err)

	slog.Info("recording...")

	for {
		select {
		case <-sigs:
			slog.Info("stopping recording...")

			var err error

			err = rec.Stop()
			checkErr(err)

			err = client.Disconnect(ctx, true)
			checkErr(err)

			os.Exit(0)
		default:
			b, err := rec.Read()
			checkErr(err)

			// Send partial audio samples.
			err = client.Send(ctx, b)
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		slog.Error("Something bad happened", "err", err)
		os.Exit(1)
	}
}
