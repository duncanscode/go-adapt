package llm

// imported as anthropic
import (
	"context"
	"encoding/json"
	"fmt"
	"go-adapt/internal/content"
	"go-adapt/internal/selection"
	"regexp"
	"strconv"

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
		  systemPrompt: LLMGuidedPrompt,
      }
}

/*   Methods (for now):
  // Returns question ID to ask next
  func (c *LLMClient) SelectNextQuestion(
      questionBank []content.Question,
      answeredHistory []AnswerRecord,
  ) (int, error) */

func (client *LLMClient) SelectNextQuestion(questionBank []content.Question, answeredHistory []selection.AnswerRecord) (int, error){
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
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: client.systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(inputPrompt)),
		},
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)

	//extract ID
	responseText := message.Content[0].Text

	questionID, err := parseQuestionID(responseText)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

func parseQuestionID(response string) (int,error){
	re := regexp.MustCompile(`<next_question_id>\s*(\d+)\s*</next_question_id>`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return 0, fmt.Errorf("could not find question ID in response")
	}
	return strconv.Atoi(matches[1])
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

