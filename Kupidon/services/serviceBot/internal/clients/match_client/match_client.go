package clientsMatch

import (
	"fmt"
	"net/http"
)

type HTTPmatchServiseClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPMatchServiseClient(baseURL string) *HTTPmatchServiseClient {
	return &HTTPmatchServiseClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *HTTPmatchServiseClient) LikeUser(fromUserID, toUserID int64) error {
	url := fmt.Sprintf("%s/like/%d/%d", c.baseURL, fromUserID, toUserID)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
