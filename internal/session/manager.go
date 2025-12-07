package session

import (
	"go-adapt/internal/bkt"
	"go-adapt/internal/content"
	"go-adapt/internal/llm"
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
	answerHistory []content.AnswerRecord
	mode string
	lastUserModel *llm.UserModel // Latest LLM user model (nil for BKT mode)
}

type QuestionResult struct {
	Question           *content.Question
	Feedback           string
	SelectionReasoning string
}

func NewSessionManager(questionBank content.QuestionBank, mode string, llmClient *llm.LLMClient, l0, t, s, g float64) *SessionManager{
	var selector selection.Selector
	if mode == "llm" {
		selector = selection.NewLLMSelector(questionBank, llmClient)
	} else {
		selector = selection.NewRuleBased(questionBank)
	}

	return &SessionManager{
		bktModel: bkt.InitializeBKTModel(l0,t,s,g),
		questionBank: questionBank,
		selector: selector,
		mode: mode,
	}
}

func (sm *SessionManager) GetNextQuestion() (*QuestionResult, error){
	ctx := selection.SelectionContext{
		PL0: sm.bktModel.GetCurrentKnowledge(),
		Answered: sm.answeredIDs,
		History: sm.answerHistory,
	}
	result, err := sm.selector.SelectQuestion(ctx)
	if err != nil {
		return nil, err
	}
	return &QuestionResult{
		Question:           result.Question,
		Feedback:           result.Feedback,
		SelectionReasoning: result.SelectionReasoning,
	}, nil
}

type SubmitAnswerResult struct {
	CurrentKnowledge float64
	Feedback         string
}

func (sm *SessionManager) SubmitAnswer(questionID int, correct bool) *SubmitAnswerResult {
	if sm.mode == "bkt" {
		if !correct {
			sm.bktModel.UpdateIncorrect()
		} else {
			sm.bktModel.UpdateCorrect()
		}
	}

	sm.answeredIDs = append(sm.answeredIDs, questionID)
	sm.answerHistory = append(sm.answerHistory, content.AnswerRecord{
		QuestionID: questionID,
		Correct:    correct,
	})

	// Prepare next question (LLM analyzes performance here)
	ctx := selection.SelectionContext{
		PL0:      sm.bktModel.GetCurrentKnowledge(),
		Answered: sm.answeredIDs,
		History:  sm.answerHistory,
	}

	feedback := ""
	// Only get feedback from PrepareNextQuestion if in LLM mode
	if sm.mode == "llm" {
		sm.selector.PrepareNextQuestion(ctx)
		// Peek at cached result to get feedback and user model without consuming it
		if llmSelector, ok := sm.selector.(*selection.LLMSelector); ok {
			if llmSelector.GetCachedResult() != nil {
				cached := llmSelector.GetCachedResult()
				feedback = cached.Feedback
				sm.lastUserModel = cached.UserModel // Store latest user model
			}
		}
	} else {
		sm.selector.PrepareNextQuestion(ctx)
	}

	knowledge := 0.0
	if sm.mode == "bkt" {
		knowledge = sm.bktModel.GetCurrentKnowledge()
	}

	return &SubmitAnswerResult{
		CurrentKnowledge: knowledge,
		Feedback:         feedback,
	}
}

func (sm *SessionManager) GetAnsweredCount() int{
	return len(sm.answeredIDs)
}

func (sm *SessionManager) GetAnsweredIDs() []int{
	return sm.answeredIDs
}

func (sm *SessionManager) GetCurrentKnowledge() float64{
	return sm.bktModel.GetCurrentKnowledge()
}

func (sm *SessionManager) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Get difficulty history from answered questions
	difficultyHistory := make([]float64, 0, len(sm.answeredIDs))
	for _, qid := range sm.answeredIDs {
		question, err := sm.questionBank.GetQuestionByID(qid)
		if err == nil {
			difficultyHistory = append(difficultyHistory, question.Metadata.Difficulty)
		}
	}

	metrics["difficulty_history"] = difficultyHistory
	metrics["mode"] = sm.mode

	if sm.mode == "bkt" {
		// BKT-specific metrics
		l0, t, s, g := sm.bktModel.GetParameters()

		metrics["knowledge_history"] = sm.bktModel.GetKnowledgeHistory()
		metrics["answer_history"] = sm.bktModel.GetAnswerHistory()
		metrics["current_knowledge"] = sm.bktModel.GetCurrentKnowledge()
		metrics["parameters"] = map[string]float64{
			"l0": l0,
			"t":  t,
			"s":  s,
			"g":  g,
		}
	} else if sm.mode == "llm" {
		// LLM-specific metrics
		if sm.lastUserModel != nil {
			metrics["user_model"] = map[string]float64{
				"knowledge_level":      sm.lastUserModel.KnowledgeLevel,
				"confidence":           sm.lastUserModel.Confidence,
				"learning_rate":        sm.lastUserModel.LearningRate,
				"pattern_consistency":  sm.lastUserModel.PatternConsistency,
				"difficulty_tolerance": sm.lastUserModel.DifficultyTolerance,
			}
		}

		// Also include answer history for LLM mode
		answerHistory := make([]bool, len(sm.answerHistory))
		for i, record := range sm.answerHistory {
			answerHistory[i] = record.Correct
		}
		metrics["answer_history"] = answerHistory
	}

	return metrics
}