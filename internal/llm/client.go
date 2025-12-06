package llm

// imported as anthropic
import (
	"context"
	"encoding/json"
	"go-adapt/internal/content"
	"go-adapt/internal/selection"

	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
	"github.com/anthropics/anthropic-sdk-go/option"
)

type LLMClient struct {
	*anthropic.Client
	systemPrompt string
}

func newLLMClient(a string) *LLMClient {
	client := anthropic.NewClient(option.WithAPIKey(a))
      return &LLMClient{
          Client: &client,
      }
}

/*   Methods (for now):
  // Returns question ID to ask next
  func (c *LLMClient) SelectNextQuestion(
      questionBank []content.Question,
      answeredHistory []AnswerRecord,
  ) (int, error) */

func (client *LLMClient) SelectNextQuestion(questionBank []content.Question, answeredHistory []selection.AnswerRecord){
	prompt := toJSONString(questionBank)
	history := toJSONString(answeredHistory)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model: anthropic.ModelClaudeHaiku4_5_20251001,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: "hey"},
		},
		Messages: questionBank,
	})
}

func toJSONString(data any) (string, error) {
	jsonBytes, err := json.MarshalIndent(
		data, "", "  "
	)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

