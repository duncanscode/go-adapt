package main

import (
	"fmt"
	"go-adapt/internal/content"
	"go-adapt/internal/session"
)

/*
cmd/server/main.go

  Purpose: Simple CLI loop to test the adaptive learning system

  Flow:
  1. Create SessionManager with initial BKT parameters (0.01, 0.1, 0.05, 0.33)
  2. Print welcome message
  3. Loop:
    - Get next question from manager
    - Display question text
    - Get user input (their answer)
    - Compare to correct answer (case-insensitive)
    - Submit answer result to manager
    - Display if correct/incorrect and current knowledge level
    - Continue until all questions answered or user quits
  4. Print final stats (total answered, final knowledge)

  User Input Handling:
  - Accept typed answers
  - Optional: allow "quit" to exit early
  - Trim whitespace, case-insensitive comparison
*/

func main() {
	bank := content.NewStaticBank()
	manager := session.NewSessionManager(bank, 0.01, 0.1, 0.05, 0.33)
	var input string
	fmt.Println("Welcome to the adaptive learning test tool.")

	for{
		question, err := manager.GetNextQuestion()
		if err != nil {
    		break
		}
		fmt.Println(question.Text)
		fmt.Scanln(&input)
		result := learnerAnswer(input, question.Answer)
		manager.SubmitAnswer(question.ID, result)



		fmt.Printf("Your answer was %s. The correct answer was %s. You were %t \n",input, question.Answer, result )


	}

}

func learnerAnswer(input string, answer string) bool{
	if input == answer{
		println("Correct")
		return true
	} else{
		println("Incorrect")
		return false
	}
}