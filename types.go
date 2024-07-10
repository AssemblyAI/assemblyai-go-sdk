package assemblyai

// Either success, or unavailable in the rare case that the model failed
type AudioIntelligenceModelStatus string

type AutoHighlightResult struct {
	// The total number of times the key phrase appears in the audio file
	Count *int64 `json:"count,omitempty"`

	// The total relevancy to the overall audio file of this key phrase - a greater number means more relevant
	Rank *float64 `json:"rank,omitempty"`

	// The text itself of the key phrase
	Text *string `json:"text,omitempty"`

	// The timestamp of the of the key phrase
	Timestamps []Timestamp `json:"timestamps,omitempty"`
}

// An array of results for the Key Phrases model, if it is enabled.
// See [Key phrases](https://www.assemblyai.com/docs/models/key-phrases) for more information.
type AutoHighlightsResult struct {
	// A temporally-sequential array of Key Phrases
	Results []AutoHighlightResult `json:"results,omitempty"`

	// The status of the Key Phrases model. Either success, or unavailable in the rare case that the model failed.
	Status AudioIntelligenceModelStatus `json:"status,omitempty"`
}

// Chapter of the audio file
type Chapter struct {
	// The starting time, in milliseconds, for the chapter
	End *int64 `json:"end,omitempty"`

	// An ultra-short summary (just a few words) of the content spoken in the chapter
	Gist *string `json:"gist,omitempty"`

	// A single sentence summary of the content spoken during the chapter
	Headline *string `json:"headline,omitempty"`

	// The starting time, in milliseconds, for the chapter
	Start *int64 `json:"start,omitempty"`

	// A one paragraph summary of the content spoken during the chapter
	Summary *string `json:"summary,omitempty"`
}

type ContentSafetyLabel struct {
	// The confidence score for the topic being discussed, from 0 to 1
	Confidence *float64 `json:"confidence,omitempty"`

	// The label of the sensitive topic
	Label *string `json:"label,omitempty"`

	// How severely the topic is discussed in the section, from 0 to 1
	Severity *float64 `json:"severity,omitempty"`
}

type ContentSafetyLabelResult struct {
	// An array of safety labels, one per sensitive topic that was detected in the section
	Labels []ContentSafetyLabel `json:"labels,omitempty"`

	// The sentence index at which the section ends
	SentencesIDxEnd *int64 `json:"sentences_idx_end,omitempty"`

	// The sentence index at which the section begins
	SentencesIDxStart *int64 `json:"sentences_idx_start,omitempty"`

	// The transcript of the section flagged by the Content Moderation model
	Text *string `json:"text,omitempty"`

	// Timestamp information for the section
	Timestamp Timestamp `json:"timestamp,omitempty"`
}

// An array of results for the Content Moderation model, if it is enabled.
// See [Content moderation](https://www.assemblyai.com/docs/models/content-moderation) for more information.
type ContentSafetyLabelsResult struct {
	Results []ContentSafetyLabelResult `json:"results,omitempty"`

	// A summary of the Content Moderation severity results for the entire audio file
	SeverityScoreSummary map[string]SeverityScoreSummary `json:"severity_score_summary,omitempty"`

	// The status of the Content Moderation model. Either success, or unavailable in the rare case that the model failed.
	Status AudioIntelligenceModelStatus `json:"status,omitempty"`

	// A summary of the Content Moderation confidence results for the entire audio file
	Summary map[string]float64 `json:"summary,omitempty"`
}

type CreateRealtimeTemporaryTokenParams struct {
	// The amount of time until the token expires in seconds
	ExpiresIn *int64 `json:"expires_in,omitempty"`
}

// A detected entity
type Entity struct {
	// The ending time, in milliseconds, for the detected entity in the audio file
	End *int64 `json:"end,omitempty"`

	// The type of entity for the detected entity
	EntityType EntityType `json:"entity_type,omitempty"`

	// The starting time, in milliseconds, at which the detected entity appears in the audio file
	Start *int64 `json:"start,omitempty"`

	// The text for the detected entity
	Text *string `json:"text,omitempty"`
}

// The type of entity for the detected entity
type EntityType string

type Error struct {
	// Error message
	Error *string `json:"error,omitempty"`

	Status *string `json:"status,omitempty"`
}

