package session

import (
	"go-adapt/internal/bkt"
	"go-adapt/internal/content"
	"go-adapt/internal/selection"
)

/*Purpose: Coordinates BKT model and question selector for a learning session

  Struct: SessionManager
  - Fields:
    - bktModel *bkt.BKTModel - tracks knowledge
    - selector selection.Selector - chooses questions
    - questionBank content.QuestionBank - source of questions
    - answeredIDs []int - tracks which questions answered
    - answerHistory []selection.AnswerRecord - full history

  Methods:
  - NewSessionManager(l0, t, s, g float64) *SessionManager
    - Initialize BKT model with parameters
    - Create StaticBank
    - Create RuleBased selector
    - Return new manager
  - GetNextQuestion() (*content.Question, error)
    - Build SelectionContext from current state (P(L), answered IDs, history)
    - Call selector.SelectQuestion()
    - Retrieve full question from bank by ID
    - Return question
  - SubmitAnswer(questionID int, correct bool) float64
    - Update BKT model (UpdateCorrect or UpdateIncorrect)
    - Add questionID to answeredIDs
    - Add to answerHistory
    - Return new currentKnowledge
  - GetCurrentKnowledge() float64
    - Return bktModel.currentKnowledge
  - GetAnsweredCount() int
    - Return len(answeredIDs)*/

type SessionManager struct{
	bktModel *bkt.BKTModel
	selector selection.Selector
	questionBank content.QuestionBank
	answeredIDs []int
	answerHistory []selection.AnswerRecord
}

func NewSessionManager(l0, t, s, g float64) *SessionManager{
	return &SessionManager{
		bktModel: bkt.InitializeBKTModel(l0,t,s,g),
		questionBank: content.NewStaticBank()
	)
	}

}