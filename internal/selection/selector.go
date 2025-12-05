package selection

import (
	"go-adapt/internal/content"
	"math"
)

// This package selects the next question for the learner

type Selector interface {
	SelectQuestion(ctx SelectionContext) (int, error)
}

type SelectionContext struct {
	PL0 float64
	Answered []int
	History []AnswerRecord
}

//Rule based
type RuleBased struct {
	questionBank content.QuestionBank
}

func (rb *RuleBased) SelectQuestion(ctx SelectionContext) (int, error) {
	allQuestions, err := rb.questionBank.GetAll()
	if err != nil {
    return 0, err
	}

	unanswered := filterUnanswered(allQuestions, ctx.Answered)
	bestQuestion := findClosestDifficulty(unanswered, ctx.PL0)
}

func findClosestDifficulty(unanswered []content.Question, targetPL float64) content.Question {
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
	return closestQuestion
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