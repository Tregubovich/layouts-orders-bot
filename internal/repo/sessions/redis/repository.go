package redis

import (
	"context"
	"errors"
	"fmt"
	"layouts-orders-bot/internal/entity"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNoSession = errors.New("no session for current chat")
)

type Repository struct {
	client *redis.Client
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func sessionKey(chatID int) string {
	return fmt.Sprintf("session:%d", chatID)
}

func answersKey(chatID int) string {
	return fmt.Sprintf("session:%d:answers", chatID)
}

func (s *Repository) StartSession(chatID int) error {
	ctx := context.Background()

	return s.client.HSet(ctx, sessionKey(chatID), "current_question", 0, "current_message", 0).Err()
}

func (s *Repository) sessionExists(ctx context.Context, chatID int) (bool, error) {
	exists, err := s.client.Exists(ctx, sessionKey(chatID)).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (s *Repository) CurQuestion(chatID int) (int, error) {
	ctx := context.Background()

	value, err := s.client.HGet(ctx, sessionKey(chatID), "current_question").Int()

	if errors.Is(err, redis.Nil) {
		return 0, ErrNoSession
	}
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s *Repository) NextQuestion(chatID int) error {
	ctx := context.Background()

	exists, err := s.client.Exists(ctx, sessionKey(chatID)).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return ErrNoSession
	}

	return s.client.HIncrBy(ctx, sessionKey(chatID), "current_question", 1).Err()
}
func (s *Repository) SetMessageID(chatID int, messageID int) error {
	ctx := context.Background()

	exists, err := s.client.Exists(ctx, sessionKey(chatID)).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return ErrNoSession
	}

	return s.client.HSet(ctx, sessionKey(chatID), "current_message", messageID).Err()
}

func (s *Repository) CurMessageID(chatID int) (int, error) {
	ctx := context.Background()

	value, err := s.client.HGet(ctx, sessionKey(chatID), "current_message").Int()

	if errors.Is(err, redis.Nil) {
		return 0, ErrNoSession
	}
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s *Repository) SaveAnswer(chatID int, key *entity.State, value string) error {
	ctx := context.Background()

	exists, err := s.client.Exists(ctx, sessionKey(chatID)).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return ErrNoSession
	}

	return s.client.HSet(ctx, answersKey(chatID), key.ID, value).Err()
}

func (s *Repository) GetAnswers(chatID int) (map[*entity.State]string, error) {
	ctx := context.Background()

	exists, err := s.client.Exists(ctx, sessionKey(chatID)).Result()
	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, ErrNoSession
	}

	values, err := s.client.HGetAll(ctx, answersKey(chatID)).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[*entity.State]string)
	for key, value := range values {
		state := convertStringToState(key)
		result[state] = value
	}

	return result, nil
}

func convertStringToState(key string) *entity.State {
	switch key {
	case entity.StateLayoutType.ID:
		return entity.StateLayoutType
	case entity.StatePurpose.ID:
		return entity.StatePurpose
	case entity.StateScale.ID:
		return entity.StateScale
	case entity.StateSize.ID:
		return entity.StateSize
	case entity.StateMaterial.ID:
		return entity.StateMaterial
	case entity.StateDetails.ID:
		return entity.StateDetails
	case entity.State3DPrint.ID:
		return entity.State3DPrint
	case entity.StateLandscape.ID:
		return entity.StateLandscape
	case entity.StateDrawings.ID:
		return entity.StateDrawings
	case entity.StateDeadline.ID:
		return entity.StateDeadline
	case entity.StateDelivery.ID:
		return entity.StateDelivery
	default:
		panic("unknown state key")
	}
}

func (s *Repository) FinishSession(chatID int) error {
	ctx := context.Background()

	_, err := s.client.Del(ctx, sessionKey(chatID), answersKey(chatID)).Result()

	return err
}
