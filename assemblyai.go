package assemblyai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	version              = "1.8.0"
	defaultBaseURLScheme = "https"
	defaultBaseURLHost   = "api.assemblyai.com"
)

// Client manages the communication with the AssemblyAI API.
type Client struct {
	baseURL   *url.URL
	userAgent string
	apiKey    string

	httpClient *http.Client

	Transcripts *TranscriptService
	LeMUR       *LeMURService
	RealTime    *RealTimeService
}

// NewClientWithOptions returns a new configurable AssemblyAI client. If you provide client
// options, they override the default values. Most users will want to use
// [NewClientWithAPIKey].
func NewClientWithOptions(opts ...ClientOption) *Client {
	defaultAPIKey := os.Getenv("ASSEMBLYAI_API_KEY")

	c := &Client{
		baseURL: &url.URL{
			Scheme: defaultBaseURLScheme,
			Host:   defaultBaseURLHost,
		},
		userAgent:  fmt.Sprintf("AssemblyAI/1.0 (sdk=Go/%s)", version),
		httpClient: &http.Client{},
		apiKey:     defaultAPIKey,
	}

	for _, f := range opts {
		f(c)
	}

	c.Transcripts = &TranscriptService{client: c}
	c.LeMUR = &LeMURService{client: c}
	c.RealTime = &RealTimeService{client: c}

	return c
}

// NewClient returns a new authenticated AssemblyAI client.
func NewClient(apiKey string) *Client {
	return NewClientWithOptions(WithAPIKey(apiKey))
}

// ClientOption lets you configure the AssemblyAI client.
type ClientOption func(*Client)

// WithHTTPClient sets the http.Client used for making requests to the API.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the API endpoint used by the client. Mainly used for testing.
func WithBaseURL(rawurl string) ClientOption {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// WithUserAgent sets the user agent used by the client.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithAPIKey sets the API key used for authentication.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)

		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := c.newRequest(ctx, method, path, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	rawurl := c.baseURL.ResolveReference(rel).String()

	req, err := http.NewRequestWithContext(ctx, method, rawurl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiKey)

	return req, err
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apierr APIError

		if err := json.NewDecoder(resp.Body).Decode(&apierr); err != nil {
			return nil, err
		}

		apierr.Status = resp.StatusCode

		return nil, apierr
	}

	if v != nil {
		switch val := v.(type) {
		case *[]byte:
			*val, err = io.ReadAll(resp.Body)
		default:
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}
