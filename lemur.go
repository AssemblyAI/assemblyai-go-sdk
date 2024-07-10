package assemblyai

import (
	"context"
)

const (
	// Claude 3.5 Sonnet is the most intelligent model to date, outperforming
	// Claude 3 Opus on a wide range of evaluations, with the speed and cost of
	// Claude 3 Sonnet.
	LeMURModelAnthropicClaude3_5_Sonnet LeMURModel = "anthropic/claude-3-5-sonnet"

	// Claude 3 Opus is good at handling complex analysis, longer tasks with
	// many steps, and higher-order math and coding tasks.
	LeMURModelAnthropicClaude3_Opus LeMURModel = "anthropic/claude-3-opus"

	// Claude 3 Haiku is the fastest model that can execute lightweight actions.
	LeMURModelAnthropicClaude3_Haiku LeMURModel = "anthropic/claude-3-haiku"

	// Claude 3 Sonnet is a legacy model with a balanced combination of
	// performance and speed for efficient, high-throughput tasks.
	LeMURModelAnthropicClaude3_Sonnet LeMURModel = "anthropic/claude-3-sonnet"

	// Claude 2.1 is a legacy model similar to Claude 2.0. The key difference is
	// that it minimizes model hallucination and system prompts, has a larger
	// context window, and performs better in citations.
	LeMURModelAnthropicClaude2_1 LeMURModel = "anthropic/claude-2-1"

	// Claude 2.0 is a legacy model that has good complex reasoning. It offers
	// more nuanced responses and improved contextual comprehension.
	LeMURModelAnthropicClaude2 LeMURModel = "anthropic/claude-2"

	// Legacy model. The same as [LeMURModelAnthropicClaude2].
	LeMURModelDefault LeMURModel = "default"

	// Claude Instant is a legacy model that is optimized for speed and cost.
	// Claude Instant can complete requests up to 20% faster than Claude 2.0.
	LeMURModelAnthropicClaudeInstant1_2 LeMURModel = "anthropic/claude-instant-1-2"

	// Legacy model. The same as [LeMURModelAnthropicClaudeInstant1_2].
	LeMURModelBasic LeMURModel = "basic"

	// Mistral 7B is an open source model that works well for summarization and
	// answering questions.
	LeMURModelAssemblyAIMistral7B LeMURModel = "assemblyai/mistral-7b"
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
