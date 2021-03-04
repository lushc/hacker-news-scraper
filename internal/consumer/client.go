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

type Client interface {
	TopStories(ctx context.Context) (TopStoriesResponse, error)
	Item(ctx context.Context, id int) (ItemResponse, error)
}

type HNClient struct {
	client *http.Client
}

type NotSuccessfulError struct {
	statusCode int
}

type TopStoriesResponse []int

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

func NewHNClient() *HNClient {
	client := retryablehttp.NewClient()
	client.Logger = nil

	return &HNClient{client: client.StandardClient()}
}

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
