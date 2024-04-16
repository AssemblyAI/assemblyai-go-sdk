package assemblyai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	// ErrSessionClosed is returned when attempting to write to a closed
	// session.
	ErrSessionClosed = errors.New("session closed")

	// ErrDisconnected is returned when attempting to write to a disconnected
	// client.
	ErrDisconnected = errors.New("client is disconnected")
)

type MessageType string

const (
	MessageTypeSessionBegins      MessageType = "SessionBegins"
	MessageTypeSessionTerminated  MessageType = "SessionTerminated"
	MessageTypePartialTranscript  MessageType = "PartialTranscript"
	MessageTypeFinalTranscript    MessageType = "FinalTranscript"
	MessageTypeSessionInformation MessageType = "SessionInformation"
)

type AudioData struct {
	// Base64 encoded raw audio data
	AudioData string `json:"audio_data,omitempty"`
}

type TerminateSession struct {
	// Set to true to end your real-time session forever
	TerminateSession bool `json:"terminate_session"`
}

type endUtteranceSilenceThreshold struct {
	// Set to true to configure the silence threshold for ending utterances.
	EndUtteranceSilenceThreshold int64 `json:"end_utterance_silence_threshold"`
}

type forceEndUtterance struct {
	// Set to true to manually end the current utterance.
	ForceEndUtterance bool `json:"force_end_utterance"`
}

type RealTimeBaseMessage struct {
	// Describes the type of the message
	MessageType MessageType `json:"message_type"`
}

type RealTimeBaseTranscript struct {
	// End time of audio sample relative to session start, in milliseconds
	AudioEnd int64 `json:"audio_end"`

	// Start time of audio sample relative to session start, in milliseconds
	AudioStart int64 `json:"audio_start"`

	// The confidence score of the entire transcription, between 0 and 1
	Confidence float64 `json:"confidence"`

	// The timestamp for the partial transcript
	Created string `json:"created"`

	// The partial transcript for your audio
	Text string `json:"text"`

	// An array of objects, with the information for each word in the
	// transcription text. Includes the start and end time of the word in
	// milliseconds, the confidence score of the word, and the text, which is
	// the word itself.
	Words []Word `json:"words"`
}

type FinalTranscript struct {
	RealTimeBaseTranscript

	// Describes the type of message
	MessageType MessageType `json:"message_type"`

	// Whether the text is punctuated and cased
	Punctuated bool `json:"punctuated"`

	// Whether the text is formatted, for example Dollar -> $
	TextFormatted bool `json:"text_formatted"`
}

type PartialTranscript struct {
	RealTimeBaseTranscript

	// Describes the type of message
	MessageType MessageType `json:"message_type"`
}

type Word struct {
	// Confidence score of the word
	Confidence float64 `json:"confidence"`

	// End time of the word in milliseconds
	End int64 `json:"end"`

	// Start time of the word in milliseconds
	Start int64 `json:"start"`

	// The word itself
	Text string `json:"text"`
}

var DefaultSampleRate = 16_000

type RealTimeClient struct {
	baseURL *url.URL
	apiKey  string
	token   string

	conn       *websocket.Conn
	httpClient *http.Client

	mtx         sync.RWMutex
	sessionOpen bool

	// done is used to clean up resources when the client disconnects.
	done chan bool

	transcriber *RealTimeTranscriber

	sampleRate int
	encoding   RealTimeEncoding
	wordBoost  []string
}

func (c *RealTimeClient) isSessionOpen() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.sessionOpen
}

func (c *RealTimeClient) setSessionOpen(open bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	c.sessionOpen = open
}

type RealTimeError struct {
	Error string `json:"error"`
}

type RealTimeClientOption func(*RealTimeClient)

// WithRealTimeBaseURL sets the API endpoint used by the client. Mainly used for
// testing.
func WithRealTimeBaseURL(rawurl string) RealTimeClientOption {
	return func(c *RealTimeClient) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// WithRealTimeAuthToken configures the client to authenticate using an
// AssemblyAI API key.
func WithRealTimeAPIKey(apiKey string) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.apiKey = apiKey
	}
}

// WithRealTimeAuthToken configures the client to authenticate using a temporary
// token generated using [CreateTemporaryToken].
func WithRealTimeAuthToken(token string) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.token = token
	}
}

// WithHandler configures the client to use the provided handler to handle
// real-time events.
//
// Deprecated: WithHandler is deprecated. Use [WithRealTimeTranscriber] instead.
func WithHandler(handler RealTimeHandler) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.transcriber = &RealTimeTranscriber{
			OnSessionBegins:     handler.SessionBegins,
			OnSessionTerminated: handler.SessionTerminated,
			OnPartialTranscript: handler.PartialTranscript,
			OnFinalTranscript:   handler.FinalTranscript,
			OnError:             handler.Error,
		}
	}
}

