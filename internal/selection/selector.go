package selection

import (
	"go-adapt/internal/content"
	"go-adapt/internal/llm"
	"math"
)

// This package selects the next question for the learner

type SelectionResult struct {
	Question           *content.Question
	Feedback           string // LLM feedback, empty for rule-based
	SelectionReasoning string // LLM reasoning for selection, empty for rule-based
}

type Selector interface {
	SelectQuestion(ctx SelectionContext) (*SelectionResult, error)
	PrepareNextQuestion(ctx SelectionContext) error // Prepare next question (LLM analyzes here)
}

type SelectionContext struct {
	PL0 float64
	Answered []int
	History  []content.AnswerRecord
}

//RULE BASED SELECTION

type RuleBased struct {
	questionBank content.QuestionBank
}

func NewRuleBased(bank content.QuestionBank) *RuleBased {
    return &RuleBased{
        questionBank: bank,
    }
}

func (rb *RuleBased) SelectQuestion(ctx SelectionContext) (*SelectionResult, error) {
	allQuestions, err := rb.questionBank.GetAll()
	if err != nil {
		return nil, err
	}

	unanswered := filterUnanswered(allQuestions, ctx.Answered)
	bestQuestion := findClosestDifficulty(unanswered, ctx.PL0)

	return &SelectionResult{
		Question: bestQuestion,
		Feedback: "", // No feedback for rule-based
	}, nil
}

// PrepareNextQuestion is a no-op for rule-based (no pre-computation needed)
func (rb *RuleBased) PrepareNextQuestion(ctx SelectionContext) error {
	return nil // Rule-based doesn't need preparation
}

func findClosestDifficulty(unanswered []content.Question, targetPL float64) *content.Question {
    // implementation
	var closestQuestion = unanswered[0]
	var minDiff = math.Abs(float64(closestQuestion.Metadata.Difficulty) - targetPL)

	for i, u := range unanswered{
		diff := math.Abs(float64(u.Metadata.Difficulty) - targetPL)
		if diff < minDiff{
			minDiff = diff
			closestQuestion = unanswered[i]
		}
	}
	return &closestQuestion
}

// LLM based selector

type LLMSelector struct{
	questionBank content.QuestionBank
	llmClient *llm.LLMClient
	cachedResult *SelectionResult // Cache for next question
}

func NewLLMSelector(qb content.QuestionBank, client *llm.LLMClient) *LLMSelector{
	return & LLMSelector{
		questionBank: qb,
		llmClient: client,
	}
}

func (ls *LLMSelector) SelectQuestion(ctx SelectionContext) (*SelectionResult, error){
	// If we have a cached result, return it
	if ls.cachedResult != nil {
		result := ls.cachedResult
		ls.cachedResult = nil // Clear cache after use
		return result, nil
	}

	// First question - pick easiest without LLM call
	if len(ctx.History) == 0 {
		allQuestions, err := ls.questionBank.GetAll()
		if err != nil {
			return nil, err
		}

		// Find question with difficulty closest to 0.1
		unanswered := filterUnanswered(allQuestions, ctx.Answered)
		firstQuestion := findClosestDifficulty(unanswered, 0.1)

		return &SelectionResult{
			Question: firstQuestion,
			Feedback: "", // No feedback for first question
		}, nil
	}

	// This shouldn't happen in normal flow (cache should be populated)
	// But fallback to LLM call if needed
	allQuestions, err := ls.questionBank.GetAll()
	if err != nil {
		return nil, err
	}

	llmResponse, err := ls.llmClient.SelectNextQuestion(allQuestions, ctx.History)
	if err != nil {
		return nil, err
	}

	question, err := ls.questionBank.GetQuestionByID(llmResponse.QuestionID)
	if err != nil {
		return nil, err
	}

	return &SelectionResult{
		Question: question,
		Feedback: llmResponse.Feedback,
	}, nil
}

// PrepareNextQuestion calls LLM to analyze performance and cache next question
func (ls *LLMSelector) PrepareNextQuestion(ctx SelectionContext) error {
	allQuestions, err := ls.questionBank.GetAll()
	if err != nil {
		return err
	}

	llmResponse, err := ls.llmClient.SelectNextQuestion(allQuestions, ctx.History)
	if err != nil {
		return err
	}

	question, err := ls.questionBank.GetQuestionByID(llmResponse.QuestionID)
	if err != nil {
		return err
	}

	// Cache the result for next SelectQuestion call
	ls.cachedResult = &SelectionResult{
		Question:           question,
		Feedback:           llmResponse.Feedback,
		SelectionReasoning: llmResponse.SelectionReasoning,
	}

	return nil
}

// GetCachedResult returns the cached result without consuming it
func (ls *LLMSelector) GetCachedResult() *SelectionResult {
	return ls.cachedResult
}

// Private helper functions (lowercase)
func filterUnanswered(questions []content.Question, answeredIDs []int) []content.Question {
	var unanswered []content.Question

    for _, q := range questions{
		answered := false
		for _, a := range answeredIDs{
			if q.ID == a{
				answered = true
				break
			}
		}
		if !answered {
			unanswered = append(unanswered, q)
		}
	}
	return unanswered
}