type LeMURActionItemsParams struct {
	LeMURBaseParams
	// How you want the action items to be returned. This can be any text.
	// Defaults to "Bullet Points".
	AnswerFormat *string `json:"answer_format,omitempty"`
}

type LeMURActionItemsResponse struct {
	LeMURBaseResponse
	// The response generated by LeMUR
	Response *string `json:"response,omitempty"`
}

type LeMURBaseParams struct {
	// Context to provide the model. This can be a string or a free-form JSON value.
	Context interface{} `json:"context,omitempty"`

	// The model that is used for the final prompt after compression is performed.
	FinalModel LeMURModel `json:"final_model,omitempty"`

	// Custom formatted transcript data. Maximum size is the context limit of the selected model, which defaults to 100000.
	// Use either transcript_ids or input_text as input into LeMUR.
	InputText *string `json:"input_text,omitempty"`

	// Max output size in tokens, up to 4000
	MaxOutputSize *int64 `json:"max_output_size,omitempty"`

	// The temperature to use for the model.
	// Higher values result in answers that are more creative, lower values are more conservative.
	// Can be any value between 0.0 and 1.0 inclusive.
	Temperature *float64 `json:"temperature,omitempty"`

	// A list of completed transcripts with text. Up to a maximum of 100 files or 100 hours, whichever is lower.
	// Use either transcript_ids or input_text as input into LeMUR.
	TranscriptIDs []string `json:"transcript_ids,omitempty"`
}

type LeMURBaseResponse struct {
	// The ID of the LeMUR request
	RequestID *string `json:"request_id,omitempty"`

	// The usage numbers for the LeMUR request
	Usage LeMURUsage `json:"usage,omitempty"`
}

// The model that is used for the final prompt after compression is performed.
type LeMURModel string

type LeMURQuestion struct {
	// How you want the answer to be returned. This can be any text. Can't be used with answer_options. Examples: "short sentence", "bullet points"
	AnswerFormat *string `json:"answer_format,omitempty"`

	// What discrete options to return. Useful for precise responses. Can't be used with answer_format. Example: ["Yes", "No"]
	AnswerOptions []string `json:"answer_options,omitempty"`

	// Any context about the transcripts you wish to provide. This can be a string or any object.
	Context interface{} `json:"context,omitempty"`

	// The question you wish to ask. For more complex questions use default model.
	Question *string `json:"question,omitempty"`
}

// An answer generated by LeMUR and its question
type LeMURQuestionAnswer struct {
	// The answer generated by LeMUR
	Answer *string `json:"answer,omitempty"`

	// The question for LeMUR to answer
	Question *string `json:"question,omitempty"`
}

type LeMURQuestionAnswerParams struct {
	LeMURBaseParams
	// A list of questions to ask
	Questions []LeMURQuestion `json:"questions,omitempty"`
}

type LeMURQuestionAnswerResponse struct {
	LeMURBaseResponse
	// The answers generated by LeMUR and their questions
	Response []LeMURQuestionAnswer `json:"response,omitempty"`
}

type LeMURSummaryParams struct {
	LeMURBaseParams
	// How you want the summary to be returned. This can be any text. Examples: "TLDR", "bullet points"
	AnswerFormat *string `json:"answer_format,omitempty"`
}

type LeMURSummaryResponse struct {
	LeMURBaseResponse
	// The response generated by LeMUR
	Response *string `json:"response,omitempty"`
}

type LeMURTaskParams struct {
	// Your text to prompt the model to produce a desired output, including any context you want to pass into the model.
	Prompt *string `json:"prompt,omitempty"`

	LeMURBaseParams
}

type LeMURTaskResponse struct {
	// The response generated by LeMUR.
	Response *string `json:"response,omitempty"`

	LeMURBaseResponse
}

// The usage numbers for the LeMUR request
type LeMURUsage struct {
	// The number of input tokens used by the model
	InputTokens *int64 `json:"input_tokens,omitempty"`

	// The number of output tokens generated by the model
	OutputTokens *int64 `json:"output_tokens,omitempty"`
}

type ListTranscriptParams struct {
	// Get transcripts that were created after this transcript ID
	AfterID *string `url:"after_id,omitempty"`

	// Get transcripts that were created before this transcript ID
	BeforeID *string `url:"before_id,omitempty"`

	// Only get transcripts created on this date
	CreatedOn *string `url:"created_on,omitempty"`

	// Maximum amount of transcripts to retrieve
	Limit *int64 `url:"limit,omitempty"`

	// Filter by transcript status
	Status TranscriptStatus `url:"status,omitempty"`

	// Only get throttled transcripts, overrides the status filter
	ThrottledOnly *bool `url:"throttled_only,omitempty"`
}

