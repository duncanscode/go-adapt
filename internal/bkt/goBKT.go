package bkt


type BKTModel struct {
	L0 float64
	T float64
	S float64
	G float64

	answerHistory []bool
	currentKnowledge float64
	knowledgeHistory []float64
}

func InitializeBKTModel(l0, t, s, g float64) *BKTModel {

	return &BKTModel{
		L0: l0,
		T: t,
		S: s,
		G: g,
		currentKnowledge: l0,
	}
}

func (bkt *BKTModel) GetCurrentKnowledge() float64{
	return bkt.currentKnowledge
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

}