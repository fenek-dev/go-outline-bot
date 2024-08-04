package payment_service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

//go:generate mockery --name HTTPClient --output ./mocks --filename http_client.go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(*Client)

type Client struct {
	BaseUrl    string
	log        *slog.Logger
	httpClient HTTPClient
}

type ErrorResponse struct {
	Message string   `json:"-"`
	Errors  []string `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	return r.Message
}

func NewClient(baseUrl string, opts ...Option) *Client {
	client := &Client{
		BaseUrl:    baseUrl,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *Client) {
		s.log = logger
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(s *Client) {
		s.httpClient = client
	}
}

// Send makes a request to the API, the response body will be
// unmarshalled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err = c.HTTPClient.Do(req)

	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) error {
		return Body.Close()
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResp *ErrorResponse
		data, err = io.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			err := json.Unmarshal(data, errResp)
			if err != nil {
				return err
			}
		}

		return errResp
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *Client) NewRequest(ctx context.Context, method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequestWithContext(ctx, method, c.Options.BaseUrl+"/"+url, buf)
}
