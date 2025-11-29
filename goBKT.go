package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
questions := []string{
    "A chemistry tutorial uses spinning 3D molecule animations with background music and flashing text. The cognitive load from the background music is:",

    "A calculus textbook explains derivatives using the formal epsilon-delta definition in the first chapter for complete beginners. The cognitive load from this mathematical complexity is:",

    "A biology lesson asks students to compare a diagram of a cell to their prior knowledge of factories and assembly lines. The cognitive load from making these connections is:",

    "An online course displays the instructor's video, slides, transcript, and chat window simultaneously, all requiring attention. The cognitive load from managing multiple information streams is:",

    "A physics problem requires students to understand both vector mathematics AND apply it to a novel engineering scenario in the same question. The cognitive load from the vector mathematics itself is:",

    "A history lesson provides a graphic organizer that helps students map cause-and-effect relationships between World War I events. The cognitive load from using this organizer is:",

    "A programming tutorial uses red text for errors, yellow for warnings, and green for success, requiring learners to remember this color code while debugging. The cognitive load from the color-coding system is:",

    "A medical training module asks students to construct their own mnemonic devices to remember the cranial nerves. The cognitive load from creating these mnemonics is:",

    "An architecture software tutorial shows each tool button with decorative icons, shadows, gradients, and animations that don't aid function identification. The cognitive load from these visual embellishments is:",

    "A statistics course requires students to simultaneously learn probability notation AND understand the conceptual meaning of probability for the first time. The cognitive load from the notation system itself is:",

    "A language learning app asks students to reflect on grammar patterns they've noticed and articulate the rules in their own words. The cognitive load from this metacognitive reflection is:",

    "A geometry lesson presents the Pythagorean theorem with a proof, worked examples, practice problems, historical context, and real-world applications all on one densely-packed page. The cognitive load from the page layout density is:",
}

answers := []string{
    "Extraneous",
    "Intrinsic",
    "Germane",
    "Extraneous",
    "Intrinsic",
    "Germane",
    "Extraneous",
    "Germane",
    "Extraneous",
    "Extraneous",
    "Germane",
    "Extraneous",
}

	model := initializeBKTModel(0.01,0.1,0.05 ,0.33)

	for i, q := range questions{
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(q,"\n")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == answers[i] {
			model.UpdateCorrect()
		} else {
			model.UpdateIncorrect()
		}
	}
}

type BKTModel struct {
	L0 float64
	T float64
	S float64
	G float64

	answerHistory []bool
	currentKnowledge float64
	knowledgeHistory []float64
}

func initializeBKTModel(l0, t, s, g float64) *BKTModel {

	return &BKTModel{
		L0: l0,
		T: t,
		S: s,
		G: g,
		currentKnowledge: l0,
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

func (bkt *BKTModel) UpdateIncorrect(){
	//probability they knew it beforehand * probability of slip
	var pLn = bkt.currentKnowledge * bkt.S
	//probability that they didn't know * probability that they didn't guess it
	var pLd = (1-bkt.currentKnowledge)*(1-bkt.G)
	//they got it wrong, this is the chance that they might have known it but slipped
	var actual = pLn/(pLn+pLd)
	//probablility they know it = probability they migh thave slipped plus probability they didn't know it * probability they learned it
	//update current knowledge for next question
	bkt.currentKnowledge = actual + ((1-bkt.G)*(bkt.T))
	bkt.knowledgeHistory = append(bkt.knowledgeHistory, bkt.currentKnowledge)
	bkt.answerHistory = append(bkt.answerHistory, false)
	fmt.Printf("Incorrect, knowledge set to = %.2f\n", bkt.currentKnowledge)
}

func (bkt *BKTModel) UpdateCorrect(){
	//probability they knew it and didn't splip
	var pLn = bkt.currentKnowledge * (1-bkt.S)
	//probability that they didn't know it and they guessed
	var pLg = (1-bkt.currentKnowledge) * (bkt.G)
	//=probability they knew it before they got it right
	var actual = pLn/(pLn+pLg)

	//set current knowledge for next question
	bkt.currentKnowledge = actual + ((1-actual)*(bkt.T))
	bkt.knowledgeHistory = append(bkt.knowledgeHistory, bkt.currentKnowledge)
	bkt.answerHistory = append(bkt.answerHistory, true)
	fmt.Printf("Correct, knowledge set to = %.2f\n", bkt.currentKnowledge)

}