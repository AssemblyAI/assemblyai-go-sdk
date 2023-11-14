package assemblyai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTranscripts_Submit(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("io.ReadAll returned error: %v", err)
		}

		// Check that we're not marshaling zero values.
		wantBody := `{"audio_url":"https://example.com/wildfires.mp3"}
`
		if string(b) != wantBody {
			t.Errorf("Request body = %v, want = %v", string(b), wantBody)
		}

		var body CreateTranscriptParameters
		json.Unmarshal(b, &body)

		want := CreateTranscriptParameters{AudioURL: String(fakeAudioURL)}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/transcript/queued.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.Submit(ctx, fakeAudioURL, nil)
	if err != nil {
		t.Errorf("Transcripts.Submit returned error: %v", err)
	}

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

	if !cmp.Equal(transcript, want) {
		t.Errorf(cmp.Diff(want, transcript))
	}
}

func TestTranscripts_Delete(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript/"+fakeTranscriptID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		writeFileResponse(t, w, "testdata/transcript/deleted.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.Delete(ctx, fakeTranscriptID)
	if err != nil {
		t.Errorf("Transcripts.Delete returned error: %v", err)
	}

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

	if !cmp.Equal(transcript, want) {
		t.Errorf(cmp.Diff(want, transcript))
	}
}

func TestTranscripts_Get(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript/"+fakeTranscriptID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		writeFileResponse(t, w, "testdata/transcript/completed.json")
	})

	ctx := context.Background()

	transcript, err := client.Transcripts.Get(ctx, fakeTranscriptID)
	if err != nil {
		t.Errorf("Transcripts.Get returned error: %v", err)
	}

	b, err := os.ReadFile("testdata/transcript/wildfires.txt")
	if err != nil {
		t.Errorf("os.ReadFile returned error: %v", err)
	}

	if *transcript.Text != strings.TrimSpace(string(b)) {
		t.Errorf("transcript.Text = %v,\nwant = %v", transcript.Text, string(b))
	}
}

func TestTranscripts_List(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/transcript", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQuery(t, r, "limit=2")

		writeFileResponse(t, w, "testdata/transcript/list.json")
	})

	ctx := context.Background()

	results, err := client.Transcripts.List(ctx, TranscriptListParameters{
		Limit: Int64(2),
	})
	if err != nil {
		t.Errorf("Transcripts.List returned error: %v", err)
	}

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

	if !cmp.Equal(results, want) {
		t.Errorf(cmp.Diff(want, results))
	}
}