// Details of the transcript page. Transcripts are sorted from newest to oldest. The previous URL always points to a page with older transcripts.
type PageDetails struct {
	// The URL used to retrieve the current page of transcripts
	CurrentURL *string `json:"current_url,omitempty"`

	// The number of results this page is limited to
	Limit *int64 `json:"limit,omitempty"`

	// The URL to the next page of transcripts. The next URL always points to a page with newer transcripts.
	NextURL *string `json:"next_url,omitempty"`

	// The URL to the next page of transcripts. The previous URL always points to a page with older transcripts.
	PrevURL *string `json:"prev_url,omitempty"`

	// The actual number of results in the page
	ResultCount *int64 `json:"result_count,omitempty"`
}

type ParagraphsResponse struct {
	AudioDuration *float64 `json:"audio_duration,omitempty"`

	Confidence *float64 `json:"confidence,omitempty"`

	ID *string `json:"id,omitempty"`

	Paragraphs []TranscriptParagraph `json:"paragraphs,omitempty"`
}

// The type of PII to redact
type PIIPolicy string

type PurgeLeMURRequestDataResponse struct {
	// Whether the request data was deleted
	Deleted *bool `json:"deleted,omitempty"`

	// The ID of the deletion request of the LeMUR request
	RequestID *string `json:"request_id,omitempty"`

	// The ID of the LeMUR request to purge the data for
	RequestIDToPurge *string `json:"request_id_to_purge,omitempty"`
}

type RealtimeTemporaryTokenResponse struct {
	// The temporary authentication token for Streaming Speech-to-Text
	Token *string `json:"token,omitempty"`
}

// Controls the filetype of the audio created by redact_pii_audio. Currently supports mp3 (default) and wav. See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
type RedactPIIAudioQuality string

// The notification when the redacted audio is ready.
type RedactedAudioNotification struct {
	RedactedAudioResponse
}

type RedactedAudioResponse struct {
	// The URL of the redacted audio file
	RedactedAudioURL *string `json:"redacted_audio_url,omitempty"`

	// The status of the redacted audio
	Status RedactedAudioStatus `json:"status,omitempty"`
}

// The status of the redacted audio
type RedactedAudioStatus string

type SentencesResponse struct {
	AudioDuration *float64 `json:"audio_duration,omitempty"`

	Confidence *float64 `json:"confidence,omitempty"`

	ID *string `json:"id,omitempty"`

	Sentences []TranscriptSentence `json:"sentences,omitempty"`
}

type Sentiment string

// The result of the Sentiment Analysis model
type SentimentAnalysisResult struct {
	// The confidence score for the detected sentiment of the sentence, from 0 to 1
	Confidence *float64 `json:"confidence,omitempty"`

	// The ending time, in milliseconds, of the sentence
	End *int64 `json:"end,omitempty"`

	// The detected sentiment for the sentence, one of POSITIVE, NEUTRAL, NEGATIVE
	Sentiment Sentiment `json:"sentiment,omitempty"`

	// The speaker of the sentence if [Speaker Diarization](https://www.assemblyai.com/docs/models/speaker-diarization) is enabled, else null
	Speaker *string `json:"speaker,omitempty"`

	// The starting time, in milliseconds, of the sentence
	Start *int64 `json:"start,omitempty"`

	// The transcript of the sentence
	Text *string `json:"text,omitempty"`
}

type SeverityScoreSummary struct {
	High *float64 `json:"high,omitempty"`

	Low *float64 `json:"low,omitempty"`

	Medium *float64 `json:"medium,omitempty"`
}

// The speech model to use for the transcription.
type SpeechModel string

// The replacement logic for detected PII, can be "entity_name" or "hash". See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
type SubstitutionPolicy string

// Format of the subtitles
type SubtitleFormat string

// The model to summarize the transcript
type SummaryModel string

// The type of summary
type SummaryType string

// Timestamp containing a start and end property in milliseconds
type Timestamp struct {
	// The end time in milliseconds
	End *int64 `json:"end,omitempty"`

	// The start time in milliseconds
	Start *int64 `json:"start,omitempty"`
}

