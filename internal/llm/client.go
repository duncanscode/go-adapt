package llm

// imported as anthropic
import (
	"context"
	"encoding/json"
	"fmt"
	"go-adapt/internal/content"
	"regexp"
	"strconv"

	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
	"github.com/anthropics/anthropic-sdk-go/option"
)

type LLMClient struct {
	*anthropic.Client
	systemPrompt string
}

func NewLLMClient(a string) *LLMClient {
	client := anthropic.NewClient(option.WithAPIKey(a))
      return &LLMClient{
          Client: &client,
		  systemPrompt: LLMGuidedPrompt,
      }
}

/*   Methods (for now):
  // Returns question ID to ask next
  func (c *LLMClient) SelectNextQuestion(
      questionBank []content.Question,
      answeredHistory []AnswerRecord,
  ) (int, error) */

type LLMResponse struct {
	QuestionID         int
	Feedback           string
	SelectionReasoning string
}

func (client *LLMClient) SelectNextQuestion(questionBank []content.Question, answeredHistory []content.AnswerRecord) (*LLMResponse, error){
	questions, _ := toJSONString(questionBank)
	history, _ := toJSONString(answeredHistory)

	inputPrompt := fmt.Sprintf(
		`
		<question_bank>
		%s
		</question_bank>

		<answer_history>
		%s
		</answer_history>

		Select the next question ID.
		`, questions, history)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model: anthropic.ModelClaudeHaiku4_5_20251001,
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: client.systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(inputPrompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM API: %w", err)
	}
	fmt.Printf("%+v\n", message.Content)

	//extract ID and feedback
	responseText := message.Content[0].Text

	questionID, err := parseQuestionID(responseText)
	if err != nil {
		return nil, err
	}

	feedback := parseFeedback(responseText)
	reasoning := parseSelectionReasoning(responseText)

	return &LLMResponse{
		QuestionID:         questionID,
		Feedback:           feedback,
		SelectionReasoning: reasoning,
	}, nil
}

func parseQuestionID(response string) (int,error){
	re := regexp.MustCompile(`<next_question_id>\s*(\d+)\s*</next_question_id>`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return 0, fmt.Errorf("could not find question ID in response")
	}
	return strconv.Atoi(matches[1])
}

func parseFeedback(response string) string {
	re := regexp.MustCompile(`(?s)<feedback>\s*(.*?)\s*</feedback>`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return "" // No feedback found, return empty string
	}
	return matches[1]
}

func parseSelectionReasoning(response string) string {
	re := regexp.MustCompile(`(?s)<selection_reasoning>\s*(.*?)\s*</selection_reasoning>`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return "" // No reasoning found, return empty string
	}
	return matches[1]
}

func toJSONString(data any) (string, error) {
	jsonBytes, err := json.MarshalIndent(
		data, "", "  ",
	)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

