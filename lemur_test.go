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

		var body LeMURSummaryParameters
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURSummaryParameters{
			LeMURBaseParameters: LeMURBaseParameters{
				TranscriptIDs: []string{"transcript_id"},
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/summarize.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParameters{
		LeMURBaseParameters: LeMURBaseParameters{
			TranscriptIDs: []string{"transcript_id"},
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

		var body LeMURQuestionAnswerParameters
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURQuestionAnswerParameters{
			LeMURBaseParameters: LeMURBaseParameters{
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

	answers, err := client.LeMUR.Question(ctx, LeMURQuestionAnswerParameters{
		LeMURBaseParameters: LeMURBaseParameters{
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

		var body LeMURActionItemsParameters
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURActionItemsParameters{
			LeMURBaseParameters: LeMURBaseParameters{
				TranscriptIDs: []string{"transcript_id"},
			},
		}

		if !cmp.Equal(body, want) {
			t.Errorf("Request body = %+v, want = %+v", body, want)
		}

		writeFileResponse(t, w, "testdata/lemur/action-items.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.ActionItems(ctx, LeMURActionItemsParameters{
		LeMURBaseParameters: LeMURBaseParameters{
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

		var body LeMURTaskParameters
		json.NewDecoder(r.Body).Decode(&body)

		want := LeMURTaskParameters{
			LeMURBaseParameters: LeMURBaseParameters{
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

	response, err := client.LeMUR.Task(ctx, LeMURTaskParameters{
		LeMURBaseParameters: LeMURBaseParameters{
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