// The result of the Topic Detection model, if it is enabled.
// See [Topic Detection](https://www.assemblyai.com/docs/models/topic-detection) for more information.
type TopicDetectionModelResult struct {
	// An array of results for the Topic Detection model
	Results []TopicDetectionResult `json:"results,omitempty"`

	// The status of the Topic Detection model. Either success, or unavailable in the rare case that the model failed.
	Status AudioIntelligenceModelStatus `json:"status,omitempty"`

	// The overall relevance of topic to the entire audio file
	Summary map[string]float64 `json:"summary,omitempty"`
}

// The result of the topic detection model
type TopicDetectionResult struct {
	Labels []struct {
		// The IAB taxonomical label for the label of the detected topic, where > denotes supertopic/subtopic relationship
		Label *string `json:"label,omitempty"`

		// How relevant the detected topic is of a detected topic
		Relevance *float64 `json:"relevance,omitempty"`
	} `json:"labels,omitempty"`

	// The text in the transcript in which a detected topic occurs
	Text *string `json:"text,omitempty"`

	Timestamp Timestamp `json:"timestamp,omitempty"`
}

// A transcript object
type Transcript struct {
	// Deprecated: The acoustic model that was used for the transcript
	AcousticModel *string `json:"acoustic_model,omitempty"`

	// The duration of this transcript object's media file, in seconds
	AudioDuration *float64 `json:"audio_duration,omitempty"`

	// The point in time, in milliseconds, in the file at which the transcription was terminated
	AudioEndAt *int64 `json:"audio_end_at,omitempty"`

	// The point in time, in milliseconds, in the file at which the transcription was started
	AudioStartFrom *int64 `json:"audio_start_from,omitempty"`

	// The URL of the media that was transcribed
	AudioURL *string `json:"audio_url,omitempty"`

	// Whether [Auto Chapters](https://www.assemblyai.com/docs/models/auto-chapters) is enabled, can be true or false
	AutoChapters *bool `json:"auto_chapters,omitempty"`

	// Whether Key Phrases is enabled, either true or false
	AutoHighlights *bool `json:"auto_highlights,omitempty"`

	// An array of results for the Key Phrases model, if it is enabled.
	// See [Key Phrases](https://www.assemblyai.com/docs/models/key-phrases) for more information.
	AutoHighlightsResult AutoHighlightsResult `json:"auto_highlights_result,omitempty"`

	// The word boost parameter value
	BoostParam *string `json:"boost_param,omitempty"`

	// An array of temporally sequential chapters for the audio file
	Chapters []Chapter `json:"chapters,omitempty"`

	// The confidence score for the transcript, between 0.0 (low confidence) and 1.0 (high confidence)
	Confidence *float64 `json:"confidence,omitempty"`

	// Whether [Content Moderation](https://www.assemblyai.com/docs/models/content-moderation) is enabled, can be true or false
	ContentSafety *bool `json:"content_safety,omitempty"`

	// An array of results for the Content Moderation model, if it is enabled.
	// See [Content moderation](https://www.assemblyai.com/docs/models/content-moderation) for more information.
	ContentSafetyLabels ContentSafetyLabelsResult `json:"content_safety_labels,omitempty"`

	// Customize how words are spelled and formatted using to and from values
	CustomSpelling []TranscriptCustomSpelling `json:"custom_spelling,omitempty"`

	// Whether custom topics is enabled, either true or false
	CustomTopics *bool `json:"custom_topics,omitempty"`

	// Transcribe Filler Words, like "umm", in your media file; can be true or false
	Disfluencies *bool `json:"disfluencies,omitempty"`

	// Whether [Dual channel transcription](https://www.assemblyai.com/docs/models/speech-recognition#dual-channel-transcription) was enabled in the transcription request, either true or false
	DualChannel *bool `json:"dual_channel,omitempty"`

	// An array of results for the Entity Detection model, if it is enabled.
	// See [Entity detection](https://www.assemblyai.com/docs/models/entity-detection) for more information.
	Entities []Entity `json:"entities,omitempty"`

	// Whether [Entity Detection](https://www.assemblyai.com/docs/models/entity-detection) is enabled, can be true or false
	EntityDetection *bool `json:"entity_detection,omitempty"`

	// Error message of why the transcript failed
	Error *string `json:"error,omitempty"`

	// Whether [Profanity Filtering](https://www.assemblyai.com/docs/models/speech-recognition#profanity-filtering) is enabled, either true or false
	FilterProfanity *bool `json:"filter_profanity,omitempty"`

	// Whether Text Formatting is enabled, either true or false
	FormatText *bool `json:"format_text,omitempty"`

	// Whether [Topic Detection](https://www.assemblyai.com/docs/models/topic-detection) is enabled, can be true or false
	IABCategories *bool `json:"iab_categories,omitempty"`

	// The result of the Topic Detection model, if it is enabled.
	// See [Topic Detection](https://www.assemblyai.com/docs/models/topic-detection) for more information.
	IABCategoriesResult TopicDetectionModelResult `json:"iab_categories_result,omitempty"`

	// The unique identifier of your transcript
	ID *string `json:"id,omitempty"`

	// The language of your audio file.
	// Possible values are found in [Supported Languages](https://www.assemblyai.com/docs/concepts/supported-languages).
	// The default value is 'en_us'.
	LanguageCode TranscriptLanguageCode `json:"language_code,omitempty"`

	// Whether [Automatic language detection](https://www.assemblyai.com/docs/models/speech-recognition#automatic-language-detection) is enabled, either true or false
	LanguageDetection *bool `json:"language_detection,omitempty"`

	// Deprecated: The language model that was used for the transcript
	LanguageModel *string `json:"language_model,omitempty"`

	// Whether Automatic Punctuation is enabled, either true or false
	Punctuate *bool `json:"punctuate,omitempty"`

	// Whether [PII Redaction](https://www.assemblyai.com/docs/models/pii-redaction) is enabled, either true or false
	RedactPII *bool `json:"redact_pii,omitempty"`

	// Whether a redacted version of the audio file was generated,
	// either true or false. See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more information.
	RedactPIIAudio *bool `json:"redact_pii_audio,omitempty"`

	// The audio quality of the PII-redacted audio file, if redact_pii_audio is enabled.
	// See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more information.
	RedactPIIAudioQuality RedactPIIAudioQuality `json:"redact_pii_audio_quality,omitempty"`

	// The list of PII Redaction policies that were enabled, if PII Redaction is enabled.
	// See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more information.
	RedactPIIPolicies []PIIPolicy `json:"redact_pii_policies,omitempty"`

	// The replacement logic for detected PII, can be "entity_type" or "hash". See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
	RedactPIISub SubstitutionPolicy `json:"redact_pii_sub,omitempty"`

	// Whether [Sentiment Analysis](https://www.assemblyai.com/docs/models/sentiment-analysis) is enabled, can be true or false
	SentimentAnalysis *bool `json:"sentiment_analysis,omitempty"`

	// An array of results for the Sentiment Analysis model, if it is enabled.
	// See [Sentiment Analysis](https://www.assemblyai.com/docs/models/sentiment-analysis) for more information.
	SentimentAnalysisResults []SentimentAnalysisResult `json:"sentiment_analysis_results,omitempty"`

	// Whether [Speaker diarization](https://www.assemblyai.com/docs/models/speaker-diarization) is enabled, can be true or false
	SpeakerLabels *bool `json:"speaker_labels,omitempty"`

	// Tell the speaker label model how many speakers it should attempt to identify, up to 10. See [Speaker diarization](https://www.assemblyai.com/docs/models/speaker-diarization) for more details.
	SpeakersExpected *int64 `json:"speakers_expected,omitempty"`

	// The speech model used for the transcription. When `null`, the default model is used.
	SpeechModel SpeechModel `json:"speech_model,omitempty"`

	// Defaults to null. Reject audio files that contain less than this fraction of speech.
	// Valid values are in the range [0, 1] inclusive.
	SpeechThreshold *float64 `json:"speech_threshold,omitempty"`

	// Deprecated: Whether speed boost is enabled
	SpeedBoost *bool `json:"speed_boost,omitempty"`

	// The status of your transcript. Possible values are queued, processing, completed, or error.
	Status TranscriptStatus `json:"status,omitempty"`

	// Whether [Summarization](https://www.assemblyai.com/docs/models/summarization) is enabled, either true or false
	Summarization *bool `json:"summarization,omitempty"`

	// The generated summary of the media file, if [Summarization](https://www.assemblyai.com/docs/models/summarization) is enabled
	Summary *string `json:"summary,omitempty"`

	// The Summarization model used to generate the summary,
	// if [Summarization](https://www.assemblyai.com/docs/models/summarization) is enabled
	SummaryModel *string `json:"summary_model,omitempty"`

	// The type of summary generated, if [Summarization](https://www.assemblyai.com/docs/models/summarization) is enabled
	SummaryType *string `json:"summary_type,omitempty"`

	// The textual transcript of your media file
	Text *string `json:"text,omitempty"`

	// True while a request is throttled and false when a request is no longer throttled
	Throttled *bool `json:"throttled,omitempty"`

	// The list of custom topics provided if custom topics is enabled
	Topics []string `json:"topics,omitempty"`

	// When dual_channel or speaker_labels is enabled, a list of turn-by-turn utterance objects.
	// See [Speaker diarization](https://www.assemblyai.com/docs/models/speaker-diarization) for more information.
	Utterances []TranscriptUtterance `json:"utterances,omitempty"`

	// Whether webhook authentication details were provided
	WebhookAuth *bool `json:"webhook_auth,omitempty"`

	// The header name to be sent with the transcript completed or failed webhook requests
	WebhookAuthHeaderName *string `json:"webhook_auth_header_name,omitempty"`

	// The status code we received from your server when delivering the transcript completed or failed webhook request, if a webhook URL was provided
	WebhookStatusCode *int64 `json:"webhook_status_code,omitempty"`

	// The URL to which we send webhook requests.
	// We sends two different types of webhook requests.
	// One request when a transcript is completed or failed, and one request when the redacted audio is ready if redact_pii_audio is enabled.
	WebhookURL *string `json:"webhook_url,omitempty"`

	// The list of custom vocabulary to boost transcription probability for
	WordBoost []string `json:"word_boost,omitempty"`

	// An array of temporally-sequential word objects, one for each word in the transcript.
	// See [Speech recognition](https://www.assemblyai.com/docs/models/speech-recognition) for more information.
	Words []TranscriptWord `json:"words,omitempty"`
}

