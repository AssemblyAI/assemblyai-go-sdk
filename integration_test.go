package assemblyai

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	integration = flag.Bool("integration", false, "Enable integration tests.")
	apiKey      = os.Getenv("ASSEMBLYAI_API_KEY")
)

func TestMain(m *testing.M) {
	flag.Parse()

	if *integration && apiKey == "" {
		fmt.Println("missing api key")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func isIntegration(t *testing.T) {
	t.Helper()

	if !*integration {
		t.Skip("skipping integration test...")
	}
}

func TestIntegration(t *testing.T) {
	isIntegration(t)

	client := NewClient(apiKey)

	ctx := context.Background()

	path := "./testdata/wildfires.mp3"

	// Transcript
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open returned error: %v", err)
	}
	defer f.Close()

	t.Logf("uploading %q...", path)

	audioURL, err := client.Upload(ctx, f)
	if err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}

	t.Log("uploaded file")
	t.Logf("submitting %q...", audioURL)

	transcript, err := client.Transcripts.Submit(ctx, audioURL, nil)
	if err != nil {
		t.Fatalf("Transcripts.Submit returned error: %v", err)
	}

	if transcript.Status != TranscriptStatusQueued {
		t.Fatalf("transcript.Status = %v, want %v", transcript.Status, err)
	}

	t.Log("submitted transcription")
	t.Log("waiting for transcription...")

	transcript, err = client.Transcripts.Wait(ctx, ToString(transcript.ID))
	if err != nil {
		t.Fatalf("Transcripts.Wait returned error: %v", err)
	}

	if transcript.Status != TranscriptStatusCompleted {
		t.Errorf("transcript.Status = %v, want %v", transcript.Status, err)
	}

	if *transcript.Text == "" {
		t.Errorf("transcript is missing text")
	}

	t.Log("completed transcription")
	t.Logf("transcript.Text = %v", ToString(transcript.Text))

	// LeMUR

	t.Logf("summarizing transcript (%v)...", transcript.ID)

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParameters{
		LeMURBaseParameters: LeMURBaseParameters{
			TranscriptIDs: []string{ToString(transcript.ID)},
		},
	})
	if err != nil {
		t.Fatalf("LeMUR.Summarize returned error: %v", err)
	}

	t.Logf("LeMUR.Summarize = %v", ToString(response.Response))
}
