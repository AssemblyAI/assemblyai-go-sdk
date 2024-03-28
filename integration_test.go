package assemblyai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		t.Skip("ASSEMBLYAI_API_KEY not set")
	}

	client := NewClient(apiKey)

	ctx := context.Background()

	path := "./testdata/wildfires.mp3"

	// Transcript
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	t.Logf("uploading %q...", path)

	audioURL, err := client.Upload(ctx, f)
	require.NoError(t, err)

	t.Log("uploaded file")
	t.Logf("submitting %q...", audioURL)

	transcript, err := client.Transcripts.SubmitFromURL(ctx, audioURL, nil)
	require.NoError(t, err)

	require.Equal(t, TranscriptStatusQueued, transcript.Status)

	t.Log("submitted transcription")
	t.Log("waiting for transcription...")

	transcript, err = client.Transcripts.Wait(ctx, ToString(transcript.ID))
	require.NoError(t, err)

	require.Equal(t, TranscriptStatusCompleted, transcript.Status)

	require.NotEmpty(t, *transcript.Text)

	t.Log("completed transcription")
	t.Logf("transcript.Text = %v", ToString(transcript.Text))

	// LeMUR

	t.Logf("summarizing transcript (%v)...", transcript.ID)

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{ToString(transcript.ID)},
		},
	})
	require.NoError(t, err)

	t.Logf("LeMUR.Summarize = %v", ToString(response.Response))
}
