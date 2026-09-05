package answers

import (
	"errors"
	"fmt"
	"layouts-orders-bot/internal/calculator"
	"layouts-orders-bot/internal/entity"
	"layouts-orders-bot/internal/handlers/answers/mocks"
	sessionsrepo "layouts-orders-bot/internal/repo/sessions/inmemory"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	meta = entity.Meta{
		ChatID:   123,
		Username: "some_user",
	}

	sampleProps = map[*entity.State]string{
		entity.StateLayoutType: entity.LayoutTypeArchitectural,
		entity.StatePurpose:    entity.PurposeEducational,
		entity.StateScale:      entity.Scale1To100,
		entity.StateSize:       entity.Size30x30,
		entity.StateMaterial:   entity.MaterialPlastic,
		entity.StateDetails:    entity.DetailsBasic,
		entity.State3DPrint:    entity.Print3DYes,
		entity.StateLandscape:  entity.LandscapeYes,
		entity.StateDrawings:   entity.DrawingsPartial,
		entity.StateDeadline:   entity.Deadline7Days,
		entity.StateDelivery:   entity.DeliveryMoscow,
	}
)

func TestNewSession(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().StartSession(meta.ChatID)

	msg, err := handler.NewSession(meta)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[0])
}

func TestNewSessionWithError(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().StartSession(meta.ChatID).Return(errors.New("some error"))

	msg, err := handler.NewSession(meta)
	require.Error(t, err)
	require.ErrorContains(t, err, "can't start session")
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func TestFinishSession(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetAnswers(meta.ChatID).Return(sampleProps, nil)
	repo.EXPECT().FinishSession(meta.ChatID)

	props, err := handler.FinishSession(meta)
	require.NoError(t, err)
	require.Equal(t, props, sampleProps)
}

func TestFinishSessionWithError(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetAnswers(meta.ChatID).Return(nil, errors.New("another error"))
	repo.EXPECT().FinishSession(meta.ChatID).Return(errors.New("some error"))

	props, err := handler.FinishSession(meta)
	require.Error(t, err)
	require.ErrorContains(t, err, "can't get answers")
	require.ErrorContains(t, err, "another error")
	require.Nil(t, props)
}

func TestMiddleState(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 7

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0])
	repo.EXPECT().NextQuestion(meta.ChatID)

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[questionIdx+1])
}

func TestMiddleStateWithWrongOption(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[0].State.Options[0], meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0])
	repo.EXPECT().NextQuestion(meta.ChatID)

	msg, err = handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[questionIdx+1])
}

func TestMiddleStateWithNumber(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0])
	repo.EXPECT().NextQuestion(meta.ChatID)

	msg, err := handler.HandleAnswer("1", meta)
	require.NoError(t, err)
	require.Equal(t, msg.Text, Questions[questionIdx+1].Text)
}

func TestMiddleStateWithWrongNumber(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil).Times(2)

	msg, err := handler.HandleAnswer("123", meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)

	msg, err = handler.HandleAnswer("0", meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)
}

func TestMiddleStateWithErrorInRepo(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0])
	repo.EXPECT().NextQuestion(meta.ChatID).Return(errors.New("some error"))

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.Error(t, err)
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func TestNoSession(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(0, ErrNoSession)

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoSession)
	require.Nil(t, msg)
}

func TestCantSaveAnswer(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0]).Return(errors.New("can't save answer"))

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.Error(t, err)
	require.ErrorContains(t, err, "can't save answer")
	require.Nil(t, msg)
}

func TestAcceptQuestion(t *testing.T) {
	handler, repo, calc := setup(t)

	questionIdx := len(Questions) - 2
	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(meta.ChatID, Questions[questionIdx].State, Questions[questionIdx].State.Options[0])
	repo.EXPECT().NextQuestion(meta.ChatID)

	repo.EXPECT().GetAnswers(meta.ChatID).Return(sampleProps, nil)

	minCost, maxCost := 12000, 25000
	calc.EXPECT().Calculate(sampleProps).Return(minCost, maxCost)

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.NoError(t, err)
	require.Contains(t, msg.Text, fmt.Sprintf("%d-%d", minCost, maxCost))
}

func TestAcceptOrder(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := len(Questions) - 1

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[0], meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFinishSession)
	require.Nil(t, msg)
}

func TestCancelOrder(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := len(Questions) - 1

	repo.EXPECT().CurQuestion(meta.ChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[questionIdx].State.Options[1], meta)
	require.NoError(t, err)
	require.Equal(t, cancelOrderMsg, msg.Text)
}

func TestFullOrder(t *testing.T) {
	repo := sessionsrepo.New()
	calc := calculator.New()

	handler := NewHandler(repo, calc)
	msg, err := handler.NewSession(meta)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[0])

	answers := make(map[*entity.State]string, len(Questions))

	for i := 0; i < len(Questions)-2; i++ {
		question := Questions[i]

		msg, err := handler.HandleAnswer(question.State.Options[0], meta)
		require.NoError(t, err)
		checkMsg(t, msg, Questions[i+1])

		answers[question.State] = question.State.Options[0]
	}

	lastQuestion := Questions[len(Questions)-2]

	msg, err = handler.HandleAnswer(lastQuestion.State.Options[0], meta)
	require.NoError(t, err)
	require.Contains(t, msg.Text, "Подтвердить заказ?")

	answers[lastQuestion.State] = lastQuestion.State.Options[0]

	msg, err = handler.HandleAnswer(entity.AcceptMessage, meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFinishSession)

	props, err := handler.FinishSession(meta)
	require.NoError(t, err)
	require.Equal(t, answers, props)
}

func setup(t *testing.T) (*Handler, *mocks.MockRepository, *mocks.MockCalculator) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockRepository(ctrl)
	calc := mocks.NewMockCalculator(ctrl)

	return NewHandler(repo, calc), repo, calc
}

func checkMsg(t *testing.T, msg *entity.Message, question *Question) {
	t.Helper()

	require.Equal(t, question.Text, msg.Text)
	require.Equal(t, question.State.GetOptions(), msg.Options)
}
