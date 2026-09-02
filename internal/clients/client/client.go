package client

import (
	"encoding/json"
	"fmt"
	"io"
	"layouts-orders-bot/internal/entity"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod          = "getUpdates"
	sendMessageMethod         = "sendMessage"
	answerCallbackQueryMethod = "answerCallbackQuery"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) GetUpdates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, message string, keyboard [][]entity.Option) error {
	q := url.Values{}
	q.Set("chat_id", strconv.Itoa(chatID))
	q.Set("text", message)

	if keyboard != nil {
		data, err := json.Marshal(FromArrayToMarkup(keyboard))
		if err != nil {
			return fmt.Errorf("can't marshal reply query: %w", err)
		}
		q.Set("reply_markup", string(data))
	}

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *Client) AnswerCallbackQuery(callbackQueryID string) error {
	q := url.Values{}
	q.Set("callback_query_id", callbackQueryID)
	_, err := c.doRequest(answerCallbackQueryMethod, q)
	if err != nil {
		return fmt.Errorf("can't answer callback query: %w", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	return body, nil
}
