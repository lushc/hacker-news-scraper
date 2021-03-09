package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseUrl = "https://hacker-news.firebaseio.com/v0/"
)

// Client is an interface for consuming the Hacker News API
type Client interface {
	TopStories(ctx context.Context) (TopStoriesResponse, error)
	Item(ctx context.Context, id int) (ItemResponse, error)
}

// HNClient is a client for the Hacker News API
type HNClient struct {
	client *http.Client
}

// NotSuccessfulError is used when encountering a non-200 response
type NotSuccessfulError struct {
	statusCode int
}

// TopStoriesResponse is the response returned when fetching top stories
type TopStoriesResponse []int

// ItemResponse is the response returned when fetching an item
type ItemResponse struct {
	ID          int    `json:"id"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Poll        int    `json:"poll"`
	Kids        []int  `json:"kids"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

// NewHNClient creates a new client for querying the Hacker News API
func NewHNClient() *HNClient {
	client := retryablehttp.NewClient()
	client.Logger = nil

	return &HNClient{client: client.StandardClient()}
}

// doGet performs a GET request for the given endpoint
func (h HNClient) doGet(ctx context.Context, endpoint string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, NotSuccessfulError{statusCode: response.StatusCode}
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// TopStories fetches the item IDs of the current top stories
func (h HNClient) TopStories(ctx context.Context) (res TopStoriesResponse, err error) {
	body, err := h.doGet(ctx, "topstories.json")
	if err != nil {
		return nil, fmt.Errorf("fetch top stories: %w", err)
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("unmarshal top stories: %w", err)
	}

	return res, nil
}

// Item fetches an item by the given ID
func (h HNClient) Item(ctx context.Context, id int) (res ItemResponse, err error) {
	body, err := h.doGet(ctx, fmt.Sprintf("item/%d.json", id))
	if err != nil {
		return ItemResponse{}, fmt.Errorf("fetch item %d: %w", id, err)
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return ItemResponse{}, fmt.Errorf("unmarshal item %d: %w", id, err)
	}

	return res, nil
}

func (e NotSuccessfulError) Error() string {
	return fmt.Sprintf("non-2xx status code returned: %d", e.statusCode)
}
