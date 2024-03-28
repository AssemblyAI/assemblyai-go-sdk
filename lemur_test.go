package assemblyai

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeMUR_Summarize(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/summary", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var body LeMURSummaryParams
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)

		want := LeMURSummaryParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
				Context:       "Additional context",
			},
		}

		require.Equal(t, want, body)

		writeFileResponse(t, w, "testdata/lemur/summarize.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
			Context:       "Additional context",
		},
	})
	require.NoError(t, err)

	require.Equal(t, lemurSummaryWildfires, *response.Response)
}

func TestLeMUR_SummarizeWithStructContext(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/summary", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var body LeMURSummaryParams
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)

		want := LeMURSummaryParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
				Context:       map[string]interface{}{"key": "value"},
			},
		}

		require.Equal(t, want, body)

		writeFileResponse(t, w, "testdata/lemur/summarize.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Summarize(ctx, LeMURSummaryParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
			Context:       map[string]interface{}{"key": "value"},
		},
	})
	require.NoError(t, err)

	require.Equal(t, lemurSummaryWildfires, *response.Response)
}

func TestLeMUR_QuestionAnswer(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/question-answer", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var body LeMURQuestionAnswerParams
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)

		want := LeMURQuestionAnswerParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
			Questions: []LeMURQuestion{
				{Question: String("What's causing the wildfires?")},
			},
		}

		require.Equal(t, want, body)

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
	require.NoError(t, err)

	want := []LeMURQuestionAnswer{
		{
			Question: String("What's causing the wildfires?"),
			Answer:   String("The dry conditions and weather systems channeling smoke from the Canadian wildfires into the US are causing the wildfires."),
		},
	}

	require.Equal(t, want, answers.Response)
}

func TestLeMUR_ActionItems(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/generate/action-items", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var body LeMURActionItemsParams
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)

		want := LeMURActionItemsParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
		}

		require.Equal(t, want, body)

		writeFileResponse(t, w, "testdata/lemur/action-items.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.ActionItems(ctx, LeMURActionItemsParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
		},
	})
	require.NoError(t, err)

	require.Equal(t, lemurActionItemsWildfires, *response.Response)
}

func TestLeMUR_Task(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	prompt := `
You are a helpful coach. Provide an analysis of the transcripts and offer areas
to improve with exact quotes. Include no preamble. Start with an overall summary
then get into the examples with feedback.
`

	handler.HandleFunc("/lemur/v3/generate/task", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var body LeMURTaskParams
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)

		want := LeMURTaskParams{
			LeMURBaseParams: LeMURBaseParams{
				TranscriptIDs: []string{"transcript_id"},
			},
			Prompt: String(prompt),
		}

		require.Equal(t, want, body)

		writeFileResponse(t, w, "testdata/lemur/task.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.Task(ctx, LeMURTaskParams{
		LeMURBaseParams: LeMURBaseParams{
			TranscriptIDs: []string{"transcript_id"},
		},
		Prompt: String(prompt),
	})
	require.NoError(t, err)

	require.Equal(t, lemurTaskWildfires, *response.Response)
}

func TestLeMUR_PurgeRequestData(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/lemur/v3/23f1485d-b3ba-4bba-8910-c16085e1afa5", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "DELETE", r.Method)

		writeFileResponse(t, w, "testdata/lemur/purge-request-data.json")
	})

	ctx := context.Background()

	response, err := client.LeMUR.PurgeRequestData(ctx, "23f1485d-b3ba-4bba-8910-c16085e1afa5")
	require.NoError(t, err)

	require.True(t, ToBool(response.Deleted))
}
