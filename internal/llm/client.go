package llm

// imported as anthropic
import (
	"go-adapt/internal/content"
	"go-adapt/internal/selection"
)


type LLMClient struct {
	apiKey string
}

func newLLMClient(a string) *LLMClient {
	return &LLMClient {
		apiKey: a,
	}
}

/*   Methods (for now):
  // Returns question ID to ask next
  func (c *LLMClient) SelectNextQuestion(
      questionBank []content.Question,
      answeredHistory []AnswerRecord,
  ) (int, error) */

  func (c *LLMClient) SelectNextQuestion(
	questionBank []content.Question,
	answerHistory []selection.AnswerRecord

  )
