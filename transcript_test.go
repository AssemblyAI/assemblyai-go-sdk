package assemblyai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranscripts_Submit(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		// Check that we're not marshaling zero values.
		wantBody := `{"audio_url":"https://example.com/wildfires.mp3"}`

		require.Equal(t, wantBody, string(bytes.TrimSpace(b)))

		var body TranscriptParams
		err = json.Unmarshal(b, &body)
		require.NoError(t, err)

		want := TranscriptParams{
			AudioURL: String(fakeAudioURL),
		}

		require.Equal(t, want, body)

		writeFileResponse(t, w, "testdata/transcript/queued.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.SubmitFromURL(ctx, fakeAudioURL, nil)
	require.NoError(t, err)

	want := Transcript{
		ID:            String(fakeTranscriptID),
		AudioURL:      String(fakeAudioURL),
		LanguageCode:  "en_us",
		LanguageModel: String("assemblyai_default"),
		AcousticModel: String("assemblyai_default"),
		Punctuate:     Bool(true),
		FormatText:    Bool(true),
		WordBoost:     []string{},
		Topics:        []string{},
		Status:        TranscriptStatusQueued,

		// Disabled models
		AutoChapters:      Bool(false),
		AutoHighlights:    Bool(false),
		ContentSafety:     Bool(false),
		CustomTopics:      Bool(false),
		Disfluencies:      Bool(false),
		EntityDetection:   Bool(false),
		FilterProfanity:   Bool(false),
		IABCategories:     Bool(false),
		LanguageDetection: Bool(false),
		RedactPII:         Bool(false),
		RedactPIIAudio:    Bool(false),
		SentimentAnalysis: Bool(false),
		SpeakerLabels:     Bool(false),
		SpeedBoost:        Bool(false),
		Summarization:     Bool(false),
		WebhookAuth:       Bool(false),

		ContentSafetyLabels: ContentSafetyLabelsResult{},
		IABCategoriesResult: TopicDetectionModelResult{},
	}

	require.Equal(t, want, transcript)
}

func TestTranscripts_Delete(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript/"+fakeTranscriptID, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "DELETE", r.Method)

		writeFileResponse(t, w, "testdata/transcript/deleted.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.Delete(ctx, fakeTranscriptID)
	require.NoError(t, err)

	want := Transcript{
		ID:            String(fakeTranscriptID),
		Status:        TranscriptStatusCompleted,
		AudioURL:      String("http://deleted_by_user"),
		WebhookURL:    String("http://deleted_by_user"),
		LanguageModel: String("assemblyai_default"),
		AcousticModel: String("assemblyai_default"),
		Text:          String("Deleted by user."),
		AudioDuration: Float64(281),

		// Disabled models
		AutoChapters:      Bool(false),
		AutoHighlights:    Bool(false),
		Disfluencies:      Bool(false),
		EntityDetection:   Bool(false),
		RedactPII:         Bool(false),
		SentimentAnalysis: Bool(false),
		Summarization:     Bool(false),
		WebhookAuth:       Bool(false),
	}

	require.Equal(t, want, transcript)
}

func TestTranscripts_Get(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript/"+fakeTranscriptID, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)

		writeFileResponse(t, w, "testdata/transcript/completed.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.Get(ctx, fakeTranscriptID)
	require.NoError(t, err)

	b, err := os.ReadFile("testdata/transcript/wildfires.txt")
	require.NoError(t, err)

	require.Equal(t, string(bytes.TrimSpace(b)), *transcript.Text)
}

func TestTranscripts_List(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "limit=2", r.URL.RawQuery)

		writeFileResponse(t, w, "testdata/transcript/list.json")
	})

	ctx := context.Background()

	results, err := client.Transcripts.List(ctx, ListTranscriptParams{
		Limit: Int64(2),
	})
	require.NoError(t, err)

	want := TranscriptList{
		PageDetails: PageDetails{
			Limit:       Int64(2),
			ResultCount: Int64(2),
			CurrentURL:  String("https://api.assemblyai.com/v2/transcript?limit=2"),
			PrevURL:     String("https://api.assemblyai.com/v2/transcript?limit=2&before_id=TRANSCRIPT_ID_2"),
			NextURL:     String("https://api.assemblyai.com/v2/transcript?limit=2&after_id=TRANSCRIPT_ID_1"),
		},
		Transcripts: []TranscriptListItem{
			{
				ID:          String("TRANSCRIPT_ID_1"),
				ResourceURL: String("https://api.assemblyai.com/v2/transcript/TRANSCRIPT_ID_1"),
				Status:      TranscriptStatusCompleted,
				AudioURL:    String("http://deleted_by_user"),
				Created:     String("2023-09-04T17:02:23.506984"),
				Completed:   String("2023-09-04T17:03:02.678334"),
			},
			{
				ID:          String("TRANSCRIPT_ID_2"),
				ResourceURL: String("https://api.assemblyai.com/v2/transcript/TRANSCRIPT_ID_2"),
				Status:      TranscriptStatusCompleted,
				AudioURL:    String("http://deleted_by_user"),
				Created:     String("2023-09-04T17:02:23.504615"),
				Completed:   String("2023-09-04T17:03:00.296060"),
			},
		},
	}

	require.Equal(t, want, results)
}

func TestTranscripts_SearchWords(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript/"+fakeTranscriptID+"/word-search", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "words=hopkins%2Cwildfires", r.URL.RawQuery)

		writeFileResponse(t, w, "testdata/transcript/word-search.json")
	})

	ctx := context.Background()

	results, err := client.Transcripts.WordSearch(ctx, fakeTranscriptID, []string{"hopkins", "wildfires"})
	require.NoError(t, err)

	want := WordSearchResponse{
		ID: String("bfc3622e-8c69-4497-9a84-fb65b30dcb07"),
		Matches: []WordSearchMatch{
			{
				Count:      Int64(2),
				Indexes:    []int64{68, 835},
				Text:       String("hopkins"),
				Timestamps: []WordSearchTimestamp{{24298, 24714}, {273498, 274090}},
			},
			{
				Count:      Int64(4),
				Indexes:    []int64{4, 90, 140, 716},
				Text:       String("wildfires"),
				Timestamps: []WordSearchTimestamp{{1668, 2346}, {33852, 34546}, {50118, 51110}, {231356, 232354}},
			},
		},
		TotalCount: Int64(6),
	}

	require.Equal(t, want, results)
}
