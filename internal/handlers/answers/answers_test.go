package answers

import (
	"errors"
	"fmt"
	"layouts-orders-bot/internal/calculator"
	"layouts-orders-bot/internal/entity"
	"layouts-orders-bot/internal/handlers/answers/mocks"
	sessionsrepo "layouts-orders-bot/internal/repo/inmemory/sessions"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const sampleChatID = 123

var sampleProps = map[entity.State]string{
	"layout_type": "архитектурный",
	"purpose":     "учебный/студенческий",
	"scale":       "1:100",
	"size":        "30x30",
	"material":    "пластик",
	"details":     "базовая",
	"3d_print":    "да",
	"landscape":   "нет",
	"drawings":    "да, с подписанными размерами",
	"deadline":    "7 дней",
	"delivery":    "москва",
}

func TestNewSession(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().StartSession(sampleChatID)

	msg, err := handler.NewSession(sampleChatID, "some_user")
	require.NoError(t, err)
	checkMsg(t, msg, Questions[0])
}

func TestNewSessionWithError(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().StartSession(sampleChatID).Return(errors.New("some error"))

	msg, err := handler.NewSession(sampleChatID, "some_user")
	require.Error(t, err)
	require.ErrorContains(t, err, "can't start session")
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func TestFinishSession(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetAnswers(sampleChatID).Return(sampleProps, nil)
	repo.EXPECT().FinishSession(sampleChatID)

	props, err := handler.FinishSession(sampleChatID)
	require.NoError(t, err)
	require.Equal(t, props, sampleProps)
}

func TestFinishSessionWithError(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetAnswers(sampleChatID).Return(sampleProps, nil)
	repo.EXPECT().FinishSession(sampleChatID).Return(errors.New("some error"))

	props, err := handler.FinishSession(sampleChatID)
	require.Error(t, err)
	require.ErrorContains(t, err, "can't finish session")
	require.ErrorContains(t, err, "some error")
	require.Nil(t, props)
}

func TestMiddleState(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 7

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0])
	repo.EXPECT().NextQuestion(sampleChatID)

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[questionIdx+1])
}

func TestMiddleStateWithWrongOption(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[0].Options[0], sampleChatID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0])
	repo.EXPECT().NextQuestion(sampleChatID)

	msg, err = handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.NoError(t, err)
	checkMsg(t, msg, Questions[questionIdx+1])
}

func TestMiddleStateWithNumber(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0])
	repo.EXPECT().NextQuestion(sampleChatID)

	msg, err := handler.HandleAnswer("1", sampleChatID)
	require.NoError(t, err)
	require.Equal(t, msg.Text, Questions[questionIdx+1].Text)
}

func TestMiddleStateWithWrongNumber(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil).Times(2)

	msg, err := handler.HandleAnswer("123", sampleChatID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)

	msg, err = handler.HandleAnswer("0", sampleChatID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWrongOption)
	require.Nil(t, msg)
}

func TestMiddleStateWithErrorInRepo(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0])
	repo.EXPECT().NextQuestion(sampleChatID).Return(errors.New("some error"))

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.Error(t, err)
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func TestNoSession(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(0, errors.New("no session"))

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.Error(t, err)
	require.ErrorContains(t, err, "no session")
	require.Nil(t, msg)
}

func TestCantSaveAnswer(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := 4

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0]).Return(errors.New("can't save answer"))

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.Error(t, err)
	require.ErrorContains(t, err, "can't save answer")
	require.Nil(t, msg)
}

func TestAcceptQuestion(t *testing.T) {
	handler, repo, calc := setup(t)

	questionIdx := len(Questions) - 2
	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)
	repo.EXPECT().SaveAnswer(sampleChatID, Questions[questionIdx].State, Questions[questionIdx].Options[0])
	repo.EXPECT().NextQuestion(sampleChatID)

	repo.EXPECT().GetAnswers(sampleChatID).Return(sampleProps, nil)

	minCost, maxCost := 12000, 25000
	calc.EXPECT().Calculate(sampleProps).Return(minCost, maxCost)

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.NoError(t, err)
	require.Contains(t, msg.Text, fmt.Sprintf("%d-%d", minCost, maxCost))
}

func TestAcceptOrder(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := len(Questions) - 1

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[0], sampleChatID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFinishSession)
	require.Nil(t, msg)
}

func TestCancelOrder(t *testing.T) {
	handler, repo, _ := setup(t)

	questionIdx := len(Questions) - 1

	repo.EXPECT().CurQuestion(sampleChatID).Return(questionIdx, nil)

	msg, err := handler.HandleAnswer(Questions[questionIdx].Options[1], sampleChatID)
	require.NoError(t, err)
	require.Equal(t, cancelOrderMsg, msg.Text)
}

func TestFullOrder(t *testing.T) {
	repo := sessionsrepo.New()
	calc := calculator.New()

	handler := NewHandler(repo, calc)
	msg, err := handler.NewSession(sampleChatID, "some_user")
	require.NoError(t, err)
	checkMsg(t, msg, Questions[0])

	answers := make(map[entity.State]string, len(Questions))

	for i := 0; i < len(Questions)-2; i++ {
		question := Questions[i]

		msg, err := handler.HandleAnswer(question.Options[0], sampleChatID)
		require.NoError(t, err)
		checkMsg(t, msg, Questions[i+1])

		answers[question.State] = question.Options[0]
	}

	lastQuestion := Questions[len(Questions)-2]

	msg, err = handler.HandleAnswer(lastQuestion.Options[0], sampleChatID)
	require.NoError(t, err)
	require.Contains(t, msg.Text, "Подтвердить заказ?")

	answers[lastQuestion.State] = lastQuestion.Options[0]

	msg, err = handler.HandleAnswer(AcceptMessage, sampleChatID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrFinishSession)

	props, err := handler.FinishSession(sampleChatID)
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
	require.Equal(t, FromStringToOptions(question.Options), msg.Options)
}