// The word boost parameter value
type TranscriptBoostParam string

// Object containing words or phrases to replace, and the word or phrase to replace with
type TranscriptCustomSpelling struct {
	// Words or phrases to replace
	From []string `json:"from,omitempty"`

	// Word or phrase to replace with
	To *string `json:"to,omitempty"`
}

// The language of your audio file. Possible values are found in [Supported Languages](https://www.assemblyai.com/docs/concepts/supported-languages).
// The default value is 'en_us'.
type TranscriptLanguageCode string

// A list of transcripts. Transcripts are sorted from newest to oldest. The previous URL always points to a page with older transcripts.
type TranscriptList struct {
	PageDetails PageDetails `json:"page_details,omitempty"`

	Transcripts []TranscriptListItem `json:"transcripts,omitempty"`
}

type TranscriptListItem struct {
	AudioURL *string `json:"audio_url,omitempty"`

	Completed *string `json:"completed,omitempty"`

	Created *string `json:"created,omitempty"`

	// Error message of why the transcript failed
	Error *string `json:"error,omitempty"`

	ID *string `json:"id,omitempty"`

	ResourceURL *string `json:"resource_url,omitempty"`

	Status TranscriptStatus `json:"status,omitempty"`
}

// The parameters for creating a transcript
type TranscriptOptionalParams struct {
	// The point in time, in milliseconds, to stop transcribing in your media file
	AudioEndAt *int64 `json:"audio_end_at,omitempty"`

	// The point in time, in milliseconds, to begin transcribing in your media file
	AudioStartFrom *int64 `json:"audio_start_from,omitempty"`

	// Enable [Auto Chapters](https://www.assemblyai.com/docs/models/auto-chapters), can be true or false
	AutoChapters *bool `json:"auto_chapters,omitempty"`

	// Enable Key Phrases, either true or false
	AutoHighlights *bool `json:"auto_highlights,omitempty"`

	// The word boost parameter value
	BoostParam TranscriptBoostParam `json:"boost_param,omitempty"`

	// Enable [Content Moderation](https://www.assemblyai.com/docs/models/content-moderation), can be true or false
	ContentSafety *bool `json:"content_safety,omitempty"`

	// The confidence threshold for the Content Moderation model. Values must be between 25 and 100.
	ContentSafetyConfidence *int64 `json:"content_safety_confidence,omitempty"`

	// Customize how words are spelled and formatted using to and from values
	CustomSpelling []TranscriptCustomSpelling `json:"custom_spelling,omitempty"`

	// Enable custom topics, either true or false
	CustomTopics *bool `json:"custom_topics,omitempty"`

	// Transcribe Filler Words, like "umm", in your media file; can be true or false
	Disfluencies *bool `json:"disfluencies,omitempty"`

	// Enable [Dual Channel](https://www.assemblyai.com/docs/models/speech-recognition#dual-channel-transcription) transcription, can be true or false.
	DualChannel *bool `json:"dual_channel,omitempty"`

	// Enable [Entity Detection](https://www.assemblyai.com/docs/models/entity-detection), can be true or false
	EntityDetection *bool `json:"entity_detection,omitempty"`

	// Filter profanity from the transcribed text, can be true or false
	FilterProfanity *bool `json:"filter_profanity,omitempty"`

	// Enable Text Formatting, can be true or false
	FormatText *bool `json:"format_text,omitempty"`

	// Enable [Topic Detection](https://www.assemblyai.com/docs/models/topic-detection), can be true or false
	IABCategories *bool `json:"iab_categories,omitempty"`

	// The language of your audio file. Possible values are found in [Supported Languages](https://www.assemblyai.com/docs/concepts/supported-languages).
	// The default value is 'en_us'.
	LanguageCode TranscriptLanguageCode `json:"language_code,omitempty"`

	// Enable [Automatic language detection](https://www.assemblyai.com/docs/models/speech-recognition#automatic-language-detection), either true or false.
	LanguageDetection *bool `json:"language_detection,omitempty"`

	// Enable Automatic Punctuation, can be true or false
	Punctuate *bool `json:"punctuate,omitempty"`

	// Redact PII from the transcribed text using the Redact PII model, can be true or false
	RedactPII *bool `json:"redact_pii,omitempty"`

	// Generate a copy of the original media file with spoken PII "beeped" out, can be true or false. See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
	RedactPIIAudio *bool `json:"redact_pii_audio,omitempty"`

	// Controls the filetype of the audio created by redact_pii_audio. Currently supports mp3 (default) and wav. See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
	RedactPIIAudioQuality RedactPIIAudioQuality `json:"redact_pii_audio_quality,omitempty"`

	// The list of PII Redaction policies to enable. See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
	RedactPIIPolicies []PIIPolicy `json:"redact_pii_policies,omitempty"`

	// The replacement logic for detected PII, can be "entity_type" or "hash". See [PII redaction](https://www.assemblyai.com/docs/models/pii-redaction) for more details.
	RedactPIISub SubstitutionPolicy `json:"redact_pii_sub,omitempty"`

	// Enable [Sentiment Analysis](https://www.assemblyai.com/docs/models/sentiment-analysis), can be true or false
	SentimentAnalysis *bool `json:"sentiment_analysis,omitempty"`

	// Enable [Speaker diarization](https://www.assemblyai.com/docs/models/speaker-diarization), can be true or false
	SpeakerLabels *bool `json:"speaker_labels,omitempty"`

	// Tells the speaker label model how many speakers it should attempt to identify, up to 10. See [Speaker diarization](https://www.assemblyai.com/docs/models/speaker-diarization) for more details.
	SpeakersExpected *int64 `json:"speakers_expected,omitempty"`

	// The speech model to use for the transcription. When `null`, the default model is used.
	SpeechModel SpeechModel `json:"speech_model,omitempty"`

	// Reject audio files that contain less than this fraction of speech.
	// Valid values are in the range [0, 1] inclusive.
	SpeechThreshold *float64 `json:"speech_threshold,omitempty"`

	// Enable [Summarization](https://www.assemblyai.com/docs/models/summarization), can be true or false
	Summarization *bool `json:"summarization,omitempty"`

	// The model to summarize the transcript
	SummaryModel SummaryModel `json:"summary_model,omitempty"`

	// The type of summary
	SummaryType SummaryType `json:"summary_type,omitempty"`

	// The list of custom topics
	Topics []string `json:"topics,omitempty"`

	// The header name to be sent with the transcript completed or failed webhook requests
	WebhookAuthHeaderName *string `json:"webhook_auth_header_name,omitempty"`

	// The header value to send back with the transcript completed or failed webhook requests for added security
	WebhookAuthHeaderValue *string `json:"webhook_auth_header_value,omitempty"`

	// The URL to which we send webhook requests. We sends two different types of webhook requests. One request when a transcript is completed or failed, and one request when the redacted audio is ready if redact_pii_audio is enabled.
	WebhookURL *string `json:"webhook_url,omitempty"`

	// The list of custom vocabulary to boost transcription probability for
	WordBoost []string `json:"word_boost,omitempty"`
}

