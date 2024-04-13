package assemblyai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_SpeechToText(t *testing.T) {
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

const framesPerBufferTelephone = 8000

func TestIntegration_RealTime(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		t.Skip("ASSEMBLYAI_API_KEY not set")
	}

	sampleRate := 8_000

	var partialTranscriptInvoked, finalTranscriptInvoked bool
	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey(apiKey),
		WithRealTimeTranscriber(&RealTimeTranscriber{
			OnPartialTranscript: func(_ PartialTranscript) {
				partialTranscriptInvoked = true
			},
			OnFinalTranscript: func(_ FinalTranscript) {
				finalTranscriptInvoked = true
			},
			OnError: func(err error) {
				require.NoError(t, err)
			},
		}),
		WithRealTimeSampleRate(sampleRate),
	)

	ctx := context.Background()

	path := "./testdata/gore-short.wav"

	// Transcript
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	err = client.Connect(ctx)
	t.Log("connected to real-time API")
	require.NoError(t, err)

	buf := make([]byte, framesPerBufferTelephone)

	for {
		_, err := f.Read(buf)
		if err != nil {
			break
		}

		err = client.Send(ctx, buf)
		require.NoError(t, err)
	}

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	require.True(t, partialTranscriptInvoked)
	require.True(t, finalTranscriptInvoked)

	t.Log("disconnected from real-time API")
}

func TestIntegration_RealTime_WithoutPartialTranscripts(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		t.Skip("ASSEMBLYAI_API_KEY not set")
	}

	sampleRate := 8_000

	var partialTranscriptInvoked, finalTranscriptInvoked bool

	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey(apiKey),
		WithRealTimeTranscriber(&RealTimeTranscriber{
			OnPartialTranscript: func(_ PartialTranscript) {
				partialTranscriptInvoked = true
			},
			OnFinalTranscript: func(_ FinalTranscript) {
				finalTranscriptInvoked = true
			},
			OnError: func(err error) {
				require.NoError(t, err)
			},
		}),
		WithRealTimeSampleRate(sampleRate),
		WithRealTimeDisablePartialTranscripts(true),
	)

	ctx := context.Background()

	path := "./testdata/gore-short.wav"

	// Transcript
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	err = client.Connect(ctx)
	t.Log("connected to real-time API")
	require.NoError(t, err)

	buf := make([]byte, framesPerBufferTelephone)

	for {
		_, err := f.Read(buf)
		if err != nil {
			break
		}

		err = client.Send(ctx, buf)
		require.NoError(t, err)
	}

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	require.False(t, partialTranscriptInvoked)
	require.True(t, finalTranscriptInvoked)

	t.Log("disconnected from real-time API")
}

func TestIntegration_RealTime_WithExtraSessionInfo(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		t.Skip("ASSEMBLYAI_API_KEY not set")
	}

	sampleRate := 8_000

	var sessionInformationInvoked bool
	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey(apiKey),
		WithRealTimeTranscriber(&RealTimeTranscriber{
			OnSessionInformation: func(event SessionInformation) {
				sessionInformationInvoked = true

				require.Greater(t, event.AudioDurationSeconds, 0.0)
			},
			OnError: func(err error) {
				require.NoError(t, err)
			},
		}),
		WithRealTimeSampleRate(sampleRate),
		WithRealTimeDisablePartialTranscripts(true),
	)

	ctx := context.Background()

	path := "./testdata/gore-short.wav"

	// Transcript
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	err = client.Connect(ctx)
	t.Log("connected to real-time API")
	require.NoError(t, err)

	buf := make([]byte, framesPerBufferTelephone)

	for {
		_, err := f.Read(buf)
		if err != nil {
			break
		}

		err = client.Send(ctx, buf)
		require.NoError(t, err)
	}

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	t.Log("disconnected from real-time API")

	require.True(t, sessionInformationInvoked)
}
