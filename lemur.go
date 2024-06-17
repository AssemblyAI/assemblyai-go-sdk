package assemblyai

import (
	"context"
)

const (
	// LeMUR Default is best at complex reasoning. It offers more nuanced
	// responses and improved contextual comprehension.
	LeMURModelDefault LeMURModel = "default"

	// LeMUR Basic is a simplified model optimized for speed and cost. LeMUR
	// Basic can complete requests up to 20% faster than Default.
	LeMURModelBasic LeMURModel = "basic"

	// Claude 2.1 is similar to Default, with key improvements: it minimizes
	// model hallucination and system prompts, has a larger context window, and
	// performs better in citations.
	LeMURModelAssemblyAIMistral7B LeMURModel = "assemblyai/mistral-7b"

	// LeMUR Mistral 7B is an LLM self-hosted by AssemblyAI. It's the fastest
	// and cheapest of the LLM options. We recommend it for use cases like basic
	// summaries and factual Q&A.
	LeMURModelAnthropicClaude2_1 LeMURModel = "anthropic/claude-2-1"
)

// LeMURService groups the operations related to LeMUR.
type LeMURService struct {
	client *Client
}

// Question returns answers to free-form questions about one or more transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#question--answer
func (s *LeMURService) Question(ctx context.Context, params LeMURQuestionAnswerParams) (LeMURQuestionAnswerResponse, error) {
	var response LeMURQuestionAnswerResponse

	req, err := s.client.newJSONRequest(ctx, "POST", "/lemur/v3/generate/question-answer", params)
	if err != nil {
		return LeMURQuestionAnswerResponse{}, err
	}

	if _, err := s.client.do(req, &response); err != nil {
		return LeMURQuestionAnswerResponse{}, err
	}

	return response, nil
}

// Summarize returns a custom summary of a set of transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#action-items
func (s *LeMURService) Summarize(ctx context.Context, params LeMURSummaryParams) (LeMURSummaryResponse, error) {
	req, err := s.client.newJSONRequest(ctx, "POST", "/lemur/v3/generate/summary", params)
	if err != nil {
		return LeMURSummaryResponse{}, err
	}

	var response LeMURSummaryResponse

	if _, err := s.client.do(req, &response); err != nil {
		return LeMURSummaryResponse{}, err
	}

	return response, nil
}

// ActionItems returns a set of action items based on a set of transcripts.
//
// https://www.assemblyai.com/docs/Models/lemur#action-items
func (s *LeMURService) ActionItems(ctx context.Context, params LeMURActionItemsParams) (LeMURActionItemsResponse, error) {
	req, err := s.client.newJSONRequest(ctx, "POST", "/lemur/v3/generate/action-items", params)
	if err != nil {
		return LeMURActionItemsResponse{}, err
	}

	var response LeMURActionItemsResponse

	if _, err := s.client.do(req, &response); err != nil {
		return LeMURActionItemsResponse{}, err
	}

	return response, nil
}

// Task lets you submit a custom prompt to LeMUR.
//
// https://www.assemblyai.com/docs/Models/lemur#task
func (s *LeMURService) Task(ctx context.Context, params LeMURTaskParams) (LeMURTaskResponse, error) {
	req, err := s.client.newJSONRequest(ctx, "POST", "/lemur/v3/generate/task", params)
	if err != nil {
		return LeMURTaskResponse{}, err
	}

	var response LeMURTaskResponse

	if _, err := s.client.do(req, &response); err != nil {
		return LeMURTaskResponse{}, err
	}

	return response, nil
}

func (s *LeMURService) PurgeRequestData(ctx context.Context, requestID string) (PurgeLeMURRequestDataResponse, error) {
	req, err := s.client.newJSONRequest(ctx, "DELETE", "/lemur/v3/"+requestID, nil)
	if err != nil {
		return PurgeLeMURRequestDataResponse{}, err
	}

	var response PurgeLeMURRequestDataResponse

	if _, err := s.client.do(req, &response); err != nil {
		return PurgeLeMURRequestDataResponse{}, err
	}

	return response, nil
}

// Retrieve a previously generated LeMUR response.
func (s *LeMURService) GetResponseData(ctx context.Context, requestID string, response interface{}) error {
	req, err := s.client.newJSONRequest(ctx, "GET", "/lemur/v3/"+requestID, nil)
	if err != nil {
		return err
	}

	if _, err := s.client.do(req, response); err != nil {
		return err
	}

	return nil
}