type TranscriptParagraph struct {
	Confidence *float64 `json:"confidence,omitempty"`

	End *int64 `json:"end,omitempty"`

	// The speaker of the sentence if [Speaker Diarization](https://www.assemblyai.com/docs/models/speaker-diarization) is enabled, else null
	Speaker *string `json:"speaker,omitempty"`

	Start *int64 `json:"start,omitempty"`

	Text *string `json:"text,omitempty"`

	Words []TranscriptWord `json:"words,omitempty"`
}

// The parameters for creating a transcript
type TranscriptParams struct {
	// The URL of the audio or video file to transcribe.
	AudioURL *string `json:"audio_url,omitempty"`

	TranscriptOptionalParams
}

// The notification when the transcript status is completed or error.
type TranscriptReadyNotification struct {
	// The status of the transcript. Either completed or error.
	Status TranscriptReadyStatus `json:"status,omitempty"`

	// The ID of the transcript
	TranscriptID *string `json:"transcript_id,omitempty"`
}

// The status of the transcript. Either completed or error.
type TranscriptReadyStatus string

type TranscriptSentence struct {
	Confidence *float64 `json:"confidence,omitempty"`

	End *int64 `json:"end,omitempty"`

	// The speaker of the sentence if [Speaker Diarization](https://www.assemblyai.com/docs/models/speaker-diarization) is enabled, else null
	Speaker *string `json:"speaker,omitempty"`

	Start *int64 `json:"start,omitempty"`

	Text *string `json:"text,omitempty"`

	Words []TranscriptWord `json:"words,omitempty"`
}