func WithRealTimeTranscriber(transcriber *RealTimeTranscriber) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.transcriber = transcriber
	}
}

// WithRealTimeSampleRate sets the sample rate for the audio data. Default is
// 16000.
func WithRealTimeSampleRate(sampleRate int) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.sampleRate = sampleRate
	}
}

// WithRealTimeWordBoost sets the word boost for the real-time transcription.
func WithRealTimeWordBoost(wordBoost []string) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.wordBoost = wordBoost
	}
}

// RealTimeEncoding is the encoding format for the audio data.
type RealTimeEncoding string

const (
	// PCM signed 16-bit little-endian (default)
	RealTimeEncodingPCMS16LE RealTimeEncoding = "pcm_s16le"

	// PCM Mu-law
	RealTimeEncodingPCMMulaw RealTimeEncoding = "pcm_mulaw"
)

// WithRealTimeEncoding specifies the encoding of the audio data.
func WithRealTimeEncoding(encoding RealTimeEncoding) RealTimeClientOption {
	return func(rtc *RealTimeClient) {
		rtc.encoding = encoding
	}
}

// NewRealTimeClientWithOptions returns a new instance of [RealTimeClient].
func NewRealTimeClientWithOptions(options ...RealTimeClientOption) *RealTimeClient {
	client := &RealTimeClient{
		baseURL: &url.URL{
			Scheme: "wss",
			Host:   "api.assemblyai.com",
			Path:   "/v2/realtime/ws",
		},
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	client.baseURL.RawQuery = client.queryFromOptions()

	return client
}

type SessionBegins struct {
	RealTimeBaseMessage

	// Timestamp when this session will expire
	ExpiresAt string `json:"expires_at"`

	// Describes the type of the message
	MessageType string `json:"message_type"`

	// Unique identifier for the established session
	SessionID string `json:"session_id"`
}

type SessionInformation struct {
	RealTimeBaseMessage

	// The duration of the audio in seconds.
	AudioDurationSeconds float64 `json:"audio_duration_seconds"`
}

type SessionTerminated struct {
	// Describes the type of the message
	MessageType MessageType `json:"message_type"`
}

// Deprecated.
type RealTimeHandler interface {
	SessionBegins(ev SessionBegins)
	SessionTerminated(ev SessionTerminated)
	FinalTranscript(transcript FinalTranscript)
	PartialTranscript(transcript PartialTranscript)
	Error(err error)
}

// NewRealTimeClient returns a new instance of [RealTimeClient] with default
// values. Use [NewRealTimeClientWithOptions] for more configuration options.
func NewRealTimeClient(apiKey string, handler RealTimeHandler) *RealTimeClient {
	return NewRealTimeClientWithOptions(WithRealTimeAPIKey(apiKey), WithHandler(handler))
}

type RealTimeTranscriber struct {
	OnSessionBegins      func(event SessionBegins)
	OnSessionTerminated  func(event SessionTerminated)
	OnSessionInformation func(event SessionInformation)
	OnPartialTranscript  func(event PartialTranscript)
	OnFinalTranscript    func(event FinalTranscript)
	OnError              func(err error)
}

// Connects opens a WebSocket connection and waits for a session to begin.
// Closes the any open WebSocket connection in case of errors.
func (c *RealTimeClient) Connect(ctx context.Context) error {
	header := make(http.Header)

	if c.apiKey != "" {
		header.Set("Authorization", c.apiKey)
	}

	opts := &websocket.DialOptions{
		HTTPHeader: header,
		HTTPClient: &http.Client{},
	}

	conn, _, err := websocket.Dial(ctx, c.baseURL.String(), opts)
	if err != nil {
		return err
	}

	c.conn = conn

	var msg json.RawMessage
	if err := wsjson.Read(ctx, c.conn, &msg); err != nil {
		return err
	}

	var realtimeError RealTimeError
	if err := json.Unmarshal(msg, &realtimeError); err != nil {
		return err
	}
	if realtimeError.Error != "" {
		return errors.New(realtimeError.Error)
	}

	var session SessionBegins
	if err := json.Unmarshal(msg, &session); err != nil {
		return err
	}

	c.setSessionOpen(true)

	if c.transcriber.OnSessionBegins != nil {
		c.transcriber.OnSessionBegins(session)
	}

	c.done = make(chan bool)

	go func() {
		for {
			if !c.isSessionOpen() {
				return
			}

			var msg json.RawMessage

			if err := wsjson.Read(ctx, c.conn, &msg); err != nil {
				if c.transcriber.OnError != nil {
					c.transcriber.OnError(err)
				}
				return
			}

			var messageType struct {
				MessageType MessageType `json:"message_type"`
			}

			if err := json.Unmarshal(msg, &messageType); err != nil {
				if c.transcriber.OnError != nil {
					c.transcriber.OnError(err)
				}
				return
			}

			switch messageType.MessageType {
			case MessageTypeFinalTranscript:
				var transcript FinalTranscript
				if err := json.Unmarshal(msg, &transcript); err != nil {
					if c.transcriber.OnError != nil {
						c.transcriber.OnError(err)
					}
					continue
				}

				if transcript.Text != "" && c.transcriber.OnFinalTranscript != nil {
					c.transcriber.OnFinalTranscript(transcript)
				}
			case MessageTypePartialTranscript:
				var transcript PartialTranscript
				if err := json.Unmarshal(msg, &transcript); err != nil {
					if c.transcriber.OnError != nil {
						c.transcriber.OnError(err)
					}
					continue
				}

				if transcript.Text != "" && c.transcriber.OnPartialTranscript != nil {
					c.transcriber.OnPartialTranscript(transcript)
				}
			case MessageTypeSessionTerminated:
				var session SessionTerminated
				if err := json.Unmarshal(msg, &session); err != nil {
					if c.transcriber.OnError != nil {
						c.transcriber.OnError(err)
					}
					continue
				}

				c.setSessionOpen(false)

				if c.transcriber.OnSessionTerminated != nil {
					c.transcriber.OnSessionTerminated(session)
				}

				c.done <- true
			case MessageTypeSessionInformation:
				var info SessionInformation
				if err := json.Unmarshal(msg, &info); err != nil {
					if c.transcriber.OnError != nil {
						c.transcriber.OnError(err)
					}
					continue
				}

				if c.transcriber.OnSessionInformation != nil {
					c.transcriber.OnSessionInformation(info)
				}
			}
		}
	}()

	return nil
}

func (c *RealTimeClient) queryFromOptions() string {
	values := url.Values{}

	// Temporary token
	if c.token != "" {
		values.Set("token", c.token)
	}

	// Sample rate
	if c.sampleRate > 0 {
		values.Set("sample_rate", strconv.Itoa(c.sampleRate))
	}

	// Encoding
	if c.encoding != "" {
		values.Set("encoding", string(c.encoding))
	}

	// Word boost
	if len(c.wordBoost) > 0 {
		b, _ := json.Marshal(c.wordBoost)
		values.Set("word_boost", string(b))
	}

	// Disable partial transcripts
	if c.transcriber.OnPartialTranscript == nil {
		values.Set("disable_partial_transcripts", "true")
	}

	// Extra session information.
	if c.transcriber.OnSessionInformation != nil {
		values.Set("enable_extra_session_information", "true")
	}

	return values.Encode()
}

// Disconnect sends the terminate_session message and waits for the server to
// send a SessionTerminated message before closing the connection.
func (c *RealTimeClient) Disconnect(ctx context.Context, waitForSessionTermination bool) error {
	terminate := TerminateSession{TerminateSession: true}

	if err := wsjson.Write(ctx, c.conn, terminate); err != nil {
		return err
	}

	if waitForSessionTermination {
		<-c.done
	}

	return c.conn.Close(websocket.StatusNormalClosure, "")
}

// Send sends audio samples to be transcribed.
//
// Expected audio format:
//
// - 16-bit signed integers
// - PCM-encoded
// - Single-channel
func (c *RealTimeClient) Send(ctx context.Context, samples []byte) error {
	if c.conn == nil || !c.isSessionOpen() {
		return ErrSessionClosed
	}

	return c.conn.Write(ctx, websocket.MessageBinary, samples)
}

// ForceEndUtterance manually ends an utterance.
func (c *RealTimeClient) ForceEndUtterance(ctx context.Context) error {
	return wsjson.Write(ctx, c.conn, forceEndUtterance{
		ForceEndUtterance: true,
	})
}

// SetEndUtteranceSilenceThreshold configures the threshold for how long to wait
// before ending an utterance. Default is 700ms.
func (c *RealTimeClient) SetEndUtteranceSilenceThreshold(ctx context.Context, threshold int64) error {
	return wsjson.Write(ctx, c.conn, endUtteranceSilenceThreshold{
		EndUtteranceSilenceThreshold: threshold,
	})
}

// RealTimeService groups operations related to the real-time transcription.
type RealTimeService struct {
	client *Client
}

// CreateTemporaryToken creates a temporary token that can be used to
// authenticate a real-time client.
func (svc *RealTimeService) CreateTemporaryToken(ctx context.Context, expiresIn int64) (*RealtimeTemporaryTokenResponse, error) {
	params := &CreateRealtimeTemporaryTokenParams{
		ExpiresIn: Int64(expiresIn),
	}

	req, err := svc.client.newJSONRequest(ctx, "POST", "/v2/realtime/token", params)
	if err != nil {
		return nil, err
	}

	var tokenResponse RealtimeTemporaryTokenResponse
	resp, err := svc.client.do(req, &tokenResponse)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &tokenResponse, nil
}
