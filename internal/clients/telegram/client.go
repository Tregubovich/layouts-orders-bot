package telegram

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
	getUpdatesMethod             = "getUpdates"
	sendMessageMethod            = "sendMessage"
	answerCallbackQueryMethod    = "answerCallbackQuery"
	editMessageReplyMarkupMethod = "editMessageReplyMarkup"
	deleteMessageMethod          = "deleteMessage"
	setMyCommandsMethod          = "setMyCommands"
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
	if token == "" {
		panic("telegram client must have a token")
	}
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

func (c *Client) SendMessage(chatID int, message string, keyboard [][]entity.Option) (int, error) {
	q := url.Values{}
	q.Set("chat_id", strconv.Itoa(chatID))
	q.Set("text", message)

	if keyboard != nil {
		data, err := json.Marshal(FromArrayToMarkup(keyboard))
		if err != nil {
			return 0, fmt.Errorf("can't marshal reply query: %w", err)
		}
		q.Set("reply_markup", string(data))
	}

	data, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return 0, fmt.Errorf("can't send message: %w", err)
	}

	var res SendMessageResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res.Result.ID, nil
}

func (c *Client) AnswerCallbackQuery(callbackQueryID string) error {
	q := url.Values{}
	q.Set("callback_query_id", callbackQueryID)
	_, err := c.doRequest(answerCallbackQueryMethod, q)
	return err
}

func (c *Client) RemoveKeyboard(chatID, messageID int) error {
	q := url.Values{}
	q.Set("chat_id", strconv.Itoa(chatID))
	q.Set("message_id", strconv.Itoa(messageID))
	q.Set("reply_markup", `{"inline_keyboard":[]}`)

	_, err := c.doRequest(editMessageReplyMarkupMethod, q)
	return err
}

func (c *Client) SetCommands(commands []*entity.Command) error {
	data, err := json.Marshal(FromEntityToCommands(commands))
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Set("commands", string(data))

	_, err = c.doRequest(setMyCommandsMethod, q)
	return err
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram API: %s: %s", resp.Status, body)
	}

	return body, nil
}