// The status of your transcript. Possible values are queued, processing, completed, or error.
type TranscriptStatus string

type TranscriptUtterance struct {
	// The confidence score for the transcript of this utterance
	Confidence *float64 `json:"confidence,omitempty"`

	// The ending time, in milliseconds, of the utterance in the audio file
	End *int64 `json:"end,omitempty"`

	// The speaker of this utterance, where each speaker is assigned a sequential capital letter - e.g. "A" for Speaker A, "B" for Speaker B, etc.
	Speaker *string `json:"speaker,omitempty"`

	// The starting time, in milliseconds, of the utterance in the audio file
	Start *int64 `json:"start,omitempty"`

	// The text for this utterance
	Text *string `json:"text,omitempty"`

	// The words in the utterance.
	Words []TranscriptWord `json:"words,omitempty"`
}

// The notifications sent to the webhook URL.
type TranscriptWebhookNotification struct{}

type TranscriptWord struct {
	Confidence *float64 `json:"confidence,omitempty"`

	End *int64 `json:"end,omitempty"`

	// The speaker of the sentence if [Speaker Diarization](https://www.assemblyai.com/docs/models/speaker-diarization) is enabled, else null
	Speaker *string `json:"speaker,omitempty"`

	Start *int64 `json:"start,omitempty"`

	Text *string `json:"text,omitempty"`
}

type UploadedFile struct {
	// A URL that points to your audio file, accessible only by AssemblyAI's servers
	UploadURL *string `json:"upload_url,omitempty"`
}

type WordSearchMatch struct {
	// The total amount of times the word is in the transcript
	Count *int64 `json:"count,omitempty"`

	// An array of all index locations for that word within the `words` array of the completed transcript
	Indexes []int64 `json:"indexes,omitempty"`

	// The matched word
	Text *string `json:"text,omitempty"`

	// An array of timestamps
	Timestamps []WordSearchTimestamp `json:"timestamps,omitempty"`
}

type WordSearchResponse struct {
	// The ID of the transcript
	ID *string `json:"id,omitempty"`

	// The matches of the search
	Matches []WordSearchMatch `json:"matches,omitempty"`

	// The total count of all matched instances. For e.g., word 1 matched 2 times, and word 2 matched 3 times, `total_count` will equal 5.
	TotalCount *int64 `json:"total_count,omitempty"`
}

// An array of timestamps structured as [`start_time`, `end_time`] in milliseconds
type WordSearchTimestamp []int64
