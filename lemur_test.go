package assemblyai

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLeMUR_Summarize(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body LeMURSummaryParams
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURSummaryParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
				Context:       "Additional context",
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/summarize.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
			Context:       "Additional context",
		},
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	want := lemurSummaryWildfires

	if *response.Response != want {
		t.Errorf("LeMUR.Summarize = %v, want = %v", response, want)
	}
}

func TestLeMUR_SummarizeWithStructContext(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/summary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body LeMURSummaryParams
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURSummaryParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
				Context:       map[string]interface{}{"key": "value"},
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/summarize.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
			Context:       map[string]interface{}{"key": "value"},
		},
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	want := lemurSummaryWildfires

	if *response.Response != want {
		t.Errorf("LeMUR.Summarize = %v, want = %v", response, want)
	}
}

func TestLeMUR_QuestionAnswer(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/question-answer", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body LeMURQuestionAnswerParams
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURQuestionAnswerParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
			Questions: []LeMURQuestion{
				{Question: String("What's causing the wildfires?")},
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/question-answer.json")
	})

	ctx := context.Background()

	questions := []LeMURQuestion{
		{Question: String("What's causing the wildfires?")},
	}

	answers, err := client.LeMUR.Question(ctx, LeMURQuestionAnswerParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
		},
		Questions: questions,
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	want := []LeMURQuestionAnswer{
		{
			Question: String("What's causing the wildfires?"),
			Answer:   String("The dry conditions and weather systems channeling smoke from the Canadian wildfires into the US are causing the wildfires."),
		},
	}

	if !cmp.Equal(answers.Response, want) {
		t.Errorf("LeMUR.Question = %+v, want = %+v", answers, want)
	}
}

func TestLeMUR_ActionItems(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/action-items", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body LeMURActionItemsParams
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURActionItemsParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/action-items.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.ActionItems(ctx, LeMURActionItemsParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
		},
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	want := lemurActionItemsWildfires

	if *response.Response != want {
		t.Errorf("LeMUR.ActionItems = %v, want = %v", response, want)
	}
}

func TestLeMUR_Task(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	prompt := `
You are a helpful coach. Provide an analysis of the transcripts and offer areas
to improve with exact quotes. Include no preamble. Start with an overall summary
then get into the examples with feedback.
`

	handler.HandleFunc("/lemur/v3/generate/task", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var body LeMURTaskParams
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURTaskParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
			Prompt: String(prompt),
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/task.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Task(ctx, LeMURTaskParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
		},
		Prompt: String(prompt),
	})
	if err != nil {
		t.Errorf("Submit returned error: %v", err)
	}

	want := lemurTaskWildfires

	if *response.Response != want {
		t.Errorf("LeMUR.ActionItems = %v, want = %v", response, want)
	}
}
